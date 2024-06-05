package commands

import (
	"bytes"
	"errors"
	"fmt"
	"klf200/binary"
	"klf200/transport"
)

type CommandSendReq struct {
	SessionID                 int
	CommandOriginator         CommandOriginator
	PriorityLevel             PriorityLevel
	ParameterActive           FunctionalParameter
	FunctionalParameterValues map[FunctionalParameter]int
	NodeIndexes               []int
	PriorityLevelLock         PriorityLevelLock
	PriorityLevelInfo         PriorityLevelInfo
	LockTime                  LockTime
}

var _ Request = (*CommandSendReq)(nil)

func (req *CommandSendReq) Code() transport.Command {
	return transport.GW_COMMAND_SEND_REQ
}

func (req *CommandSendReq) NewConfirm() Confirm {
	return &CommandSendCfm{}
}

func (req *CommandSendReq) Write() ([]byte, error) {

	buff := &bytes.Buffer{}
	writer := binary.MakeBinaryWriter(buff)

	writer.WriteU16(uint16(req.SessionID))
	writer.WriteU8(uint8(req.CommandOriginator))
	writer.WriteU8(uint8(req.PriorityLevel))
	writer.WriteU8(uint8(req.ParameterActive))

	var bitmap uint16 = 0

	for index := 0; index < 17; index++ {
		param := FunctionalParameter(index)
		_, defined := req.FunctionalParameterValues[param]
		if !defined && param == FunctionalParameterMP {
			return nil, errors.New("missing value for FunctionalParameterMP")
		}

		if param != FunctionalParameterMP {
			// TODO: test this
			pos := 16 - param
			if defined {
				bitmap |= (1 << pos)
			} else {
				bitmap &= ^(1 << pos)
			}
		}
	}

	writer.WriteU16(bitmap)

	for index := 0; index < 17; index++ {
		param := FunctionalParameter(index)
		val, defined := req.FunctionalParameterValues[param]

		if !defined {
			val = 0
		}

		writer.WriteU16(uint16(val))
	}

	if len(req.NodeIndexes) < 1 || len(req.NodeIndexes) > 20 {
		return nil, fmt.Errorf("bad node indexes len (got %d, expected > 0 && <= 20)", len(req.NodeIndexes))
	}

	writer.WriteU8(uint8(len(req.NodeIndexes)))

	for index := 0; index < 20; index++ {
		var value uint8 = 0
		if index < len(req.NodeIndexes) {
			value = uint8(req.NodeIndexes[index])
		}

		writer.WriteU8(value)
	}

	writer.WriteU8((uint8(req.PriorityLevelLock)))
	writer.WriteU16((uint16(req.PriorityLevelInfo)))
	writer.WriteU8(uint8(req.LockTime))

	return buff.Bytes(), nil
}

type CommandSendCfm struct {
	SessionID int
	Success   bool
}

var _ Confirm = (*CommandSendCfm)(nil)

func (cfm *CommandSendCfm) Code() transport.Command {
	return transport.GW_COMMAND_SEND_CFM
}

func (cfm *CommandSendCfm) Read(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8
	var u16 uint16

	u16, _ = reader.ReadU16()
	cfm.SessionID = int(u16)

	u8, _ = reader.ReadU8()
	switch u8 {
	case 0:
		cfm.Success = false
	case 1:
		cfm.Success = true
	default:
		return fmt.Errorf("bad status")
	}

	return nil
}
