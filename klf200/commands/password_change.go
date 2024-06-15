package commands

import (
	"bytes"
	"fmt"

	"github.com/mylife-home/klf200-go/klf200/transport"
)

type PasswordChangeReq struct {
	Password string
}

var _ Request = (*PasswordChangeReq)(nil)

func (req *PasswordChangeReq) Code() transport.Command {
	return transport.GW_PASSWORD_CHANGE_REQ
}

func (req *PasswordChangeReq) NewConfirm() Confirm {
	return &PasswordChangeCfm{}
}

func (req *PasswordChangeReq) Write() ([]byte, error) {
	array := []byte(req.Password)
	if len(array) > 32 {
		return nil, fmt.Errorf("password too long")
	}

	remain := 32 - len(array)

	if remain > 0 {
		pad := make([]byte, remain)
		array = append(array, pad...)
	}

	return array, nil
}

type PasswordChangeCfm struct {
	Success bool
}

var _ Confirm = (*PasswordChangeCfm)(nil)

func (cfm *PasswordChangeCfm) Code() transport.Command {
	return transport.GW_PASSWORD_CHANGE_CFM
}

func (cfm *PasswordChangeCfm) Read(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("bad length")
	}

	switch data[0] {
	case 0:
		cfm.Success = true
	case 1:
		cfm.Success = false
	default:
		return fmt.Errorf("bad status")
	}
	return nil
}

type PasswordChangeNtf struct {
	NewPassword string
}

var _ Notify = (*PasswordChangeNtf)(nil)

func init() {
	registerNotify(func() Notify { return &PasswordChangeNtf{} })
}

func (ntf *PasswordChangeNtf) Code() transport.Command {
	return transport.GW_PASSWORD_CHANGE_NTF
}

func (ntf *PasswordChangeNtf) Read(data []byte) error {
	if len(data) != 32 {
		return fmt.Errorf("bad length")
	}

	data = bytes.TrimRight(data, "\x00")

	ntf.NewPassword = string(data)

	return nil
}
