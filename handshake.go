package klf200

import (
	"context"
	"errors"
	"fmt"

	"github.com/mylife-home/klf200-go/commands"
	"github.com/mylife-home/klf200-go/transport"
)

type handshakeData struct {
	ctx      context.Context
	conn     *connection
	password string
}

func handshake(ctx context.Context, conn *connection, password string) error {
	handshake := &handshakeData{
		ctx:      ctx,
		conn:     conn,
		password: password,
	}

	return handshake.execute()
}

func (handshake *handshakeData) execute() error {
	if err := handshake.handshakeAuthenticate(); err != nil {
		return err
	}

	if err := handshake.handshakeCheckVersion(); err != nil {
		return err
	}

	return nil
}

func (handshake *handshakeData) handshakeAuthenticate() error {

	req := &commands.PasswordEnterReq{
		Password: handshake.password,
	}

	cfm, err := handshake.sendCommandWithResponse(req)
	if err != nil {
		return err
	}

	tcfm := cfm.(*commands.PasswordEnterCfm)
	if !tcfm.Success {
		return errors.New("authentication failed")
	}

	return nil
}

func (handshake *handshakeData) handshakeCheckVersion() error {
	req := &commands.GetProtocolVersionReq{}

	cfm, err := handshake.sendCommandWithResponse(req)
	if err != nil {
		return err
	}

	tcfm := cfm.(*commands.GetProtocolVersionCfm)

	if tcfm.MajorVersion == 3 && tcfm.MinorVersion == 14 {
		handshake.conn.log.Infof("Protocol version : %d.%d", tcfm.MajorVersion, tcfm.MinorVersion)
	} else {
		handshake.conn.log.Warnf("Protocol version : %d.%d, expected 3.14", tcfm.MajorVersion, tcfm.MinorVersion)
	}

	return nil
}

func (handshake *handshakeData) sendCommandWithResponse(req commands.Request) (commands.Confirm, error) {

	data, err := req.Write()
	if err != nil {
		return nil, err
	}

	frame := &transport.Frame{
		Cmd:  req.Code(),
		Data: data,
	}

	handshake.send(frame)

	frame, err = handshake.receive()
	if err != nil {
		return nil, err
	}

	cfm := req.NewConfirm()

	if frame.Cmd != cfm.Code() {
		return nil, fmt.Errorf("received unexpected frame %d (expected %d)", frame.Cmd, cfm.Code())
	}

	if err := cfm.Read(frame.Data); err != nil {
		return nil, err
	}

	return cfm, nil
}

func (handshake *handshakeData) send(frame *transport.Frame) {
	handshake.conn.Write(frame)
}

func (handshake *handshakeData) receive() (*transport.Frame, error) {
	select {
	case <-handshake.ctx.Done():
		return nil, errors.New("client closing")

	case frame := <-handshake.conn.Read():
		return frame, nil

	case err := <-handshake.conn.Errors():
		return nil, err
	}
}
