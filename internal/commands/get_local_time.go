package commands

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mylife-home/klf200-go/internal/binary"
	"github.com/mylife-home/klf200-go/internal/transport"
)

type GetLocalTimeReq struct {
}

var _ Request = (*GetLocalTimeReq)(nil)

func (req *GetLocalTimeReq) Code() transport.Command {
	return transport.GW_GET_LOCAL_TIME_REQ
}

func (req *GetLocalTimeReq) NewConfirm() Confirm {
	return &GetLocalTimeCfm{}
}

func (req *GetLocalTimeReq) Write() ([]byte, error) {
	return emptyData, nil
}

type DaylightSavingFlag int8

const DaylightSavingUnknown DaylightSavingFlag = -1
const DaylightSavingNotActive DaylightSavingFlag = 0
const DaylightSavingActive DaylightSavingFlag = 1

type LocalTime struct {
	// Seconds after the minute (local time), range 0-61
	Second int

	// Minutes after the hour (local time), range 0-59
	Minute int

	// Hours since midnight (local time), range 0-23
	Hour int

	// Day of the month, range 1-31
	DayOfMonth int

	// Months since January, range 0-11
	Month int

	Year int

	// Days since Sunday, range 0-6
	WeekDay int

	// Days since January 1, range 0-365
	DayOfYear int

	DaylightSavingFlag DaylightSavingFlag
}

type GetLocalTimeCfm struct {
	UtcTime   time.Time
	LocalTime LocalTime
}

var _ Confirm = (*GetLocalTimeCfm)(nil)

func (cfm *GetLocalTimeCfm) Code() transport.Command {
	return transport.GW_GET_LOCAL_TIME_CFM
}

func (cfm *GetLocalTimeCfm) Read(data []byte) error {
	if len(data) != 15 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))

	utc, _ := reader.ReadU32()
	cfm.UtcTime = time.Unix(int64(utc), 0)

	cfm.LocalTime.Second = cfm.readU8AsInt(reader)
	cfm.LocalTime.Minute = cfm.readU8AsInt(reader)
	cfm.LocalTime.Hour = cfm.readU8AsInt(reader)
	cfm.LocalTime.DayOfMonth = cfm.readU8AsInt(reader)
	cfm.LocalTime.Month = cfm.readU8AsInt(reader)
	cfm.LocalTime.Year = cfm.readU16AsInt(reader) + 1900
	cfm.LocalTime.WeekDay = cfm.readU8AsInt(reader)
	cfm.LocalTime.DayOfYear = cfm.readU16AsInt(reader)
	cfm.LocalTime.DaylightSavingFlag = DaylightSavingFlag(cfm.readU8AsInt(reader))

	return nil
}

func (cfm *GetLocalTimeCfm) readU8AsInt(reader binary.BinaryReader) int {
	value, _ := reader.ReadU8()
	return int(value)
}

func (cfm *GetLocalTimeCfm) readU16AsInt(reader binary.BinaryReader) int {
	value, _ := reader.ReadU16()
	return int(value)
}
