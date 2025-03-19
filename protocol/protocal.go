package protocol

const VERSION = 1

func NewProtocolkit() *Protocolkit {
	protocolkit := Protocolkit{}
	return &protocolkit
}

type Protocolkit struct {
}
