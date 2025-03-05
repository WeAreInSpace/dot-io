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
	ReadJson() (val any, err error)
	ReadBytes() ([]byte, error)
	ReadStreamBytes() (*bytes.Buffer, error)
}

type PacketReaderTo interface {
	ReadTo(len int64, data []byte) error
	ReadStreamTo(len int64, buffer *bytes.Buffer) error

	ReadInt32To(*int32) error
	ReadInt64To(*int64) error
	ReadStringTo(*string) error
	ReadStreamStringTo(*bytes.Buffer) error
	ReadJsonTo(val any) error
	ReadBytesTo(data []byte) error
	ReadStreamBytesTo(buffer *bytes.Buffer) error
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

func (ipk *InPacket) ReadTo(len int64, data []byte) error {
	readData, err := ipk.Read(len)
	if err != nil {
		return err
	}

	copy(data, readData)

	return nil
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

func (ipk *InPacket) ReadStreamTo(len int64, buffer *bytes.Buffer) error {
	err := ipk.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return err
	}

	written, err := io.CopyN(buffer, ipk.conn, len)
	if (written < len) || (err != nil) {
		return err
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return errors.New("read timeout")
	}

	return nil
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

func (ipk *InPacket) ReadInt32To(data *int32) error {
	readData, err := ipk.ReadInt32()
	if err != nil {
		return err
	}

	*data = readData

	return nil
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

func (ipk *InPacket) ReadInt64To(data *int64) error {
	readData, err := ipk.ReadInt64()
	if err != nil {
		return err
	}

	*data = readData

	return nil
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

func (ipk *InPacket) ReadStringTo(data *string) error {
	readData, err := ipk.ReadString()
	if err != nil {
		return err
	}

	*data = readData

	return nil
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

func (ipk *InPacket) ReadStreamStringTo(buffer *bytes.Buffer) error {
	length, err := ipk.ReadInt64()
	if err != nil {
		return err
	}

	err = ipk.ReadStreamTo(length, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (ipk *InPacket) ReadJson() (val any, err error) {
	jsonString, err := ipk.ReadStreamString()
	if err != nil {
		return
	}

	jsonDecoder := sonic.ConfigDefault.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return
}

func (ipk *InPacket) ReadJsonTo(val any) error {
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

func (ipk *InPacket) ReadBytesTo(data []byte) error {
	readData, err := ipk.ReadBytes()
	if err != nil {
		return err
	}

	copy(data, readData)

	return nil
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

func (ipk *InPacket) ReadStreamBytesTo(buffer *bytes.Buffer) error {
	length, err := ipk.ReadInt64()
	if err != nil {
		return err
	}

	err = ipk.ReadStreamTo(length, buffer)
	if err != nil {
		return err
	}

	return nil
}
