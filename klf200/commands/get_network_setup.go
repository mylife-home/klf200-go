package commands

import (
	"bytes"
	"fmt"
	"klf200/binary"
	"klf200/transport"
	"net"
)

type GetNetworkSetupReq struct {
}

var _ Request = (*GetNetworkSetupReq)(nil)

func (req *GetNetworkSetupReq) Code() transport.Command {
	return transport.GW_GET_NETWORK_SETUP_REQ
}

func (req *GetNetworkSetupReq) NewConfirm() Confirm {
	return &GetNetworkSetupCfm{}
}

func (req *GetNetworkSetupReq) Write() ([]byte, error) {
	return emptyData, nil
}

type GetNetworkSetupCfm struct {
	IpAddress net.IP
	Mask      net.IPMask
	DefGW     net.IP
	DHCP      bool
}

var _ Confirm = (*GetNetworkSetupCfm)(nil)

func (cfm *GetNetworkSetupCfm) Code() transport.Command {
	return transport.GW_GET_NETWORK_SETUP_CFM
}

func (cfm *GetNetworkSetupCfm) Read(data []byte) error {
	if len(data) != 13 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var value uint8

	cfm.IpAddress = cfm.readIp(reader)
	cfm.Mask = cfm.readMask(reader)
	cfm.DefGW = cfm.readIp(reader)

	value, _ = reader.ReadU8()

	switch value {
	case 0:
		cfm.DHCP = false
	case 1:
		cfm.DHCP = true
	default:
		return fmt.Errorf("bad DHCP flag %d", value)
	}

	return nil
}

func (cfm *GetNetworkSetupCfm) readIp(reader binary.BinaryReader) net.IP {
	data := make([]byte, 4)
	reader.Read(data)
	return net.IPv4(data[0], data[1], data[2], data[3])
}

func (cfm *GetNetworkSetupCfm) readMask(reader binary.BinaryReader) net.IPMask {
	data := make([]byte, 4)
	reader.Read(data)
	return net.IPv4Mask(data[0], data[1], data[2], data[3])
}
