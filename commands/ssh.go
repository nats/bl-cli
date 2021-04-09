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
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/pkg/ssh"
)

var (
	sshHostRE = regexp.MustCompile(`^((?P<m1>\w+)@)?(?P<m2>.*?)(:(?P<m3>\d+))?$`)
)

// SSH creates the ssh commands hierarchy
func SSH(parent *Command) *Command {
	usr, err := user.Current()
	checkErr(err)

	path := filepath.Join(usr.HomeDir, ".ssh", "id_rsa")

	sshDesc := fmt.Sprintf(`Access a Server using SSH by providing its ID or name.

You may specify the user to login with by passing the `+"`"+`--%s`+"`"+` flag. To access the Server on a non-default port, use the `+"`"+`--%s`+"`"+` flag. By default, the connection will be made to the Server's public IP address. In order access it using its private IP address, use the `+"`"+`--%s`+"`"+` flag.
`, blcli.ArgSSHUser, blcli.ArgsSSHPort, blcli.ArgsSSHPrivateIP)

	cmdSSH := CmdBuilder(parent, RunSSH, "ssh <server-id|name>", "Access a Server using SSH", sshDesc, Writer)
	AddStringFlag(cmdSSH, blcli.ArgSSHUser, "", "root", "SSH user for connection")
	AddStringFlag(cmdSSH, blcli.ArgsSSHKeyPath, "", path, "Path to SSH private key")
	AddIntFlag(cmdSSH, blcli.ArgsSSHPort, "", 22, "The remote port sshd is running on")
	AddBoolFlag(cmdSSH, blcli.ArgsSSHAgentForwarding, "", false, "Enable SSH agent forwarding")
	AddBoolFlag(cmdSSH, blcli.ArgsSSHPrivateIP, "", false, "SSH to Server's private IP address")
	AddStringFlag(cmdSSH, blcli.ArgSSHCommand, "", "", "Command to execute on Server")

	return cmdSSH
}

// RunSSH finds a server to ssh to given input parameters (name or id).
func RunSSH(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	serverID := c.Args[0]

	if serverID == "" {
		return blcli.NewMissingArgsErr(c.NS)
	}

	user, err := c.Doit.GetString(c.NS, blcli.ArgSSHUser)
	if err != nil {
		return err
	}

	keyPath, err := c.Doit.GetString(c.NS, blcli.ArgsSSHKeyPath)
	if err != nil {
		return err
	}

	port, err := c.Doit.GetInt(c.NS, blcli.ArgsSSHPort)
	if err != nil {
		return err
	}

	var opts = make(ssh.Options)
	opts[blcli.ArgsSSHAgentForwarding], err = c.Doit.GetBool(c.NS, blcli.ArgsSSHAgentForwarding)
	if err != nil {
		return err
	}

	opts[blcli.ArgSSHCommand], err = c.Doit.GetString(c.NS, blcli.ArgSSHCommand)
	if err != nil {
		return nil
	}

	privateIPChoice, err := c.Doit.GetBool(c.NS, blcli.ArgsSSHPrivateIP)
	if err != nil {
		return err
	}

	var server *bl.Server

	ss := c.Servers()
	if id, err := strconv.Atoi(serverID); err == nil {
		// serverID is an integer

		blServer, err := ss.Get(id)
		if err != nil {
			return err
		}

		server = blServer
	} else {
		// serverID is a string
		servers, err := ss.List()
		if err != nil {
			return err
		}

		shi := extractHostInfo(serverID)

		if shi.user != "" {
			user = shi.user
		}

		if i, err := strconv.Atoi(shi.port); shi.port != "" && err != nil {
			port = i
		}

		for _, s := range servers {
			if s.Name == shi.host {
				server = &s
				break
			}
			if strconv.Itoa(s.ID) == shi.host {
				server = &s
				break
			}
		}

		if server == nil {
			return errors.New("Could not find Server")
		}

	}

	if user == "" {
		user = defaultSSHUser(server)
	}

	ip, err := privateIPElsePub(server, privateIPChoice)
	if err != nil {
		return err
	}

	if ip == "" {
		return errors.New("Could not find Server address")
	}

	runner := c.Doit.SSH(user, ip, keyPath, port, opts)
	return runner.Run()
}

func defaultSSHUser(server *bl.Server) string {
	slug := strings.ToLower(server.Image.Slug)
	if strings.Contains(slug, "coreos") {
		return "core"
	}

	return "root"
}

type sshHostInfo struct {
	user string
	host string
	port string
}

func extractHostInfo(in string) sshHostInfo {
	m := sshHostRE.FindStringSubmatch(in)
	r := map[string]string{}
	for i, n := range sshHostRE.SubexpNames() {
		r[n] = m[i]
	}

	return sshHostInfo{
		user: r["m1"],
		host: r["m2"],
		port: r["m3"],
	}
}

func privateIPElsePub(server *bl.Server, choice bool) (string, error) {
	if choice {
		return server.PrivateIPv4()
	}
	return server.PublicIPv4()
}
