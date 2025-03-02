package in

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	"github.com/bytedance/sonic"
)

func ToInt32(data *bytes.Buffer) (number int32, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

func ToInt64(data *bytes.Buffer) (number int64, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

func NewInPacket(conn *net.TCPConn) *InPacket {
	return &InPacket{
		conn: conn,
	}
}

type PacketReader interface {
	Read(len int64) ([]byte, error)
	ReadStream(len int64) (*bytes.Buffer, error)

	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
	ReadStreamString() (*bytes.Buffer, error)
	ReadJson(val any) error
	ReadBytes() ([]byte, error)
	ReadStreamBytes() (*bytes.Buffer, error)
}

type InPacket struct {
	conn *net.TCPConn
}

func (ipk *InPacket) Read(len int64) ([]byte, error) {
	err := ipk.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return nil, err
	}

	byteBuffer := make([]byte, len)
	written, err := ipk.conn.Read(byteBuffer)
	if written < int(len) {
		return nil, errors.New("there is no data left")
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, errors.New("read timeout")
	}

	return byteBuffer, nil
}

func (ipk *InPacket) ReadStream(len int64) (*bytes.Buffer, error) {
	err := ipk.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return nil, err
	}

	byteBuffer := new(bytes.Buffer)
	written, err := io.CopyN(byteBuffer, ipk.conn, len)
	if (written < len) || (err != nil) {
		return nil, err
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, errors.New("read timeout")
	}

	return byteBuffer, nil
}

func (ipk *InPacket) ReadInt32() (int32, error) {
	rawData, err := ipk.ReadStream(4)
	if err != nil {
		return 0, err
	}

	number, err := ToInt32(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ipk *InPacket) ReadInt64() (int64, error) {
	rawData, err := ipk.ReadStream(8)
	if err != nil {
		return 0, err
	}

	number, err := ToInt64(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ipk *InPacket) ReadString() (string, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return "", err
	}

	rawData, err := ipk.Read(length)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func (ipk *InPacket) ReadStreamString() (*bytes.Buffer, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	rawData, err := ipk.ReadStream(length)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (ipk *InPacket) ReadJson(val any) error {
	jsonString, err := ipk.ReadStreamString()
	if err != nil {
		return err
	}

	jsonDecoder := sonic.ConfigDefault.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return nil
}

func (ipk *InPacket) ReadBytes() ([]byte, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	byteBuf, err := ipk.Read(length)
	if err != nil {
		return nil, err
	}

	return byteBuf, nil
}

func (ipk *InPacket) ReadStreamBytes() (*bytes.Buffer, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	byteBuf, err := ipk.ReadStream(length)
	if err != nil {
		return nil, err
	}

	return byteBuf, nil
}
