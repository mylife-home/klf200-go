package commands

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/mylife-home/klf200-go/internal/binary"
	"github.com/mylife-home/klf200-go/internal/transport"
)

type StatusRequestReq struct {
	SessionID            int
	NodeIndexes          []int
	StatusType           StatusRequestStatusType
	FunctionalParameters map[FunctionalParameter]bool
}

type StatusRequestStatusType int

// Request Target position
const StatusRequestTargetPosition StatusRequestStatusType = 0

// Request Current position
const StatusRequestCurrentPosition StatusRequestStatusType = 1

// Request Remaining time
const StatusRequestRemainingTime StatusRequestStatusType = 2

// Request Main info
const StatusRequestMainInfo StatusRequestStatusType = 3

func (t StatusRequestStatusType) String() string {
	switch t {
	case StatusRequestTargetPosition:
		return "StatusRequestTargetPosition"
	case StatusRequestCurrentPosition:
		return "StatusRequestCurrentPosition"
	case StatusRequestRemainingTime:
		return "StatusRequestRemainingTime"
	case StatusRequestMainInfo:
		return "StatusRequestMainInfo"
	default:
		return fmt.Sprintf("<%d>", t)
	}
}

var _ Request = (*StatusRequestReq)(nil)

func (req *StatusRequestReq) Code() transport.Command {
	return transport.GW_STATUS_REQUEST_REQ
}

func (req *StatusRequestReq) NewConfirm() Confirm {
	return &StatusRequestCfm{}
}

func (req *StatusRequestReq) Write() ([]byte, error) {

	buff := &bytes.Buffer{}
	writer := binary.MakeBinaryWriter(buff)

	writer.WriteU16(uint16(req.SessionID))

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

	writer.WriteU8(uint8(req.StatusType))

	var bitmap uint16 = 0
	parametersCount := 0

	for index := 1; index < 17; index++ {
		param := FunctionalParameter(index)
		value := req.FunctionalParameters[param]

		if value {
			parametersCount++
		}

		// TODO: test this
		pos := 16 - param
		if value {
			bitmap |= (1 << pos)
		} else {
			bitmap &= ^(1 << pos)
		}
	}

	if parametersCount > 7 {
		return nil, errors.New("too many functional parameters")
	}

	writer.WriteU16(bitmap)

	return buff.Bytes(), nil
}

type StatusRequestCfm struct {
	SessionID int
	Success   bool
}

var _ Confirm = (*StatusRequestCfm)(nil)

func (cfm *StatusRequestCfm) Code() transport.Command {
	return transport.GW_STATUS_REQUEST_CFM
}

func (cfm *StatusRequestCfm) Read(data []byte) error {
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

type StatusRequestNtf struct {
	SessionID   int
	StatusID    CommandRunOwner
	NodeIndex   int
	RunStatus   CommandRunStatus
	StatusReply CommandRunStatusReply
	StatusType  StatusRequestStatusType
	StatusData  StatusData
}

var _ Notify = (*StatusRequestNtf)(nil)

func init() {
	registerNotify(func() Notify { return &StatusRequestNtf{} })
}

func (ntf *StatusRequestNtf) Code() transport.Command {
	return transport.GW_STATUS_REQUEST_NTF
}

func (ntf *StatusRequestNtf) Read(data []byte) error {
	if len(data) < 6 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8
	var u16 uint16

	u16, _ = reader.ReadU16()
	ntf.SessionID = int(u16)

	u8, _ = reader.ReadU8()
	ntf.StatusID = CommandRunOwner(u8)

	u8, _ = reader.ReadU8()
	ntf.NodeIndex = int(u8)

	u8, _ = reader.ReadU8()
	ntf.RunStatus = CommandRunStatus(u8)

	u8, _ = reader.ReadU8()
	ntf.StatusReply = CommandRunStatusReply(u8)

	u8, _ = reader.ReadU8()
	ntf.StatusType = StatusRequestStatusType(u8)

	switch ntf.StatusType {
	case StatusRequestTargetPosition, StatusRequestCurrentPosition, StatusRequestRemainingTime:
		if len(data) != 59 {
			return fmt.Errorf("bad length")
		}

		ntf.StatusData = &StatusDataParameters{}

	case StatusRequestMainInfo:
		if len(data) != 18 {
			return fmt.Errorf("bad length")
		}

		ntf.StatusData = &StatusDataMainInfo{}
	}

	// Note: may be 255 if the status reply indicates an error.
	// In this case there is no data (only 0s).

	if ntf.StatusData != nil {
		ntf.StatusData.read(reader)
	}

	return nil
}

// StatusDataParameters or StatusDataMainInfo
type StatusData interface {
	read(reader binary.BinaryReader)
}

// Filled when StatusType = "Target Position" or StatusType = "Current Position" or StatusType = "Remaining Time"
type StatusDataParameters struct {
	ParameterValues map[FunctionalParameter]int
}

func (data *StatusDataParameters) read(reader binary.BinaryReader) {
	var count uint8
	var typ uint16
	var value uint16

	count, _ = reader.ReadU8()

	data.ParameterValues = make(map[FunctionalParameter]int)

	for index := 0; index < int(count); index++ {
		typ, _ = reader.ReadU16()
		value, _ = reader.ReadU16()

		data.ParameterValues[FunctionalParameter(typ)] = int(value)
	}
}

// Filled when StatusType = "Main info"
type StatusDataMainInfo struct {
	TargetPosition             MPValue
	CurrentPosition            MPValue
	RemainingTime              time.Duration
	LastMasterExecutionAddress uint32
	LastCommandOriginator      CommandOriginator
}

func (data *StatusDataMainInfo) read(reader binary.BinaryReader) {
	var u8 uint8
	var u16 uint16
	var u32 uint32

	u16, _ = reader.ReadU16()
	data.TargetPosition = MPValue(u16)

	u16, _ = reader.ReadU16()
	data.CurrentPosition = MPValue(u16)

	u16, _ = reader.ReadU16()
	data.RemainingTime = time.Second * time.Duration(u16)

	u32, _ = reader.ReadU32()
	data.LastMasterExecutionAddress = u32

	u8, _ = reader.ReadU8()
	data.LastCommandOriginator = CommandOriginator(u8)
}
