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
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/binarylane/bl-cli/bl"
)

type Server struct {
	Servers bl.Servers
}

var _ Displayable = &Server{}

func (d *Server) JSON(out io.Writer) error {
	return writeJSON(d.Servers, out)
}

func (d *Server) Cols() []string {
	cols := []string{
		"ID", "Name", "PublicIPv4", "PrivateIPv4", "PublicIPv6", "Memory", "VCPUs", "Disk", "Region", "Image", "VPCID", "Status", "Tags", "Features", "Volumes",
	}
	return cols
}

func (d *Server) ColMap() map[string]string {
	return map[string]string{
		"ID": "ID", "Name": "Name", "PublicIPv4": "Public IPv4", "PrivateIPv4": "Private IPv4", "PublicIPv6": "Public IPv6",
		"Memory": "Memory", "VCPUs": "VCPUs", "Disk": "Disk",
		"Region": "Region", "Image": "Image", "VPCID": "VPC ID", "Status": "Status",
		"Tags": "Tags", "Features": "Features", "Volumes": "Volumes",
		"SizeSlug": "Size Slug",
	}
}

func (d *Server) KV() []map[string]interface{} {
	out := []map[string]interface{}{}
	for _, d := range d.Servers {
		sort.Strings(d.Tags)
		tags := strings.Join(d.Tags, ",")
		image := fmt.Sprintf("%s %s", d.Image.Distribution, d.Image.Name)
		ip, _ := d.PublicIPv4()
		privIP, _ := d.PrivateIPv4()
		ip6, _ := d.PublicIPv6()
		features := strings.Join(d.Features, ",")
		volumes := strings.Join(d.VolumeIDs, ",")
		vpcID := ""
		if d.VPCID != 0 {
			vpcID = strconv.Itoa(d.VPCID)
		}
		m := map[string]interface{}{
			"ID": d.ID, "Name": d.Name, "PublicIPv4": ip, "PrivateIPv4": privIP, "PublicIPv6": ip6,
			"Memory": d.Memory, "VCPUs": d.Vcpus, "Disk": d.Disk,
			"Region": d.Region.Slug, "Image": image, "VPCID": vpcID, "Status": d.Status,
			"Tags": tags, "Features": features, "Volumes": volumes,
			"SizeSlug": d.SizeSlug,
		}
		out = append(out, m)
	}

	return out
}
