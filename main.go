package main

import (
	"net"
	"os"
)

func main() {
	con, err := net.Dial("tcp4", os.Getenv("KLF200_ADDRESS"))
	if err != nil {
		panic(err)
	}

	con.Close()
}
