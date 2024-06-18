package klf200

import (
	"context"
	"errors"
	"sync"

	"github.com/mylife-home/klf200-go/transport"
)

type connection struct {
	sock       *socket
	write      chan *transport.Frame
	read       chan *transport.Frame
	errors     chan error
	exit       chan struct{}
	decoder    transport.SlipDecoder
	workerSync sync.WaitGroup
	log        Logger
}

var errConnectionRemotelyClosed = errors.New("connection closed by remote side")

func makeConnection(ctx context.Context, address string, log Logger) (*connection, error) {
	sock, err := makeSocket(ctx, address)
	if err != nil {
		return nil, err
	}

	conn := &connection{
		sock:   sock,
		write:  make(chan *transport.Frame, 10),
		read:   make(chan *transport.Frame, 10),
		errors: make(chan error, 10),
		exit:   make(chan struct{}, 1),
		log:    log,
	}

	conn.workerSync.Add(1)
	go conn.worker()

	return conn, nil
}

func (conn *connection) worker() {
	defer conn.workerSync.Done()

	for {
		select {
		case <-conn.exit:
			return

		case err := <-conn.sock.Errors():
			conn.errors <- err

		case data := <-conn.sock.Read():
			conn.processRead(data)

		case frame := <-conn.write:
			conn.processWrite(frame)
		}
	}
}

func (conn *connection) processRead(data []byte) {
	if len(data) == 0 {
		conn.errors <- errConnectionRemotelyClosed
		return
	}

	if err := conn.decoder.AddRaw(data); err != nil {
		conn.errors <- err
		return
	}

	for {
		buff := conn.decoder.NextFrame()
		if buff == nil {
			break
		}

		frame, err := transport.FrameRead(buff)
		if err != nil {
			conn.errors <- err
			return
		}

		// conn.log.Debugf("Recv frame %v", frame)
		conn.read <- frame
	}
}

func (conn *connection) processWrite(frame *transport.Frame) {
	// conn.log.Debugf("Send frame %v", frame)

	buffer := transport.SlipEncode(frame.Write())
	conn.sock.Write(buffer.Bytes())
}

func (conn *connection) Write(frame *transport.Frame) {
	conn.write <- frame
}

func (conn *connection) Read() <-chan *transport.Frame {
	return conn.read
}

func (conn *connection) Errors() <-chan error {
	return conn.errors
}

func (conn *connection) Close() {
	close(conn.exit)
	conn.workerSync.Wait()

	conn.sock.Close()
}
