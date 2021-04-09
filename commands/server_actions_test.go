/*
Copyright 2018 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/stretchr/testify/assert"
)

func TestServerActionCommand(t *testing.T) {
	cmd := ServerAction()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "change-kernel", "enable-backups", "disable-backups", "enable-ipv6", "enable-private-networking", "get", "power-cycle", "power-off", "power-on", "password-reset", "reboot", "rebuild", "rename", "resize", "restore", "shutdown", "snapshot")
}

func TestServerActionsChangeKernel(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().ChangeKernel(1, 2).Return(&testAction, nil)

		config.Doit.Set(config.NS, blcli.ArgKernelID, 2)
		config.Args = append(config.Args, "1")

		err := RunServerActionChangeKernel(config)
		assert.NoError(t, err)
	})
}
func TestServerActionsEnableBackups(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().EnableBackups(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionEnableBackups(config)
		assert.NoError(t, err)
	})

}
func TestServerActionsDisableBackups(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().DisableBackups(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionDisableBackups(config)
		assert.NoError(t, err)
	})

}
func TestServerActionsEnableIPv6(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().EnableIPv6(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionEnableIPv6(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsEnablePrivateNetworking(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().EnablePrivateNetworking(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionEnablePrivateNetworking(config)
		assert.NoError(t, err)
	})
}
func TestServerActionsGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Get(1, 2).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgActionID, 2)

		err := RunServerActionGet(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsPasswordReset(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().PasswordReset(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionPasswordReset(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsPowerCycle(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().PowerCycle(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionPowerCycle(config)
		assert.NoError(t, err)
	})

}
func TestServerActionsPowerOff(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().PowerOff(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionPowerOff(config)
		assert.NoError(t, err)
	})
}
func TestServerActionsPowerOn(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().PowerOn(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionPowerOn(config)
		assert.NoError(t, err)
	})

}
func TestServerActionsReboot(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Reboot(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionReboot(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsRebuildByImageID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().RebuildByImageID(1, 2).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgImage, "2")

		err := RunServerActionRebuild(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsRebuildByImageSlug(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().RebuildByImageSlug(1, "slug").Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgImage, "slug")

		err := RunServerActionRebuild(config)
		assert.NoError(t, err)
	})

}
func TestServerActionsRename(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Rename(1, "name").Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgServerName, "name")

		err := RunServerActionRename(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsResize(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Resize(1, "1gb", true).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgSizeSlug, "1gb")
		config.Doit.Set(config.NS, blcli.ArgResizeDisk, true)

		err := RunServerActionResize(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsRestore(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Restore(1, 2).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgImageID, 2)

		err := RunServerActionRestore(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsShutdown(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Shutdown(1).Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActionShutdown(config)
		assert.NoError(t, err)
	})
}

func TestServerActionsSnapshot(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.serverActions.EXPECT().Snapshot(1, "name").Return(&testAction, nil)

		config.Args = append(config.Args, "1")

		config.Doit.Set(config.NS, blcli.ArgSnapshotName, "name")

		err := RunServerActionSnapshot(config)
		assert.NoError(t, err)
	})
}
