package commands

import (
	"fmt"

	"github.com/mylife-home/klf200-go/transport"
)

type RtcSetTimeZoneReq struct {
	TimeZoneString string
}

var _ Request = (*RtcSetTimeZoneReq)(nil)

func (req *RtcSetTimeZoneReq) Code() transport.Command {
	return transport.GW_RTC_SET_TIME_ZONE_REQ
}

func (req *RtcSetTimeZoneReq) NewConfirm() Confirm {
	return &RtcSetTimeZoneCfm{}
}

func (req *RtcSetTimeZoneReq) Write() ([]byte, error) {
	array := []byte(req.TimeZoneString)
	if len(array) > 64 {
		return nil, fmt.Errorf("timezone string too long")
	}

	remain := 64 - len(array)

	if remain > 0 {
		pad := make([]byte, remain)
		array = append(array, pad...)
	}

	return array, nil
}

type RtcSetTimeZoneCfm struct {
	Success bool
}

var _ Confirm = (*RtcSetTimeZoneCfm)(nil)

func (cfm *RtcSetTimeZoneCfm) Code() transport.Command {
	return transport.GW_RTC_SET_TIME_ZONE_CFM
}

func (cfm *RtcSetTimeZoneCfm) Read(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("bad length")
	}

	switch data[0] {
	case 0:
		cfm.Success = false
	case 1:
		cfm.Success = true
	default:
		return fmt.Errorf("bad status")
	}
	return nil

}
