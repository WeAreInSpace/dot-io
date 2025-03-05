package connection

import (
	"errors"
	"net"
	"sync"

	"github.com/WeAreInSpace/dot-io/packet/in"
	"github.com/WeAreInSpace/dot-io/packet/out"
	"github.com/google/uuid"
)

func NewConnectionManager() (*ConnectionManager, error) {
	devive := make(deviceMap)
	mutex := new(sync.RWMutex)

	connMgr := &ConnectionManager{
		Device: devive,
		Mx:     mutex,
	}

	return connMgr, nil
}

type deviceMap map[uuid.UUID]*ConnectionData

type ConnectionManager struct {
	Device deviceMap

	Mx *sync.RWMutex
}

type ConnectionData struct {
	Authentication ClientAuthentication

	Conn *net.TCPConn

	Ipk *in.InPacket
	Opk *out.OutPacket
}

func (mgr *ConnectionManager) Add(conn *ConnectionData) (*uuid.UUID, error) {
	mgr.Mx.Lock()
	defer mgr.Mx.Unlock()

	var clientUuid uuid.UUID
	for {
		uuid, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		if _, ok := mgr.Device[uuid]; !ok {
			clientUuid = uuid
			break
		}
	}

	mgr.Device[clientUuid] = conn

	return &clientUuid, nil
}

func (mgr *ConnectionManager) Get(uuid uuid.UUID) (*ConnectionData, error) {
	mgr.Mx.RLock()
	defer mgr.Mx.RUnlock()

	if _, ok := mgr.Device[uuid]; !ok {
		return nil, errors.New("connection data not found")
	}

	connectionData := mgr.Device[uuid]
	return connectionData, nil
}

func (mgr *ConnectionManager) Remove(uuid uuid.UUID) error {
	mgr.Mx.RLock()
	defer mgr.Mx.RUnlock()

	if _, ok := mgr.Device[uuid]; !ok {
		return errors.New("connection data not found")
	}

	delete(mgr.Device, uuid)

	return nil
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
	clientUuid, err := mgr.Add(connData)
	if err != nil {
		return err
	}

	serverConnectionHeader := &ServerConnectionHeader{
		ConnectionUUID: clientUuid.String(),
	}
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
