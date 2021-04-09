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

// FloatingIP wraps a binarylane FloatingIP.
type FloatingIP struct {
	*binarylane.FloatingIP
}

// FloatingIPs is a slice of FloatingIP.
type FloatingIPs []FloatingIP

// FloatingIPsService is the binarylane FloatingIPsService interface.
type FloatingIPsService interface {
	List() (FloatingIPs, error)
	Get(ip string) (*FloatingIP, error)
	Create(ficr *binarylane.FloatingIPCreateRequest) (*FloatingIP, error)
	Delete(ip string) error
}

type floatingIPsService struct {
	client *binarylane.Client
}

var _ FloatingIPsService = &floatingIPsService{}

// NewFloatingIPsService builds an instance of FloatingIPsService.
func NewFloatingIPsService(client *binarylane.Client) FloatingIPsService {
	return &floatingIPsService{
		client: client,
	}
}

func (fis *floatingIPsService) List() (FloatingIPs, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := fis.client.FloatingIPs.List(context.TODO(), opt)
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

	var list FloatingIPs
	for _, x := range si {
		fip := x.(binarylane.FloatingIP)
		list = append(list, FloatingIP{FloatingIP: &fip})
	}

	return list, nil
}

func (fis *floatingIPsService) Get(ip string) (*FloatingIP, error) {
	fip, _, err := fis.client.FloatingIPs.Get(context.TODO(), ip)
	if err != nil {
		return nil, err
	}

	return &FloatingIP{FloatingIP: fip}, nil
}

func (fis *floatingIPsService) Create(ficr *binarylane.FloatingIPCreateRequest) (*FloatingIP, error) {
	fip, _, err := fis.client.FloatingIPs.Create(context.TODO(), ficr)
	if err != nil {
		return nil, err
	}

	return &FloatingIP{FloatingIP: fip}, nil
}

func (fis *floatingIPsService) Delete(ip string) error {
	_, err := fis.client.FloatingIPs.Delete(context.TODO(), ip)
	return err
}
