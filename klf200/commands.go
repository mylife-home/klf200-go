package klf200

import (
	"errors"
	"klf200/commands"
	"sync/atomic"
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

type Session struct {
	id int
}

func (cmds *Commands) newSessionId() int {
	const maxUint16 = uint64(^uint16(0))
	sessionId := cmds.sessionIdGen.Add(1) % maxUint16
	return int(sessionId)
}

func (cmds *Commands) ChangePosition(nodeIndexes []int, position commands.MPValue) (*Session, error) {
	sessionId := cmds.newSessionId()

	// TODO: customize parameters
	req := &commands.CommandSendReq{
		SessionID:                 sessionId,
		CommandOriginator:         commands.CommandOriginatorUser,
		PriorityLevel:             commands.PriorityUserLevel2,
		ParameterActive:           commands.FunctionalParameterMP,
		FunctionalParameterValues: map[commands.FunctionalParameter]int{commands.FunctionalParameterMP: int(position)},
		NodeIndexes:               nodeIndexes,
		PriorityLevelLock:         commands.PriorityLevelLockNoNewLock,
		PriorityLevelInfo:         commands.NewPriorityLevelInfo(),
		LockTime:                  commands.NewLockTimeUnlimited(),
	}

	b, _ := req.Write()

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

	return &Session{id: sessionId}, nil
}

// TODO: missing API
