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

// Tag is a wrapper for binarylane.Tag
type Tag struct {
	*binarylane.Tag
}

// Tags is a slice of Tag.
type Tags []Tag

// TagsService is an interface for interacting with BinaryLane's tags api.
type TagsService interface {
	List() (Tags, error)
	Get(string) (*Tag, error)
	Create(*binarylane.TagCreateRequest) (*Tag, error)
	Delete(string) error
	TagResources(string, *binarylane.TagResourcesRequest) error
	UntagResources(string, *binarylane.UntagResourcesRequest) error
}

type tagsService struct {
	client *binarylane.Client
}

var _ TagsService = (*tagsService)(nil)

// NewTagsService builds a TagsService instance.
func NewTagsService(client *binarylane.Client) TagsService {
	return &tagsService{
		client: client,
	}
}

func (ts *tagsService) List() (Tags, error) {
	f := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		list, resp, err := ts.client.Tags.List(context.TODO(), opt)
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

	list := make(Tags, len(si))
	for i := range si {
		a := si[i].(binarylane.Tag)
		list[i] = Tag{Tag: &a}
	}

	return list, nil
}

func (ts *tagsService) Get(name string) (*Tag, error) {
	t, _, err := ts.client.Tags.Get(context.TODO(), name)
	if err != nil {
		return nil, err
	}

	return &Tag{Tag: t}, nil
}

func (ts *tagsService) Create(tcr *binarylane.TagCreateRequest) (*Tag, error) {
	t, _, err := ts.client.Tags.Create(context.TODO(), tcr)
	if err != nil {
		return nil, err
	}

	return &Tag{Tag: t}, nil
}

func (ts *tagsService) Delete(name string) error {
	_, err := ts.client.Tags.Delete(context.TODO(), name)
	return err
}

func (ts *tagsService) TagResources(name string, trr *binarylane.TagResourcesRequest) error {
	_, err := ts.client.Tags.TagResources(context.TODO(), name, trr)
	return err
}

func (ts *tagsService) UntagResources(name string, urr *binarylane.UntagResourcesRequest) error {
	_, err := ts.client.Tags.UntagResources(context.TODO(), name, urr)
	return err
}
