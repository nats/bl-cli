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
	"strconv"
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/pkg/runner"
	"github.com/binarylane/bl-cli/pkg/ssh"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSSHComand(t *testing.T) {
	parent := &Command{
		Command: &cobra.Command{
			Use:   "compute",
			Short: "compute commands",
			Long:  "compute commands are for controlling and managing infrastructure",
		},
	}
	cmd := SSH(parent)
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd)
}

func TestSSH_ID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Get(testServer.ID).Return(&testServer, nil)

		config.Args = append(config.Args, strconv.Itoa(testServer.ID))

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_InvalidID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		err := RunSSH(config)
		assert.Error(t, err)
	})
}

func TestSSH_UnknownServer(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Args = append(config.Args, "missing")

		err := RunSSH(config)
		assert.EqualError(t, err, "Could not find Server")
	})
}

func TestSSH_ServerWithNoPublic(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().List().Return(testPrivateServerList, nil)

		config.Args = append(config.Args, testPrivateServer.Name)

		err := RunSSH(config)
		assert.EqualError(t, err, "Could not find Server address")
	})
}

func TestSSH_CustomPort(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.sshRunner.EXPECT().Run().Return(nil)

		tc := config.Doit.(*blcli.TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, 2222, port)
			return tm.sshRunner
		}

		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Doit.Set(config.NS, blcli.ArgsSSHPort, "2222")
		config.Args = append(config.Args, testServer.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_CustomUser(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.sshRunner.EXPECT().Run().Return(nil)

		tc := config.Doit.(*blcli.TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, "foobar", user)
			return tm.sshRunner
		}

		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Doit.Set(config.NS, blcli.ArgSSHUser, "foobar")
		config.Args = append(config.Args, testServer.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_AgentForwarding(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.sshRunner.EXPECT().Run().Return(nil)

		tc := config.Doit.(*blcli.TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, true, opts[blcli.ArgsSSHAgentForwarding])
			return tm.sshRunner
		}

		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Doit.Set(config.NS, blcli.ArgsSSHAgentForwarding, true)
		config.Args = append(config.Args, testServer.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func TestSSH_CommandExecuting(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.sshRunner.EXPECT().Run().Return(nil)

		tc := config.Doit.(*blcli.TestConfig)
		tc.SSHFn = func(user, host, keyPath string, port int, opts ssh.Options) runner.Runner {
			assert.Equal(t, "uptime", opts[blcli.ArgSSHCommand])
			return tm.sshRunner
		}

		tm.servers.EXPECT().List().Return(testServerList, nil)
		config.Doit.Set(config.NS, blcli.ArgSSHCommand, "uptime")
		config.Args = append(config.Args, testServer.Name)

		err := RunSSH(config)
		assert.NoError(t, err)
	})
}

func Test_extractHostInfo(t *testing.T) {
	cases := []struct {
		s string
		e sshHostInfo
	}{
		{s: "host", e: sshHostInfo{host: "host"}},
		{s: "root@host", e: sshHostInfo{user: "root", host: "host"}},
		{s: "root@host:22", e: sshHostInfo{user: "root", host: "host", port: "22"}},
		{s: "host:22", e: sshHostInfo{host: "host", port: "22"}},
		{s: "dokku@simple-task-02efb9c544", e: sshHostInfo{host: "simple-task-02efb9c544", user: "dokku"}},
	}

	for _, c := range cases {
		i := extractHostInfo(c.s)
		assert.Equal(t, c.e, i)
	}
}
