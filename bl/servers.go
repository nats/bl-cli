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

package bl

import (
	"context"

	"github.com/binarylane/go-binarylane"
	"github.com/binarylane/go-binarylane/util"
)

// ServerIPTable is a table of interface IPS.
type ServerIPTable map[InterfaceType]string

// InterfaceType is a an interface type.
type InterfaceType string

const (
	// InterfacePublic is a public interface.
	InterfacePublic InterfaceType = "public"
	// InterfacePrivate is a private interface.
	InterfacePrivate InterfaceType = "private"
)

// Server is a wrapper for binarylane.Server
type Server struct {
	*binarylane.Server
}

// Servers is a slice of Server.
type Servers []Server

// Kernel is a wrapper for binarylane.Kernel
type Kernel struct {
	*binarylane.Kernel
}

// Kernels is a slice of Kernel.
type Kernels []Kernel

// ServersService is an interface for interacting with BinaryLane's server api.
type ServersService interface {
	List() (Servers, error)
	ListByTag(string) (Servers, error)
	Get(int) (*Server, error)
	Create(*binarylane.ServerCreateRequest, bool) (*Server, error)
	CreateMultiple(*binarylane.ServerMultiCreateRequest) (Servers, error)
	Delete(int) error
	DeleteByTag(string) error
	Kernels(int) (Kernels, error)
	Snapshots(int) (Images, error)
	Backups(int) (Images, error)
	Actions(int) (Actions, error)
	Neighbors(int) (Servers, error)
}

type serversService struct {
	client *binarylane.Client
}

var _ ServersService = &serversService{}

// NewServersService builds a ServersService instance.
func NewServersService(client *binarylane.Client) ServersService {
	return &serversService{
		client: client,
	}
}

func (ss *serversService) List() (Servers, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ss.client.Servers.List(context.TODO(), opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Servers, len(si))
	for i := range si {
		a := si[i].(binarylane.Server)
		list[i] = Server{Server: &a}
	}

	return list, nil
}

func (ss *serversService) ListByTag(tagName string) (Servers, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ss.client.Servers.ListByTag(context.TODO(), tagName, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Servers, len(si))
	for i := range si {
		a := si[i].(binarylane.Server)
		list[i] = Server{Server: &a}
	}

	return list, nil
}

func (ss *serversService) Get(id int) (*Server, error) {
	d, _, err := ss.client.Servers.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	return &Server{Server: d}, nil
}

func (ss *serversService) Create(dcr *binarylane.ServerCreateRequest, wait bool) (*Server, error) {
	d, resp, err := ss.client.Servers.Create(context.TODO(), dcr)
	if err != nil {
		return nil, err
	}

	if wait {
		var action *binarylane.LinkAction
		for _, a := range resp.Links.Actions {
			if a.Rel == "create" {
				action = &a
				break
			}
		}

		if action != nil {
			_ = util.WaitForActive(context.TODO(), ss.client, action.HREF)
			server, err := ss.Get(d.ID)
			if err != nil {
				return nil, err
			}
			d = server.Server
		}
	}

	return &Server{Server: d}, nil
}

func (ss *serversService) CreateMultiple(dmcr *binarylane.ServerMultiCreateRequest) (Servers, error) {
	binarylaneServers, _, err := ss.client.Servers.CreateMultiple(context.TODO(), dmcr)
	if err != nil {
		return nil, err
	}

	var servers Servers
	for _, d := range binarylaneServers {
		servers = append(servers, Server{Server: &d})
	}

	return servers, nil
}

func (ss *serversService) Delete(id int) error {
	_, err := ss.client.Servers.Delete(context.TODO(), id)
	return err
}

func (ss *serversService) DeleteByTag(tag string) error {
	_, err := ss.client.Servers.DeleteByTag(context.TODO(), tag)
	return err
}

func (ss *serversService) Kernels(id int) (Kernels, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ss.client.Servers.Kernels(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Kernels, len(si))
	for i := range si {
		a := si[i].(binarylane.Kernel)
		list[i] = Kernel{Kernel: &a}
	}

	return list, nil
}

func (ss *serversService) Snapshots(id int) (Images, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ss.client.Servers.Snapshots(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Images, len(si))
	for i := range si {
		a := si[i].(binarylane.Image)
		list[i] = Image{Image: &a}
	}

	return list, nil
}

func (ss *serversService) Backups(id int) (Images, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ss.client.Servers.Backups(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Images, len(si))
	for i := range si {
		a := si[i].(binarylane.Image)
		list[i] = Image{Image: &a}
	}

	return list, nil
}

func (ss *serversService) Actions(id int) (Actions, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ss.client.Servers.Actions(context.TODO(), id, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(Actions, len(si))
	for i := range si {
		a := si[i].(binarylane.Action)
		list[i] = Action{Action: &a}
	}

	return list, nil
}

func (ss *serversService) Neighbors(id int) (Servers, error) {
	list, _, err := ss.client.Servers.Neighbors(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	servers := make(Servers, len(list))
	for i := range list {
		servers[i] = Server{&list[i]}
	}

	return servers, nil
}
