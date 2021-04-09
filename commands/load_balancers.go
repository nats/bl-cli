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

package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/binarylane/go-binarylane"
	"github.com/spf13/cobra"
)

// LoadBalancer creates the load balancer command.
func LoadBalancer() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "load-balancer",
			Short: "Display commands to manage load balancers",
			Long: `The sub-commands of ` + "`" + `bl compute load-balancer` + "`" + ` manage your load balancers.

With the load-balancer command, you can list, create, or delete load balancers, and manage their configuration details.`,
		},
	}
	lbDetail := `

- The load balancer's ID
- The load balancer's name
- The load balancer's IP address
- The load balancer's traffic algorithm. Must
  be either ` + "`" + `round_robin` + "`" + ` or ` + "`" + `least_connections` + "`" + `
- The current state of the load balancer. This can be ` + "`" + `new` + "`" + `, ` + "`" + `active` + "`" + `, or ` + "`" + `errored` + "`" + `.
- The load balancer's creation date, in ISO8601 combined date and time format.
- The load balancer's forwarding rules. See ` + "`" + `bl compute load-balancer add-forwarding-rules --help` + "`" + ` for a list.
- The ` + "`" + `health_check` + "`" + ` settings for the load balancer.
- The ` + "`" + `sticky_sessions` + "`" + ` settings for the load balancer.
- The datacenter region the load balancer is located in.
- The Server tag corresponding to the Servers assigned to the load balancer.
- The IDs of the Servers assigned to the load balancer.
- Whether HTTP request to the load balancer on port 80 will be redirected to HTTPS on port 443.
- Whether the PROXY protocol is in use on the load balancer.
`
	forwardingDetail := `

- ` + "`" + `entry_protocol` + "`" + `: The entry protocol used for traffic to the load balancer. Possible values are: ` + "`" + `http` + "`" + `, ` + "`" + `https` + "`" + `, ` + "`" + `http2` + "`" + `, or ` + "`" + `tcp` + "`" + `.
- ` + "`" + `entry_port` + "`" + `: The entry port used for incoming traffic to the load balancer.
- ` + "`" + `target_protocol` + "`" + `: The target protocol used for traffic from the load balancer to the backend Servers. Possible values are: ` + "`" + `http` + "`" + `, ` + "`" + `https` + "`" + `, ` + "`" + `http2` + "`" + `, or ` + "`" + `tcp` + "`" + `.
- ` + "`" + `target_port` + "`" + `: The target port used to send traffic from the load balancer to the backend Servers.
- ` + "`" + `certificate_id` + "`" + `: The ID of the TLS certificate used for SSL termination, if enabled. Can be obtained with ` + "`" + `bl certificate list` + "`" + `
- ` + "`" + `tls_passthrough` + "`" + `: Whether SSL passthrough is enabled on the load balancer.
`
	forwardingRulesTxt := "A comma-separated list of key-value pairs representing forwarding rules, which define how traffic is routed, e.g.: `entry_protocol:tcp, entry_port:3306, target_protocol:tcp, target_port:3306`. Use a quoted string of space-separated values for multiple rules"
	CmdBuilder(cmd, RunLoadBalancerGet, "get <id>", "Retrieve a load balancer", "Use this command to retrieve information about a load balancer instance, including:"+lbDetail, Writer,
		aliasOpt("g"), displayerType(&displayers.LoadBalancer{}))

	cmdRecordCreate := CmdBuilder(cmd, RunLoadBalancerCreate, "create",
		"Create a new load balancer", "Use this command to create a new load balancer on your account. Valid forwarding rules are:"+forwardingDetail, Writer, aliasOpt("c"))
	AddStringFlag(cmdRecordCreate, blcli.ArgLoadBalancerName, "", "",
		"The load balancer's name", requiredOpt())
	AddStringFlag(cmdRecordCreate, blcli.ArgRegionSlug, "", "",
		"The load balancer's region, e.g.: `syd`", requiredOpt())
	AddStringFlag(cmdRecordCreate, blcli.ArgSizeSlug, "", "lb-small",
		"The load balancer's size, e.g.: `lb-small`", requiredOpt())
	AddStringFlag(cmdRecordCreate, blcli.ArgVPCID, "", "", "The ID of the VPC to create the load balancer in")
	AddStringFlag(cmdRecordCreate, blcli.ArgLoadBalancerAlgorithm, "",
		"round_robin", "The algorithm to use when traffic is distributed across your Servers; possible values: `round_robin` or `least_connections`")
	AddBoolFlag(cmdRecordCreate, blcli.ArgRedirectHTTPToHTTPS, "", false,
		"Redirects HTTP requests to the load balancer on port 80 to HTTPS on port 443")
	AddBoolFlag(cmdRecordCreate, blcli.ArgEnableProxyProtocol, "", false,
		"enable proxy protocol")
	AddBoolFlag(cmdRecordCreate, blcli.ArgEnableBackendKeepalive, "", false,
		"enable keepalive connections to backend target servers")
	AddStringFlag(cmdRecordCreate, blcli.ArgTagName, "", "", "server tag name")
	AddStringSliceFlag(cmdRecordCreate, blcli.ArgServerIDs, "", []string{},
		"A comma-separated list of Server IDs to add to the load balancer, e.g.: `12,33`")
	AddStringFlag(cmdRecordCreate, blcli.ArgStickySessions, "", "",
		"A comma-separated list of key-value pairs representing a list of active sessions, e.g.: `type:cookies, cookie_name:DO-LB, cookie_ttl_seconds:5`")
	AddStringFlag(cmdRecordCreate, blcli.ArgHealthCheck, "", "",
		"A comma-separated list of key-value pairs representing recent health check results, e.g.: `protocol:http, port:80, path:/index.html, check_interval_seconds:10, response_timeout_seconds:5, healthy_threshold:5, unhealthy_threshold:3`")
	AddStringFlag(cmdRecordCreate, blcli.ArgForwardingRules, "", "",
		forwardingRulesTxt)

	cmdRecordUpdate := CmdBuilder(cmd, RunLoadBalancerUpdate, "update <id>",
		"Update a load balancer's configuration", `Use this command to update the configuration of a specified load balancer.`, Writer, aliasOpt("u"))
	AddStringFlag(cmdRecordUpdate, blcli.ArgLoadBalancerName, "", "",
		"The load balancer's name", requiredOpt())
	AddStringFlag(cmdRecordUpdate, blcli.ArgRegionSlug, "", "",
		"The load balancer's region, e.g.: `syd`", requiredOpt())
	AddStringFlag(cmdRecordUpdate, blcli.ArgSizeSlug, "", "",
		"The load balancer's size, e.g.: `lb-small`", requiredOpt())
	AddStringFlag(cmdRecordUpdate, blcli.ArgVPCID, "", "", "The ID of the VPC to create the load balancer in")
	AddStringFlag(cmdRecordUpdate, blcli.ArgLoadBalancerAlgorithm, "",
		"round_robin", "The algorithm to use when traffic is distributed across your Servers; possible values: `round_robin` or `least_connections`")
	AddBoolFlag(cmdRecordUpdate, blcli.ArgRedirectHTTPToHTTPS, "", false,
		"Flag to redirect HTTP requests to the load balancer on port 80 to HTTPS on port 443")
	AddBoolFlag(cmdRecordUpdate, blcli.ArgEnableProxyProtocol, "", false,
		"enable proxy protocol")
	AddBoolFlag(cmdRecordUpdate, blcli.ArgEnableBackendKeepalive, "", false,
		"enable keepalive connections to backend target servers")
	AddStringFlag(cmdRecordUpdate, blcli.ArgTagName, "", "", "Assigns Servers with the specified tag to the load balancer")
	AddStringSliceFlag(cmdRecordUpdate, blcli.ArgServerIDs, "", []string{},
		"A comma-separated list of Server IDs, e.g.: `215,378`")
	AddStringFlag(cmdRecordUpdate, blcli.ArgStickySessions, "", "",
		"A comma-separated list of key-value pairs representing a list of active sessions, e.g.: `type:cookies, cookie_name:DO-LB, cookie_ttl_seconds:5`")
	AddStringFlag(cmdRecordUpdate, blcli.ArgHealthCheck, "", "",
		"A comma-separated list of key-value pairs representing recent health check results, e.g.: `protocol:http, port:80, path:/index.html, check_interval_seconds:10, response_timeout_seconds:5, healthy_threshold:5, unhealthy_threshold:3`")
	AddStringFlag(cmdRecordUpdate, blcli.ArgForwardingRules, "", "", forwardingRulesTxt)

	CmdBuilder(cmd, RunLoadBalancerList, "list", "List load balancers", "Use this command to get a list of the load balancers on your account, including the following information for each:"+lbDetail, Writer,
		aliasOpt("ls"), displayerType(&displayers.LoadBalancer{}))

	cmdRunRecordDelete := CmdBuilder(cmd, RunLoadBalancerDelete, "delete <id>",
		"Permanently delete a load balancer", `Use this command to permanently delete the speicified load balancer. This is irreversable.`, Writer, aliasOpt("d", "rm"))
	AddBoolFlag(cmdRunRecordDelete, blcli.ArgForce, blcli.ArgShortForce, false,
		"Delete the load balancer without a confirmation prompt")

	cmdAddServers := CmdBuilder(cmd, RunLoadBalancerAddServers, "add-servers <id>",
		"Add Servers to a load balancer", `Use this command to add Servers to a load balancer.`, Writer)
	AddStringSliceFlag(cmdAddServers, blcli.ArgServerIDs, "", []string{},
		"A comma-separated list of IDs of Server to add to the load balancer, example value: `12,33`")

	cmdRemoveServers := CmdBuilder(cmd, RunLoadBalancerRemoveServers,
		"remove-servers <id>", "Remove Servers from a load balancer", `Use this command to remove Servers from a load balancer. This command does not destroy any Servers.`, Writer)
	AddStringSliceFlag(cmdRemoveServers, blcli.ArgServerIDs, "", []string{},
		"A comma-separated list of IDs of Servers to remove from the load balancer, example value: `12,33`")

	cmdAddForwardingRules := CmdBuilder(cmd, RunLoadBalancerAddForwardingRules,
		"add-forwarding-rules <id>", "Add forwarding rules to a load balancer", "Use this command to add forwarding rules to a load balancer, specified with the `--forwarding-rules` flag. Valid rules include:"+forwardingDetail, Writer)
	AddStringFlag(cmdAddForwardingRules, blcli.ArgForwardingRules, "", "", forwardingRulesTxt)

	cmdRemoveForwardingRules := CmdBuilder(cmd, RunLoadBalancerRemoveForwardingRules,
		"remove-forwarding-rules <id>", "Remove forwarding rules from a load balancer", "Use this command to remove forwarding rules from a load balancer, specified with the `--forwarding-rules` flag. Valid rules include:"+forwardingDetail, Writer)
	AddStringFlag(cmdRemoveForwardingRules, blcli.ArgForwardingRules, "", "", forwardingRulesTxt)

	return cmd
}

// RunLoadBalancerGet retrieves an existing load balancer by its identifier.
func RunLoadBalancerGet(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	lbs := c.LoadBalancers()
	lb, err := lbs.Get(id)
	if err != nil {
		return err
	}

	item := &displayers.LoadBalancer{LoadBalancers: bl.LoadBalancers{*lb}}
	return c.Display(item)
}

// RunLoadBalancerList lists load balancers.
func RunLoadBalancerList(c *CmdConfig) error {
	lbs := c.LoadBalancers()
	list, err := lbs.List()
	if err != nil {
		return err
	}

	item := &displayers.LoadBalancer{LoadBalancers: list}
	return c.Display(item)
}

// RunLoadBalancerCreate creates a new load balancer with a given configuration.
func RunLoadBalancerCreate(c *CmdConfig) error {
	r := new(binarylane.LoadBalancerRequest)
	if err := buildRequestFromArgs(c, r); err != nil {
		return err
	}

	lbs := c.LoadBalancers()
	lb, err := lbs.Create(r)
	if err != nil {
		return err
	}

	item := &displayers.LoadBalancer{LoadBalancers: bl.LoadBalancers{*lb}}
	return c.Display(item)
}

// RunLoadBalancerUpdate updates an existing load balancer with new configuration.
func RunLoadBalancerUpdate(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return blcli.NewMissingArgsErr(c.NS)
	}
	lbID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	r := new(binarylane.LoadBalancerRequest)
	if err := buildRequestFromArgs(c, r); err != nil {
		return err
	}

	lbs := c.LoadBalancers()
	lb, err := lbs.Update(lbID, r)
	if err != nil {
		return err
	}

	item := &displayers.LoadBalancer{LoadBalancers: bl.LoadBalancers{*lb}}
	return c.Display(item)
}

// RunLoadBalancerDelete deletes a load balancer by its identifier.
func RunLoadBalancerDelete(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	lbID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	force, err := c.Doit.GetBool(c.NS, blcli.ArgForce)
	if err != nil {
		return err
	}

	if force || AskForConfirmDelete("load balancer", 1) == nil {
		lbs := c.LoadBalancers()
		if err := lbs.Delete(lbID); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Operation aborted.")
	}

	return nil
}

// RunLoadBalancerAddServers adds servers to a load balancer.
func RunLoadBalancerAddServers(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	lbID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	serverIDsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgServerIDs)
	if err != nil {
		return err
	}

	serverIDs, err := extractServerIDs(serverIDsList)
	if err != nil {
		return err
	}

	return c.LoadBalancers().AddServers(lbID, serverIDs...)
}

// RunLoadBalancerRemoveServers removes servers from a load balancer.
func RunLoadBalancerRemoveServers(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	lbID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	serverIDsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgServerIDs)
	if err != nil {
		return err
	}

	serverIDs, err := extractServerIDs(serverIDsList)
	if err != nil {
		return err
	}

	return c.LoadBalancers().RemoveServers(lbID, serverIDs...)
}

// RunLoadBalancerAddForwardingRules adds forwarding rules to a load balancer.
func RunLoadBalancerAddForwardingRules(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	lbID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	fra, err := c.Doit.GetString(c.NS, blcli.ArgForwardingRules)
	if err != nil {
		return err
	}

	forwardingRules, err := extractForwardingRules(fra)
	if err != nil {
		return err
	}

	return c.LoadBalancers().AddForwardingRules(lbID, forwardingRules...)
}

// RunLoadBalancerRemoveForwardingRules removes forwarding rules from a load balancer.
func RunLoadBalancerRemoveForwardingRules(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	lbID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	fra, err := c.Doit.GetString(c.NS, blcli.ArgForwardingRules)
	if err != nil {
		return err
	}

	forwardingRules, err := extractForwardingRules(fra)
	if err != nil {
		return err
	}

	return c.LoadBalancers().RemoveForwardingRules(lbID, forwardingRules...)
}

func extractForwardingRules(s string) (forwardingRules []binarylane.ForwardingRule, err error) {
	if len(s) == 0 {
		return forwardingRules, err
	}

	list := strings.Split(s, " ")

	for _, v := range list {
		forwardingRule := new(binarylane.ForwardingRule)
		if err := fillStructFromStringSliceArgs(forwardingRule, v); err != nil {
			return nil, err
		}

		forwardingRules = append(forwardingRules, *forwardingRule)
	}

	return forwardingRules, err
}

func fillStructFromStringSliceArgs(obj interface{}, s string) error {
	if len(s) == 0 {
		return nil
	}

	kvs := strings.Split(s, ",")
	m := map[string]string{}

	for _, v := range kvs {
		p := strings.Split(v, ":")
		if len(p) == 2 {
			m[p[0]] = p[1]
		} else {
			return fmt.Errorf("Unexpected input value %v: must be a key:value pair", p)
		}
	}

	structValue := reflect.Indirect(reflect.ValueOf(obj))
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		f := structValue.Field(i)
		jv := strings.Split(structType.Field(i).Tag.Get("json"), ",")[0]

		if val, exists := m[jv]; exists {
			switch f.Kind() {
			case reflect.Bool:
				if v, err := strconv.ParseBool(val); err == nil {
					f.Set(reflect.ValueOf(v))
				}
			case reflect.Int:
				if v, err := strconv.Atoi(val); err == nil {
					f.Set(reflect.ValueOf(v))
				}
			case reflect.String:
				f.Set(reflect.ValueOf(val))
			default:
				return fmt.Errorf("Unexpected type for struct field %v", val)
			}
		}
	}

	return nil
}

func buildRequestFromArgs(c *CmdConfig, r *binarylane.LoadBalancerRequest) error {
	name, err := c.Doit.GetString(c.NS, blcli.ArgLoadBalancerName)
	if err != nil {
		return err
	}
	r.Name = name

	region, err := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	if err != nil {
		return err
	}
	r.Region = region

	size, err := c.Doit.GetString(c.NS, blcli.ArgSizeSlug)
	if err != nil {
		return err
	}
	r.SizeSlug = size

	algorithm, err := c.Doit.GetString(c.NS, blcli.ArgLoadBalancerAlgorithm)
	if err != nil {
		return err
	}
	r.Algorithm = algorithm

	tag, err := c.Doit.GetString(c.NS, blcli.ArgTagName)
	if err != nil {
		return err
	}
	r.Tag = tag

	vpcID, err := c.Doit.GetInt(c.NS, blcli.ArgVPCID)
	if err != nil {
		return err
	}
	r.VPCID = vpcID

	redirectHTTPToHTTPS, err := c.Doit.GetBool(c.NS, blcli.ArgRedirectHTTPToHTTPS)
	if err != nil {
		return err
	}
	r.RedirectHttpToHttps = redirectHTTPToHTTPS

	enableProxyProtocol, err := c.Doit.GetBool(c.NS, blcli.ArgEnableProxyProtocol)
	if err != nil {
		return err
	}
	r.EnableProxyProtocol = enableProxyProtocol

	enableBackendKeepalive, err := c.Doit.GetBool(c.NS, blcli.ArgEnableBackendKeepalive)
	if err != nil {
		return err
	}
	r.EnableBackendKeepalive = enableBackendKeepalive

	serverIDsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgServerIDs)
	if err != nil {
		return err
	}

	serverIDs, err := extractServerIDs(serverIDsList)
	if err != nil {
		return err
	}
	r.ServerIDs = serverIDs

	ssa, err := c.Doit.GetString(c.NS, blcli.ArgStickySessions)
	if err != nil {
		return err
	}

	stickySession := new(binarylane.StickySessions)
	if err := fillStructFromStringSliceArgs(stickySession, ssa); err != nil {
		return err
	}
	r.StickySessions = stickySession

	hca, err := c.Doit.GetString(c.NS, blcli.ArgHealthCheck)
	if err != nil {
		return err
	}

	healthCheck := new(binarylane.HealthCheck)
	if err := fillStructFromStringSliceArgs(healthCheck, hca); err != nil {
		return err
	}
	r.HealthCheck = healthCheck

	fra, err := c.Doit.GetString(c.NS, blcli.ArgForwardingRules)
	if err != nil {
		return err
	}

	forwardingRules, err := extractForwardingRules(fra)
	if err != nil {
		return err
	}
	r.ForwardingRules = forwardingRules

	return nil
}
