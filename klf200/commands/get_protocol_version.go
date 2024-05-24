package commands

import (
	"bytes"
	"fmt"
	"klf200/binary"
	"klf200/transport"
)

type GetProtocolVersionReq struct {
}

var _ Request = (*GetProtocolVersionReq)(nil)

func (req *GetProtocolVersionReq) Code() transport.Command {
	return GW_GET_PROTOCOL_VERSION_REQ
}

func (req *GetProtocolVersionReq) NewConfirm() Confirm {
	return &GetProtocolVersionCfm{}
}

func (req *GetProtocolVersionReq) Write() ([]byte, error) {
	return emptyData, nil
}

type GetProtocolVersionCfm struct {
	MajorVersion int
	MinorVersion int
}

var _ Confirm = (*GetProtocolVersionCfm)(nil)

func (cfm *GetProtocolVersionCfm) Code() transport.Command {
	return GW_GET_PROTOCOL_VERSION_CFM
}

func (cfm *GetProtocolVersionCfm) Read(data []byte) error {
	if len(data) != 4 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var value uint16

	value, _ = reader.ReadU16()
	cfm.MajorVersion = int(value)

	value, _ = reader.ReadU16()
	cfm.MinorVersion = int(value)

	return nil
}
