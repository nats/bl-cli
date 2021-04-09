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

package displayers

import (
	"io"

	"github.com/binarylane/bl-cli/bl"
)

type Tag struct {
	Tags bl.Tags
}

var _ Displayable = &Tag{}

func (t *Tag) JSON(out io.Writer) error {
	return writeJSON(t.Tags, out)
}

func (t *Tag) Cols() []string {
	return []string{"Name", "ServerCount"}
}

func (t *Tag) ColMap() map[string]string {
	return map[string]string{
		"Name":        "Name",
		"ServerCount": "Server Count",
	}
}

func (t *Tag) KV() []map[string]interface{} {
	out := []map[string]interface{}{}

	for _, x := range t.Tags {
		serverCount := x.Resources.Servers.Count
		o := map[string]interface{}{
			"Name":        x.Name,
			"ServerCount": serverCount,
		}
		out = append(out, o)
	}

	return out
}
