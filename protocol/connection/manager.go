package connection

import (
	"net"
	"sync"

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
	err := ipk.ReadJson(clientConnectionHeader)
	if err != nil {
		return err
	}

	clientConnectionStatus := &Status{}
	err = ipk.ReadJson(clientConnectionStatus)
	if err != nil {
		return err
	}

	connData := &ConnectionData{
		Authentication: clientConnectionHeader.Authentication,
		Conn:           conn,

		Ipk: ipk,
		Opk: opk,
	}

	serverConnectionHeader := &ServerConnectionHeader{}
	err = opk.WriteJson(serverConnectionHeader)
	if err != nil {
		return err
	}

	serverConnectionStatus := &Status{}
	err = opk.WriteJson(serverConnectionStatus)
	if err != nil {
		return err
	}

	go handleFunc(connData)
	return nil
}
