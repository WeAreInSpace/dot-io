package dotio

import (
	"log"
	"net"
	"sync"

	"github.com/WeAreInSpace/dot-io/packet"
	"github.com/WeAreInSpace/dot-io/protocol/connection"
)

type ListenerConfig struct {
	Address string
	Network string

	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpListener *net.TCPListener
}

func validateListenerConfig(conf *ListenerConfig) error {
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

func NewListener(conf *ListenerConfig) (*Listener, error) {
	if conf == nil {
		conf = &ListenerConfig{}
	}
	err := validateListenerConfig(conf)
	if err != nil {
		return nil, err
	}

	connectionMgr, err := connection.NewConnectionManager()
	if err != nil {
		return nil, err
	}

	feildManager := packet.NewFieldManager()

	listener := &Listener{
		Wg: conf.Wg,
		Mx: conf.Mx,

		TcpListener:   conf.TcpListener,
		ConnectionMgr: connectionMgr,
		FeildMgr:      feildManager,
	}

	return listener, nil
}

type Listener struct {
	Wg *sync.WaitGroup
	Mx *sync.RWMutex

	TcpListener   *net.TCPListener
	ConnectionMgr *connection.ConnectionManager
	FeildMgr      *packet.FieldManager
}

func (l *Listener) HandleConnection(cbOnConnect func(cdt *connection.ConnectionData)) {
	for {
		conn, err := l.TcpListener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		go l.ConnectionMgr.HandleConnection(
			conn,
			func(cdt *connection.ConnectionData) {
				cbOnConnect(cdt)
			},
		)
	}
}
