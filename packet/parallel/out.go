package parallel

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"
)

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

func NewParallelOut(w io.Writer) *ParallelOut {
	errChan := make(chan error)

	return &ParallelOut{
		errChan: errChan,
		W:       w,
	}
}

type ParallelOut struct {
	wg sync.WaitGroup
	mx sync.Mutex

	errChan    chan error
	feildCount uint8

	W io.Writer
}

func (po *ParallelOut) write(r io.Reader) {
	buffer := []byte{}
	for {
		buff := make([]byte, 1)
		_, err := r.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				po.errChan <- err
				break
			}
		}
		buffer = append(buffer, buff...)
	}

	for dataPosition, buff := range buffer {
		po.wg.Add(1)

		go func() {
			defer po.wg.Done()

			po.mx.Lock()
			fakeByteArr := []byte{buff}
			_, err := io.Copy(po.W, rawToUint8(po.feildCount))
			if err != nil {
				po.errChan <- err
				po.mx.Unlock()
				return
			}

			_, err = io.Copy(po.W, rawToInt64(int64(dataPosition)))
			if err != nil {
				po.errChan <- err
				po.mx.Unlock()
				return
			}

			_, err = po.W.Write(fakeByteArr)
			if err != nil {
				po.errChan <- err
				po.mx.Unlock()
				return
			}
			po.mx.Unlock()
		}()
	}
}

func (po *ParallelOut) WriteBytes(r io.Reader) {
	po.write(r)
	po.feildCount++
}

func (p *ParallelOut) Wait() error {
	p.wg.Wait()
	close(p.errChan)

	for err := range p.errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
