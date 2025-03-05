package main

import (
	"log"

	dotio "github.com/WeAreInSpace/dot-io"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
)

func main() {
	conn := client()
	conn.Call(
		func(cdt *dotio.ConnectionData) {
			callShow(cdt, "Hello Dot.IO\n")
			callHello(cdt)
		},
	)
}

type ShowDataSchema struct {
	Message string
}

func client() *dotio.Connection {
	connection, err := dotio.NewConnection(nil, connection.ClientConnectionHeader{ProtocolVersion: 1.0})
	if err != nil {
		log.Fatalln(err)
	}

	return connection
}

func callShow(cdt *dotio.ConnectionData, message string) {
	err := cdt.Opk.WriteString("show")
	if err != nil {
		log.Fatalln(err)
	}

	dataToShow := &ShowDataSchema{
		Message: message,
	}
	err = cdt.Opk.WriteJson(dataToShow)
	if err != nil {
		log.Fatalln(err)
	}
}

func callHello(cdt *dotio.ConnectionData) {
	err := cdt.Opk.WriteString("sayHello")
	if err != nil {
		log.Fatalln(err)
	}

	str, err := cdt.Ipk.ReadString()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(str)
}
