package connection

import (
	"net"
	"sync"

	"github.com/WeAreInSpace/dot-io/packet"
	"github.com/WeAreInSpace/dot-io/packet/in"
	"github.com/WeAreInSpace/dot-io/packet/out"
)

func NewConnectionManager() (*ConnectionManager, error) {
	mutex := new(sync.RWMutex)

	connMgr := &ConnectionManager{
		Mx: mutex,
	}

	return connMgr, nil
}

type ConnectionManager struct {
	Mx *sync.RWMutex
}

type ConnectionData struct {
	Authentication ClientAuthentication

	Conn *net.TCPConn

	Ipk *in.InPacket
	Opk *out.OutPacket
}

func (mgr *ConnectionManager) HandleConnection(conn *net.TCPConn, handleFunc func(cdt *ConnectionData)) error {
	opk := out.NewOutPacket(conn)
	ipk := in.NewInPacket(conn)

	clientConnectionHeader := &ClientConnectionHeader{}
	clientConnectionStatus := &Status{}
	err := packet.TryAndRuturnThis(
		ipk.ReadJsonTo(clientConnectionHeader),
		ipk.ReadJsonTo(clientConnectionStatus),
	)
	if err != nil {
		return err
	}

	serverConnectionHeader := &ServerConnectionHeader{}
	serverConnectionStatus := &Status{}
	err = packet.TryAndRuturnThis(
		opk.WriteJson(serverConnectionHeader),
		opk.WriteJson(serverConnectionStatus),
	)
	if err != nil {
		return err
	}

	connData := &ConnectionData{
		Authentication: clientConnectionHeader.Authentication,
		Conn:           conn,

		Ipk: ipk,
		Opk: opk,
	}
	go handleFunc(connData)

	return nil
}
