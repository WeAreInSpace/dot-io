package pin

import (
	"io"
	"sync"
)

type ParallelIn struct {
	maxPipe int8
	pipeMap map[int8]io.Reader
	mx      sync.Mutex
}

func (pi *ParallelIn) ReadInt32() {

}
