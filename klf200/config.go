package klf200

import (
	"context"
	"klf200/commands"
	"klf200/utils"
	"reflect"
)

type Config struct {
	client        *Client
	sysTableTrans utils.Mutex
}

func newConfig(client *Client) *Config {
	return &Config{
		client:        client,
		sysTableTrans: utils.NewMutex(),
	}
}

func (config *Config) GetSystemTable(ctx context.Context) ([]commands.SystemtableObject, error) {
	// Permits only one request at a time to avoid notifications mismatchs
	if !config.sysTableTrans.TryLockWithContext(ctx) {
		return nil, ctx.Err()
	}

	defer config.sysTableTrans.Unlock()

	n := config.client.RegisterNotifications([]reflect.Type{reflect.TypeOf(&commands.CsGetSystemtableDataNtf{})})

	_, err := config.client.execute(&commands.CsGetSystemtableDataReq{})
	if err != nil {
		return nil, err
	}

	objects := make([]commands.SystemtableObject, 0)

	for {
		// todo timeout
		notif := <-n.Stream()

		packet := notif.(*commands.CsGetSystemtableDataNtf)

		for index := 0; index < packet.NumberOfEntry; index++ {
			object := packet.Objects[index]
			objects = append(objects, object)
		}

		if packet.RemainingNumberOfEntry == 0 {
			break
		}
	}

	n.Close()

	return objects, nil
}

// TODO: missing API