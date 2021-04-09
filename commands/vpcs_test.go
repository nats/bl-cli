package commands

import (
	"strconv"
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/go-binarylane"

	"github.com/stretchr/testify/assert"
)

var (
	testVPC = bl.VPC{
		VPC: &binarylane.VPC{
			Name:        "vpc-name",
			RegionSlug:  "nyc1",
			Description: "vpc description",
			IPRange:     "10.116.0.0/20",
		}}

	testVPCList = bl.VPCs{
		testVPC,
	}
)

func TestVPCsCommand(t *testing.T) {
	cmd := VPCs()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "get", "list", "create", "update", "delete")
}

func TestVPCGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		vpcID := 2
		tm.vpcs.EXPECT().Get(vpcID).Return(&testVPC, nil)

		config.Args = append(config.Args, strconv.Itoa(vpcID))

		err := RunVPCGet(config)
		assert.NoError(t, err)
	})
}

func TestVPCGetNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunVPCGet(config)
		assert.Error(t, err)
	})
}

func TestVPCList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.vpcs.EXPECT().List().Return(testVPCList, nil)

		err := RunVPCList(config)
		assert.NoError(t, err)
	})
}

func TestVPCCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		r := binarylane.VPCCreateRequest{
			Name:        "vpc-name",
			RegionSlug:  "nyc1",
			Description: "vpc description",
			IPRange:     "10.116.0.0/20",
		}
		tm.vpcs.EXPECT().Create(&r).Return(&testVPC, nil)

		config.Doit.Set(config.NS, blcli.ArgVPCName, "vpc-name")
		config.Doit.Set(config.NS, blcli.ArgRegionSlug, "nyc1")
		config.Doit.Set(config.NS, blcli.ArgVPCDescription, "vpc description")
		config.Doit.Set(config.NS, blcli.ArgVPCIPRange, "10.116.0.0/20")

		err := RunVPCCreate(config)
		assert.NoError(t, err)
	})
}

func TestVPCUpdate(t *testing.T) {
	tests := []struct {
		desc            string
		setup           func(*CmdConfig)
		expectedVPCId   int
		expectedRequest *binarylane.VPCUpdateRequest
	}{
		{
			desc: "update vpc name",
			setup: func(in *CmdConfig) {
				in.Args = append(in.Args, "2")
				in.Doit.Set(in.NS, blcli.ArgVPCName, "update-vpc-name-test")

			},
			expectedVPCId: 2,
			expectedRequest: &binarylane.VPCUpdateRequest{
				Name: "update-vpc-name-test",
			},
		},

		{
			desc: "update vpc name and description",
			setup: func(in *CmdConfig) {
				in.Args = append(in.Args, "2")
				in.Doit.Set(in.NS, blcli.ArgVPCName, "update-vpc-name-test")
				in.Doit.Set(in.NS, blcli.ArgVPCDescription, "i am a new desc")

			},
			expectedVPCId: 2,
			expectedRequest: &binarylane.VPCUpdateRequest{
				Name:        "update-vpc-name-test",
				Description: "i am a new desc",
			},
		},

		{
			desc: "update vpc name and description and set to default",
			setup: func(in *CmdConfig) {
				in.Args = append(in.Args, "2")
				in.Doit.Set(in.NS, blcli.ArgVPCName, "update-vpc-name-test")
				in.Doit.Set(in.NS, blcli.ArgVPCDescription, "i am a new desc")
				in.Doit.Set(in.NS, blcli.ArgVPCDefault, true)
			},
			expectedVPCId: 2,
			expectedRequest: &binarylane.VPCUpdateRequest{
				Name:        "update-vpc-name-test",
				Description: "i am a new desc",
				Default:     boolPtr(true),
			},
		},
	}

	for _, tt := range tests {
		withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
			if tt.setup != nil {
				tt.setup(config)
			}

			tm.vpcs.EXPECT().Update(tt.expectedVPCId, tt.expectedRequest).Return(&testVPC, nil)
			err := RunVPCUpdate(config)

			assert.NoError(t, err)
		})
	}
}

func TestVPCUpdatefNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunVPCUpdate(config)
		assert.Error(t, err)
	})
}

func TestVPCDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		vpcID := 2
		tm.vpcs.EXPECT().Delete(vpcID).Return(nil)

		config.Args = append(config.Args, strconv.Itoa(vpcID))
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunVPCDelete(config)
		assert.NoError(t, err)
	})
}

func TestVPCDeleteNoID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunVPCDelete(config)
		assert.Error(t, err)
	})
}
