package commands

import (
	"bytes"
	"fmt"
	"klf200/binary"
	"klf200/transport"
	"time"
)

type SetUtcReq struct {
	Timestamp time.Time
}

var _ Request = (*SetUtcReq)(nil)

func (req *SetUtcReq) Code() transport.Command {
	return transport.GW_SET_UTC_REQ
}

func (req *SetUtcReq) NewConfirm() Confirm {
	return &SetUtcCfm{}
}

func (req *SetUtcReq) Write() ([]byte, error) {
	buff := &bytes.Buffer{}
	writer := binary.MakeBinaryWriter(buff)
	writer.WriteU32(uint32(req.Timestamp.Unix()))

	return buff.Bytes(), nil
}

type SetUtcCfm struct {
}

var _ Confirm = (*SetUtcCfm)(nil)

func (cfm *SetUtcCfm) Code() transport.Command {
	return transport.GW_SET_UTC_CFM
}

func (cfm *SetUtcCfm) Read(data []byte) error {
	if len(data) != 0 {
		return fmt.Errorf("bad length")
	}

	return nil
}
