package commands

import (
	"fmt"
	"time"
)

type CommandOriginator int

// User Remote control causing action on actuator
const CommandOriginatorUser CommandOriginator = 1

// Rain sensor
const CommandOriginatorRain CommandOriginator = 2

// Timer controlled
const CommandOriginatorTimer CommandOriginator = 3

// UPS unit
const CommandOriginatorUps CommandOriginator = 5

// Stand Alone Automatic Controls
const CommandOriginatorSaac CommandOriginator = 8

// Wind sensor
const CommandOriginatorWind CommandOriginator = 9

// Managers for requiring a particular electric load shed
const CommandOriginatorLoadShedding CommandOriginator = 11

// Local light sensor
const CommandOriginatorLocalLight CommandOriginator = 12

// Used in context with commands transmitted
// on basis of an unknown sensor for protection
// of an end-product or house goods
const CommandOriginatorUnspecificEnvironmentSensor CommandOriginator = 13

// Used in context with emergency or security commands
const CommandOriginatorEmergency CommandOriginator = 255

type PriorityLevel int

// Provide the most secured level.
//
// Since consequences of misusing this level can deeply impact the
// system behaviour, and therefore the io-homecontrol image, it
// is mandatory for the manufacturer that wants to use this level
// of priority to receive an agreement from io-homecontrol®
//
// In any case the reception of such a command will disable all
// categories (Level 0 to 7).
const PriorityProtectionHuman PriorityLevel = 0

// Used by local sensors that are relative to goods protection: end-
// product protection, house goods protection.
//
// Examples: wind sensor on a terrace awning, rain sensor on a roof window, etc.
const PriorityProtectionEnvironment PriorityLevel = 1

// Used by controller to send one (or a set of one shot) immediate
// action commands when user manually requested for this.
//
// Controllers prescribed as having a higher level of priority than
// others use this level.
//
// For example, this level can be used in combination with a lock
// command on other levels of priority, for providing an exclusive
// access to actuators control.   e.g Parents/Children different
// access rights, ...
const PriorityUserLevel1 PriorityLevel = 2

// Used by controller to send one (or a set of one shot) immediate
// action commands when user manually requested for this.
// This level is the default level used by controllers.
const PriorityUserLevel2 PriorityLevel = 3

// TBD. Don't use
const PriorityComfortLevel1 PriorityLevel = 4

// TBD. Don't use
const PriorityComfortLevel2 PriorityLevel = 5

// TBD. Don't use
const PriorityComfortLevel3 PriorityLevel = 6

// TBD. Don't use
const PriorityComfortLevel4 PriorityLevel = 7

type FunctionalParameter int

// Main Parameter
const FunctionalParameterMP FunctionalParameter = 0

// Functional Parameter number 1
const FunctionalParameterFP1 FunctionalParameter = 1

// Functional Parameter number 2
const FunctionalParameterFP2 FunctionalParameter = 2

// Functional Parameter number 3
const FunctionalParameterFP3 FunctionalParameter = 3

// Functional Parameter number 4
const FunctionalParameterFP4 FunctionalParameter = 4

// Functional Parameter number 5
const FunctionalParameterFP5 FunctionalParameter = 5

// Functional Parameter number 6
const FunctionalParameterFP6 FunctionalParameter = 6

// Functional Parameter number 7
const FunctionalParameterFP7 FunctionalParameter = 7

// Functional Parameter number 8
const FunctionalParameterFP8 FunctionalParameter = 8

// Functional Parameter number 9
const FunctionalParameterFP9 FunctionalParameter = 9

// Functional Parameter number 10
const FunctionalParameterFP10 FunctionalParameter = 10

// Functional Parameter number 11
const FunctionalParameterFP11 FunctionalParameter = 11

// Functional Parameter number 12
const FunctionalParameterFP12 FunctionalParameter = 12

// Functional Parameter number 13
const FunctionalParameterFP13 FunctionalParameter = 13

// Functional Parameter number 14
const FunctionalParameterFP14 FunctionalParameter = 14

// Functional Parameter number 15
const FunctionalParameterFP15 FunctionalParameter = 15

// Functional Parameter number 16
const FunctionalParameterFP16 FunctionalParameter = 16

type MPValue uint16

// 0 = fully opened, 100 = fully closed
func NewMPValueAbsolute(percent int) MPValue {
	// 0000 -> C800
	return MPValue(percent * 0xC800 / 100)
}

// -100% -> go fully opened
// +100% -> go fully closed
func NewMPValueRelative(percent int) MPValue {
	// C900 -> D0D0
	const size = 0xD0D0 - 0xC900
	return MPValue(((percent + 100) * size / 200) + 0xC900)
}

func NewMPValueTarget() MPValue {
	return MPValue(0xD100)
}

func NewMPValueCurrent() MPValue {
	return MPValue(0xD200)
}

func NewMPValueDefault() MPValue {
	return MPValue(0xD300)
}

func NewMPValueIgnore() MPValue {
	return MPValue(0xD400)
}

func (value MPValue) Absolute() (bool, int) {
	if int(value) > 0xC800 {
		return false, 0
	}

	return true, int(value) * 100 / 0xC800
}

func (value MPValue) Relative() (bool, int) {
	if int(value) < 0xC900 || int(value) > 0xD0D0 {
		return false, 0
	}

	const size = 0xD0D0 - 0xC900
	return true, ((int(value) - 0xC900) * 100 / size) - 100
}

func (value MPValue) Target() bool {
	return int(value) == 0xD100
}

func (value MPValue) Current() bool {
	return int(value) == 0xD200
}

func (value MPValue) Default() bool {
	return int(value) == 0xD300
}

func (value MPValue) Ignore() bool {
	return int(value) == 0xD400
}

func (value MPValue) String() string {
	if ok, val := value.Absolute(); ok {
		return fmt.Sprintf("Absolute(%d%%)", val)
	} else if ok, val := value.Relative(); ok {
		return fmt.Sprintf("Relative(%d%%)", val)
	} else if value.Target() {
		return "Target"
	} else if value.Current() {
		return "Current"
	} else if value.Default() {
		return "Default"
	} else if value.Ignore() {
		return "Ignore"
	} else {
		return fmt.Sprintf("?? %d", value)
	}
}

type PriorityLevelLock int

// Do not set a new lock on priority level. Information in the parameters PL_0_3, PL_4_7
// and LockTime are not used. This is the one typically used.
const PriorityLevelLockNoNewLock PriorityLevelLock = 0

// Information in the parameters PL_0_3, PL_4_7 and LockTime are used to lock one or
// more priority level
const PriorityLevelLockNewLock PriorityLevelLock = 1

type PriorityLevelInfo uint16

// TODO: test this accesses
/*
Bit 7-6 = PLI 0
Bit 5-4 = PLI 1
Bit 3-2 = PLI 2
Bit 1-0 = PLI 3

Bit 7-6 = PLI 4
Bit 5-4 = PLI 5
Bit 3-2 = PLI 6
Bit 1-0 = PLI 7
*/

func NewPriorityLevelInfo() PriorityLevelInfo {
	return PriorityLevelInfo(0)
}

func (info PriorityLevelInfo) PLI0() PriorityLevelLevel {
	return PriorityLevelLevel(info << 14 & 0x03)
}

func (info PriorityLevelInfo) PLI1() PriorityLevelLevel {
	return PriorityLevelLevel(info << 12 & 0x03)
}

func (info PriorityLevelInfo) PLI2() PriorityLevelLevel {
	return PriorityLevelLevel(info << 10 & 0x03)
}

func (info PriorityLevelInfo) PLI3() PriorityLevelLevel {
	return PriorityLevelLevel(info << 8 & 0x03)
}

func (info PriorityLevelInfo) PLI4() PriorityLevelLevel {
	return PriorityLevelLevel(info << 6 & 0x03)
}

func (info PriorityLevelInfo) PLI5() PriorityLevelLevel {
	return PriorityLevelLevel(info << 4 & 0x03)
}

func (info PriorityLevelInfo) PLI6() PriorityLevelLevel {
	return PriorityLevelLevel(info << 2 & 0x03)
}

func (info PriorityLevelInfo) PLI7() PriorityLevelLevel {
	return PriorityLevelLevel(info << 0 & 0x03)
}

type PriorityLevelLevel int

// Disable the priority related to the Master
const PriorityLevelLevelDisable PriorityLevelLevel = 0

// Enable the priority related to the Master
const PriorityLevelLevelEnable PriorityLevelLevel = 1

// Enable all pool entry for the specified priority level
// Must be used with caution!
const PriorityLevelLevelEnableAll PriorityLevelLevel = 2

// Do not make any action. When used, the priority setting
// for the specific level will be kept in its current state
const PriorityLevelLevelKeepCurrent PriorityLevelLevel = 3

type LockTime int

func NewLockTimeUnlimited() LockTime {
	return LockTime(255)
}

// Note: should be >= 1 && <= 7650 seconds
func NewLockTime(duration time.Duration) LockTime {
	return LockTime(duration.Seconds() - 1)
}

func (lockTime LockTime) Unlimited() bool {
	return int(lockTime) == 255
}

func (lockTime LockTime) Duration() time.Duration {
	if lockTime.Unlimited() {
		return time.Duration(0)
	}

	return time.Second * time.Duration(int(lockTime)+1)
}
