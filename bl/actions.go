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

//Action is a wrapper for binarylane.Action
type Action struct {
	*binarylane.Action
}

// Actions is a slice of Action.
type Actions []Action

// ActionsService is an interface for interacting with BinaryLane's action api.
type ActionsService interface {
	List() (Actions, error)
	Get(int) (*Action, error)
}

type actionsService struct {
	client *binarylane.Client
}

var _ ActionsService = &actionsService{}

// NewActionsService builds an ActionsService instance.
func NewActionsService(client *binarylane.Client) ActionsService {
	return &actionsService{
		client: client,
	}
}

func (as *actionsService) List() (Actions, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := as.client.Actions.List(context.TODO(), opt)
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

func (as *actionsService) Get(id int) (*Action, error) {
	a, _, err := as.client.Actions.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}
