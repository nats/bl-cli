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

// ServerActionsService is an interface for interacting with BinaryLane's server action api.
type ServerActionsService interface {
	Shutdown(int) (*Action, error)
	ShutdownByTag(string) (Actions, error)
	PowerOff(int) (*Action, error)
	PowerOffByTag(string) (Actions, error)
	PowerOn(int) (*Action, error)
	PowerOnByTag(string) (Actions, error)
	PowerCycle(int) (*Action, error)
	PowerCycleByTag(string) (Actions, error)
	Reboot(int) (*Action, error)
	Restore(int, int) (*Action, error)
	Resize(int, string, bool) (*Action, error)
	Rename(int, string) (*Action, error)
	Snapshot(int, string) (*Action, error)
	SnapshotByTag(string, string) (Actions, error)
	EnableBackups(int) (*Action, error)
	EnableBackupsByTag(string) (Actions, error)
	DisableBackups(int) (*Action, error)
	DisableBackupsByTag(string) (Actions, error)
	PasswordReset(int) (*Action, error)
	RebuildByImageID(int, int) (*Action, error)
	RebuildByImageSlug(int, string) (*Action, error)
	ChangeKernel(int, int) (*Action, error)
	EnableIPv6(int) (*Action, error)
	EnableIPv6ByTag(string) (Actions, error)
	EnablePrivateNetworking(int) (*Action, error)
	EnablePrivateNetworkingByTag(string) (Actions, error)
	Get(int, int) (*Action, error)
	GetByURI(string) (*Action, error)
}

type serverActionsService struct {
	client *binarylane.Client
}

var _ ServerActionsService = &serverActionsService{}

// NewServerActionsService builds an instance of ServerActionsService.
func NewServerActionsService(client *binarylane.Client) ServerActionsService {
	return &serverActionsService{
		client: client,
	}
}

func (sas *serverActionsService) handleActionResponse(a *binarylane.Action, err error) (*Action, error) {
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}

func (sas *serverActionsService) handleTagActionResponse(a []binarylane.Action, err error) (Actions, error) {
	if err != nil {
		return nil, err
	}

	actions := make([]Action, 0, len(a))

	for _, action := range a {
		actions = append(actions, Action{Action: &action})
	}

	return actions, nil
}

func (sas *serverActionsService) Shutdown(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.Shutdown(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) ShutdownByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.ShutdownByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) PowerOff(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.PowerOff(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) PowerOffByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.PowerOffByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) PowerOn(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.PowerOn(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) PowerOnByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.PowerOnByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) PowerCycle(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.PowerCycle(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) PowerCycleByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.PowerCycleByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) Reboot(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.Reboot(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) Restore(id, imageID int) (*Action, error) {
	a, _, err := sas.client.ServerActions.Restore(context.TODO(), id, imageID)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) Resize(id int, sizeSlug string, resizeDisk bool) (*Action, error) {
	a, _, err := sas.client.ServerActions.Resize(context.TODO(), id, sizeSlug, resizeDisk)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) Rename(id int, name string) (*Action, error) {
	a, _, err := sas.client.ServerActions.Rename(context.TODO(), id, name)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) Snapshot(id int, name string) (*Action, error) {
	a, _, err := sas.client.ServerActions.Snapshot(context.TODO(), id, name)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) SnapshotByTag(tag string, name string) (Actions, error) {
	a, _, err := sas.client.ServerActions.SnapshotByTag(context.TODO(), tag, name)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) EnableBackups(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.EnableBackups(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) EnableBackupsByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.EnableBackupsByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) DisableBackups(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.DisableBackups(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) DisableBackupsByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.DisableBackupsByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) PasswordReset(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.PasswordReset(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) RebuildByImageID(id, imageID int) (*Action, error) {
	a, _, err := sas.client.ServerActions.RebuildByImageID(context.TODO(), id, imageID)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) RebuildByImageSlug(id int, slug string) (*Action, error) {
	a, _, err := sas.client.ServerActions.RebuildByImageSlug(context.TODO(), id, slug)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) ChangeKernel(id, kernelID int) (*Action, error) {
	a, _, err := sas.client.ServerActions.ChangeKernel(context.TODO(), id, kernelID)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) EnableIPv6(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.EnableIPv6(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) EnableIPv6ByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.EnableIPv6ByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) EnablePrivateNetworking(id int) (*Action, error) {
	a, _, err := sas.client.ServerActions.EnablePrivateNetworking(context.TODO(), id)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) EnablePrivateNetworkingByTag(tag string) (Actions, error) {
	a, _, err := sas.client.ServerActions.EnablePrivateNetworkingByTag(context.TODO(), tag)
	return sas.handleTagActionResponse(a, err)
}

func (sas *serverActionsService) Get(id int, actionID int) (*Action, error) {
	a, _, err := sas.client.ServerActions.Get(context.TODO(), id, actionID)
	return sas.handleActionResponse(a, err)
}

func (sas *serverActionsService) GetByURI(uri string) (*Action, error) {
	a, _, err := sas.client.ServerActions.GetByURI(context.TODO(), uri)
	return sas.handleActionResponse(a, err)
}
