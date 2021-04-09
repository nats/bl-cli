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

// FloatingIPAction creates the floating IP action command.
func FloatingIPAction() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "floating-ip-action",
			Short:   "Display commands to associate floating IP addresses with Servers",
			Long:    "Floating IP actions are commands that are used to manage BinaryLane floating IP addresses.",
			Aliases: []string{"fipa"},
		},
	}
	flipactionDetail := `

	- The unique numeric ID used to identify and reference a floating IP action.
	- The status of the floating IP action. This will be either "in-progress", "completed", or "errored".
	- A time value given in ISO8601 combined date and time format that represents when the action was initiated.
	- A time value given in ISO8601 combined date and time format that represents when the action was completed.
	- The resource ID, which is a unique identifier for the resource that the action is associated with.
	- The type of resource that the action is associated with.
	- The region where the action occurred.
	- The slug for the region where the action occurred.
`
	CmdBuilder(cmd, RunFloatingIPActionsGet,
		"get <floating-ip> <action-id>", "Retrieve the status of a floating IP action", `Use this command to retrieve the status of a floating IP action. Outputs the following information:`+flipactionDetail, Writer,
		displayerType(&displayers.Action{}))

	CmdBuilder(cmd, RunFloatingIPActionsAssign,
		"assign <floating-ip> <server-id>", "Assign a floating IP address to a Server", "Use this command to assign a floating IP address to a Server by specifying the `server_id`.", Writer,
		displayerType(&displayers.Action{}))

	CmdBuilder(cmd, RunFloatingIPActionsUnassign,
		"unassign <floating-ip>", "Unassign a floating IP address from a Server", `Use this command to unassign a floating IP address from a Server. The floating IP address will be reserved in the region but not assigned to a Server.`, Writer,
		displayerType(&displayers.Action{}))

	return cmd
}

// RunFloatingIPActionsGet retrieves an action for a floating IP.
func RunFloatingIPActionsGet(c *CmdConfig) error {
	if len(c.Args) != 2 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	ip := c.Args[0]

	fia := c.FloatingIPActions()

	actionID, err := strconv.Atoi(c.Args[1])
	if err != nil {
		return err
	}

	a, err := fia.Get(ip, actionID)
	if err != nil {
		return err
	}

	item := &displayers.Action{Actions: bl.Actions{*a}}
	return c.Display(item)
}

// RunFloatingIPActionsAssign assigns a floating IP to a server.
func RunFloatingIPActionsAssign(c *CmdConfig) error {
	if len(c.Args) != 2 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	ip := c.Args[0]

	fia := c.FloatingIPActions()

	serverID, err := strconv.Atoi(c.Args[1])
	if err != nil {
		return err
	}

	a, err := fia.Assign(ip, serverID)
	if err != nil {
		checkErr(fmt.Errorf("could not assign IP to server: %v", err))
	}

	item := &displayers.Action{Actions: bl.Actions{*a}}
	return c.Display(item)
}

// RunFloatingIPActionsUnassign unassigns a floating IP to a server.
func RunFloatingIPActionsUnassign(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}

	ip := c.Args[0]

	fia := c.FloatingIPActions()

	a, err := fia.Unassign(ip)
	if err != nil {
		checkErr(fmt.Errorf("could not unassign IP to server: %v", err))
	}

	item := &displayers.Action{Actions: bl.Actions{*a}}
	return c.Display(item)
}
