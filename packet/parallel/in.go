package parallel

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

func binaryToUint8(data io.Reader) (number uint8, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

func binaryToInt64(data io.Reader) (number int64, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

type ParallelIn struct {
	wg sync.WaitGroup
	mx sync.Mutex

	errChan    chan error
	feildIO    map[uint8]io.Writer
	feildCount uint8
	feildBuff  map[uint8][]byte

	R io.Reader
}

func (pi *ParallelIn) readAll() {
	for {
		go func() {
			pi.mx.Lock()
			feildIndex, err := binaryToUint8(pi.R)
			if err != nil {
				pi.mx.Unlock()
				if err == io.EOF {
					return
				} else {
					pi.errChan <- err
					return
				}
			}

			dataPosition, err := binaryToInt64(pi.R)
			if err != nil {
				pi.mx.Unlock()
				pi.errChan <- err
				return
			}

			fmt.Println(feildIndex, dataPosition)
		}()
	}
}

func (pi *ParallelIn) read(w io.Writer) {
}
