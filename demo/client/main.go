package main

import (
	"log"

	dotio "github.com/WeAreInSpace/dot-io"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
)

func main() {
	conn := client()
	conn.Call(
		func(cdt *dotio.ConnectionIO) {
			callShow(cdt, "Hello Dot.IO\n")
			callHello(cdt)
		},
	)
	conn.Close()
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

func callShow(cdt *dotio.ConnectionIO, message string) {
	err := cdt.Ob.WriteString("show")
	if err != nil {
		log.Fatalln(err)
	}

	dataToShow := &ShowDataSchema{
		Message: message,
	}
	err = cdt.Ob.WriteJson(dataToShow)
	if err != nil {
		log.Fatalln(err)
	}
}

func callHello(cdt *dotio.ConnectionIO) {
	err := cdt.Ob.WriteString("sayHello")
	if err != nil {
		log.Fatalln(err)
	}

	str, err := cdt.Ib.ReadString()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(str)
}
