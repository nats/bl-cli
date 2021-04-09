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

// BillingHistory is a wrapper for binarylane.BillingHistory
type BillingHistory struct {
	*binarylane.BillingHistory
}

// BillingHistoryService is an interface for interacting with BinaryLane's invoices api.
type BillingHistoryService interface {
	List() (*BillingHistory, error)
}

type billingHistoryService struct {
	client *binarylane.Client
}

var _ BillingHistoryService = &billingHistoryService{}

// NewBillingHistoryService builds an BillingHistoryService instance.
func NewBillingHistoryService(client *binarylane.Client) BillingHistoryService {
	return &billingHistoryService{
		client: client,
	}
}

func (is *billingHistoryService) List() (*BillingHistory, error) {
	listFn := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		historyList, resp, err := is.client.BillingHistory.List(context.Background(), opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(historyList.BillingHistory))
		for i := range historyList.BillingHistory {
			si[i] = historyList.BillingHistory[i]
		}
		return si, resp, err
	}

	paginatedList, err := PaginateResp(listFn)
	if err != nil {
		return nil, err
	}
	list := make([]binarylane.BillingHistoryEntry, len(paginatedList))
	for i := range paginatedList {
		list[i] = paginatedList[i].(binarylane.BillingHistoryEntry)
	}

	return &BillingHistory{BillingHistory: &binarylane.BillingHistory{BillingHistory: list}}, nil
}
