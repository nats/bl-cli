#!/bin/bash

# regenerated generated mocks

set -euo pipefail

cd "bl"

GO111MODULE=off go get -u github.com/golang/mock/mockgen

mockgen -source account.go -package=mocks AccountService > mocks/AccountService.go
mockgen -source actions.go -package=mocks ActionService > mocks/ActionService.go
mockgen -source balance.go -package=mocks BalanceService > mocks/BalanceService.go
mockgen -source billing_history.go -package=mocks BillingHistoryService > mocks/BillingHistoryService.go
mockgen -source domains.go -package=mocks DomainService > mocks/DomainService.go
mockgen -source server_actions.go -package=mocks ServerActionsService > mocks/ServerActionService.go
mockgen -source servers.go -package=mocks ServersService > mocks/ServersService.go
mockgen -source firewalls.go -package=mocks FirewallsService > mocks/FirewallsService.go
mockgen -source floating_ip_actions.go -package=mocks FloatingIPActionsService > mocks/FloatingIPActionsService.go
mockgen -source floating_ips.go -package=mocks FloatingIPsService > mocks/FloatingIPsService.go
mockgen -source image_actions.go -package=mocks ImageActionsService > mocks/ImageActionsService.go
mockgen -source images.go -package=mocks ImageService > mocks/ImageService.go
mockgen -source invoices.go -package=mocks InvoicesService > mocks/InvoicesService.go
mockgen -source load_balancers.go -package=mocks LoadBalancersService > mocks/LoadBalancersService.go
mockgen -source projects.go -package=mocks ProjectsService > mocks/ProjectsService.go
mockgen -source regions.go -package=mocks RegionsService > mocks/RegionsService.go
mockgen -source snapshots.go -package=mocks SnapshotsService > mocks/SnapshotsService.go
mockgen -source sizes.go -package=mocks SizesService > mocks/SizesService.go
mockgen -source sshkeys.go -package=mocks KeysService > mocks/KeysService.go
mockgen -source tags.go -package=mocks TagsService > mocks/TagsService.go
mockgen -source vpcs.go -package=mocks VPCsService > mocks/VPCsService.go
mockgen -source ../pkg/runner/runner.go -package=mocks Runner > mocks/Runner.go
