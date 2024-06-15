package klf200

import (
	"errors"
	"time"

	"github.com/mylife-home/klf200-go/klf200/commands"
)

type Device struct {
	client *Client
}

func (dev *Device) ChangePassword(newPassword string) error {
	req := &commands.PasswordChangeReq{Password: newPassword}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return err
	}

	if !cfm.(*commands.PasswordChangeCfm).Success {
		return errors.New("the request failed")
	}

	return nil
}

func (dev *Device) GetVersion() (*commands.GetVersionCfm, error) {
	req := &commands.GetVersionReq{}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return nil, err
	}

	return cfm.(*commands.GetVersionCfm), nil
}

func (dev *Device) GetProtocolVersion() (*commands.GetProtocolVersionCfm, error) {
	req := &commands.GetProtocolVersionReq{}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return nil, err
	}

	return cfm.(*commands.GetProtocolVersionCfm), nil
}

func (dev *Device) GetState() (*commands.GetStateCfm, error) {
	req := &commands.GetStateReq{}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return nil, err
	}

	return cfm.(*commands.GetStateCfm), nil
}

func (dev *Device) LeaveLearnState() error {
	req := &commands.LeaveLearnStateReq{}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return err
	}

	if !cfm.(*commands.LeaveLearnStateCfm).Success {
		return errors.New("the request failed")
	}

	return nil
}

func (dev *Device) SetUtc(timestamp time.Time) error {
	req := &commands.SetUtcReq{Timestamp: timestamp}
	_, err := dev.client.execute(req)
	if err != nil {
		return err
	}

	return nil
}

func (dev *Device) SetTimeZone(tzstr string) error {
	// TODO: help create tzstr
	req := &commands.RtcSetTimeZoneReq{TimeZoneString: tzstr}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return err
	}

	if !cfm.(*commands.RtcSetTimeZoneCfm).Success {
		return errors.New("the request failed")
	}

	return nil
}

func (dev *Device) GetLocalTime() (*commands.GetLocalTimeCfm, error) {
	req := &commands.GetLocalTimeReq{}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return nil, err
	}

	return cfm.(*commands.GetLocalTimeCfm), nil
}

func (dev *Device) Reboot() error {
	req := &commands.RebootReq{}
	_, err := dev.client.execute(req)
	if err != nil {
		return err
	}

	return nil
}

func (dev *Device) SetFactoryDefault() error {
	req := &commands.SetFactoryDefaultReq{}
	_, err := dev.client.execute(req)
	if err != nil {
		return err
	}

	return nil
}

func (dev *Device) GetNetworkSetup() (*commands.GetNetworkSetupCfm, error) {
	req := &commands.GetNetworkSetupReq{}
	cfm, err := dev.client.execute(req)
	if err != nil {
		return nil, err
	}

	return cfm.(*commands.GetNetworkSetupCfm), nil
}
