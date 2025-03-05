package main

import (
	"fmt"
	"log"
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
				if err != nil {
					log.Println("ggg", err)
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
	err := cdt.Ipk.ReadJson(dataToShow)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s %s", cdt.Conn.RemoteAddr().String(), dataToShow.Message)
}
