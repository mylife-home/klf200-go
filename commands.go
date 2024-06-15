package klf200

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/mylife-home/klf200-go/commands"
)

type Commands struct {
	client       *Client
	sessionIdGen atomic.Uint64
}

func newCommands(client *Client) *Commands {
	return &Commands{
		client: client,
	}
}

func (cmds *Commands) newSessionId() int {
	const maxUint16 = uint64(^uint16(0))
	sessionId := cmds.sessionIdGen.Add(1) % maxUint16
	return int(sessionId)
}

type Session struct {
	id       int
	notifier Notifier
	ctx      context.Context
	events   chan Event
}

func newSession(client *Client, id int, ctx context.Context) *Session {
	notifier := client.RegisterNotifications([]reflect.Type{
		reflect.TypeOf(&commands.CommandRemainingTimeNtf{}),
		reflect.TypeOf(&commands.CommandRunStatusNtf{}),
		reflect.TypeOf(&commands.SessionFinishedNtf{}),
	})

	sess := &Session{
		id:       id,
		notifier: notifier,
		ctx:      ctx,
		events:   make(chan Event, 100),
	}

	go sess.worker()

	return sess
}

func (sess *Session) Events() <-chan Event {
	return sess.events
}

func (sess *Session) worker() {

	for {

		notif, err := sess.selectNotif()
		if err != nil {
			sess.events <- &RunError{err}
			break
		}

		finished := false

		switch notif := notif.(type) {
		case *commands.CommandRunStatusNtf:
			if sess.id == notif.SessionID {
				sess.events <- &RunStatus{
					StatusID:       notif.StatusID,
					ParameterValue: commands.MPValue(notif.ParameterValue),
					RunStatus:      notif.RunStatus,
					StatusReply:    notif.StatusReply,
				}
			}

		case *commands.CommandRemainingTimeNtf:
			if sess.id == notif.SessionID {
				sess.events <- &RunRemainingTime{
					Duration: notif.Duration,
				}
			}

		case *commands.SessionFinishedNtf:
			if sess.id == notif.SessionID {
				finished = true
			}
		}

		if finished {
			break
		}
	}

	close(sess.events)
	sess.notifier.Close()
}

type Event interface {
}

type RunStatus struct {
	StatusID       commands.CommandRunOwner
	ParameterValue commands.MPValue
	RunStatus      commands.CommandRunStatus
	StatusReply    commands.CommandRunStatusReply
}

type RunRemainingTime struct {
	Duration time.Duration
}

type RunError struct {
	Err error
}

func (sess *Session) selectNotif() (commands.Notify, error) {
	select {
	// TODO: handle disconnection
	case <-sess.ctx.Done():
		return nil, sess.ctx.Err()

	case notif := <-sess.notifier.Stream():
		return notif, nil
	}
}

func (cmds *Commands) ChangePosition(ctx context.Context, nodeIndex int, position commands.MPValue) (*Session, error) {
	sessionId := cmds.newSessionId()

	// TODO: customize parameters
	req := &commands.CommandSendReq{
		SessionID:                 sessionId,
		CommandOriginator:         commands.CommandOriginatorUser,
		PriorityLevel:             commands.PriorityUserLevel2,
		ParameterActive:           commands.FunctionalParameterMP,
		FunctionalParameterValues: map[commands.FunctionalParameter]int{commands.FunctionalParameterMP: int(position)},
		NodeIndexes:               []int{nodeIndex},
		PriorityLevelLock:         commands.PriorityLevelLockNoNewLock,
		PriorityLevelInfo:         commands.NewPriorityLevelInfo(),
		LockTime:                  commands.NewLockTimeUnlimited(),
	}

	cfm, err := cmds.client.execute(req)
	if err != nil {
		return nil, err
	}

	tcfm := cfm.(*commands.CommandSendCfm)

	if !tcfm.Success {
		return nil, errors.New("the request failed")
	}

	if tcfm.SessionID != sessionId {
		return nil, errors.New("session id mismatch")
	}

	return newSession(cmds.client, sessionId, ctx), nil
}

func (cmds *Commands) Mode(ctx context.Context, nodeIndex int) (*Session, error) {
	sessionId := cmds.newSessionId()

	// TODO: customize parameters
	req := &commands.ModeSendReq{
		SessionID:         sessionId,
		CommandOriginator: commands.CommandOriginatorUser,
		PriorityLevel:     commands.PriorityUserLevel2,
		ModeNumber:        0,
		ModeParameter:     0,
		NodeIndexes:       []int{nodeIndex},
		PriorityLevelLock: commands.PriorityLevelLockNoNewLock,
		PriorityLevelInfo: commands.NewPriorityLevelInfo(),
		LockTime:          commands.NewLockTimeUnlimited(),
	}

	cfm, err := cmds.client.execute(req)
	if err != nil {
		return nil, err
	}

	tcfm := cfm.(*commands.ModeSendCfm)

	if tcfm.Status != commands.ModeSendStatusSuccess {
		return nil, fmt.Errorf("error : '%s'", tcfm.Status)
	}

	if tcfm.SessionID != sessionId {
		return nil, errors.New("session id mismatch")
	}

	return newSession(cmds.client, sessionId, ctx), nil
}

func (cmds *Commands) Status(ctx context.Context, nodeIndexes []int) (map[int]StatusData, error) {
	sessionId := cmds.newSessionId()

	// TODO: customize parameters
	req := &commands.StatusRequestReq{
		SessionID:            sessionId,
		NodeIndexes:          nodeIndexes,
		StatusType:           commands.StatusRequestMainInfo,
		FunctionalParameters: make(map[commands.FunctionalParameter]bool),
	}

	n := cmds.client.RegisterNotifications([]reflect.Type{
		reflect.TypeOf(&commands.StatusRequestNtf{}),
		reflect.TypeOf(&commands.SessionFinishedNtf{}),
	})

	cfm, err := cmds.client.execute(req)
	if err != nil {
		return nil, err
	}

	tcfm := cfm.(*commands.StatusRequestCfm)

	if !tcfm.Success {
		return nil, errors.New("the request failed")
	}

	if tcfm.SessionID != sessionId {
		return nil, errors.New("session id mismatch")
	}

	data := make(map[int]StatusData)

	for {
		notif, err := cmds.selectStatusNotif(ctx, n)
		if err != nil {
			return nil, err
		}

		finished := false

		switch notif := notif.(type) {
		case *commands.StatusRequestNtf:
			if sessionId == notif.SessionID {

				statusData := &StatusData{
					StatusID:    notif.StatusID,
					RunStatus:   notif.RunStatus,
					StatusReply: notif.StatusReply,
				}

				// May be nil if status represents an error
				if notif.StatusData != nil {
					mainInfo := notif.StatusData.(*commands.StatusDataMainInfo)

					statusData.TargetPosition = mainInfo.TargetPosition
					statusData.CurrentPosition = mainInfo.CurrentPosition
					statusData.RemainingTime = mainInfo.RemainingTime
					statusData.LastMasterExecutionAddress = mainInfo.LastMasterExecutionAddress
					statusData.LastCommandOriginator = mainInfo.LastCommandOriginator
				}

				data[notif.NodeIndex] = *statusData
			}

		case *commands.SessionFinishedNtf:
			if sessionId == notif.SessionID {
				finished = true
			}
		}

		if finished {
			break
		}
	}

	n.Close()

	return data, nil
}

type StatusData struct {
	StatusID                   commands.CommandRunOwner
	RunStatus                  commands.CommandRunStatus
	StatusReply                commands.CommandRunStatusReply
	TargetPosition             commands.MPValue
	CurrentPosition            commands.MPValue
	RemainingTime              time.Duration
	LastMasterExecutionAddress uint32
	LastCommandOriginator      commands.CommandOriginator
}

func (cmds *Commands) selectStatusNotif(ctx context.Context, n Notifier) (commands.Notify, error) {
	select {
	// TODO: handle disconnection
	case <-ctx.Done():
		return nil, ctx.Err()

	case notif := <-n.Stream():
		return notif, nil
	}
}

// TODO: missing API
