package main

import "net"

func main() {
	con, err := net.Dial("tcp4", "klf200.mti-team2.dyndns.org:51200")
	if err != nil {
		panic(err)
	}

	con.Close()
}
