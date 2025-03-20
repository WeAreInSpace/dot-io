package protocol

import (
	"encoding/json"
	"io"

	"github.com/WeAreInSpace/dot-io/data"
)

const VERSION = 1

type ProtocolData struct {
	ProtocolVersion int
	KeepAlivePeriod int64
	FeildGroup      *data.FieldkitManager
}

func (p *ProtocolData) Export(w io.Writer) {
	protocol := ProtocolSchema{
		ProtocolVersion: p.ProtocolVersion,
		KeepAlivePeriod: p.KeepAlivePeriod,
		FeildGroup:      p.FeildGroup.Export(),
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(protocol)
}

type ProtocolSchema struct {
	ProtocolVersion int                     `json:"version"`
	KeepAlivePeriod int64                   `json:"keepalive-period"`
	FeildGroup      []data.FeildGroupSchema `json:"feild-groups"`
}
