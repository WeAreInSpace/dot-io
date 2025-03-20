package out

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

func ToInt32(number int32) *bytes.Buffer {
	binaryBuffer := new(bytes.Buffer)
	binary.Write(binaryBuffer, binary.BigEndian, number)
	return binaryBuffer
}

func ToInt64(number int64) *bytes.Buffer {
	binaryBuffer := new(bytes.Buffer)
	binary.Write(binaryBuffer, binary.BigEndian, number)
	return binaryBuffer
}

func NewOutbound(w io.Writer) *Outbound {
	return &Outbound{
		w: w,
	}
}

type Writer interface {
	Write([]byte) error
	WriteStream(io.Reader) error

	WriteInt32(int32) error
	WriteInt64(int64) error
	WriteString(string) error
	WriteStreamString(int64, io.Reader) error
	WriteJson(any) error
	WriteBytes([]byte) error
	WriteStreamBytes(int64, io.Reader) error
}

type Outbound struct {
	w io.Writer
}

func (opk *Outbound) Write(data []byte) error {
	_, err := opk.w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteStream(data io.Reader) error {
	_, err := io.Copy(opk.w, data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteInt32(data int32) error {
	err := opk.WriteStream(ToInt32(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteInt64(data int64) error {
	err := opk.WriteStream(ToInt64(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteString(data string) error {
	dataLen := len(data)

	err := opk.WriteInt64(int64(dataLen))
	if err != nil {
		return err
	}

	err = opk.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteStreamString(len int64, data io.Reader) error {
	err := opk.WriteInt64(len)
	if err != nil {
		return err
	}

	err = opk.WriteStream(data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteJson(data any) error {
	jsonBuffer := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(jsonBuffer)

	jsonEncoder.Encode(data)

	err := opk.WriteStreamString(int64(jsonBuffer.Len()), jsonBuffer)
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteBytes(data []byte) error {
	err := opk.WriteInt64(int64(len(data)))
	if err != nil {
		return err
	}

	err = opk.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *Outbound) WriteStreamBytes(len int64, data io.Reader) error {
	err := opk.WriteInt64(len)
	if err != nil {
		return err
	}

	err = opk.WriteStream(data)
	if err != nil {
		return err
	}

	return nil
}
