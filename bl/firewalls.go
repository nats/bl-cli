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
	"errors"

	"github.com/binarylane/go-binarylane"
)

// Firewall wraps a binarylane Firewall.
type Firewall struct {
	*binarylane.Firewall
}

// Firewalls is a slice of Firewall.
type Firewalls []Firewall

// FirewallsService is the binarylane FirewallsService interface.
type FirewallsService interface {
	Get(fID string) (*Firewall, error)
	Create(fr *binarylane.FirewallRequest) (*Firewall, error)
	Update(fID string, fr *binarylane.FirewallRequest) (*Firewall, error)
	List() (Firewalls, error)
	ListByServer(sID int) (Firewalls, error)
	Delete(fID string) error
	AddServers(fID string, sIDs ...int) error
	RemoveServers(fID string, sIDs ...int) error
	AddTags(fID string, tags ...string) error
	RemoveTags(fID string, tags ...string) error
	AddRules(fID string, rr *binarylane.FirewallRulesRequest) error
	RemoveRules(fID string, rr *binarylane.FirewallRulesRequest) error
}

var _ FirewallsService = &firewallsService{}

type firewallsService struct {
	client *binarylane.Client
}

// NewFirewallsService builds an instance of FirewallsService.
func NewFirewallsService(client *binarylane.Client) FirewallsService {
	return &firewallsService{client: client}
}

func (fs *firewallsService) Get(fID string) (*Firewall, error) {
	f, _, err := fs.client.Firewalls.Get(context.TODO(), fID)
	if err != nil {
		return nil, err
	}

	return &Firewall{Firewall: f}, nil
}

func (fs *firewallsService) Create(fr *binarylane.FirewallRequest) (*Firewall, error) {
	f, _, err := fs.client.Firewalls.Create(context.TODO(), fr)
	if err != nil {
		return nil, err
	}

	return &Firewall{Firewall: f}, nil
}

func (fs *firewallsService) Update(fID string, fr *binarylane.FirewallRequest) (*Firewall, error) {
	f, _, err := fs.client.Firewalls.Update(context.TODO(), fID, fr)
	if err != nil {
		return nil, err
	}

	return &Firewall{Firewall: f}, nil
}

func (fs *firewallsService) List() (Firewalls, error) {
	listFn := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := fs.client.Firewalls.List(context.TODO(), opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	return paginatedListHelper(listFn)
}

func (fs *firewallsService) ListByServer(sID int) (Firewalls, error) {
	listFn := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := fs.client.Firewalls.ListByServer(context.TODO(), sID, opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	return paginatedListHelper(listFn)
}

func (fs *firewallsService) Delete(fID string) error {
	_, err := fs.client.Firewalls.Delete(context.TODO(), fID)
	return err
}

func (fs *firewallsService) AddServers(fID string, sIDs ...int) error {
	_, err := fs.client.Firewalls.AddServers(context.TODO(), fID, sIDs...)
	return err
}

func (fs *firewallsService) RemoveServers(fID string, sIDs ...int) error {
	_, err := fs.client.Firewalls.RemoveServers(context.TODO(), fID, sIDs...)
	return err
}

func (fs *firewallsService) AddTags(fID string, tags ...string) error {
	_, err := fs.client.Firewalls.AddTags(context.TODO(), fID, tags...)
	return err
}

func (fs *firewallsService) RemoveTags(fID string, tags ...string) error {
	_, err := fs.client.Firewalls.RemoveTags(context.TODO(), fID, tags...)
	return err
}

func (fs *firewallsService) AddRules(fID string, rr *binarylane.FirewallRulesRequest) error {
	_, err := fs.client.Firewalls.AddRules(context.TODO(), fID, rr)
	return err
}

func (fs *firewallsService) RemoveRules(fID string, rr *binarylane.FirewallRulesRequest) error {
	_, err := fs.client.Firewalls.RemoveRules(context.TODO(), fID, rr)
	return err
}

func paginatedListHelper(listFn func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error)) (Firewalls, error) {
	si, err := PaginateResp(listFn)
	if err != nil {
		return nil, err
	}

	list := make([]Firewall, len(si))
	for i := range si {
		a, ok := si[i].(binarylane.Firewall)
		if !ok {
			return nil, errors.New("unexpected value in response")
		}

		list[i] = Firewall{Firewall: &a}
	}

	return list, nil
}
