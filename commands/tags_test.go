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
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/go-binarylane"
	"github.com/stretchr/testify/assert"
)

var (
	testTag = bl.Tag{
		Tag: &binarylane.Tag{
			Name: "mytag",
			Resources: &binarylane.TaggedResources{
				Count:         5,
				LastTaggedURI: fmt.Sprintf("https://api.binarylane.com.au/v2/servers/%d", testServer.ID),
				Servers: &binarylane.TaggedServersResources{
					Count:      5,
					LastTagged: testServer.Server,
				},
				Images: &binarylane.TaggedImagesResources{
					Count: 0,
				},
			}}}
	testTagList = bl.Tags{
		testTag,
	}
)

func TestTTagCommand(t *testing.T) {
	cmd := Tags()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "create", "get", "delete", "list")
}

func TestTagGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tags.EXPECT().Get("mytag").Return(&testTag, nil)

		config.Args = append(config.Args, "mytag")

		err := RunCmdTagGet(config)
		assert.NoError(t, err)
	})
}

func TestTagList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tags.EXPECT().List().Return(testTagList, nil)

		err := RunCmdTagList(config)
		assert.NoError(t, err)
	})
}

func TestTagCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tcr := binarylane.TagCreateRequest{Name: "new-tag"}
		tm.tags.EXPECT().Create(&tcr).Return(&testTag, nil)
		config.Args = append(config.Args, "new-tag")

		err := RunCmdTagCreate(config)
		assert.NoError(t, err)
	})
}

func TestTagDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tags.EXPECT().Delete("my-tag").Return(nil)
		config.Args = append(config.Args, "my-tag")

		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunCmdTagDelete(config)
		assert.NoError(t, err)
	})
}

func TestTagDeleteMultiple(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.tags.EXPECT().Delete("my-tag").Return(nil)
		tm.tags.EXPECT().Delete("my-tag-secondary").Return(nil)
		config.Args = append(config.Args, "my-tag", "my-tag-secondary")

		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunCmdTagDelete(config)
		assert.NoError(t, err)
	})
}
