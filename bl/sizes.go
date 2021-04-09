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

// Size wraps binarylane Size.
type Size struct {
	*binarylane.Size
}

// Sizes is a slice of Size.
type Sizes []Size

// SizesService is the binarylane SizesService interface.
type SizesService interface {
	List() (Sizes, error)
}

type sizesService struct {
	client *binarylane.Client
}

var _ SizesService = &sizesService{}

// NewSizesService builds an instance of SizesService.
func NewSizesService(client *binarylane.Client) SizesService {
	return &sizesService{
		client: client,
	}
}

func (rs *sizesService) List() (Sizes, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := rs.client.Sizes.List(context.TODO(), opt)
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

	list := make(Sizes, len(si))
	for i := range si {
		r := si[i].(binarylane.Size)
		list[i] = Size{Size: &r}
	}

	return list, nil
}
