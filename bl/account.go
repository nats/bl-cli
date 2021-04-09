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

// Account is a wrapper for binarylane.Account.
type Account struct {
	*binarylane.Account
}

// RateLimit is a wrapper for binarylane.Rate.
type RateLimit struct {
	*binarylane.Rate
}

// AccountService is an interface for interacting with BinaryLane's account api.
type AccountService interface {
	Get() (*Account, error)
	RateLimit() (*RateLimit, error)
}

type accountService struct {
	client *binarylane.Client
}

var _ AccountService = &accountService{}

// NewAccountService builds an AccountService instance.
func NewAccountService(client *binarylane.Client) AccountService {
	return &accountService{
		client: client,
	}
}

func (as *accountService) Get() (*Account, error) {
	binarylaneAccount, _, err := as.client.Account.Get(context.TODO())
	if err != nil {
		return nil, err
	}

	account := &Account{Account: binarylaneAccount}
	return account, nil
}

func (as *accountService) RateLimit() (*RateLimit, error) {
	_, resp, err := as.client.Account.Get(context.TODO())
	if err != nil {
		return nil, err
	}

	rateLimit := &RateLimit{Rate: &resp.Rate}
	return rateLimit, nil
}
