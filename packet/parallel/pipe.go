package parallel

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"
)

type PipeHeader struct {
	PipeOrder uint8
}

type Pipe struct {
	Wg         *sync.WaitGroup
	ErrChan    chan error
	PipeNumber uint8

	stopChan     chan bool
	dataPosition int64
	mx           sync.Mutex
}

/*
Convert binary to int32

# Example

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(255))
	ToInt32(buffer)
*/
func binaryToUint8(data io.Reader) (number uint8, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

/*
Convert binary to int64

# Example

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(255))
	ToInt64(buffer)
*/
func binaryToInt64(data io.Reader) (number int64, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}

	return
}

type ReadPipe struct {
	R      io.Reader
	Buffer *bytes.Buffer
	*Pipe
}

func rawToUint8(number uint8) *bytes.Buffer {
	binaryBuffer := new(bytes.Buffer)
	binary.Write(binaryBuffer, binary.BigEndian, number)
	return binaryBuffer
}

func rawToInt64(number int64) *bytes.Buffer {
	binaryBuffer := new(bytes.Buffer)
	binary.Write(binaryBuffer, binary.BigEndian, number)
	return binaryBuffer
}

func NewWritePipe(pipe *Pipe, r io.Reader, w io.Writer) *writePipe {
	pipe.stopChan = make(chan bool)
	return &writePipe{
		W:    w,
		R:    r,
		Pipe: pipe,
	}
}

/*
pipeNumber uint8, dataPosition: uint64, data byte[1]{}
*/

type writePipe struct {
	W io.Writer
	R io.Reader
	*Pipe
}

func (w *writePipe) write() {
	w.Wg.Add(1)
	go func() {
		defer w.Wg.Done()

		for {
			buff := make([]byte, 1)
			_, err := w.R.Read(buff)
			if err != nil {
				if err == io.EOF {
					return
				} else {
					w.ErrChan <- err
					return
				}
			}

			w.mx.Lock()
			_, err = io.Copy(w.W, rawToUint8(w.PipeNumber))
			if err != nil {
				w.ErrChan <- err
				w.mx.Unlock()
				return
			}

			w.dataPosition++
			_, err = io.Copy(w.W, rawToInt64(w.dataPosition))
			if err != nil {
				w.ErrChan <- err
				w.mx.Unlock()
				return
			}

			_, err = w.W.Write(buff)
			if err != nil {
				w.ErrChan <- err
				w.mx.Unlock()
				return
			}
			w.mx.Unlock()
		}
	}()
}
