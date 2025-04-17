package protocol

import (
	"encoding/json"
	"io"

	"github.com/WeAreInSpace/dot-io/packet"
)

const VERSION = 1
const DEFAULT_KEEP_ALIVE = 10

type ProtocolData struct {
	ProtocolVersion int
	KeepAlivePeriod int64
	Feildkit        *packet.FieldkitManager
}

func (p *ProtocolData) Export(w io.Writer) {
	protocol := ProtocolSchema{
		ProtocolVersion: p.ProtocolVersion,
		KeepAlivePeriod: p.KeepAlivePeriod,
		FeildGroup:      p.Feildkit.Export(),
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(protocol)
}

type ProtocolSchema struct {
	ProtocolVersion int                       `json:"version"`
	KeepAlivePeriod int64                     `json:"keepalive-period"`
	FeildGroup      []packet.FeildGroupSchema `json:"feild-groups"`
}
