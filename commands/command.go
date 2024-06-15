package commands

import "github.com/mylife-home/klf200-go/transport"

type Command interface {
	Code() transport.Command
}

type Request interface {
	Command
	NewConfirm() Confirm
	Write() ([]byte, error)
}

type Confirm interface {
	Command
	Read(data []byte) error
}

type Notify interface {
	Command
	Read(data []byte) error
}

var notifyRegistry = make(map[transport.Command]func() Notify)

func registerNotify(builder func() Notify) {
	code := builder().Code()
	notifyRegistry[code] = builder
}

func GetNotify(code transport.Command) Notify {
	builder, ok := notifyRegistry[code]
	if !ok {
		return nil
	}

	return builder()
}

var emptyData = make([]byte, 0)
