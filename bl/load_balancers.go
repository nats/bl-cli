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

// LoadBalancer wraps a binarylane LoadBalancer.
type LoadBalancer struct {
	*binarylane.LoadBalancer
}

// LoadBalancers is a slice of LoadBalancer.
type LoadBalancers []LoadBalancer

// LoadBalancersService is the binarylane LoadBalancersService interface.
type LoadBalancersService interface {
	Get(lbID int) (*LoadBalancer, error)
	List() (LoadBalancers, error)
	Create(lbr *binarylane.LoadBalancerRequest) (*LoadBalancer, error)
	Update(lbID int, lbr *binarylane.LoadBalancerRequest) (*LoadBalancer, error)
	Delete(lbID int) error
	AddServers(lbID int, sIDs ...int) error
	RemoveServers(lbID int, sIDs ...int) error
	AddForwardingRules(lbID int, rules ...binarylane.ForwardingRule) error
	RemoveForwardingRules(lbID int, rules ...binarylane.ForwardingRule) error
}

var _ LoadBalancersService = &loadBalancersService{}

type loadBalancersService struct {
	client *binarylane.Client
}

// NewLoadBalancersService builds an instance of LoadBalancersService.
func NewLoadBalancersService(client *binarylane.Client) LoadBalancersService {
	return &loadBalancersService{
		client: client,
	}
}

func (lbs *loadBalancersService) Get(lbID int) (*LoadBalancer, error) {
	lb, _, err := lbs.client.LoadBalancers.Get(context.TODO(), lbID)
	if err != nil {
		return nil, err
	}

	return &LoadBalancer{LoadBalancer: lb}, nil
}

func (lbs *loadBalancersService) List() (LoadBalancers, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := lbs.client.LoadBalancers.List(context.TODO(), opt)
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

	list := make([]LoadBalancer, len(si))
	for i := range si {
		a := si[i].(binarylane.LoadBalancer)
		list[i] = LoadBalancer{LoadBalancer: &a}
	}

	return list, nil
}

func (lbs *loadBalancersService) Create(lbr *binarylane.LoadBalancerRequest) (*LoadBalancer, error) {
	lb, _, err := lbs.client.LoadBalancers.Create(context.TODO(), lbr)
	if err != nil {
		return nil, err
	}

	return &LoadBalancer{LoadBalancer: lb}, nil
}

func (lbs *loadBalancersService) Update(lbID int, lbr *binarylane.LoadBalancerRequest) (*LoadBalancer, error) {
	lb, _, err := lbs.client.LoadBalancers.Update(context.TODO(), lbID, lbr)
	if err != nil {
		return nil, err
	}

	return &LoadBalancer{LoadBalancer: lb}, nil
}

func (lbs *loadBalancersService) Delete(lbID int) error {
	_, err := lbs.client.LoadBalancers.Delete(context.TODO(), lbID)
	return err
}

func (lbs *loadBalancersService) AddServers(lbID int, sIDs ...int) error {
	_, err := lbs.client.LoadBalancers.AddServers(context.TODO(), lbID, sIDs...)
	return err
}

func (lbs *loadBalancersService) RemoveServers(lbID int, sIDs ...int) error {
	_, err := lbs.client.LoadBalancers.RemoveServers(context.TODO(), lbID, sIDs...)
	return err
}

func (lbs *loadBalancersService) AddForwardingRules(lbID int, rules ...binarylane.ForwardingRule) error {
	_, err := lbs.client.LoadBalancers.AddForwardingRules(context.TODO(), lbID, rules...)
	return err
}

func (lbs *loadBalancersService) RemoveForwardingRules(lbID int, rules ...binarylane.ForwardingRule) error {
	_, err := lbs.client.LoadBalancers.RemoveForwardingRules(context.TODO(), lbID, rules...)
	return err
}
