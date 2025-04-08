package packet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

func rawToInt32(number int32) *bytes.Buffer {
	binaryBuffer := new(bytes.Buffer)
	binary.Write(binaryBuffer, binary.BigEndian, number)
	return binaryBuffer
}

func rawToInt64(number int64) *bytes.Buffer {
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
	WriteInt32(int32) error
	WriteInt64(int64) error
	WriteString(string) error
	WriteJson(any) error
	WriteBytes([]byte) error
	WriteStreamBytes(int64, io.Reader) error
}

type Outbound struct {
	w io.Writer
}

func (ob *Outbound) write(data []byte) error {
	_, err := ob.w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) writeStream(data io.Reader) error {
	_, err := io.Copy(ob.w, data)
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) WriteInt32(data int32) error {
	err := ob.writeStream(rawToInt32(data))
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) WriteInt64(data int64) error {
	err := ob.writeStream(rawToInt64(data))
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) WriteString(data string) error {
	dataLen := len(data)

	err := ob.WriteInt64(int64(dataLen))
	if err != nil {
		return err
	}

	err = ob.write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) WriteJson(data any) error {
	jsonBuffer := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(jsonBuffer)

	jsonEncoder.Encode(data)

	err := ob.WriteStreamBytes(int64(jsonBuffer.Len()), jsonBuffer)
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) WriteBytes(data []byte) error {
	err := ob.WriteInt64(int64(len(data)))
	if err != nil {
		return err
	}

	err = ob.write(data)
	if err != nil {
		return err
	}

	return nil
}

func (ob *Outbound) WriteStreamBytes(len int64, data io.Reader) error {
	err := ob.WriteInt64(len)
	if err != nil {
		return err
	}

	err = ob.writeStream(data)
	if err != nil {
		return err
	}

	return nil
}
