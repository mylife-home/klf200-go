package commands

import (
	"fmt"
	"klf200/transport"
)

type RebootReq struct {
}

var _ Request = (*RebootReq)(nil)

func (req *RebootReq) Code() transport.Command {
	return transport.GW_REBOOT_REQ
}

func (req *RebootReq) NewConfirm() Confirm {
	return &RebootCfm{}
}

func (req *RebootReq) Write() ([]byte, error) {
	return emptyData, nil
}

type RebootCfm struct {
}

var _ Confirm = (*RebootCfm)(nil)

func (cfm *RebootCfm) Code() transport.Command {
	return transport.GW_REBOOT_CFM
}

func (cfm *RebootCfm) Read(data []byte) error {
	if len(data) != 0 {
		return fmt.Errorf("bad length")
	}

	return nil
}
