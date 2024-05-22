package commands

import (
	"bytes"
	"fmt"
	"klf200/pkg/binary"
	"klf200/pkg/transport"
)

type GetVersionReq struct {
}

var _ Request = (*GetVersionReq)(nil)

func (req *GetVersionReq) Code() transport.Command {
	return GW_GET_VERSION_REQ
}

func (req *GetVersionReq) NewConfirm() Confirm {
	return &GetVersionCfm{}
}

func (req *GetVersionReq) Write() ([]byte, error) {
	return emptyData, nil
}

type GetVersionCfm struct {
	SoftwareVersion SoftWareVersionData
	HardwareVersion int
	ProductGroup    int
	ProductType     int
}

type SoftWareVersionData struct {
	CommandVersionNumber int
	VersionWholeNumber   int
	VersionSubNumber     int
	BranchID             int
	BuildNumber          int
	MicroBuild           int
}

var _ Confirm = (*GetVersionCfm)(nil)

func (cfm *GetVersionCfm) Code() transport.Command {
	return GW_GET_VERSION_CFM
}

func (cfm *GetVersionCfm) Read(data []byte) error {
	if len(data) != 9 {
		return fmt.Errorf("bad length")
	}

	reader := binary.MakeBinaryReader(bytes.NewBuffer(data))

	cfm.SoftwareVersion.CommandVersionNumber = cfm.readAsInt(reader)
	cfm.SoftwareVersion.VersionWholeNumber = cfm.readAsInt(reader)
	cfm.SoftwareVersion.VersionSubNumber = cfm.readAsInt(reader)
	cfm.SoftwareVersion.BranchID = cfm.readAsInt(reader)
	cfm.SoftwareVersion.BuildNumber = cfm.readAsInt(reader)
	cfm.SoftwareVersion.MicroBuild = cfm.readAsInt(reader)
	cfm.HardwareVersion = cfm.readAsInt(reader)
	cfm.ProductGroup = cfm.readAsInt(reader)
	cfm.ProductType = cfm.readAsInt(reader)

	return nil
}

func (cfm *GetVersionCfm) readAsInt(reader binary.BinaryReader) int {
	value, _ := reader.ReadU8()
	return int(value)
}
