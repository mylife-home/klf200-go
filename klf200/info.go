package klf200

import (
	"context"
	"errors"
	"fmt"
	"klf200/commands"
	"klf200/utils"
	"reflect"
)

type Info struct {
	client          *Client
	getAllInfoTrans utils.Mutex
}

func newInfo(client *Client) *Info {
	return &Info{
		client:          client,
		getAllInfoTrans: utils.NewMutex(),
	}
}

func (info *Info) GetAllNodesInformation(ctx context.Context) ([]*commands.GetAllNodesInformationNtf, error) {
	// Permits only one request at a time to avoid notifications mismatchs
	if !info.getAllInfoTrans.TryLockWithContext(ctx) {
		return nil, ctx.Err()
	}

	defer info.getAllInfoTrans.Unlock()

	n := info.client.RegisterNotifications([]reflect.Type{
		reflect.TypeOf(&commands.GetAllNodesInformationNtf{}),
		reflect.TypeOf(&commands.GetAllNodesInformationFinishedNtf{}),
	})

	cfm, err := info.client.execute(&commands.GetAllNodesInformationReq{})
	if err != nil {
		return nil, err
	}

	tcfm := cfm.(*commands.GetAllNodesInformationCfm)

	if !tcfm.Success {
		return nil, errors.New("system table empty")
	}

	nodes := make([]*commands.GetAllNodesInformationNtf, 0, tcfm.TotalNumberOfNodes)

	for {
		notif, err := info.selectNotif(ctx, n)
		if err != nil {
			return nil, err
		}

		exit := false

		switch notif := notif.(type) {
		case *commands.GetAllNodesInformationNtf:
			nodes = append(nodes, notif)
		case *commands.GetAllNodesInformationFinishedNtf:
			exit = true
		}

		if exit {
			break
		}
	}

	if len(nodes) != tcfm.TotalNumberOfNodes {
		return nil, fmt.Errorf("nodes count mismatch (ntf=%d, cfm=%d)", len(nodes), tcfm.TotalNumberOfNodes)
	}

	n.Close()

	return nodes, nil
}

func (info *Info) selectNotif(ctx context.Context, n Notifier) (commands.Notify, error) {
	select {
	// TODO: handle disconnection
	case <-ctx.Done():
		return nil, ctx.Err()

	case notif := <-n.Stream():
		return notif, nil
	}
}

// GW_GET_ALL_NODES_INFORMATION_REQ
// GW_GET_ALL_NODES_INFORMATION_CFM
// GW_GET_ALL_NODES_INFORMATION_NTF
// GW_GET_ALL_NODES_INFORMATION_FINISHED_NTF

// TODO: missing API
