package in

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

func ToInt32(data io.Reader) (number int32, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

func ToInt64(data io.Reader) (number int64, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

func NewInbound(r io.Reader) *Inbound {
	return &Inbound{
		r: r,
	}
}

type Reader interface {
	ReturnableReader
	ThrowableReader
}

type ReturnableReader interface {
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
	ReadStreamString() (io.ReadWriter, error)
	ReadJson() (any, error)
	ReadBytes() ([]byte, error)
	ReadStreamBytes() (io.ReadWriter, error)
}

type ThrowableReader interface {
	ReadInt32To(*int32) error
	ReadInt64To(*int64) error
	ReadStringTo(*string) error
	ReadStreamStringTo(io.Writer) error
	ReadJsonTo(any) error
	ReadBytesTo([]byte) error
	ReadStreamBytesTo(io.Writer) error
}

type Inbound struct {
	r io.Reader
}

func (ipk *Inbound) read(len int64) ([]byte, error) {
	byteBuffer := make([]byte, len)
	written, err := ipk.r.Read(byteBuffer)
	if written < int(len) {
		return nil, errors.New("there is no data left")
	}
	if err != nil {
		return nil, err
	}

	return byteBuffer, nil
}

func (ipk *Inbound) readStream(len int64) (io.ReadWriter, error) {
	byteBuffer := new(bytes.Buffer)
	written, err := io.CopyN(byteBuffer, ipk.r, len)
	if written < len {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return byteBuffer, nil
}

func (ipk *Inbound) readStreamTo(len int64, buffer io.Writer) error {
	written, err := io.CopyN(buffer, ipk.r, len)
	if (written < len) || (err != nil) {
		return err
	}

	return nil
}

func (ipk *Inbound) ReadInt32() (int32, error) {
	rawData, err := ipk.readStream(4)
	if err != nil {
		return 0, err
	}

	number, err := ToInt32(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ipk *Inbound) ReadInt32To(data *int32) error {
	readData, err := ipk.ReadInt32()
	if err != nil {
		return err
	}

	*data = readData

	return nil
}

func (ipk *Inbound) ReadInt64() (int64, error) {
	rawData, err := ipk.readStream(8)
	if err != nil {
		return 0, err
	}

	number, err := ToInt64(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (ipk *Inbound) ReadInt64To(data *int64) error {
	readData, err := ipk.ReadInt64()
	if err != nil {
		return err
	}

	*data = readData

	return nil
}

func (ipk *Inbound) ReadString() (string, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return "", err
	}

	rawData, err := ipk.read(length)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func (ipk *Inbound) ReadStringTo(data *string) error {
	readData, err := ipk.ReadString()
	if err != nil {
		return err
	}

	*data = readData

	return nil
}

func (ipk *Inbound) ReadStreamString() (io.ReadWriter, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	rawData, err := ipk.readStream(length)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (ipk *Inbound) ReadStreamStringTo(buffer io.Writer) error {
	length, err := ipk.ReadInt64()
	if err != nil {
		return err
	}

	err = ipk.readStreamTo(length, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (ipk *Inbound) ReadJson() (val any, err error) {
	jsonString, err := ipk.ReadStreamString()
	if err != nil {
		return
	}

	jsonDecoder := json.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return
}

func (ipk *Inbound) ReadJsonTo(val any) error {
	jsonString, err := ipk.ReadStreamString()
	if err != nil {
		return err
	}

	jsonDecoder := json.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return nil
}

func (ipk *Inbound) ReadBytes() ([]byte, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	byteBuf, err := ipk.read(length)
	if err != nil {
		return nil, err
	}

	return byteBuf, nil
}

func (ipk *Inbound) ReadBytesTo(data []byte) error {
	readData, err := ipk.ReadBytes()
	if err != nil {
		return err
	}

	copy(data, readData)

	return nil
}

func (ipk *Inbound) ReadStreamBytes() (io.ReadWriter, error) {
	length, err := ipk.ReadInt64()
	if err != nil {
		return nil, err
	}

	byteBuf, err := ipk.readStream(length)
	if err != nil {
		return nil, err
	}

	return byteBuf, nil
}

func (ipk *Inbound) ReadStreamBytesTo(buffer io.Writer) error {
	length, err := ipk.ReadInt64()
	if err != nil {
		return err
	}

	err = ipk.readStreamTo(length, buffer)
	if err != nil {
		return err
	}

	return nil
}
