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
)

// FloatingIPActionsService is an interface for interacting with
// BinaryLane's floating ip action api.
type FloatingIPActionsService interface {
	Assign(ip string, serverID int) (*Action, error)
	Unassign(ip string) (*Action, error)
	Get(ip string, actionID int) (*Action, error)
	List(ip string, opt *binarylane.ListOptions) ([]Action, error)
}

type floatingIPActionsService struct {
	client *binarylane.Client
}

var _ FloatingIPActionsService = &floatingIPActionsService{}

// NewFloatingIPActionsService builds a FloatingIPActionsService instance.
func NewFloatingIPActionsService(client *binarylane.Client) FloatingIPActionsService {
	return &floatingIPActionsService{
		client: client,
	}
}

func (fia *floatingIPActionsService) Assign(ip string, serverID int) (*Action, error) {
	a, _, err := fia.client.FloatingIPActions.Assign(context.TODO(), ip, serverID)
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}

func (fia *floatingIPActionsService) Unassign(ip string) (*Action, error) {
	a, _, err := fia.client.FloatingIPActions.Unassign(context.TODO(), ip)
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}

func (fia *floatingIPActionsService) Get(ip string, actionID int) (*Action, error) {
	a, _, err := fia.client.FloatingIPActions.Get(context.TODO(), ip, actionID)
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}

func (fia *floatingIPActionsService) List(ip string, opt *binarylane.ListOptions) ([]Action, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := fia.client.FloatingIPActions.List(context.TODO(), ip, opt)
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
