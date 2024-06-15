package commands

import (
	"bytes"
	"fmt"

	"github.com/mylife-home/klf200-go/binary"
	"github.com/mylife-home/klf200-go/transport"
)

type ErrorNumber uint8

type ErrorNtf struct {
	ErrorNumber ErrorNumber
}

var _ Notify = (*ErrorNtf)(nil)

func init() {
	registerNotify(func() Notify { return &ErrorNtf{} })
}

func (ntf *ErrorNtf) Code() transport.Command {
	return transport.GW_ERROR_NTF
}

func (ntf *ErrorNtf) Read(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))

	e, _ := reader.ReadU8()
	ntf.ErrorNumber = ErrorNumber(e)

	return nil
}

// Not further defined error.
const ErrorUnknown ErrorNumber = 0

// Unknown Command or command is not accepted at this state.
const ErrorBadCommand ErrorNumber = 1

// ERROR on Frame Structure.
const ErrorBadFrame ErrorNumber = 2

// Busy. Try again later.
const ErrorBusy ErrorNumber = 7

// Bad system table index.
const ErrorBadIndex ErrorNumber = 8

// Not authenticated
const ErrorNotAuthenticated ErrorNumber = 12

func (e ErrorNumber) String() string {
	switch e {
	case ErrorUnknown:
		return "ErrorUnknown"
	case ErrorBadCommand:
		return "ErrorBadCommand"
	case ErrorBadFrame:
		return "ErrorBadFrame"
	case ErrorBusy:
		return "ErrorBusy"
	case ErrorBadIndex:
		return "ErrorBadIndex"
	case ErrorNotAuthenticated:
		return "ErrorNotAuthenticated"
	default:
		return fmt.Sprintf("<%d>", e)
	}
}
