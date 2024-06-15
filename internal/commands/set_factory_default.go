package commands

import (
	"fmt"

	"github.com/mylife-home/klf200-go/internal/transport"
)

type SetFactoryDefaultReq struct {
}

var _ Request = (*SetFactoryDefaultReq)(nil)

func (req *SetFactoryDefaultReq) Code() transport.Command {
	return transport.GW_SET_FACTORY_DEFAULT_REQ
}

func (req *SetFactoryDefaultReq) NewConfirm() Confirm {
	return &SetFactoryDefaultCfm{}
}

func (req *SetFactoryDefaultReq) Write() ([]byte, error) {
	return emptyData, nil
}

type SetFactoryDefaultCfm struct {
}

var _ Confirm = (*SetFactoryDefaultCfm)(nil)

func (cfm *SetFactoryDefaultCfm) Code() transport.Command {
	return transport.GW_SET_FACTORY_DEFAULT_CFM
}

func (cfm *SetFactoryDefaultCfm) Read(data []byte) error {
	if len(data) != 0 {
		return fmt.Errorf("bad length")
	}

	return nil
}
