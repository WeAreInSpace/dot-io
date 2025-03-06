package dotio

import (
	"log"
	"net"
	"sync"

	"github.com/WeAreInSpace/dot-io/packet"
	"github.com/WeAreInSpace/dot-io/packet/in"
	"github.com/WeAreInSpace/dot-io/packet/out"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
)

/*
 Server Side
*/

type ServerConfig struct {
	Address string
	Network string

	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpListener *net.TCPListener
}

func validateServerConfig(conf *ServerConfig) error {
	if conf.Address == "" {
		conf.Address = "127.0.0.1:42500"
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}

	if conf.Wg == nil {
		conf.Wg = new(sync.WaitGroup)
	}

	if conf.Mx == nil {
		conf.Mx = new(sync.RWMutex)
	}

	if conf.TcpListener == nil {
		addr, err := net.ResolveTCPAddr("tcp", conf.Address)
		if err != nil {
			return err
		}

		listener, err := net.ListenTCP(conf.Network, addr)
		if err != nil {
			return err
		}

		conf.TcpListener = listener
	}

	return nil
}

type Listener struct {
	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpListener *net.TCPListener
	Connection  *connection.ConnectionManager
	Feildkit    *packet.FieldkitManager
}

func NewListener(conf *ServerConfig) (*Listener, error) {
	if conf == nil {
		conf = &ServerConfig{}
	}
	err := validateServerConfig(conf)
	if err != nil {
		return nil, err
	}

	connectionMgr, err := connection.NewConnectionManager()
	if err != nil {
		return nil, err
	}

	feildkitManager := packet.NewFieldkitManager()

	listener := &Listener{
		Wg: conf.Wg,
		Mx: conf.Mx,

		TcpListener: conf.TcpListener,
		Connection:  connectionMgr,
		Feildkit:    feildkitManager,
	}

	return listener, nil
}

func (l *Listener) OnConnection(cbOnConnect func(cdt *connection.ConnectionData)) {
	for {
		conn, err := l.TcpListener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			err := l.Connection.HandleConnection(
				conn,
				func(cdt *connection.ConnectionData) {
					cbOnConnect(cdt)
				},
			)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

/*
 Client Side
*/

type ClientConfig struct {
	Address string
	Network string

	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	Feildkit *packet.FieldkitManager

	TcpConn *net.TCPConn
}

func validateClientConfig(conf *ClientConfig) error {
	if conf.Address == "" {
		conf.Address = "127.0.0.1:42500"
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}

	if conf.Wg == nil {
		conf.Wg = new(sync.WaitGroup)
	}

	if conf.Mx == nil {
		conf.Mx = new(sync.RWMutex)
	}

	if conf.TcpConn == nil {
		addr, err := net.ResolveTCPAddr("tcp", conf.Address)
		if err != nil {
			return err
		}

		conn, err := net.DialTCP(conf.Network, nil, addr)
		if err != nil {
			return err
		}

		conf.TcpConn = conn
	}

	if conf.Feildkit == nil {
		conf.Feildkit = packet.NewFieldkitManager()
	}

	return nil
}

type Connection struct {
	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpConn  *net.TCPConn
	Feildkit *packet.FieldkitManager

	*ConnectionData

	ServerHeader *connection.ServerConnectionHeader
}

func NewConnection(conf *ClientConfig, clientConnectionHeader connection.ClientConnectionHeader) (*Connection, error) {
	if conf == nil {
		conf = &ClientConfig{}
	}
	err := validateClientConfig(conf)
	if err != nil {
		return nil, err
	}

	ipk := in.NewInPacket(conf.TcpConn)
	opk := out.NewOutPacket(conf.TcpConn)

	clientConnectionStatus := &connection.Status{}
	err = packet.TryAndRuturnThis(
		opk.WriteJson(clientConnectionHeader),
		opk.WriteJson(clientConnectionStatus),
	)
	if err != nil {
		return nil, err
	}

	serverConnectionHeader := &connection.ServerConnectionHeader{}
	serverConnectionStatus := &connection.Status{}
	err = packet.TryAndRuturnThis(
		ipk.ReadJsonTo(serverConnectionHeader),
		ipk.ReadJsonTo(serverConnectionStatus),
	)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Wg: conf.Wg,
		Mx: conf.Mx,

		TcpConn:  conf.TcpConn,
		Feildkit: conf.Feildkit,

		ConnectionData: &ConnectionData{
			Ipk: ipk,
			Opk: opk,
		},

		ServerHeader: serverConnectionHeader,
	}
	return connection, nil
}

type ConnectionData struct {
	Ipk *in.InPacket
	Opk *out.OutPacket
}

func (c *Connection) Call(cb func(cdt *ConnectionData)) {
	cb(c.ConnectionData)
}
