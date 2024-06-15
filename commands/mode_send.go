package commands

import (
	"bytes"
	"fmt"

	"github.com/mylife-home/klf200-go/binary"
	"github.com/mylife-home/klf200-go/transport"
)

type ModeSendReq struct {
	SessionID         int
	CommandOriginator CommandOriginator
	PriorityLevel     PriorityLevel

	ModeNumber    int // Set to 0
	ModeParameter int // Set to 0

	NodeIndexes       []int
	PriorityLevelLock PriorityLevelLock
	PriorityLevelInfo PriorityLevelInfo
	LockTime          LockTime
}

var _ Request = (*ModeSendReq)(nil)

func (req *ModeSendReq) Code() transport.Command {
	return transport.GW_MODE_SEND_REQ
}

func (req *ModeSendReq) NewConfirm() Confirm {
	return &ModeSendCfm{}
}

func (req *ModeSendReq) Write() ([]byte, error) {

	buff := &bytes.Buffer{}
	writer := binary.MakeBinaryWriter(buff)

	writer.WriteU16(uint16(req.SessionID))
	writer.WriteU8(uint8(req.CommandOriginator))
	writer.WriteU8(uint8(req.PriorityLevel))
	writer.WriteU8(uint8(req.ModeNumber))
	writer.WriteU8(uint8(req.ModeParameter))

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

type ModeSendCfm struct {
	SessionID int
	Status    ModeSendStatus
}

type ModeSendStatus int

// OK. Accepted by Command Handler
const ModeSendStatusSuccess ModeSendStatus = 0

// Failed. Rejected by Command Handler
const ModeSendStatusRejected ModeSendStatus = 1

// Failed with unknown Client ID
const ModeSendStatusUnknownClientId ModeSendStatus = 2

// Failed. Session ID already in use
const ModeSendStatusSessionIdAlreadyInUse ModeSendStatus = 3

// Failed. Busy – no free session slots – try again
const ModeSendStatusNoFreeSessionSlot ModeSendStatus = 4

// Failed. Illegal parameter value
const ModeSendStatusIllegalParameterValue ModeSendStatus = 5

// Failed. Not further defined error
const ModeSendStatusUnknownError ModeSendStatus = 255

func (s ModeSendStatus) String() string {
	switch s {
	case ModeSendStatusSuccess:
		return "ModeSendStatusSuccess"
	case ModeSendStatusRejected:
		return "ModeSendStatusRejected"
	case ModeSendStatusUnknownClientId:
		return "ModeSendStatusUnknownClientId"
	case ModeSendStatusSessionIdAlreadyInUse:
		return "ModeSendStatusSessionIdAlreadyInUse"
	case ModeSendStatusNoFreeSessionSlot:
		return "ModeSendStatusNoFreeSessionSlot"
	case ModeSendStatusIllegalParameterValue:
		return "ModeSendStatusIllegalParameterValue"
	case ModeSendStatusUnknownError:
		return "ModeSendStatusUnknownError"
	default:
		return fmt.Sprintf("<%d>", s)
	}
}

var _ Confirm = (*ModeSendCfm)(nil)

func (cfm *ModeSendCfm) Code() transport.Command {
	return transport.GW_MODE_SEND_CFM
}

func (cfm *ModeSendCfm) Read(data []byte) error {
	if len(data) != 3 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8
	var u16 uint16

	u16, _ = reader.ReadU16()
	cfm.SessionID = int(u16)

	u8, _ = reader.ReadU8()
	cfm.Status = ModeSendStatus(u8)

	return nil
}
