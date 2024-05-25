package commands

import (
	"fmt"
	"klf200/transport"
)

type LeaveLearnStateReq struct {
}

var _ Request = (*LeaveLearnStateReq)(nil)

func (req *LeaveLearnStateReq) Code() transport.Command {
	return transport.GW_LEAVE_LEARN_STATE_REQ
}

func (req *LeaveLearnStateReq) NewConfirm() Confirm {
	return &LeaveLearnStateCfm{}
}

func (req *LeaveLearnStateReq) Write() ([]byte, error) {
	return emptyData, nil
}

type LeaveLearnStateCfm struct {
	Success bool
}

var _ Confirm = (*LeaveLearnStateCfm)(nil)

func (cfm *LeaveLearnStateCfm) Code() transport.Command {
	return transport.GW_LEAVE_LEARN_STATE_CFM
}

func (cfm *LeaveLearnStateCfm) Read(data []byte) error {
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
