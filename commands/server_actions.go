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
	"fmt"
	"strconv"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/spf13/cobra"
)

type actionFn func(das bl.ServerActionsService) (*bl.Action, error)

func performAction(c *CmdConfig, fn actionFn) error {
	das := c.ServerActions()

	a, err := fn(das)
	if err != nil {
		return err
	}

	wait, err := c.Doit.GetBool(c.NS, blcli.ArgCommandWait)
	if err != nil {
		return err
	}

	if wait {
		a, err = actionWait(c, a.ID, 5)
		if err != nil {
			return err
		}

	}

	item := &displayers.Action{Actions: bl.Actions{*a}}
	return c.Display(item)
}

// ServerAction creates the server-action command.
func ServerAction() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "server-action",
			Aliases: []string{"da"},
			Short:   "Display server action commands",
			Long: `Use the subcommands of ` + "`" + `bl compute server-action` + "`" + ` to perform actions on servers.

Servers actions are tasks that can be executed on a server, such as rebooting, resizing, or snapshotting a server.`,
		},
	}

	cmdServerActionGet := CmdBuilder(cmd, RunServerActionGet, "get <server-id>", "Retrieve a specific Server action", `Use this command to retrieve a Server action.`, Writer,
		aliasOpt("g"), displayerType(&displayers.Action{}))
	AddIntFlag(cmdServerActionGet, blcli.ArgActionID, "", 0, "Action ID", requiredOpt())

	cmdServerActionEnableBackups := CmdBuilder(cmd, RunServerActionEnableBackups,
		"enable-backups <server-id>", "Enable backups on a Server", `Use this command to enable backups on a Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionEnableBackups, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionDisableBackups := CmdBuilder(cmd, RunServerActionDisableBackups,
		"disable-backups <server-id>", "Disable backups on a Server", `Use this command to disable backups on a Server. This does not delete existing backups.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionDisableBackups, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionReboot := CmdBuilder(cmd, RunServerActionReboot,
		"reboot <server-id>", "Reboot a Server", `Use this command to reboot a Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionReboot, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionPowerCycle := CmdBuilder(cmd, RunServerActionPowerCycle,
		"power-cycle <server-id>", "Powercycle a Server", `Use this command to powercycle a Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionPowerCycle, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionShutdown := CmdBuilder(cmd, RunServerActionShutdown,
		"shutdown <server-id>", "Shut down a Server", `Use this command to shut down a Server. Servers that are powered off are still billable. To stop billing, destroy the Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionShutdown, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionPowerOff := CmdBuilder(cmd, RunServerActionPowerOff,
		"power-off <server-id>", "Power off a Server", `Use this command to power off a Server. Servers that are powered off are still billable. To stop billing, destroy the Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionPowerOff, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionPowerOn := CmdBuilder(cmd, RunServerActionPowerOn,
		"power-on <server-id>", "Power on a Server", `Use this command to power on a Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionPowerOn, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionPasswordReset := CmdBuilder(cmd, RunServerActionPasswordReset,
		"password-reset <server-id>", "Reset the root password for a Server", `Use this command to initiate a root password reset on a Server. This also powercycles the Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionPasswordReset, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionEnableIPv6 := CmdBuilder(cmd, RunServerActionEnableIPv6,
		"enable-ipv6 <server-id>", "Enable IPv6 on a Server", `Use this command to enable IPv6 networking on a Server. BinaryLane will automatically assign an IPv6 address to the Server.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionEnableIPv6, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionEnablePrivateNetworking := CmdBuilder(cmd, RunServerActionEnablePrivateNetworking,
		"enable-private-networking <server-id>", "Enable private networking on a Server", `Use this command to enable private networking on a Server. This adds a private IPv4 address to the Server that other Servers inside the network can access. The Server will require additional internal network configuration for it to become accessible through the private network.`, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionEnablePrivateNetworking, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionRestore := CmdBuilder(cmd, RunServerActionRestore,
		"restore <server-id>", "Restore a Server from a backup", `Use this command to restore a Server from a backup.`, Writer,
		displayerType(&displayers.Action{}))
	AddIntFlag(cmdServerActionRestore, blcli.ArgImageID, "", 0, "Image ID", requiredOpt())
	AddBoolFlag(cmdServerActionRestore, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	serverResizeDesc := `Use this command to resize a Server to a different plan.

By default, this command will only increase or decrease the CPU and RAM of the Server, not its disk size. This can be reversed.

To also increase the Server's disk size, pass the ` + "`--resize-disk`" + ` flag. This is a permanent change and cannot be reversed as a Server's disk size cannot be decreased.

In order to resize a Server, it must first be powered off.`
	cmdServerActionResize := CmdBuilder(cmd, RunServerActionResize,
		"resize <server-id>", "Resize a Server", serverResizeDesc, Writer,
		displayerType(&displayers.Action{}))
	AddBoolFlag(cmdServerActionResize, blcli.ArgResizeDisk, "", false, "Resize the Server's disk size in addition to its RAM and CPU.")
	AddStringFlag(cmdServerActionResize, blcli.ArgSizeSlug, "", "", "A slug indicating the new size for the Server (e.g. `s-2vcpu-2gb`). Run `bl compute size list` for a list of valid sizes.", requiredOpt())
	AddBoolFlag(cmdServerActionResize, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionRebuild := CmdBuilder(cmd, RunServerActionRebuild,
		"rebuild <server-id>", "Rebuild a Server", `Use this command to rebuild a Server from an image.`, Writer,
		displayerType(&displayers.Action{}))
	AddStringFlag(cmdServerActionRebuild, blcli.ArgImage, "", "", "Image ID or Slug", requiredOpt())
	AddBoolFlag(cmdServerActionRebuild, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionRename := CmdBuilder(cmd, RunServerActionRename,
		"rename <server-id>", "Rename a Server", `Use this command to rename a Server. When using a fully qualified domain name (FQDN) this also updates the pointer (PTR) record.`, Writer,
		displayerType(&displayers.Action{}))
	AddStringFlag(cmdServerActionRename, blcli.ArgServerName, "", "", "Server name", requiredOpt())
	AddBoolFlag(cmdServerActionRename, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionChangeKernel := CmdBuilder(cmd, RunServerActionChangeKernel,
		"change-kernel <server-id>", "Change a Server's kernel", `Use this command to change a Server's kernel.`, Writer,
		displayerType(&displayers.Action{}))
	AddIntFlag(cmdServerActionChangeKernel, blcli.ArgKernelID, "", 0, "Kernel ID", requiredOpt())
	AddBoolFlag(cmdServerActionChangeKernel, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	cmdServerActionSnapshot := CmdBuilder(cmd, RunServerActionSnapshot,
		"snapshot <server-id>", "Take a Server snapshot", `Use this command to take a snapshot of a Server. We recommend that you power off the Server before taking a snapshot to ensure data consistency.`, Writer,
		displayerType(&displayers.Action{}))
	AddStringFlag(cmdServerActionSnapshot, blcli.ArgSnapshotName, "", "", "Snapshot name", requiredOpt())
	AddBoolFlag(cmdServerActionSnapshot, blcli.ArgCommandWait, "", false, "Wait for action to complete")

	return cmd
}

// RunServerActionGet returns a server action by id.
func RunServerActionGet(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		serverID, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		actionID, err := c.Doit.GetInt(c.NS, blcli.ArgActionID)
		if err != nil {
			return nil, err
		}

		a, err := das.Get(serverID, actionID)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionEnableBackups disables backups for a server.
func RunServerActionEnableBackups(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		a, err := das.EnableBackups(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionDisableBackups disables backups for a server.
func RunServerActionDisableBackups(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		a, err := das.DisableBackups(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionReboot reboots a server.
func RunServerActionReboot(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, err
		}

		a, err := das.Reboot(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionPowerCycle power cycles a server.
func RunServerActionPowerCycle(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PowerCycle(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionShutdown shuts a server down.
func RunServerActionShutdown(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])
		if err != nil {
			return nil, fmt.Errorf("Could not convert args into integer")
		}

		a, err := das.Shutdown(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionPowerOff turns server power off.
func RunServerActionPowerOff(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PowerOff(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionPowerOn turns server power on.
func RunServerActionPowerOn(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PowerOn(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionPasswordReset resets the server root password.
func RunServerActionPasswordReset(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.PasswordReset(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionEnableIPv6 enables IPv6 for a server.
func RunServerActionEnableIPv6(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.EnableIPv6(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionEnablePrivateNetworking enables private networking for a server.
func RunServerActionEnablePrivateNetworking(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		a, err := das.EnablePrivateNetworking(id)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionRestore restores a server using an image id.
func RunServerActionRestore(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		image, err := c.Doit.GetInt(c.NS, blcli.ArgImageID)
		if err != nil {
			return nil, err
		}

		a, err := das.Restore(id, image)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionResize resizes a server giving a size slug and
// optionally expands the disk.
func RunServerActionResize(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		size, err := c.Doit.GetString(c.NS, blcli.ArgSizeSlug)
		if err != nil {
			return nil, err
		}

		disk, err := c.Doit.GetBool(c.NS, blcli.ArgResizeDisk)
		if err != nil {
			return nil, err
		}

		a, err := das.Resize(id, size, disk)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionRebuild rebuilds a server using an image id or slug.
func RunServerActionRebuild(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		image, err := c.Doit.GetString(c.NS, blcli.ArgImage)
		if err != nil {
			return nil, err
		}

		var a *bl.Action
		if i, aerr := strconv.Atoi(image); aerr == nil {
			a, err = das.RebuildByImageID(id, i)
		} else {
			a, err = das.RebuildByImageSlug(id, image)
		}
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionRename renames a server.
func RunServerActionRename(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		name, err := c.Doit.GetString(c.NS, blcli.ArgServerName)
		if err != nil {
			return nil, err
		}

		a, err := das.Rename(id, name)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionChangeKernel changes the kernel for a server.
func RunServerActionChangeKernel(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		kernel, err := c.Doit.GetInt(c.NS, blcli.ArgKernelID)
		if err != nil {
			return nil, err
		}

		a, err := das.ChangeKernel(id, kernel)
		return a, err
	}

	return performAction(c, fn)
}

// RunServerActionSnapshot creates a snapshot for a server.
func RunServerActionSnapshot(c *CmdConfig) error {
	fn := func(das bl.ServerActionsService) (*bl.Action, error) {
		err := ensureOneArg(c)
		if err != nil {
			return nil, err
		}
		id, err := strconv.Atoi(c.Args[0])

		if err != nil {
			return nil, err
		}

		name, err := c.Doit.GetString(c.NS, blcli.ArgSnapshotName)
		if err != nil {
			return nil, err
		}

		a, err := das.Snapshot(id, name)
		return a, err
	}

	return performAction(c, fn)
}
