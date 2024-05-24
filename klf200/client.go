package klf200

import (
	"context"
	"fmt"
	"itv2-go/itv2/commands"
	"log"
	"sync"
	"time"
)

type ConnectionStatus uint8

const (
	ConnectionClosed      ConnectionStatus = 0
	ConnectionHandshaking ConnectionStatus = 1
	ConnectionOpen        ConnectionStatus = 2
)

func (status ConnectionStatus) String() string {
	switch status {
	case ConnectionClosed:
		return "Closed"
	case ConnectionHandshaking:
		return "Handshaking"
	case ConnectionOpen:
		return "Open"
	default:
		return fmt.Sprintf("<%d>", status)
	}
}

const heartbeatInterval = time.Second * 5

type Client struct {
	servAddr string
	password string

	status                    ConnectionStatus
	connectionStatusCallbacks []func(ConnectionStatus)
	notificationsCallbacks    []func(commands.Command)

	ctx        context.Context
	close      context.CancelFunc
	workerSync sync.WaitGroup

	conn *connection
}

func MakeClient(servAddr string, password string) *Client {
	ctx, close := context.WithCancel(context.Background())

	client := &Client{
		servAddr:                  servAddr,
		password:                  password,
		ctx:                       ctx,
		close:                     close,
		status:                    ConnectionClosed,
		connectionStatusCallbacks: make([]func(ConnectionStatus), 0),
		notificationsCallbacks:    make([]func(commands.Command), 0),
	}

	client.workerSync.Add(1)

	go client.worker()

	return client
}

func (client *Client) Close() {
	client.close()
	client.workerSync.Wait()
}

func (client *Client) changeStatus(newStatus ConnectionStatus) {
	if client.status == newStatus {
		return
	}

	client.status = newStatus

	for _, callback := range client.connectionStatusCallbacks {
		go callback(newStatus)
	}
}

func (client *Client) Status() ConnectionStatus {
	return client.status
}

func (client *Client) RegisterStatusChange(callback func(ConnectionStatus)) {
	client.connectionStatusCallbacks = append(client.connectionStatusCallbacks, callback)
}

func (client *Client) RegisterNotifications(callback func(commands.Command)) {
	client.notificationsCallbacks = append(client.notificationsCallbacks, callback)
}

func (client *Client) worker() {
	defer client.workerSync.Done()

	for {
		client.connection()

		select {
		case <-client.ctx.Done():
			return
		case <-time.After(time.Second * 5):
			// reconnect
		}
	}
}

func (client *Client) connection() {
	log.Printf("Dial to '%s'", client.servAddr)

	conn, err := makeConnection(client.ctx, client.servAddr)
	if err != nil {
		log.Printf("Could not connect to '%s': %s", client.servAddr, err)
		return
	}

	client.conn = conn
	defer func() {
		client.conn.Close()
		client.conn = nil
		client.changeStatus(ConnectionClosed)
		log.Printf("Connection closed")
	}()

	client.changeStatus(ConnectionHandshaking)
	log.Printf("Start handshake")

	if err := handshake(client.ctx, client.conn, client.password); err != nil {
		log.Printf("Handshake failed: %s", err)
		return
	}

	client.changeStatus(ConnectionOpen)
	log.Printf("Handshake done")
	/*
	   client.transactions = newTransactionManager()

	   	defer func() {
	   		client.transactions.CancelAll()
	   		client.transactions = nil
	   	}()

	   	for {
	   		select {
	   		case <-client.ctx.Done():
	   			return

	   		case <-time.After(heartbeatInterval):
	   			go client.heartbeat()

	   		case err := <-client.conn.Errors():
	   			log.Printf("Error on connection: %s", err)
	   			return

	   		case cmd := <-client.conn.Read():
	   			client.processCommand(cmd)
	   		}
	   	}
	*/
}

/*
func (client *Client) processCommand(cmd commands.Command) {
	if client.transactions.ProcessCommand(cmd) {
		return
	}

	// Not matched by transaction manager, consider it notification
	for _, callback := range client.notificationsCallbacks {
		go callback(cmd)
	}
}

// send command and expect basic response
func (client *Client) executeCommand(cmd commands.CommandWithAppSeq) error {
	_, err := client.execCmdInternal(cmd)
	return err
}

// send request and expect response
func (client *Client) executeRequest(req commands.RequestData) (commands.ResponseData, error) {
	cmd := &commands.Request{
		ReqCode: req.RequestCode(),
		ReqData: req,
	}

	res, err := client.execCmdInternal(cmd)

	var resData commands.ResponseData
	if err == nil {
		resData = res.(commands.ResponseData)
	}

	return resData, err
}

func (client *Client) execCmdInternal(cmd commands.CommandWithAppSeq) (commands.Command, error) {
	// Note: can be set to nil in the middle
	conn := client.conn
	transactions := client.transactions

	if conn == nil || transactions == nil {
		return nil, fmt.Errorf("not connected")
	}

	cmd.SetAppSeq(conn.NextAppSeq())

	conn.Write(cmd)

	transaction := makeTransaction(cmd)
	transactions.addTransaction(transaction)
	res, err := transaction.Wait()
	transactions.removeTransaction(transaction)

	return res, err
}

func (client *Client) heartbeat() {
	cmd := &commands.UserActivity{
		PartitionNumber: &serialization.VarBytes{},
		Type:            4,
	}

	if err := client.executeCommand(cmd); err != nil {
		log.Printf("Heartbeat error: %s", err)
	}

	log.Printf("Heartbeat OK")
}
*/
