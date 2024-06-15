package klf200

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/mylife-home/klf200-go/klf200/commands"
	"github.com/mylife-home/klf200-go/klf200/transport"
)

const executeTimeout = time.Second * 5
const reconnectTimeout = time.Second * 5
const heartbeatInterval = time.Minute

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

type Notifier interface {
	Stream() <-chan commands.Notify
	Close()
}

type Client struct {
	servAddr string
	password string
	log      Logger

	status                    ConnectionStatus
	connectionStatusCallbacks []func(ConnectionStatus)
	notifiers                 map[*notifier]struct{}
	notifiersLock             sync.Mutex

	ctx         context.Context
	close       context.CancelFunc
	workerSync  sync.WaitGroup
	trans       sync.Mutex
	pendingConf *pendingConfirm

	conn *connection

	device   *Device
	config   *Config
	info     *Info
	commands *Commands
}

func MakeClient(servAddr string, password string, log Logger) *Client {
	ctx, close := context.WithCancel(context.Background())

	client := &Client{
		servAddr:                  servAddr,
		password:                  password,
		log:                       log,
		ctx:                       ctx,
		close:                     close,
		status:                    ConnectionClosed,
		connectionStatusCallbacks: make([]func(ConnectionStatus), 0),
		notifiers:                 make(map[*notifier]struct{}),
	}

	client.device = &Device{client}
	client.config = newConfig(client)
	client.info = newInfo(client)
	client.commands = newCommands(client)

	return client
}

func (client *Client) Start() {
	client.workerSync.Add(1)
	go client.worker()
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

func (client *Client) RegisterNotifications(types []reflect.Type) Notifier {
	n := &notifier{client: client, stream: make(chan commands.Notify, 1000)}

	if types != nil {
		n.types = make(map[reflect.Type]struct{})

		for _, typ := range types {
			n.types[typ] = struct{}{}
		}
	}

	client.notifiersLock.Lock()
	defer client.notifiersLock.Unlock()

	client.notifiers[n] = struct{}{}

	return n
}

type notifier struct {
	client *Client
	types  map[reflect.Type]struct{}
	stream chan commands.Notify
}

func (n *notifier) process(notif commands.Notify) {
	if n.filter(notif) {
		n.stream <- notif

	}
}

func (n *notifier) filter(notif commands.Notify) bool {
	if n.types == nil {
		return true
	}

	_, found := n.types[reflect.TypeOf(notif)]
	return found
}

func (n *notifier) Stream() <-chan commands.Notify {
	return n.stream
}

func (n *notifier) Close() {
	n.client.notifiersLock.Lock()
	defer n.client.notifiersLock.Unlock()

	delete(n.client.notifiers, n)
	close(n.stream)
}

func (client *Client) worker() {
	defer client.workerSync.Done()

	for {
		client.connection()

		select {
		case <-client.ctx.Done():
			return
		case <-time.After(reconnectTimeout):
			// reconnect
		}
	}
}

func (client *Client) connection() {
	client.log.Infof("Dial to '%s'", client.servAddr)

	conn, err := makeConnection(client.ctx, client.servAddr, client.log)
	if err != nil {
		client.log.WithError(err).Errorf("Could not connect to '%s'", client.servAddr)
		return
	}

	client.conn = conn
	defer func() {
		client.conn.Close()
		client.conn = nil
		client.changeStatus(ConnectionClosed)
		client.log.Infof("Connection closed")
	}()

	client.changeStatus(ConnectionHandshaking)
	client.log.Debugf("Start handshake")

	if err := handshake(client.ctx, client.conn, client.password); err != nil {
		client.log.WithError(err).Error("Handshake failed")
		return
	}

	client.changeStatus(ConnectionOpen)
	client.log.Debugf("Handshake done")

	defer func() {
		pendingConf := client.pendingConf
		if pendingConf != nil {
			pendingConf.Cancel()
		}
	}()

	for {
		select {
		case <-client.ctx.Done():
			return

		case <-time.After(heartbeatInterval):
			go client.heartbeat()

		case err := <-client.conn.Errors():
			client.log.WithError(err).Error("Error on connection")
			return

		case frame := <-client.conn.Read():
			client.processFrame(frame)
		}
	}
}

func (client *Client) processFrame(frame *transport.Frame) {
	// try to read it as notify
	notify := commands.GetNotify(frame.Cmd)
	if notify != nil {
		if err := notify.Read(frame.Data); err != nil {
			client.log.WithError(err).Errorf("Cannot read frame %s", frame.Cmd)
			return
		}

		client.notifiersLock.Lock()
		defer client.notifiersLock.Unlock()

		for n := range client.notifiers {
			n.process(notify)
		}

		return
	}

	pendingConf := client.pendingConf
	if pendingConf != nil {
		pendingConf.Confirm(frame)

		return
	}

	client.log.Warnf("Got unmatched frame %s", frame.Cmd)
}

func (client *Client) send(conn *connection, req commands.Request) error {

	data, err := req.Write()
	if err != nil {
		return err
	}

	frame := &transport.Frame{
		Cmd:  req.Code(),
		Data: data,
	}

	conn.Write(frame)

	return nil
}

func (client *Client) execute(req commands.Request) (commands.Confirm, error) {
	conn := client.conn
	if conn == nil {
		return nil, errors.New("not connected")
	}

	client.trans.Lock()
	defer client.trans.Unlock()

	pendingConf := newPendingConfirm()
	client.pendingConf = pendingConf
	defer func() {
		client.pendingConf = nil
	}()

	if err := client.send(conn, req); err != nil {
		return nil, err
	}

	frame, err := pendingConf.Wait()
	if err != nil {
		return nil, err
	}

	cfm := req.NewConfirm()
	if frame.Cmd != cfm.Code() {
		return nil, fmt.Errorf("unexpected confirm: got %d, expected %d", frame.Cmd, cfm.Code())
	}

	if err := cfm.Read(frame.Data); err != nil {
		return nil, err
	}

	return cfm, nil
}

type pendingConfirm struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	cfm    chan *transport.Frame
}

func newPendingConfirm() *pendingConfirm {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, executeTimeout)
	ctx, cancel := context.WithCancelCause(ctx)
	cfm := make(chan *transport.Frame)

	return &pendingConfirm{ctx, cancel, cfm}
}

func (pc *pendingConfirm) Cancel() {
	pc.cancel(errors.New("connection closed"))
}

func (pc *pendingConfirm) Confirm(frame *transport.Frame) {
	pc.cfm <- frame
}

func (pc *pendingConfirm) Wait() (*transport.Frame, error) {
	select {
	case <-pc.ctx.Done():
		return nil, pc.ctx.Err()
	case cfm := <-pc.cfm:
		pc.cancel(nil)
		return cfm, nil
	}
}

func (client *Client) heartbeat() {
	if _, err := client.device.GetState(); err != nil {
		client.log.WithError(err).Errorf("Heartbeat error")
	}

	client.log.Debugf("Heartbeat OK")
}

func (client *Client) Device() *Device {
	return client.device
}

func (client *Client) Config() *Config {
	return client.config
}

func (client *Client) Info() *Info {
	return client.info
}

// TODO: ActivationLog

func (client *Client) Commands() *Commands {
	return client.commands
}

// TODO: Scenes

// TODO: ContactInput
