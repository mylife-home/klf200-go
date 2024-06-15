package commands

import (
	"bytes"
	"fmt"

	"github.com/mylife-home/klf200-go/internal/binary"
	"github.com/mylife-home/klf200-go/internal/transport"
)

type GetStateReq struct {
}

var _ Request = (*GetStateReq)(nil)

func (req *GetStateReq) Code() transport.Command {
	return transport.GW_GET_STATE_REQ
}

func (req *GetStateReq) NewConfirm() Confirm {
	return &GetStateCfm{}
}

func (req *GetStateReq) Write() ([]byte, error) {
	return emptyData, nil
}

type GatewayState int

// Test mode.
const GatewayStateTest GatewayState = 0

// Gateway mode, no actuator nodes in the system table.
const GatewayStateGatewayMode GatewayState = 1

// Gateway mode, with one or more actuator nodes in the system table.
const GatewayStateGatewayModeWithActuator GatewayState = 2

// Beacon mode, not configured by a remote controller.
const GatewayStateBeaconMode GatewayState = 3

// Beacon mode, has been configured by a remote controller.
const GatewayStateBeaconModeConfigured GatewayState = 4

type GatewaySubState int

// Idle state.
const GatewaySubStateIdle GatewaySubState = 0x00

// Performing task in Configuration Service handler
const GatewaySubStateConfigurationServiceHandler GatewaySubState = 0x01

// Performing Scene Configuration
const GatewaySubStateSceneConfiguration GatewaySubState = 0x02

// Performing Information Service Configuration.
const GatewaySubStateInformationServiceConfiguration GatewaySubState = 0x03

// Performing Contact input Configuration.
const GatewaySubStateContactInputConfiguration GatewaySubState = 0x04

// Performing task in Command Handler
const GatewaySubStateCommandHandler GatewaySubState = 0x80

// Performing task in Activate Group Handler
const GatewaySubStateActivateGroupHandler GatewaySubState = 0x81

// Performing task in Activate Scene Handler
const GatewaySubStateActivateSceneHandler GatewaySubState = 0x82

type GetStateCfm struct {
	GatewayState GatewayState
	SubState     GatewaySubState
	StateData    uint
}

var _ Confirm = (*GetStateCfm)(nil)

func (cfm *GetStateCfm) Code() transport.Command {
	return transport.GW_GET_STATE_CFM
}

func (cfm *GetStateCfm) Read(data []byte) error {
	if len(data) != 6 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var value uint8
	var lvalue uint32

	value, _ = reader.ReadU8()
	cfm.GatewayState = GatewayState(value)

	value, _ = reader.ReadU8()
	cfm.SubState = GatewaySubState(value)

	lvalue, _ = reader.ReadU32()
	cfm.StateData = uint(lvalue)

	return nil
}
