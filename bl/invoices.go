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

// Invoice is a wrapper for binarylane.Invoice
type Invoice struct {
	*binarylane.Invoice
}

// InvoiceItem is a wrapper for binarylane.InvoiceItem
type InvoiceItem struct {
	*binarylane.InvoiceItem
}

// InvoiceSummary is a wrapper for binarylane.InvoiceSummary
type InvoiceSummary struct {
	*binarylane.InvoiceSummary
}

// InvoiceList is a the results when listing invoices
type InvoiceList struct {
	*binarylane.InvoiceList
}

// InvoicesService is an interface for interacting with BinaryLane's invoices api.
type InvoicesService interface {
	Get(string) (*Invoice, error)
	List() (*InvoiceList, error)
	GetSummary(string) (*InvoiceSummary, error)
	GetPDF(string) ([]byte, error)
	GetCSV(string) ([]byte, error)
}

type invoicesService struct {
	client *binarylane.Client
}

var _ InvoicesService = &invoicesService{}

// NewInvoicesService builds an InvoicesService instance.
func NewInvoicesService(client *binarylane.Client) InvoicesService {
	return &invoicesService{
		client: client,
	}
}

func (is *invoicesService) List() (*InvoiceList, error) {
	var invoicePreview binarylane.InvoiceListItem

	listFn := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		invoiceList, resp, err := is.client.Invoices.List(context.Background(), opt)
		if err != nil {
			return nil, nil, err
		}
		invoicePreview = invoiceList.InvoicePreview

		si := make([]interface{}, len(invoiceList.Invoices))
		for i := range invoiceList.Invoices {
			si[i] = invoiceList.Invoices[i]
		}
		return si, resp, err
	}

	paginatedList, err := PaginateResp(listFn)
	if err != nil {
		return nil, err
	}
	list := make([]binarylane.InvoiceListItem, len(paginatedList))
	for i := range paginatedList {
		list[i] = paginatedList[i].(binarylane.InvoiceListItem)
	}

	return &InvoiceList{
		InvoiceList: &binarylane.InvoiceList{
			Invoices:       list,
			InvoicePreview: invoicePreview,
		},
	}, nil
}

func (is *invoicesService) Get(uuid string) (*Invoice, error) {
	listFn := func(opt *binarylane.ListOptions) ([]interface{}, *binarylane.Response, error) {
		invoice, resp, err := is.client.Invoices.Get(context.Background(), uuid, opt)
		if err != nil {
			return nil, nil, err
		}
		si := make([]interface{}, len(invoice.InvoiceItems))
		for i := range invoice.InvoiceItems {
			si[i] = invoice.InvoiceItems[i]
		}
		return si, resp, err
	}

	paginatedList, err := PaginateResp(listFn)
	if err != nil {
		return nil, err
	}

	list := make([]binarylane.InvoiceItem, len(paginatedList))
	for i := range paginatedList {
		list[i] = paginatedList[i].(binarylane.InvoiceItem)
	}

	return &Invoice{
		Invoice: &binarylane.Invoice{
			InvoiceItems: list,
		},
	}, nil
}

func (is *invoicesService) GetSummary(uuid string) (*InvoiceSummary, error) {
	summary, _, err := is.client.Invoices.GetSummary(context.Background(), uuid)
	if err != nil {
		return nil, err
	}

	return &InvoiceSummary{InvoiceSummary: summary}, nil
}

func (is *invoicesService) GetPDF(uuid string) ([]byte, error) {
	pdf, _, err := is.client.Invoices.GetPDF(context.Background(), uuid)
	if err != nil {
		return nil, err
	}

	return pdf, nil
}

func (is *invoicesService) GetCSV(uuid string) ([]byte, error) {
	csv, _, err := is.client.Invoices.GetCSV(context.Background(), uuid)
	if err != nil {
		return nil, err
	}

	return csv, nil
}
