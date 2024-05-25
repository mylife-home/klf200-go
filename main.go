package main

import (
	"fmt"
	"klf200"
	"klf200/commands"
	"os"
)

func main() {

	client := klf200.MakeClient(os.Getenv("KLF200_ADDRESS"), os.Getenv("KLF200_PASSWORD"))

	client.RegisterNotifications(notify)

	client.RegisterStatusChange(func(cs klf200.ConnectionStatus) {
		fmt.Printf("got status change %d\n", cs)

		if cs == klf200.ConnectionOpen {
			open(client)
		}
	})

	client.Start()

	for {
	}

	client.Close()
}

func notify(n commands.Notify) {
	fmt.Printf("got notify %d\n", n.Code())
}

func open(client *klf200.Client) {
	ver, err := client.Version()
	if err != nil {
		panic(err)
	}

	fmt.Printf("SoftwareVersion.CommandVersionNumber = %d\n", ver.SoftwareVersion.CommandVersionNumber)
	fmt.Printf("SoftwareVersion.VersionWholeNumber = %d\n", ver.SoftwareVersion.VersionWholeNumber)
	fmt.Printf("SoftwareVersion.VersionSubNumber = %d\n", ver.SoftwareVersion.VersionSubNumber)
	fmt.Printf("SoftwareVersion.BranchID = %d\n", ver.SoftwareVersion.BranchID)
	fmt.Printf("SoftwareVersion.BuildNumber = %d\n", ver.SoftwareVersion.BuildNumber)
	fmt.Printf("SoftwareVersion.MicroBuild = %d\n", ver.SoftwareVersion.MicroBuild)
	fmt.Printf("HardwareVersion = %d\n", ver.HardwareVersion)
	fmt.Printf("ProductGroup = %d\n", ver.ProductGroup)
	fmt.Printf("ProductType = %d\n", ver.ProductType)

	state, err := client.State()
	if err != nil {
		panic(err)
	}

	fmt.Printf("GatewayState = %d\n", state.GatewayState)
	fmt.Printf("SubState = %d\n", state.SubState)
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
