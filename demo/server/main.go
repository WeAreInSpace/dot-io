package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	dotio "github.com/WeAreInSpace/dot-io"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
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

	f := l.Feildkit.New("Command")
	f.ReadString("cmd")

	f = l.Feildkit.New("Show", `[ if (Command.cmd = "show") ]`)
	f.ReadJson("msg", "Message data to show")

	f = l.Feildkit.New("Say Hello", `[ if (Command.cmd = "sayHello") ]`)
	f.WriteString("msg", "Hello World")

	file, err := os.OpenFile("protocol_data/server.json", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	l.ProtocolData.Export(file)

	wg.Add(1)
	go server(l, wg)

	wg.Wait()
}

func server(l *dotio.Listener, wg *sync.WaitGroup) {
	l.OnConnection(
		func(cdt *connection.ConnectionData) {
			for {
				cmd, err := cdt.Ib.ReadString()
				if err, ok := err.(net.Error); ok && err != nil {
					log.Println("Net error: ", err)
					cdt.Conn.Close()
					break
				} else if err != nil {
					log.Println("Error: ", err)
					cdt.Conn.Close()
					break
				}

				switch cmd {
				case "show":
					{
						show(cdt)
					}
				case "sayHello":
					{
						cdt.Ob.WriteString("Hello World")
					}
				case "close":
					{
						continue
					}
				}
			}
			cdt.Conn.Close()
		},
	)

	wg.Done()
}

type ShowDataSchema struct {
	Message string
}

func show(cdt *connection.ConnectionData) {
	dataToShow := &ShowDataSchema{}
	err := cdt.Ib.ReadJsonTo(dataToShow)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s %s\n", cdt.Conn.RemoteAddr().String(), dataToShow.Message)
}
