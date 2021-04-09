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

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

// Snapshot creates the snapshot command
func Snapshot() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "snapshot",
			Aliases: []string{"s"},
			Short:   "Access and manage snapshots",
			Long:    "The subcommands of `bl compute snapshot` allow you to manage and retrieve information about Server and block storage volume snapshots.",
		},
	}

	snapshotDetail := `

  - The snapshot's ID
  - The snapshot's name
  - The date and time when the snapshot was created
  - The slugs of the datacenter regions in which the snapshot is available
  - The type of resource the snapshot was made from, Server or volume, and its ID
  - The minimum size in GB required for a Server or volume to use this snapshot
  - The compressed, billable size of the snapshot
`

	cmdRunSnapshotList := CmdBuilder(cmd, RunSnapshotList, "list [glob]",
		"List Server and volume snapshots", "List information about Server and block storage volume snapshots, including:"+snapshotDetail,
		Writer, aliasOpt("ls"), displayerType(&displayers.Snapshot{}))
	AddStringFlag(cmdRunSnapshotList, blcli.ArgResourceType, "", "", "Filter by resource type (`server` or `volume`)")
	AddStringFlag(cmdRunSnapshotList, blcli.ArgRegionSlug, "", "", "Filter by regional availability")

	CmdBuilder(cmd, RunSnapshotGet, "get <snapshot-id>...",
		"Retrieve a Server or volume snapshot", "Retrieve information about a Server or block storage volume snapshot, including:"+snapshotDetail,
		Writer, aliasOpt("g"), displayerType(&displayers.Snapshot{}))

	cmdRunSnapshotDelete := CmdBuilder(cmd, RunSnapshotDelete, "delete <snapshot-id>...",
		"Delete a snapshot of a Server or volume", "Delete a snapshot of a Server or volume by specifying its ID.",
		Writer, aliasOpt("d"), displayerType(&displayers.Snapshot{}))
	AddBoolFlag(cmdRunSnapshotDelete, blcli.ArgForce, blcli.ArgShortForce, false, "Delete the snapshot without confirmation")

	return cmd
}

// RunSnapshotList returns a list of snapshots
func RunSnapshotList(c *CmdConfig) error {
	var err error
	ss := c.Snapshots()

	restype, err := c.Doit.GetString(c.NS, blcli.ArgResourceType)
	if err != nil {
		return err
	}

	region, err := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	if err != nil {
		return err
	}

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	var matchedList []bl.Snapshot
	var list []bl.Snapshot

	if restype == "server" {
		list, err = ss.ListServer()
		if err != nil {
			return err
		}
	} else if restype == "volume" {
		list, err = ss.ListVolume()
		if err != nil {
			return err
		}
	} else {
		list, err = ss.List()
		if err != nil {
			return err
		}
	}

	for _, snapshot := range list {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(snapshot.ID) {
					skip = false
				}
				if m.Match(snapshot.Name) {
					skip = false
				}
			}
		}

		if !skip && region != "" {
			for _, snapshotRegion := range snapshot.Regions {
				if region != snapshotRegion {
					skip = true
				} else {
					skip = false
					break
				}
			}

		}

		if !skip {
			matchedList = append(matchedList, snapshot)
		}
	}

	item := &displayers.Snapshot{Snapshots: matchedList}
	return c.Display(item)
}

// RunSnapshotGet returns a snapshot
func RunSnapshotGet(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	ss := c.Snapshots()
	ids := c.Args

	var matchedList []bl.Snapshot

	for _, id := range ids {
		s, err := ss.Get(id)
		if err != nil {
			return err
		}
		matchedList = append(matchedList, *s)
	}
	item := &displayers.Snapshot{Snapshots: matchedList}
	return c.Display(item)
}

// RunSnapshotDelete destroys snapshot(s) by id
func RunSnapshotDelete(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	force, err := c.Doit.GetBool(c.NS, blcli.ArgForce)
	if err != nil {
		return err
	}

	ss := c.Snapshots()
	ids := c.Args

	if force || AskForConfirmDelete("snapshot", len(ids)) == nil {
		for _, id := range ids {
			err := ss.Delete(id)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("Operation aborted.")
	}
	return nil
}
