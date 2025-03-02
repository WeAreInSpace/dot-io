package in

import (
	"bytes"
	"encoding/binary"
	"errors"
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

type InPacket struct {
	conn *net.TCPConn
}

func (ipk InPacket) Read(len int64) (*bytes.Buffer, error) {
	err := ipk.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return nil, err
	}

	byteBuf := make([]byte, len)
	readLen, err := ipk.conn.Read(byteBuf)
	if readLen < int(len) {
		return nil, errors.New("there is no data left")
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, errors.New("read timeout")
	}
	if err != nil {
		return nil, err
	}

	byteBuffer := new(bytes.Buffer)
	_, err = byteBuffer.Write(byteBuf)
	if err != nil {
		return nil, err
	}

	return byteBuffer, nil
}

func (ipk InPacket) ReadInt32() (int32, error) {
	rawData, err := ipk.Read(4)
	if err != nil {
		return 0, err
	}

	number, err := ToInt32(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ipk InPacket) ReadInt64() (int64, error) {
	rawData, err := ipk.Read(8)
	if err != nil {
		return 0, err
	}

	number, err := ToInt64(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ipk InPacket) ReadString() (*bytes.Buffer, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	rawData, err := ipk.Read(length)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (ipk InPacket) ReadJson(val any) error {
	jsonString, err := ipk.ReadString()
	if err != nil {
		return err
	}

	jsonDecoder := sonic.ConfigDefault.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return nil
}
