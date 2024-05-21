package main

import (
	"crypto/tls"
	"fmt"
	"klf200/pkg/commands"
	"klf200/pkg/transport"
	"os"
)

func main() {

	conf := &tls.Config{}
	conf.InsecureSkipVerify = true

	con, err := tls.Dial("tcp4", os.Getenv("KLF200_ADDRESS"), conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TLS OK\n")

	req := &commands.PasswordEnterReq{
		Password: os.Getenv("KLF200_PASSWORD"),
	}

	data, err := req.Write()
	if err != nil {
		panic(err)
	}

	frame := &transport.Frame{
		Cmd:  req.Code(),
		Data: data,
	}

	buff := transport.SlipEncode(frame.Write()).Bytes()
	n, err := con.Write(buff)
	if err != nil {
		panic(err)
	}

	if n != len(buff) {
		panic("bad write len")
	}

	// dumpByteSlice(buff)

	fmt.Printf("WRITE OK\n")

	buff = make([]byte, 500)
	decoder := &transport.SlipDecoder{}

	for {
		n, err = con.Read(buff)
		if err != nil {
			panic(err)
		}

		if n == 0 {
			panic("connection closed")
		}

		fmt.Printf("READ %d\n", n)

		// dumpByteSlice(buff[0:n])

		err = decoder.AddRaw(buff[0:n])
		if err != nil {
			panic(err)
		}

		for {
			buff := decoder.NextFrame()
			if buff == nil {
				break
			}

			frame, err := transport.FrameRead(buff)
			if err != nil {
				panic(err)
			}

			cfm := &commands.PasswordEnterCfm{}
			if frame.Cmd == cfm.Code() {
				cfm.Read(frame.Data)
				fmt.Printf("cfm success %t\n", cfm.Success)
			} else {
				fmt.Printf("frame %d\n", frame.Cmd)
			}
		}
	}

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
