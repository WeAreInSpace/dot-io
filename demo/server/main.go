package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	dotio "github.com/WeAreInSpace/dot-io"
	"github.com/WeAreInSpace/dot-io/packet"
	"github.com/WeAreInSpace/dot-io/protocol"
)

func main() {
	wg := new(sync.WaitGroup)

	l, err := dotio.NewListener(
		&dotio.ServerConfig{
			Address: ":35002",
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	feildkit := packet.NewFieldkit()

	f := feildkit.New("Command")
	f.ReadString("cmd")

	f = feildkit.New("Show", `[ if (Command.cmd = "show") ]`)
	f.ReadJson("msg", "Message data to show")

	f = feildkit.New("Say Hello", `[ if (Command.cmd = "sayHello") ]`)
	f.WriteString("msg", "Hello World")

	file, err := os.OpenFile("protocol_data/server.json", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	protocolData := protocol.ProtocolData{
		ProtocolVersion: protocol.VERSION,
		KeepAlivePeriod: protocol.DEFAULT_KEEP_ALIVE,
		Feildkit:        feildkit,
	}
	protocolData.Export(file)

	wg.Add(1)
	go server(l, wg)

	wg.Wait()
}

func server(l *dotio.Listener, wg *sync.WaitGroup) {
	l.OnConnection(
		func(ctx *dotio.Ctx) {
			for {
				cmd, err := ctx.Ib.ReadString()
				if err, ok := err.(net.Error); ok && err != nil {
					log.Println("Net error: ", err)
					ctx.Conn.Close()
					break
				} else if err != nil {
					log.Println("Error: ", err)
					ctx.Conn.Close()
					break
				}

				switch cmd {
				case "show":
					{
						show(ctx)
					}
				case "sayHello":
					{
						ctx.Ob.WriteString("Hello World")
					}
				case "close":
					{
						continue
					}
				}
			}
			ctx.Conn.Close()
		},
	)

	wg.Done()
}

type ShowDataSchema struct {
	Message string
}

func show(ctx *dotio.Ctx) {
	dataToShow := &ShowDataSchema{}
	err := ctx.Ib.ReadJsonTo(dataToShow)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s %s\n", ctx.Conn.RemoteAddr().String(), dataToShow.Message)
}
