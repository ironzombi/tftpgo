package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

/*
 * Types for tftp:
 * tftp read method in golang
 * basic usage for tftp protocol
 * Read request RRQ
 * Data DATA
 * Aknowledgement ACK
 * Error ERROR
 */

const (
	DatagramSize = 516              // max datagram size
	BlockSize    = DatagramSize - 4 // datagram - 4-byte header
)

type OpCode uint16

const (
	OpRRQ OpCode = iota + 1
	_            // no WRQ support
	OpData
	OpAck
	OpErr
)

type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

type ReadReq struct {
	Filename string
	Mode     string
}

//
func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}

	// op code + filename + 0 byte + mode + 0 byte
	cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpRRQ) // write the op code
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Filename) // write the filename
	if err != nil {
		return nil, nil
	}

	err = b.WriteByte(0) // write the 0 byte
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode) // write the mode
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write the 0 byte
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read the op code
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return errors.New("invalid RRQ")
	}

	q.Filename, err = r.ReadString(0) // read the filename
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Filename = strings.TrimRight(q.Filename, "\x00") // remove the 0 byte
	if len(q.Filename) == 0 {
		return errors.New("invalid RRQ")
	}

	q.Mode, err = r.ReadString(0) // read the mode
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Mode = strings.TrimRight(q.Mode, "\x00") // remove the 0 byte
	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ")
	}

	actual := strings.ToLower(q.Mode) // enforce octet mode
	if actual != "octet" {
		return errors.New("only binary transfers supported")
	}

	return nil
}

type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	b.Grow(DatagramSize)

	d.Block++ // increment the block number

	err := binary.Write(b, binary.BigEndian, OpData) // write the op code
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, d.Block) // write the block number
	if err != nil {
		return nil, err
	}

	// write upto BlockSize of bytes
	_, err = io.CopyN(b, d.Payload, BlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {
	if l := len(p); l < 4 || l > DatagramSize {
		return errors.New("invalid data packet")
	}

	var opcode OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode) // read the op code
	if err != nil || opcode != OpData {
		return errors.New("invalid data packet")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block) // read the block number
	if err != nil {
		return errors.New("invalid data packet")
	}

	d.Payload = bytes.NewBuffer(p[4:])

	return nil
}

type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
	cap := 2 + 2 // [OpCode][Block#] 4bytes

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpAck) // write opcode
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, a) // write block number
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (a Ack) UnmarshalBinary(p []byte) error {
	var code OpCode

	r := bytes.NewReader(p)

	err := binary.Read(r, binary.BigEndian, &code) // read the opcode
	if err != nil {
		return err
	}

	if code != OpAck {
		return errors.New("INVALID ACK")
	}

	return binary.Read(r, binary.BigEndian, a)
}

type Err struct {
	Error   ErrCode
	Message string
}

func (e Err) MarshalBinary() ([]byte, error) {
	//opcode + errorcode + message + 0 byte
	cap := 2 + 2 + len(e.Message) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpErr) // write opcode
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, e.Error) // write errcode
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(e.Message) //write msg
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write 0byte end
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e Err) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read opcode
	if err != nil {
		return err
	}

	if code != OpErr {
		return errors.New("ERROR invalid")
	}

	err = binary.Read(r, binary.BigEndian, &e.Error) // read error message
	if err != nil {
		return err
	}

	e.Message, err = r.ReadString(0)
	e.Message = strings.TrimRight(e.Message, "\x00") // remove the trailing null/0-byte

	return err
}
