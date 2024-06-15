package commands

import (
	"bytes"
	"fmt"
	"time"

	"github.com/mylife-home/klf200-go/binary"
	"github.com/mylife-home/klf200-go/transport"
)

type CsGetSystemtableDataReq struct {
}

var _ Request = (*CsGetSystemtableDataReq)(nil)

func (req *CsGetSystemtableDataReq) Code() transport.Command {
	return transport.GW_CS_GET_SYSTEMTABLE_DATA_REQ
}

func (req *CsGetSystemtableDataReq) NewConfirm() Confirm {
	return &CsGetSystemtableDataCfm{}
}

func (req *CsGetSystemtableDataReq) Write() ([]byte, error) {
	return emptyData, nil
}

type CsGetSystemtableDataCfm struct {
}

var _ Confirm = (*CsGetSystemtableDataCfm)(nil)

func (cfm *CsGetSystemtableDataCfm) Code() transport.Command {
	return transport.GW_CS_GET_SYSTEMTABLE_DATA_CFM
}

func (cfm *CsGetSystemtableDataCfm) Read(data []byte) error {
	if len(data) != 0 {
		return fmt.Errorf("bad length")
	}

	return nil
}

type CsGetSystemtableDataNtf struct {
	NumberOfEntry          int
	Objects                []SystemtableObject
	RemainingNumberOfEntry int
}

type SystemtableObject struct {
	SystemTableIndex        int
	ActuatorAddress         uint
	ActuatorType            ActuatorType
	ActuatorSubType         ActuatorSubType
	PowerSaveMode           bool
	RfSupport               bool
	ActuatorTurnaroundTime  time.Duration
	IoManufacturer          IoManufacturer
	BackboneReferenceNumber uint
}

type ActuatorType int

const VenetianBlind ActuatorType = 1
const RollerShutter ActuatorType = 2
const Awning ActuatorType = 3
const WindowOpener ActuatorType = 4
const GarageOpener ActuatorType = 5
const Light ActuatorType = 6
const GateOpener ActuatorType = 7
const RollingDoorOpener ActuatorType = 8
const Lock ActuatorType = 9
const Blind ActuatorType = 10
const Beacon ActuatorType = 12
const DualShutter ActuatorType = 13
const HeatingTemperatureInterface ActuatorType = 14
const OnOffSwitch ActuatorType = 15
const HorizontalAwning ActuatorType = 16
const ExternalVenetianBlind ActuatorType = 17
const LouvreBlind ActuatorType = 18
const CurtainTrack ActuatorType = 19
const VentilationPoint ActuatorType = 20
const ExteriorHeating ActuatorType = 21
const HeatPump ActuatorType = 22
const IntrusionAlarm ActuatorType = 23
const SwingingShutter ActuatorType = 24

type ActuatorSubType int

// TODO: sub type decoding

type IoManufacturer int

const Velux IoManufacturer = 1
const Somfy IoManufacturer = 2
const Honeywell IoManufacturer = 3
const Hormann IoManufacturer = 4
const AssaAbloy IoManufacturer = 5
const Niko IoManufacturer = 6
const WindowMaster IoManufacturer = 7
const Renson IoManufacturer = 8
const Ciat IoManufacturer = 9
const Secuyou IoManufacturer = 10
const Overkiz IoManufacturer = 11
const AtlanticGroup IoManufacturer = 12

var _ Notify = (*CsGetSystemtableDataNtf)(nil)

func init() {
	registerNotify(func() Notify { return &CsGetSystemtableDataNtf{} })
}

func (ntf *CsGetSystemtableDataNtf) Code() transport.Command {
	return transport.GW_CS_GET_SYSTEMTABLE_DATA_NTF
}

func (ntf *CsGetSystemtableDataNtf) Read(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8

	u8, _ = reader.ReadU8()
	ntf.NumberOfEntry = int(u8)

	// object size = 11
	if len(data) != 2+11*ntf.NumberOfEntry {
		return fmt.Errorf("bad length")
	}

	ntf.Objects = make([]SystemtableObject, ntf.NumberOfEntry)

	for index := 0; index < ntf.NumberOfEntry; index++ {
		ntf.readObject(reader, &ntf.Objects[index])
	}

	u8, _ = reader.ReadU8()
	ntf.RemainingNumberOfEntry = int(u8)

	return nil
}

func (ntf *CsGetSystemtableDataNtf) readObject(reader binary.BinaryReader, obj *SystemtableObject) {
	var u8 uint8
	var u16 uint16

	u8, _ = reader.ReadU8()
	obj.SystemTableIndex = int(u8)

	obj.ActuatorAddress = 0
	u8, _ = reader.ReadU8()
	obj.ActuatorAddress |= uint(u8) << 16
	u8, _ = reader.ReadU8()
	obj.ActuatorAddress |= uint(u8) << 8
	u8, _ = reader.ReadU8()
	obj.ActuatorAddress |= uint(u8)

	u16, _ = reader.ReadU16()
	obj.ActuatorType = ActuatorType(u16 >> 6)
	obj.ActuatorSubType = ActuatorSubType(u16 & 0x3F)

	u8, _ = reader.ReadU8()
	obj.PowerSaveMode = u8&1 == 1
	obj.RfSupport = u8&3 == 1

	switch u8 >> 6 {
	case 0:
		obj.ActuatorTurnaroundTime = time.Millisecond * 5
	case 1:
		obj.ActuatorTurnaroundTime = time.Millisecond * 10
	case 2:
		obj.ActuatorTurnaroundTime = time.Millisecond * 20
	case 3:
		obj.ActuatorTurnaroundTime = time.Millisecond * 40
	}

	u8, _ = reader.ReadU8()
	obj.IoManufacturer = IoManufacturer(u8)

	obj.BackboneReferenceNumber = 0
	u8, _ = reader.ReadU8()
	obj.BackboneReferenceNumber |= uint(u8) << 16
	u8, _ = reader.ReadU8()
	obj.BackboneReferenceNumber |= uint(u8) << 8
	u8, _ = reader.ReadU8()
	obj.BackboneReferenceNumber |= uint(u8)
}

func (t ActuatorType) String() string {
	switch t {
	case VenetianBlind:
		return "VenetianBlind"
	case RollerShutter:
		return "RollerShutter"
	case Awning:
		return "Awning"
	case WindowOpener:
		return "WindowOpener"
	case GarageOpener:
		return "GarageOpener"
	case Light:
		return "Light"
	case GateOpener:
		return "GateOpener"
	case RollingDoorOpener:
		return "RollingDoorOpener"
	case Lock:
		return "Lock"
	case Blind:
		return "Blind"
	case Beacon:
		return "Beacon"
	case DualShutter:
		return "DualShutter"
	case HeatingTemperatureInterface:
		return "HeatingTemperatureInterface"
	case OnOffSwitch:
		return "OnOffSwitch"
	case HorizontalAwning:
		return "HorizontalAwning"
	case ExternalVenetianBlind:
		return "ExternalVenetianBlind"
	case LouvreBlind:
		return "LouvreBlind"
	case CurtainTrack:
		return "CurtainTrack"
	case VentilationPoint:
		return "VentilationPoint"
	case ExteriorHeating:
		return "ExteriorHeating"
	case HeatPump:
		return "HeatPump"
	case IntrusionAlarm:
		return "IntrusionAlarm"
	case SwingingShutter:
		return "SwingingShutter"
	default:
		return fmt.Sprintf("<%d>", t)
	}
}

func (m IoManufacturer) String() string {
	switch m {
	case Velux:
		return "Velux"
	case Somfy:
		return "Somfy"
	case Honeywell:
		return "Honeywell"
	case Hormann:
		return "Hormann"
	case AssaAbloy:
		return "AssaAbloy"
	case Niko:
		return "Niko"
	case WindowMaster:
		return "WindowMaster"
	case Renson:
		return "Renson"
	case Ciat:
		return "Ciat"
	case Secuyou:
		return "Secuyou"
	case Overkiz:
		return "Overkiz"
	case AtlanticGroup:
		return "AtlanticGroup"
	default:
		return fmt.Sprintf("<%d>", m)
	}
}
