package commands

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mylife-home/klf200-go/internal/binary"
	"github.com/mylife-home/klf200-go/internal/transport"
)

type GetAllNodesInformationReq struct {
}

var _ Request = (*GetAllNodesInformationReq)(nil)

func (req *GetAllNodesInformationReq) Code() transport.Command {
	return transport.GW_GET_ALL_NODES_INFORMATION_REQ
}

func (req *GetAllNodesInformationReq) NewConfirm() Confirm {
	return &GetAllNodesInformationCfm{}
}

func (req *GetAllNodesInformationReq) Write() ([]byte, error) {
	return emptyData, nil
}

type GetAllNodesInformationCfm struct {
	Success            bool
	TotalNumberOfNodes int
}

var _ Confirm = (*GetAllNodesInformationCfm)(nil)

func (cfm *GetAllNodesInformationCfm) Code() transport.Command {
	return transport.GW_GET_ALL_NODES_INFORMATION_CFM
}

func (cfm *GetAllNodesInformationCfm) Read(data []byte) error {
	if len(data) != 2 {
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

	cfm.TotalNumberOfNodes = int(data[1])

	return nil
}

type Velocity int

// The node operates by its default velocity
const VelocityDefault Velocity = 0

// The node operates in silent mode (slow)
const VelocitySilent Velocity = 1

// The node operates with fast velocity
const VelocityFast Velocity = 2

// Not supported by node
const VelocityNotAvailable Velocity = 255

type NodeTypeSubType int

// TODO

type ProductType int

// TODO

type NodeVariation int

// Not set
const NodeVariationNotSet NodeVariation = 0

// Window is a top hung window
const NodeVariationTopHung NodeVariation = 1

// Window is a kip window
const NodeVariationKip NodeVariation = 2

// Window is a flat roof
const NodeVariationFlatRoof NodeVariation = 3

// Window is a sky light
const NodeVariationSkyLight NodeVariation = 4

type PowerMode int

const PowerModeAlwaysAlive PowerMode = 0
const PowerModeLowPowerMode PowerMode = 1

type NodeState int

// This status information is only returned about an ACTIAVTE_FUNC, an
// ACTIVATE_MODE, an ACTIVATE_STATE or a WINK command.
//
// The parameter is unable to execute due to given conditions. An example can be that
// the temperature is too high. It indicates that the parameter could not execute per
// the contents of the present activate command.
const NodeStateNonExecuting NodeState = 0

// This status information is only returned about an ACTIVATE_STATUS_REQ
// command.
//
// An error has occurred while executing. This error information will be cleared the
// next time the parameter is going into ‘Waiting for executing’, ‘Waiting for power’ or
// ‘Executing’.
//
// A parameter can have the execute status ‘Error while executing’ only if the previous
// execute status was ‘Executing’. Note that this execute status gives information
// about the previous execution of the parameter, and gives no indication whether the
// following execution will fail.
const NodeStateErrorWhileExecution NodeState = 1

const NodeStateNotUsed NodeState = 2

// The parameter is waiting for power to proceed execution
const NodeStateWaitingForPower = 3

// Execution for the parameter is in progress
const NodeStateExecuting NodeState = 4

// The parameter is not executing and no error has been detected. No activation of the
// parameter has been initiated. The parameter is ready for activation.
const NodeStateDone NodeState = 5

// The state is unknown
const NodeStateUnknown NodeState = 255

type NodePosition int

const NodePositionMin NodePosition = 0x0000
const NodePositionMax NodePosition = 0xC800

// No feed-back value known
const NodePositionUnknown NodePosition = 0xF7FF

type NodeAliasId int

// A position a window can be opened to for getting some ventilation and
// where the window is still locked
const NodeAliasSecuredVentilation NodeAliasId = 0xD803

type GetAllNodesInformationNtf struct {
	NodeID             int
	Order              int
	Placement          int
	Name               string
	Velocity           Velocity
	NodeTypeSubType    NodeTypeSubType
	ProductGroup       int
	ProductType        ProductType
	NodeVariation      NodeVariation
	PowerMode          PowerMode
	BuildNumber        int
	SerialNumber       int
	State              NodeState
	CurrentPosition    NodePosition
	Target             NodePosition
	FP1CurrentPosition NodePosition
	FP2CurrentPosition NodePosition
	FP3CurrentPosition NodePosition
	FP4CurrentPosition NodePosition
	RemainingTime      time.Duration
	TimeStamp          time.Time
	Aliases            map[NodeAliasId]int
}

var _ Notify = (*GetAllNodesInformationNtf)(nil)

func init() {
	registerNotify(func() Notify { return &GetAllNodesInformationNtf{} })
}

func (ntf *GetAllNodesInformationNtf) Code() transport.Command {
	return transport.GW_GET_ALL_NODES_INFORMATION_NTF
}

func (ntf *GetAllNodesInformationNtf) Read(data []byte) error {
	if len(data) != 124 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8
	var u16 uint16
	var u32 uint32

	u8, _ = reader.ReadU8()
	ntf.NodeID = int(u8)

	u16, _ = reader.ReadU16()
	ntf.Order = int(u16)

	u8, _ = reader.ReadU8()
	ntf.Placement = int(u8)

	name := make([]byte, 64)
	reader.Read(name)
	ntf.Name = string(bytes.TrimRight(name, "\x00"))

	u8, _ = reader.ReadU8()
	ntf.Velocity = Velocity(u8)

	u16, _ = reader.ReadU16()
	ntf.NodeTypeSubType = NodeTypeSubType(u16)

	u8, _ = reader.ReadU8()
	ntf.ProductGroup = int(u8)

	u8, _ = reader.ReadU8()
	ntf.ProductType = ProductType(u8)

	u8, _ = reader.ReadU8()
	ntf.NodeVariation = NodeVariation(u8)

	u8, _ = reader.ReadU8()
	ntf.PowerMode = PowerMode(u8)

	u8, _ = reader.ReadU8()
	ntf.BuildNumber = int(u8)

	u16, _ = reader.ReadU16()
	ntf.BuildNumber = int(u16)

	u8, _ = reader.ReadU8()
	ntf.State = NodeState(u8)

	u16, _ = reader.ReadU16()
	ntf.CurrentPosition = NodePosition(u16)

	u16, _ = reader.ReadU16()
	ntf.Target = NodePosition(u16)

	u16, _ = reader.ReadU16()
	ntf.FP1CurrentPosition = NodePosition(u16)

	u16, _ = reader.ReadU16()
	ntf.FP2CurrentPosition = NodePosition(u16)

	u16, _ = reader.ReadU16()
	ntf.FP3CurrentPosition = NodePosition(u16)

	u16, _ = reader.ReadU16()
	ntf.FP4CurrentPosition = NodePosition(u16)

	u16, _ = reader.ReadU16()
	ntf.RemainingTime = time.Second * time.Duration(u16)

	u32, _ = reader.ReadU32()
	ntf.TimeStamp = time.Unix(int64(u32), 0)

	ntf.Aliases = make(map[NodeAliasId]int)

	nbrOfAlias, _ := reader.ReadU8()
	for index := 0; index < int(nbrOfAlias); index++ {
		typ, _ := reader.ReadU16()
		value, _ := reader.ReadU16()

		ntf.Aliases[NodeAliasId(typ)] = int(value)
	}

	return nil
}

type GetAllNodesInformationFinishedNtf struct {
}

var _ Notify = (*GetAllNodesInformationFinishedNtf)(nil)

func init() {
	registerNotify(func() Notify { return &GetAllNodesInformationFinishedNtf{} })
}

func (ntf *GetAllNodesInformationFinishedNtf) Code() transport.Command {
	return transport.GW_GET_ALL_NODES_INFORMATION_FINISHED_NTF
}

func (ntf *GetAllNodesInformationFinishedNtf) Read(data []byte) error {
	if len(data) != 0 {
		return fmt.Errorf("bad length")
	}

	return nil
}

func (v Velocity) String() string {
	switch v {
	case VelocityDefault:
		return "VelocityDefault"
	case VelocitySilent:
		return "VelocitySilent"
	case VelocityFast:
		return "VelocityFast"
	case VelocityNotAvailable:
		return "VelocityNotAvailable"
	default:
		return fmt.Sprintf("<%d>", v)
	}
}

func (v NodeVariation) String() string {
	switch v {
	case NodeVariationNotSet:
		return "NodeVariationNotSet"
	case NodeVariationTopHung:
		return "NodeVariationTopHung"
	case NodeVariationKip:
		return "NodeVariationKip"
	case NodeVariationFlatRoof:
		return "NodeVariationFlatRoof"
	case NodeVariationSkyLight:
		return "NodeVariationSkyLight"
	default:
		return fmt.Sprintf("<%d>", v)
	}
}

func (pm PowerMode) String() string {
	switch pm {
	case PowerModeAlwaysAlive:
		return "PowerModeAlwaysAlive"
	case PowerModeLowPowerMode:
		return "PowerModeLowPowerMode"
	default:
		return fmt.Sprintf("<%d>", pm)
	}
}

func (s NodeState) String() string {
	switch s {
	case NodeStateNonExecuting:
		return "NodeStateNonExecuting"
	case NodeStateErrorWhileExecution:
		return "NodeStateErrorWhileExecution"
	case NodeStateNotUsed:
		return "NodeStateNotUsed"
	case NodeStateWaitingForPower:
		return "NodeStateWaitingForPower"
	case NodeStateExecuting:
		return "NodeStateExecuting"
	case NodeStateDone:
		return "NodeStateDone"
	case NodeStateUnknown:
		return "NodeStateUnknown"
	default:
		return fmt.Sprintf("<%d>", s)
	}
}
