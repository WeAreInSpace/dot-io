package in

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

/*
Convert binary to int32

# Example

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(255))
	ToInt32(buffer)
*/
func ToInt32(data io.Reader) (number int32, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}
	return
}

/*
Convert binary to int64

# Example

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(255))
	ToInt64(buffer)
*/
func ToInt64(data io.Reader) (number int64, err error) {
	err = binary.Read(data, binary.BigEndian, &number)
	if err != nil {
		return 0, err
	}

	return
}

// Return an instance of Inbound
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

func (ib *Inbound) read(len int64) ([]byte, error) {
	byteBuffer := make([]byte, len)
	written, err := ib.r.Read(byteBuffer)
	if written < int(len) {
		return nil, errors.New("there is no data left")
	}
	if err != nil {
		return nil, err
	}

	return byteBuffer, nil
}

func (ib *Inbound) readStream(len int64) (io.ReadWriter, error) {
	byteBuffer := new(bytes.Buffer)
	written, err := io.CopyN(byteBuffer, ib.r, len)
	if written < len {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return byteBuffer, nil
}

func (ib *Inbound) readStreamTo(len int64, buffer io.Writer) error {
	written, err := io.CopyN(buffer, ib.r, len)
	if (written < len) || (err != nil) {
		return err
	}

	return nil
}

/*
Read int32 from reader and return (int32, error)

# Example

	number, err := ib.ReadInt32()
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
*/
func (ib *Inbound) ReadInt32() (int32, error) {
	rawData, err := ib.readStream(4)
	if err != nil {
		return 0, err
	}

	number, err := ToInt32(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

/*
Read int32 from reader and set value

# Example

	var number int32
	err = ib.ReadInt32To(&number)
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
*/
func (ib *Inbound) ReadInt32To(data *int32) error {
	readData, err := ib.ReadInt32()
	if err != nil {
		return err
	}

	*data = readData

	return nil
}

/*
Read int64 from reader and return (int64, error)

# Example

	number, err := ib.ReadInt64()
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
*/
func (ib *Inbound) ReadInt64() (int64, error) {
	rawData, err := ib.readStream(8)
	if err != nil {
		return 0, err
	}

	number, err := ToInt64(rawData)
	if err != nil {
		return 0, err
	}
	return number, nil
}

/*
Read int64 from reader and set value

# Example

	var number int64
	err = ib.ReadInt64To(&number)
	if err != nil {
		panic(err)
	}
	fmt.Println(number)
*/
func (ib *Inbound) ReadInt64To(data *int64) error {
	readData, err := ib.ReadInt64()
	if err != nil {
		return err
	}

	*data = readData

	return nil
}

func (ib *Inbound) ReadString() (string, error) {
	length, err := ib.ReadInt64()
	if err != nil {
		return "", err
	}

	rawData, err := ib.read(length)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func (ib *Inbound) ReadStringTo(data *string) error {
	readData, err := ib.ReadString()
	if err != nil {
		return err
	}

	*data = readData

	return nil
}

func (ib *Inbound) ReadStreamString() (io.ReadWriter, error) {
	length, err := ib.ReadInt64()
	if err != nil {
		return nil, err
	}

	rawData, err := ib.readStream(length)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (ib *Inbound) ReadStreamStringTo(buffer io.Writer) error {
	length, err := ib.ReadInt64()
	if err != nil {
		return err
	}

	err = ib.readStreamTo(length, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (ib *Inbound) ReadJson() (val any, err error) {
	jsonString, err := ib.ReadStreamString()
	if err != nil {
		return
	}

	jsonDecoder := json.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return
}

func (ib *Inbound) ReadJsonTo(val any) error {
	jsonString, err := ib.ReadStreamString()
	if err != nil {
		return err
	}

	jsonDecoder := json.NewDecoder(jsonString)
	jsonDecoder.Decode(&val)

	return nil
}

func (ib *Inbound) ReadBytes() ([]byte, error) {
	length, err := ib.ReadInt64()
	if err != nil {
		return nil, err
	}

	byteBuf, err := ib.read(length)
	if err != nil {
		return nil, err
	}

	return byteBuf, nil
}

func (ib *Inbound) ReadBytesTo(data []byte) error {
	readData, err := ib.ReadBytes()
	if err != nil {
		return err
	}

	copy(data, readData)

	return nil
}

func (ib *Inbound) ReadStreamBytes() (io.ReadWriter, error) {
	length, err := ib.ReadInt64()
	if err != nil {
		return nil, err
	}

	byteBuf, err := ib.readStream(length)
	if err != nil {
		return nil, err
	}

	return byteBuf, nil
}

func (ib *Inbound) ReadStreamBytesTo(buffer io.Writer) error {
	length, err := ib.ReadInt64()
	if err != nil {
		return err
	}

	err = ib.readStreamTo(length, buffer)
	if err != nil {
		return err
	}

	return nil
}
