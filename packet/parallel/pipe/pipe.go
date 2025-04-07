package pipe

import "io"

type Pipe struct {
	Io io.ReadWriter
}
