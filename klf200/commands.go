package klf200

import (
	"context"
	"errors"
	"klf200/commands"
	"reflect"
	"sync/atomic"
	"time"
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
	cancel   context.CancelCauseFunc
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
			finished = true
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

// Mode

// Status request

// TODO: missing API
