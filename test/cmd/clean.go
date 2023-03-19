/*
Copyright © 2023 Matthew R Kasun <mkasun@nusak.ca>

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
package cmd

import (
	"fmt"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean up network",
	Long: `cleans up network to facilitate tests
	remove all gateways and removes interface/conf file on extclients`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("clean")
		fmt.Println(cleanNetwork(&config))
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cleanNetwork(config *netmaker.Config) bool {
	pass := true
	netclient := netmaker.GetNetclient(config.Network)
	for _, machine := range netclient {
		if machine.Node.IsEgressGateway {
			slog.Info("deleting egress ", machine.Host.Name)
			netmaker.DeleteEgress(machine.Node.ID, machine.Node.Network)
		}
		if machine.Node.IsIngressGateway {
			slog.Info("deleting ingress", machine.Host.Name)
			netmaker.DeleteIngress(machine.Node.ID, machine.Node.Network)
		}
		if machine.Host.IsRelay {
			slog.Info("deleting relay", machine.Host.Name)
			netmaker.DeleteRelay(machine.Host.ID)
		}
	}
	slog.Info("reseting extclient")
	logger.Info("resteting extclient")
	if err := netmaker.RestoreExtClient(config); err != nil {
		slog.Error("restoring extclient", "err", err)
		pass = false
	}
	relayed := netmaker.GetHost("relayed", netclient)
	if relayed == nil {
		slog.Error("did not find relayed netclient")
		pass = false
	}
	egress := netmaker.GetHost("egress", netclient)
	if egress == nil {
		slog.Error("did not find egress netclient")
		pass = false
	}
	if relayed != nil && egress != nil {
		slog.Info("reseting firewall on relayed/egress")
		_, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "iptables -D OUTPUT -d "+egress.Host.EndpointIP+" -j DROP")
		if err != nil {
			slog.Warn("error resetting egress firewall" + err.Error())
		}
		_, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "iptables -D OUTPUT -d "+relayed.Host.EndpointIP+" -j DROP")
		if err != nil {
			slog.Warn("error resetting egress firewall" + err.Error())
		}
	}
	return pass
}
