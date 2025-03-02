package out

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	"github.com/bytedance/sonic"
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

type OutPacket struct {
	conn *net.TCPConn
}

func (opk *OutPacket) WriteString(data string) error {
	dataLen := len(data)

	err := opk.WriteInt64(int64(dataLen))
	if err != nil {
		return err
	}

	err = opk.WriteBytes([]byte(data))
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

	_, err = io.Copy(opk.conn, data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteInt32(data int32) error {
	_, err := io.Copy(opk.conn, ToInt32(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteInt64(data int64) error {
	_, err := io.Copy(opk.conn, ToInt64(data))
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteBytes(data []byte) error {
	_, err := opk.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (opk *OutPacket) WriteJson(data any) error {
	jsonBuffer := new(bytes.Buffer)
	jsonEncoder := sonic.ConfigDefault.NewEncoder(jsonBuffer)

	jsonEncoder.Encode(data)

	err := opk.WriteStreamString(int64(jsonBuffer.Len()), jsonBuffer)
	if err != nil {
		return err
	}

	return nil
}
