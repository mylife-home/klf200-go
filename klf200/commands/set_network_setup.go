package commands

import (
	"bytes"
	"fmt"
	"klf200/binary"
	"klf200/transport"
	"net"
)

type SetNetworkSetupReq struct {
	IpAddress net.IP
	Mask      net.IPMask
	DefGW     net.IP
	DHCP      bool
}

var _ Request = (*SetNetworkSetupReq)(nil)

func (req *SetNetworkSetupReq) Code() transport.Command {
	return transport.GW_SET_NETWORK_SETUP_REQ
}

func (req *SetNetworkSetupReq) NewConfirm() Confirm {
	return &SetNetworkSetupCfm{}
}

func (req *SetNetworkSetupReq) Write() ([]byte, error) {

	buff := &bytes.Buffer{}
	writer := binary.MakeBinaryWriter(buff)

	writer.Write(req.IpAddress.To4())
	writer.Write(req.Mask)
	writer.Write(req.DefGW.To4())

	var dhcp uint8 = 0
	if req.DHCP {
		dhcp = 1
	}

	writer.WriteU8(dhcp)

	return buff.Bytes(), nil
}

type SetNetworkSetupCfm struct {
}

var _ Confirm = (*SetNetworkSetupCfm)(nil)

func (cfm *SetNetworkSetupCfm) Code() transport.Command {
	return transport.GW_SET_NETWORK_SETUP_CFM
}

func (cfm *SetNetworkSetupCfm) Read(data []byte) error {
	if len(data) != 0 {
		return fmt.Errorf("bad length")
	}

	return nil
}
