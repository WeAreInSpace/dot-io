package out

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
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

func NewOutPacket(conn *net.TCPConn) *OutPacket {
	return &OutPacket{
		conn: conn,
	}
}

type PacketWriter interface {
	Write(data []byte) error
	WriteStream(data io.Reader) error

	WriteInt32(data int32) error
	WriteInt64(data int64) error
	WriteString(data string) error
	WriteStreamString(len int64, data io.Reader) error
	WriteJson(data any) error
	WriteBytes(data []byte) error
	WriteStreamBytes(len int64, data io.Reader) error
}

type OutPacket struct {
	conn *net.TCPConn
}

func (opk *OutPacket) Write(data []byte) error {
	_, err := opk.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteStream(data io.Reader) error {
	_, err := io.Copy(opk.conn, data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteInt32(data int32) error {
	err := opk.WriteStream(ToInt32(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteInt64(data int64) error {
	err := opk.WriteStream(ToInt64(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteString(data string) error {
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

func (opk *OutPacket) WriteStreamString(len int64, data io.Reader) error {
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

func (opk *OutPacket) WriteJson(data any) error {
	jsonBuffer := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(jsonBuffer)

	jsonEncoder.Encode(data)

	err := opk.WriteStreamString(int64(jsonBuffer.Len()), jsonBuffer)
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteBytes(data []byte) error {
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

func (opk *OutPacket) WriteStreamBytes(len int64, data io.Reader) error {
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
