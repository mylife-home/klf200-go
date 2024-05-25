package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"klf200"
	"klf200/commands"
	"net"
	"os"
	"sync"
)

type connection struct {
	conn        *tls.Conn
	workersSync sync.WaitGroup
}

func makeConnection(ctx context.Context, address string) (*connection, error) {
	dialer := net.Dialer{}
	netConn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{}
	conf.InsecureSkipVerify = true
	tlsConn := tls.Client(netConn, conf)

	conn := &connection{conn: tlsConn}

	return conn, nil
}

func main() {

	client := klf200.MakeClient(os.Getenv("KLF200_ADDRESS"), os.Getenv("KLF200_PASSWORD"))

	client.RegisterStatusChange(func(cs klf200.ConnectionStatus) {
		fmt.Printf("got status change %d\n", cs)

		if cs == klf200.ConnectionOpen {
			ver, err := client.Version()
			if err != nil {
				panic(err)
			}

			fmt.Printf("version %v\n", ver)
		}
	})

	client.RegisterNotifications(func(n commands.Notify) {
		fmt.Printf("got notify %d\n", n.Code())
	})

	for {
	}
	client.Close()
}

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
