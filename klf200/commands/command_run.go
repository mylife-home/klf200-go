package commands

import (
	"bytes"
	"fmt"
	"klf200/binary"
	"klf200/transport"
	"time"
)

type CommandRunStatusNtf struct {
	SessionID       int
	StatusID        CommandRunOwner
	NodeIndex       int
	NodeParameter   FunctionalParameter
	ParameterValue  int
	RunStatus       CommandRunStatus
	StatusReply     CommandRunStatusReply
	InformationCode uint32
}

type CommandRunOwner int

// The status is from a local user activation. (My self)
const CommandRunOwnerSelf CommandRunOwner = 0x00

// The status is from a user activation.
const CommandRunOwnerUser CommandRunOwner = 0x01

// The status is from a rain sensor activation.
const CommandRunOwnerRain CommandRunOwner = 0x02

// The status is from a timer generated action.
const CommandRunOwnerTimer CommandRunOwner = 0x03

// The status is from a UPS generated action.
const CommandRunOwnerUps CommandRunOwner = 0x05

// The status is from an automatic program generated action. (SAAC)
const CommandRunOwnerProgram CommandRunOwner = 0x08

// The status is from a Wind sensor generated action.
const CommandRunOwnerWind CommandRunOwner = 0x09

// The status is from an actuator generated action.
const CommandRunOwnerMyself CommandRunOwner = 0x0A

// The status is from a automatic cycle generated action.
const CommandRunOwnerAutomaticCycle CommandRunOwner = 0x0B

// The status is from an emergency or a security generated action.
const CommandRunOwnerEmergency CommandRunOwner = 0x0C

// The status is from an unknown command originator action
const CommandRunOwnerUnknown CommandRunOwner = 0xFF

type CommandRunStatus int

// Execution is completed with no errors.
const CommandRunStatusCompleted CommandRunStatus = 0

// Execution has failed. (Get specifics in the following error code)
const CommandRunStatusFailed CommandRunStatus = 1

// Execution is still active
const CommandRunStatusActive CommandRunStatus = 2

type CommandRunStatusReply int

// Used to indicate unknown reply.
const CommandRunStatusReplyUnknownStatusReply CommandRunStatusReply = 0x00

// Indicates no errors detected.
const CommandRunStatusReplyCommandCompletedOk CommandRunStatusReply = 0x01

// Indicates no communication to node.
const CommandRunStatusReplyNoContact CommandRunStatusReply = 0x02

// Indicates manually operated by a user.
const CommandRunStatusReplyManuallyOperated CommandRunStatusReply = 0x03

// Indicates node has been blocked by an object.
const CommandRunStatusReplyBlocked CommandRunStatusReply = 0x04

// Indicates the node contains a wrong system key.
const CommandRunStatusReplyWrongSystemkey CommandRunStatusReply = 0x05

// Indicates the node is locked on this priority level.
const CommandRunStatusReplyPriorityLevelLocked CommandRunStatusReply = 0x06

// Indicates node has stopped in another position than expected.
const CommandRunStatusReplyReachedWrongPosition CommandRunStatusReply = 0x07

// Indicates an error has occurred during execution of command.
const CommandRunStatusReplyErrorDuringExecution CommandRunStatusReply = 0x08

// Indicates no movement of the node parameter.
const CommandRunStatusReplyNoExecution CommandRunStatusReply = 0x09

// Indicates the node is calibrating the parameters.
const CommandRunStatusReplyCalibrating CommandRunStatusReply = 0x0A

// Indicates the node power consumption is too high.
const CommandRunStatusReplyPowerConsumptionTooHigh CommandRunStatusReply = 0x0B

// Indicates the node power consumption is too low.
const CommandRunStatusReplyPowerConsumptionTooLow CommandRunStatusReply = 0x0C

// Indicates door lock errors. (Door open during lock command)
const CommandRunStatusReplyLockPositionOpen CommandRunStatusReply = 0x0D

// Indicates the target was not reached in time.
const CommandRunStatusReplyMotionTimeTooLongCommunicationEnded CommandRunStatusReply = 0x0E

// Indicates the node has gone into thermal protection mode.
const CommandRunStatusReplyThermalProtection CommandRunStatusReply = 0x0F

// Indicates the node is not currently operational.
const CommandRunStatusReplyProductNotOperational CommandRunStatusReply = 0x10

// Indicates the filter needs maintenance.
const CommandRunStatusReplyFilterMaintenanceNeeded CommandRunStatusReply = 0x11

// Indicates the battery level is low.
const CommandRunStatusReplyBatteryLevel CommandRunStatusReply = 0x12

// Indicates the node has modified the target value of the command.
const CommandRunStatusReplyTargetModified CommandRunStatusReply = 0x13

// Indicates this node does not support the mode received.
const CommandRunStatusReplyModeNotImplemented CommandRunStatusReply = 0x14

// Indicates the node is unable to move in the right direction.
const CommandRunStatusReplyCommandIncompatibleToMovement CommandRunStatusReply = 0x15

// Indicates dead bolt is manually locked during unlock command.
const CommandRunStatusReplyUserAction CommandRunStatusReply = 0x16

// Indicates dead bolt error.
const CommandRunStatusReplyDeadBoltError CommandRunStatusReply = 0x17

// Indicates the node has gone into automatic cycle mode.
const CommandRunStatusReplyAutomaticCycleEngaged CommandRunStatusReply = 0x18

// Indicates wrong load on node.
const CommandRunStatusReplyWrongLoadConnected CommandRunStatusReply = 0x19

// Indicates that node is unable to reach received colour code.
const CommandRunStatusReplyColourNotReachable CommandRunStatusReply = 0x1A

// Indicates the node is unable to reach received target position.
const CommandRunStatusReplyTargetNotReachable CommandRunStatusReply = 0x1B

// Indicates io-protocol has received an invalid index.
const CommandRunStatusReplyBadIndexReceived CommandRunStatusReply = 0x1C

// Indicates that the command was overruled by a new command.
const CommandRunStatusReplyCommandOverruled CommandRunStatusReply = 0x1D

// Indicates that the node reported waiting for power.
const CommandRunStatusReplyNodeWaitingForPower CommandRunStatusReply = 0x1E

// Indicates an unknown error code received. (Hex code is shown on display)
const CommandRunStatusReplyInformationCode CommandRunStatusReply = 0xDF

// Indicates the parameter was limited by an unknown device. (Same as LIMITATION_BY_UNKNOWN_DEVICE)
const CommandRunStatusReplyParameterLimited CommandRunStatusReply = 0xE0

// Indicates the parameter was limited by local button.
const CommandRunStatusReplyLimitationByLocalUser CommandRunStatusReply = 0xE1

// Indicates the parameter was limited by a remote control.
const CommandRunStatusReplyLimitationByUser CommandRunStatusReply = 0xE2

// Indicates the parameter was limited by a rain sensor.
const CommandRunStatusReplyLimitationByRain CommandRunStatusReply = 0xE3

// Indicates the parameter was limited by a timer.
const CommandRunStatusReplyLimitationByTimer CommandRunStatusReply = 0xE4

// Indicates the parameter was limited by a power supply.
const CommandRunStatusReplyLimitationByUps CommandRunStatusReply = 0xE6

// Indicates the parameter was limited by an unknown device. (Same as PARAMETER_LIMITED)
const CommandRunStatusReplyLimitationByUnknownDevice CommandRunStatusReply = 0xE7

// Indicates the parameter was limited by a standalone automatic controller.
const CommandRunStatusReplyLimitationBySaac CommandRunStatusReply = 0xEA

// Indicates the parameter was limited by a wind sensor.
const CommandRunStatusReplyLimitationByWind CommandRunStatusReply = 0xEB

// Indicates the parameter was limited by the node itself.
const CommandRunStatusReplyLimitationByMyself CommandRunStatusReply = 0xEC

// Indicates the parameter was limited by an automatic cycle.
const CommandRunStatusReplyLimitationByAutomaticCycle CommandRunStatusReply = 0xED

// Indicates the parameter was limited by an emergency
const CommandRunStatusReplyLimitationByEmergency CommandRunStatusReply = 0xEE

var _ Notify = (*CommandRunStatusNtf)(nil)

func init() {
	registerNotify(func() Notify { return &CommandRunStatusNtf{} })
}

func (ntf *CommandRunStatusNtf) Code() transport.Command {
	return transport.GW_COMMAND_RUN_STATUS_NTF
}

func (ntf *CommandRunStatusNtf) Read(data []byte) error {
	if len(data) != 13 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8
	var u16 uint16
	var u32 uint32

	u16, _ = reader.ReadU16()
	ntf.SessionID = int(u16)

	u8, _ = reader.ReadU8()
	ntf.StatusID = CommandRunOwner(u8)

	u8, _ = reader.ReadU8()
	ntf.NodeIndex = int(u8)

	u8, _ = reader.ReadU8()
	ntf.NodeParameter = FunctionalParameter(u8)

	u16, _ = reader.ReadU16()
	ntf.ParameterValue = int(u16)

	u8, _ = reader.ReadU8()
	ntf.RunStatus = CommandRunStatus(u8)

	u8, _ = reader.ReadU8()
	ntf.StatusReply = CommandRunStatusReply(u8)

	u32, _ = reader.ReadU32()
	ntf.InformationCode = uint32(u32)

	return nil
}

type CommandRemainingTimeNtf struct {
	SessionID     int
	NodeIndex     int
	NodeParameter FunctionalParameter
	Duration      time.Duration
}

var _ Notify = (*CommandRemainingTimeNtf)(nil)

func init() {
	registerNotify(func() Notify { return &CommandRemainingTimeNtf{} })
}

func (ntf *CommandRemainingTimeNtf) Code() transport.Command {
	return transport.GW_COMMAND_REMAINING_TIME_NTF
}

func (ntf *CommandRemainingTimeNtf) Read(data []byte) error {
	if len(data) != 6 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u8 uint8
	var u16 uint16

	u16, _ = reader.ReadU16()
	ntf.SessionID = int(u16)

	u8, _ = reader.ReadU8()
	ntf.NodeIndex = int(u8)

	u8, _ = reader.ReadU8()
	ntf.NodeParameter = FunctionalParameter(u8)

	u16, _ = reader.ReadU16()
	ntf.Duration = time.Second * time.Duration(u16)

	return nil
}

type SessionFinishedNtf struct {
	SessionID int
}

var _ Notify = (*SessionFinishedNtf)(nil)

func init() {
	registerNotify(func() Notify { return &SessionFinishedNtf{} })
}

func (ntf *SessionFinishedNtf) Code() transport.Command {
	return transport.GW_SESSION_FINISHED_NTF
}

func (ntf *SessionFinishedNtf) Read(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))
	var u16 uint16

	u16, _ = reader.ReadU16()
	ntf.SessionID = int(u16)

	return nil
}

func (value CommandRunOwner) String() string {
	switch value {
	case CommandRunOwnerUser:
		return "User"
	case CommandRunOwnerRain:
		return "Rain"
	case CommandRunOwnerTimer:
		return "Timer"
	case CommandRunOwnerUps:
		return "Ups"
	case CommandRunOwnerProgram:
		return "Program"
	case CommandRunOwnerWind:
		return "Wind"
	case CommandRunOwnerMyself:
		return "Myself"
	case CommandRunOwnerAutomaticCycle:
		return "AutomaticCycle"
	case CommandRunOwnerEmergency:
		return "Emergency"
	case CommandRunOwnerUnknown:
		return "Unknown"
	default:
		return fmt.Sprintf("<%d>", value)
	}
}

func (value CommandRunStatus) String() string {
	switch value {
	case CommandRunStatusCompleted:
		return "Completed"
	case CommandRunStatusFailed:
		return "Failed"
	case CommandRunStatusActive:
		return "Active"
	default:
		return fmt.Sprintf("<%d>", value)
	}
}

func (value CommandRunStatusReply) String() string {
	switch value {
	case CommandRunStatusReplyUnknownStatusReply:
		return "UnknownStatusReply"
	case CommandRunStatusReplyCommandCompletedOk:
		return "CommandCompletedOk"
	case CommandRunStatusReplyNoContact:
		return "NoContact"
	case CommandRunStatusReplyManuallyOperated:
		return "ManuallyOperated"
	case CommandRunStatusReplyBlocked:
		return "Blocked"
	case CommandRunStatusReplyWrongSystemkey:
		return "WrongSystemkey"
	case CommandRunStatusReplyPriorityLevelLocked:
		return "PriorityLevelLocked"
	case CommandRunStatusReplyReachedWrongPosition:
		return "ReachedWrongPosition"
	case CommandRunStatusReplyErrorDuringExecution:
		return "ErrorDuringExecution"
	case CommandRunStatusReplyNoExecution:
		return "NoExecution"
	case CommandRunStatusReplyCalibrating:
		return "Calibrating"
	case CommandRunStatusReplyPowerConsumptionTooHigh:
		return "PowerConsumptionTooHigh"
	case CommandRunStatusReplyPowerConsumptionTooLow:
		return "PowerConsumptionTooLow"
	case CommandRunStatusReplyLockPositionOpen:
		return "LockPositionOpen"
	case CommandRunStatusReplyMotionTimeTooLongCommunicationEnded:
		return "MotionTimeTooLongCommunicationEnded"
	case CommandRunStatusReplyThermalProtection:
		return "ThermalProtection"
	case CommandRunStatusReplyProductNotOperational:
		return "ProductNotOperational"
	case CommandRunStatusReplyFilterMaintenanceNeeded:
		return "FilterMaintenanceNeeded"
	case CommandRunStatusReplyBatteryLevel:
		return "BatteryLevel"
	case CommandRunStatusReplyTargetModified:
		return "TargetModified"
	case CommandRunStatusReplyModeNotImplemented:
		return "ModeNotImplemented"
	case CommandRunStatusReplyCommandIncompatibleToMovement:
		return "CommandIncompatibleToMovement"
	case CommandRunStatusReplyUserAction:
		return "UserAction"
	case CommandRunStatusReplyDeadBoltError:
		return "DeadBoltError"
	case CommandRunStatusReplyAutomaticCycleEngaged:
		return "AutomaticCycleEngaged"
	case CommandRunStatusReplyWrongLoadConnected:
		return "WrongLoadConnected"
	case CommandRunStatusReplyColourNotReachable:
		return "ColourNotReachable"
	case CommandRunStatusReplyTargetNotReachable:
		return "TargetNotReachable"
	case CommandRunStatusReplyBadIndexReceived:
		return "BadIndexReceived"
	case CommandRunStatusReplyCommandOverruled:
		return "CommandOverruled"
	case CommandRunStatusReplyNodeWaitingForPower:
		return "NodeWaitingForPower"
	case CommandRunStatusReplyInformationCode:
		return "InformationCode"
	case CommandRunStatusReplyParameterLimited:
		return "ParameterLimited"
	case CommandRunStatusReplyLimitationByLocalUser:
		return "LimitationByLocalUser"
	case CommandRunStatusReplyLimitationByUser:
		return "LimitationByUser"
	case CommandRunStatusReplyLimitationByRain:
		return "LimitationByRain"
	case CommandRunStatusReplyLimitationByTimer:
		return "LimitationByTimer"
	case CommandRunStatusReplyLimitationByUps:
		return "LimitationByUps"
	case CommandRunStatusReplyLimitationByUnknownDevice:
		return "LimitationByUnknownDevice"
	case CommandRunStatusReplyLimitationBySaac:
		return "LimitationBySaac"
	case CommandRunStatusReplyLimitationByWind:
		return "LimitationByWind"
	case CommandRunStatusReplyLimitationByMyself:
		return "LimitationByMyself"
	case CommandRunStatusReplyLimitationByAutomaticCycle:
		return "LimitationByAutomaticCycle"
	case CommandRunStatusReplyLimitationByEmergency:
		return "LimitationByEmergency"
	default:
		return fmt.Sprintf("<%d>", value)
	}
}
