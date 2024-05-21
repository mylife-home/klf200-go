package commands

import (
	"fmt"
	"klf200/pkg/transport"
)

type PasswordEnterReq struct {
	Password string
}

var _ Request = (*PasswordEnterReq)(nil)

func (req *PasswordEnterReq) Code() transport.Command {
	return GW_PASSWORD_ENTER_REQ
}

func (req *PasswordEnterReq) Write() ([]byte, error) {
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

type PasswordEnterCfm struct {
	Status bool
}

var _ Confirm = (*PasswordEnterCfm)(nil)

func (cfm *PasswordEnterCfm) Code() transport.Command {
	return GW_PASSWORD_ENTER_CFM
}

func (cfm *PasswordEnterCfm) Read(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("bad length")
	}

	switch data[0] {
	case 0:
		cfm.Status = false
	case 1:
		cfm.Status = true
	default:
		return fmt.Errorf("bad status")
	}
	return nil
}
