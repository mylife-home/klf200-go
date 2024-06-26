package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mylife-home/klf200-go"
)

func main() {

	client := klf200.MakeClient(os.Getenv("KLF200_ADDRESS"), os.Getenv("KLF200_PASSWORD"), &logger{})

	//client.RegisterNotifications(notify)

	client.RegisterStatusChange(func(cs klf200.ConnectionStatus) {
		fmt.Printf("got status change %d\n", cs)

		if cs == klf200.ConnectionOpen {
			open(client)
		}
	})

	client.Start()

	for {
		time.Sleep(time.Second * 5)
	}

	client.Close()
}

/*
	func notify(n commands.Notify) {
		fmt.Printf("got notify %d\n", n.Code())
	}
*/
func open(client *klf200.Client) {
	ver, err := client.Device().GetVersion()
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

	state, err := client.Device().GetState()
	if err != nil {
		panic(err)
	}

	fmt.Printf("GatewayState = %d\n", state.GatewayState)
	fmt.Printf("SubState = %d\n", state.SubState)

	// client.SetUtc(time.Now())
	// TODO: set time zone

	devTime, err := client.Device().GetLocalTime()
	if err != nil {
		panic(err)
	}

	fmt.Printf("UtcTime = %v\n", devTime.UtcTime)
	fmt.Printf("Second = %d\n", devTime.LocalTime.Second)
	fmt.Printf("Minute = %d\n", devTime.LocalTime.Minute)
	fmt.Printf("Hour = %d\n", devTime.LocalTime.Hour)
	fmt.Printf("DayOfMonth = %d\n", devTime.LocalTime.DayOfMonth)
	fmt.Printf("Month = %d\n", devTime.LocalTime.Month)
	fmt.Printf("Year = %d\n", devTime.LocalTime.Year)
	fmt.Printf("WeekDay = %d\n", devTime.LocalTime.WeekDay)
	fmt.Printf("DayOfYear = %d\n", devTime.LocalTime.DayOfYear)
	fmt.Printf("DaylightSavingFlag = %d\n", devTime.LocalTime.DaylightSavingFlag)

	net, err := client.Device().GetNetworkSetup()
	if err != nil {
		panic(err)
	}

	fmt.Printf("DHCP = %t\n", net.DHCP)
	fmt.Printf("IpAddress = %v\n", net.IpAddress)
	fmt.Printf("Mask = %v\n", net.Mask)
	fmt.Printf("DefGW = %v\n", net.DefGW)

	objects, err := client.Config().GetSystemTable(context.Background())
	if err != nil {
		panic(err)
	}

	for _, object := range objects {
		fmt.Printf("Object = %v\n", object)
	}

	nodes, err := client.Info().GetAllNodesInformation(context.TODO())
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		fmt.Printf("Node = %v\n", node)
	}

	for _, node := range nodes {
		fmt.Printf("%s %d\n", node.Name, node.CurrentPosition)
	}
	/*
	   //sess, err := client.Commands().ChangePosition(context.Background(), 7, commands.NewMPValueRelative(-30))
	   //sess, err := client.Commands().Mode(context.Background(), 7)

	   	if err != nil {
	   		panic(err)
	   	}

	   fmt.Printf("Session = %v\n", sess)

	   	for event := range sess.Events() {
	   		switch event := event.(type) {
	   		case *klf200.RunError:
	   			fmt.Printf("run error %v\n", event)
	   		case *klf200.RunStatus:
	   			fmt.Printf("run status %v\n", event)
	   		case *klf200.RunRemainingTime:
	   			fmt.Printf("run remaining time %v\n", event)
	   		}
	   	}
	*/

	for {
		status, err := client.Commands().Status(context.Background(), []int{0, 1, 2, 3, 4, 5, 6, 7})

		if err != nil {
			panic(err)
		}

		for _, status := range status {
			fmt.Printf("Status for node %d = %v\n", status.NodeIndex, status)
		}

		time.Sleep(time.Second * 5)
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

type logger struct {
	err error
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *logger) Debug(msg string) {
	fmt.Print("DEBUG ")
	l.print(msg)
}

func (l *logger) Info(msg string) {
	fmt.Print("INFO  ")
	l.print(msg)
}

func (l *logger) Warn(msg string) {
	fmt.Print("WARN  ")
	l.print(msg)
}

func (l *logger) Error(msg string) {
	fmt.Print("ERROR ")
	l.print(msg)
}

func (l *logger) print(msg string) {
	fmt.Print(msg)
	if l.err != nil {
		fmt.Print(": ")
		fmt.Printf("%s", l.err)
	}
	fmt.Println("")
}

func (l *logger) WithError(err error) klf200.Logger {
	return &logger{err}
}
