package parallel

import (
	"errors"
	"io"
	"math"
	"sync"
)

func NewParallelOut(w io.Writer) *ParallelOut {
	pipes := []*writePipe{}
	errChan := make(chan error)

	return &ParallelOut{
		pipes: pipes,

		errChan: errChan,

		W: w,
	}
}

type ParallelOut struct {
	pipes []*writePipe

	wg      sync.WaitGroup
	errChan chan error

	W io.Writer
}

func (p *ParallelOut) write(r io.Reader) error {
	pipeNumber := len(p.pipes)
	if pipeNumber > math.MaxUint8 {
		return errors.New("there are too many pipes")
	}
	if pipeNumber < 0 {
		return errors.New("there is no pipe")
	}

	pipe := NewWritePipe(
		&Pipe{
			Wg:         &p.wg,
			ErrChan:    p.errChan,
			PipeNumber: uint8(pipeNumber),
		},
		r,
		p.W,
	)

	p.pipes = append(p.pipes, pipe)

	return nil
}

func (p *ParallelOut) WriteBytes(r io.Reader) error {
	err := p.write(r)
	if err != nil {
		return err
	}

	return nil
}

func (p *ParallelOut) Wait() error {
	for _, pipe := range p.pipes {
		pipe.write()
	}
	p.wg.Wait()
	close(p.errChan)

	for err := range p.errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
