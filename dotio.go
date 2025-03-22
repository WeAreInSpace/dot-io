package dotio

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/WeAreInSpace/dot-io/data"
	"github.com/WeAreInSpace/dot-io/data/in"
	"github.com/WeAreInSpace/dot-io/data/out"
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

	if conf.Wg == nil {
		conf.Wg = new(sync.WaitGroup)
	}

	if conf.Mx == nil {
		conf.Mx = new(sync.RWMutex)
	}

	if conf.KeepAlivePeriod <= 0 {
		conf.KeepAlivePeriod = 10
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

	KeepAlivePeriod int64

	TcpListener  *net.TCPListener
	Connection   *connection.ConnectionManager
	Feildkit     *data.FieldkitManager
	ProtocolData *protocol.ProtocolData
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

	feildkitManager := data.NewFieldkitManager()

	protocolData := &protocol.ProtocolData{
		ProtocolVersion: protocol.VERSION,
		KeepAlivePeriod: conf.KeepAlivePeriod,
		FeildGroup:      feildkitManager,
	}

	listener := &Listener{
		Wg: conf.Wg,
		Mx: conf.Mx,

		KeepAlivePeriod: conf.KeepAlivePeriod,

		TcpListener:  conf.TcpListener,
		Connection:   connectionMgr,
		Feildkit:     feildkitManager,
		ProtocolData: protocolData,
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

		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * time.Duration(l.KeepAlivePeriod))

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

	KeepAlivePeriod int64

	Feildkit *data.FieldkitManager

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
		addr, err := net.ResolveTCPAddr("tcp", ":8000")
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

	if conf.Feildkit == nil {
		conf.Feildkit = data.NewFieldkitManager()
	}

	return nil
}

type Connection struct {
	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpConn  *net.TCPConn
	Feildkit *data.FieldkitManager

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

	ib := in.NewInbound(conf.TcpConn)
	ob := out.NewOutbound(conf.TcpConn)

	clientConnectionStatus := &connection.Status{}
	err = data.TryAndRuturnThis(
		ob.WriteJson(clientConnectionHeader),
		ob.WriteJson(clientConnectionStatus),
	)
	if err != nil {
		return nil, err
	}

	serverConnectionHeader := &connection.ServerConnectionHeader{}
	serverConnectionStatus := &connection.Status{}
	err = data.TryAndRuturnThis(
		ib.ReadJsonTo(serverConnectionHeader),
		ib.ReadJsonTo(serverConnectionStatus),
	)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Wg: conf.Wg,
		Mx: conf.Mx,

		TcpConn:  conf.TcpConn,
		Feildkit: conf.Feildkit,

		ConnectionIO: &ConnectionIO{
			Ib: ib,
			Ob: ob,
		},

		ServerHeader: serverConnectionHeader,
	}
	return connection, nil
}

type ConnectionIO struct {
	Ib in.Reader
	Ob out.Writer
}

func (c *Connection) Call(cb func(cd *ConnectionIO)) {
	cb(c.ConnectionIO)
}

func (c *Connection) Close() {
	c.TcpConn.Close()
}

/*
 Local
*/

type LocalIO struct {
	Ib in.Reader
	Ob out.Writer
}

func NewLocal(lio io.ReadWriter) *LocalIO {
	ib := in.NewInbound(lio.(io.Reader))
	ob := out.NewOutbound(lio.(io.Writer))
	localIO := &LocalIO{
		Ib: ib,
		Ob: ob,
	}

	return localIO
}
