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
	"github.com/spf13/cobra"
)

// Version creates a version command.
func Version() *Command {
	return &Command{
		Command: &cobra.Command{
			Use:   "version",
			Short: "Show the current version",
			Long:  "The `bl version` command displays the version of the bl software.",
			Run: func(cmd *cobra.Command, args []string) {
				if blcli.Build != "" {
					blcli.DoitVersion.Build = blcli.Build
				}
				if blcli.Major != "" {
					i, _ := strconv.Atoi(blcli.Major)
					blcli.DoitVersion.Major = i
				}
				if blcli.Minor != "" {
					i, _ := strconv.Atoi(blcli.Minor)
					blcli.DoitVersion.Minor = i
				}
				if blcli.Patch != "" {
					i, _ := strconv.Atoi(blcli.Patch)
					blcli.DoitVersion.Patch = i
				}
				if blcli.Label != "" {
					blcli.DoitVersion.Label = blcli.Label
				}

				fmt.Println(blcli.DoitVersion.Complete(&blcli.GithubLatestVersioner{}))
			},
		},
	}
}
