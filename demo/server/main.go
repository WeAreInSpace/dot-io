package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	dotio "github.com/WeAreInSpace/dot-io"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
)

func main() {
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go server()

	wg.Wait()
}

func server() {
	l, err := dotio.NewListener(nil)
	if err != nil {
		log.Fatalln(err)
	}

	l.OnConnection(
		func(cdt *connection.ConnectionData) {
			for {
				cmd, err := cdt.Ipk.ReadString()
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
						continue
					}
				case "sayHello":
					{
						cdt.Opk.WriteString("Hello World")
						continue
					}
				}
			}
			cdt.Conn.Close()
		},
	)
}

type ShowDataSchema struct {
	Message string
}

func show(cdt *connection.ConnectionData) {
	dataToShow := &ShowDataSchema{}
	err := cdt.Ipk.ReadJsonTo(dataToShow)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s %s", cdt.Conn.RemoteAddr().String(), dataToShow.Message)
}
