package klf200

import (
	"errors"
	"fmt"
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
	dumpByteSlice(b)

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

func dumpByteSlice(b []byte) {
	var a [16]byte
	n := (len(b) + 15) &^ 15
	for i := 0; i < n; i++ {
		if i%16 == 0 {
			fmt.Printf("%4d", i)
		}
		if i%8 == 0 {
			fmt.Print(" ")
		}
		if i < len(b) {
			fmt.Printf(" %02X", b[i])
		} else {
			fmt.Print("   ")
		}
		if i >= len(b) {
			a[i%16] = ' '
		} else if b[i] < 32 || b[i] > 126 {
			a[i%16] = '.'
		} else {
			a[i%16] = b[i]
		}
		if i%16 == 15 {
			fmt.Printf("  %s\n", string(a[:]))
		}
	}
}
