package dotio

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/WeAreInSpace/dot-io/packet"
	"github.com/WeAreInSpace/dot-io/protocol"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
)

/*
 Server Side
*/

var (
	Conf_CbOnNetErr func() = nil
)

type ServerConfig struct {
	Address string
	Network string

	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	KeepAlivePeriod int64

	TcpListener *net.TCPListener
}

func validateServerConfig(conf *ServerConfig) error {
	if conf.Address == "" {
		conf.Address = "127.0.0.1"
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}

	if conf.KeepAlivePeriod <= 0 {
		conf.KeepAlivePeriod = protocol.DEFAULT_KEEP_ALIVE
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

	KeepAlivePeriod int64

	TcpListener *net.TCPListener
}

func NewListener(conf *ServerConfig) (*Listener, error) {
	if conf == nil {
		conf = &ServerConfig{}
	}
	err := validateServerConfig(conf)
	if err != nil {
		return nil, err
	}

	listener := &Listener{
		Wg: conf.Wg,
		Mx: conf.Mx,

		KeepAlivePeriod: conf.KeepAlivePeriod,

		TcpListener: conf.TcpListener,
	}

	return listener, nil
}

type Ctx struct {
	Authentication connection.ClientAuthentication

	Conn *net.TCPConn

	Ib packet.Reader
	Ob packet.Writer
}

func handleConnection(conn *net.TCPConn, handleFunc func(cdt *Ctx)) error {
	ib := packet.NewInbound(conn)
	ob := packet.NewOutbound(conn)
	clientConnectionHeader := &connection.ClientConnectionHeader{}
	clientConnectionStatus := &connection.Status{}
	err := packet.TryAndRuturnThis(
		ib.ReadJsonTo(clientConnectionHeader),
		ib.ReadJsonTo(clientConnectionStatus),
	)
	if err != nil {
		return err
	}

	serverConnectionHeader := &connection.ServerConnectionHeader{}
	serverConnectionStatus := &connection.Status{}
	err = packet.TryAndRuturnThis(
		ob.WriteJson(serverConnectionHeader),
		ob.WriteJson(serverConnectionStatus),
	)
	if err != nil {
		return err
	}

	connData := &Ctx{
		Authentication: clientConnectionHeader.Authentication,
		Conn:           conn,

		Ib: ib,
		Ob: ob,
	}

	go handleFunc(connData)

	return nil
}

func (l *Listener) OnConnection(cbOnConnect func(cdt *Ctx)) {
	for {
		conn, err := l.TcpListener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * time.Duration(l.KeepAlivePeriod))

		go func() {
			err := handleConnection(
				conn,
				func(cdt *Ctx) {
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

	KeepAlivePeriod int64

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

	if conf.KeepAlivePeriod <= 0 {
		conf.KeepAlivePeriod = 10
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

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * time.Duration(conf.KeepAlivePeriod))

		conf.TcpConn = conn
	}

	return nil
}

type Connection struct {
	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpConn *net.TCPConn

	*ConnectionIO

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

	ib := packet.NewInbound(conf.TcpConn)
	ob := packet.NewOutbound(conf.TcpConn)

	clientConnectionStatus := &connection.Status{}
	err = packet.TryAndRuturnThis(
		ob.WriteJson(clientConnectionHeader),
		ob.WriteJson(clientConnectionStatus),
	)
	if err != nil {
		return nil, err
	}

	serverConnectionHeader := &connection.ServerConnectionHeader{}
	serverConnectionStatus := &connection.Status{}
	err = packet.TryAndRuturnThis(
		ib.ReadJsonTo(serverConnectionHeader),
		ib.ReadJsonTo(serverConnectionStatus),
	)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Wg: conf.Wg,
		Mx: conf.Mx,

		TcpConn: conf.TcpConn,

		ConnectionIO: &ConnectionIO{
			Ib: ib,
			Ob: ob,
		},

		ServerHeader: serverConnectionHeader,
	}
	return connection, nil
}

type ConnectionIO struct {
	Ib packet.Reader
	Ob packet.Writer
}

func (c *Connection) Call(cb func(cd *ConnectionIO)) {
	cb(c.ConnectionIO)
}

func (c *Connection) Close() {
	c.TcpConn.Close()
}
