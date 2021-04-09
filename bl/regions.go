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

// Region wraps binarylane Region.
type Region struct {
	*binarylane.Region
}

// Regions is a slice of Region.
type Regions []Region

// RegionsService is the binarylane RegionsService interface.
type RegionsService interface {
	List() (Regions, error)
}

type regionsService struct {
	client *binarylane.Client
}

var _ RegionsService = &regionsService{}

// NewRegionsService builds an instance of RegionsService.
func NewRegionsService(client *binarylane.Client) RegionsService {
	return &regionsService{
		client: client,
	}
}

func (rs *regionsService) List() (Regions, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := rs.client.Regions.List(context.TODO(), opt)
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

	list := make(Regions, len(si))
	for i := range si {
		r := si[i].(binarylane.Region)
		list[i] = Region{Region: &r}
	}

	return list, nil
}
