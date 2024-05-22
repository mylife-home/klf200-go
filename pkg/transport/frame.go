package transport

import (
	"bytes"
	"fmt"
	"klf200/pkg/binary"
)

type Command uint16

const protocolId uint8 = 0

type Frame struct {
	Cmd  Command
	Data []byte
}

func (frame *Frame) Write() *bytes.Buffer {
	buff := &bytes.Buffer{}
	writer := binary.MakeBinaryWriter(buff)

	len := len(frame.Data)
	if len > 250 {
		panic("drame data too long")
	}

	len = len + 3 // count len itself + command

	writer.WriteU8(protocolId)
	writer.WriteU8(uint8(len))
	writer.WriteU16(uint16(frame.Cmd))
	writer.Write(frame.Data)
	writer.WriteU8(checksum(buff.Bytes()))

	return buff
}

func FrameRead(buff *bytes.Buffer) (*Frame, error) {
	reader := binary.MakeBinaryReader(buff)
	raw := buff.Bytes()
	frame := &Frame{}

	currentProtocolId, err := reader.ReadU8()
	if err != nil {
		return nil, err
	}

	if currentProtocolId != protocolId {
		return nil, fmt.Errorf("unexpected protocol ID")
	}

	len, err := reader.ReadU8()
	if err != nil {
		return nil, err
	}

	if len < 3 || len > 253 {
		return nil, fmt.Errorf("unexpected length")
	}

	cmd, err := reader.ReadU16()
	if err != nil {
		return nil, err
	}

	data := make([]byte, len-3)
	err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	currentCs, err := reader.ReadU8()
	if err != nil {
		return nil, err
	}

	expectedCs := checksum(raw[0 : len+1]) // include protocol id

	if expectedCs != currentCs {
		return nil, fmt.Errorf("wrong checksum")
	}

	frame.Cmd = Command(cmd)
	frame.Data = data

	return frame, nil
}

func checksum(buff []byte) byte {
	var cs byte

	for _, b := range buff {
		cs = cs ^ b
	}

	return cs
}
