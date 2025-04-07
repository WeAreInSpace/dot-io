package connection

import (
	"net"
	"sync"

	"github.com/WeAreInSpace/dot-io/packet"
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

	Ib packet.Reader
	Ob packet.Writer
}

func (mgr *ConnectionManager) HandleConnection(conn *net.TCPConn, handleFunc func(cdt *ConnectionData)) error {
	ob := packet.NewOutbound(conn)
	ib := packet.NewInbound(conn)

	clientConnectionHeader := &ClientConnectionHeader{}
	clientConnectionStatus := &Status{}
	err := packet.TryAndRuturnThis(
		ib.ReadJsonTo(clientConnectionHeader),
		ib.ReadJsonTo(clientConnectionStatus),
	)
	if err != nil {
		return err
	}

	serverConnectionHeader := &ServerConnectionHeader{}
	serverConnectionStatus := &Status{}
	err = packet.TryAndRuturnThis(
		ob.WriteJson(serverConnectionHeader),
		ob.WriteJson(serverConnectionStatus),
	)
	if err != nil {
		return err
	}

	connData := &ConnectionData{
		Authentication: clientConnectionHeader.Authentication,
		Conn:           conn,

		Ib: ib,
		Ob: ob,
	}

	go handleFunc(connData)

	return nil
}
