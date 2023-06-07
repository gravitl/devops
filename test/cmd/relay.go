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
	"net"
	"strings"
	"time"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// relayCmd represents the relay command
var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "run a relay test",
	Long: `creates a relay and
	verifies all other nodes have received the update`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("relay")
		fmt.Println(relaytest(&config))
	},
}

func init() {
	rootCmd.AddCommand(relayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// relayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// relayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func relaytest(config *netmaker.Config) bool {
	pass := true
	netmaker.SetCxt(config.Api, config.Masterkey)
	netclient := netmaker.GetNetclient(config.Network)
	relay := netmaker.GetHost("relay", netclient)
	if relay == nil {
		slog.Error("did not find relay netclient", "test", "relay")
		return false
	}
	relayed := netmaker.GetHost("relayed", netclient)
	if relayed == nil {
		slog.Error("did not find relayed netclient", "test", "relay")
		return false
	}
	slog.Info("creating relay")
	netmaker.CreateRelay(relay, relayed)
	slog.Info("setting firewall rules")
	egress := netmaker.GetHost("egress", netclient)
	if egress == nil {
		slog.Error("did not find egress netclient", "test", "relay")
		return false
	}
	_, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "iptables -A OUTPUT -d "+egress.Host.EndpointIP+" -j DROP")
	if err != nil {
		slog.Error("failed to set firewall rule on relayed", "test", "relay")
		return false
	}
	_, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "iptables -A OUTPUT -d "+relayed.Host.EndpointIP+" -j DROP")
	if err != nil {
		slog.Error("failed to set firewall rule on relayed", "test", "relay")
		return false
	}
	slog.Info("waiting for changes to propogate")
	time.Sleep(time.Second * 30)
	defer resetFirewall(relayed, egress)
	slog.Info("ping egress from relayed")
	ip, _, err := net.ParseCIDR(egress.Node.Address)
	if err != nil {
		slog.Error("failed to parse egress address", egress.Node.Address)
		return false
	}
	out, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "ping -c 3 "+ip.String())
	if err != nil {
		slog.Error("error connecting to relayed", "err", err)
		pass = false
	} else {
		if !strings.Contains(out, "3 received") {
			slog.Error("failed to ping egress from relayed", "output", out)
			pass = false
		}
	}
	slog.Info("ping relayed from egress")
	ip, _, err = net.ParseCIDR(relayed.Node.Address)
	if err != nil {
		slog.Error("failed to parse relayed address", "address", relayed.Node.Address, "test", "relay")
		return false
	}
	out, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "ping -c 3 "+ip.String())
	if err != nil {
		slog.Error("error connecting to egress", "test", "relay", "err", err)
		pass = true
	} else {
		if !strings.Contains(out, "3 received") {
			slog.Error("failed to ping relayed from egress", "test", "relay", "output", out)
			pass = true
		}
	}
	return pass
}

func resetFirewall(relayed, egress *netmaker.Netclient) {
	slog.Info("reseting firewall on relayed/egress")
	_, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "iptables -D OUTPUT -d "+egress.Host.EndpointIP+" -j DROP")
	if err != nil {
		slog.Error("failed to set firewall rule on relayed", "test", "relay")
	}
	_, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "iptables -D OUTPUT -d "+relayed.Host.EndpointIP+" -j DROP")
	if err != nil {
		slog.Error("failed to set firewall rule on relayed", "test", "relay")
	}
}
