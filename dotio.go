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
		conf.Address = "127.0.0.1"
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
		addr, err := net.ResolveTCPAddr("tcp", ":8000")
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

	TcpConn *net.TCPConn
}

func validateClientConfig(conf *ClientConfig) error {
	if conf.Address == "" {
		conf.Address = "127.0.0.1"
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
		addr, err := net.ResolveTCPAddr("tcp", ":8000")
		if err != nil {
			return err
		}

		conn, err := net.DialTCP(conf.Network, nil, addr)
		if err != nil {
			return err
		}

		conf.TcpConn = conn
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

func NewConnection(conf *ClientConfig, connectionHeader connection.ClientConnectionHeader) (*Connection, error) {
	if conf == nil {
		conf = &ClientConfig{}
	}
	err := validateClientConfig(conf)
	if err != nil {
		return nil, err
	}

	ipk := in.NewInPacket(conf.TcpConn)
	opk := out.NewOutPacket(conf.TcpConn)

	err = opk.WriteJson(connectionHeader)
	if err != nil {
		return nil, err
	}

	connectionStatus := &connection.Status{}
	err = opk.WriteJson(connectionStatus)
	if err != nil {
		return nil, err
	}

	serverConnectionHeader := &connection.ServerConnectionHeader{}
	err = ipk.ReadJson(serverConnectionHeader)
	if err != nil {
		return nil, err
	}

	serverConnectionStatus := &connection.Status{}
	err = ipk.ReadJson(serverConnectionStatus)
	if err != nil {
		return nil, err
	}

	feildkitManager := packet.NewFieldkitManager()

	connection := &Connection{
		Wg: conf.Wg,
		Mx: conf.Mx,

		TcpConn:  conf.TcpConn,
		Feildkit: feildkitManager,

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
