package connection

import (
	"net"
	"sync"

	"github.com/WeAreInSpace/dot-io/data"
	"github.com/WeAreInSpace/dot-io/data/in"
	"github.com/WeAreInSpace/dot-io/data/out"
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

	Ib in.Reader
	Ob out.Writer
}

func (mgr *ConnectionManager) HandleConnection(conn *net.TCPConn, handleFunc func(cdt *ConnectionData)) error {
	ob := out.NewOutbound(conn)
	ib := in.NewInbound(conn)

	clientConnectionHeader := &ClientConnectionHeader{}
	clientConnectionStatus := &Status{}
	err := data.TryAndRuturnThis(
		ib.ReadJsonTo(clientConnectionHeader),
		ib.ReadJsonTo(clientConnectionStatus),
	)
	if err != nil {
		return err
	}

	serverConnectionHeader := &ServerConnectionHeader{}
	serverConnectionStatus := &Status{}
	err = data.TryAndRuturnThis(
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
