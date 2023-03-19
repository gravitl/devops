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
	"strings"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "run a ping test",
	Long:  `ping all nodes on network and reports result`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("ping")
		fmt.Println(pingtest(&config))
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func pingtest(config *netmaker.Config) bool {
	netclient := netmaker.GetNetclient(config.Network)
	destinations, err := netmaker.GetWireGuardIPs(config.Network)
	if err != nil {
		slog.Error("unable to get wireguard IP for network", "network", config.Network, "test", "ping", "err", err)
	}
	failures := make(map[string]string)
	for _, hosts := range netclient {
		source := hosts.Host.EndpointIP
		slog.Info("ping from", "host", hosts.Host.Name, "ip", source)
		for _, destination := range destinations {
			out, err := ssh.Run([]byte(config.Key), source, "ping -c 3 "+destination.String())
			if err != nil {
				slog.Error("error connecting to host", "host", hosts.Host.Name, "test", "ping", "err", err)
				failures[hosts.Host.Name] = "unable to connect"
				break
			}
			if !strings.Contains(out, "3 received") {
				slog.Error("failed to ping", "destination", destination, "output", out)
				failures[hosts.Host.Name] = failures[hosts.Host.Name] + " " + destination.String()
				continue
			}
		}
	}
	if len(failures) > 0 {
		for k, v := range failures {
			slog.Error("ping failues", "host", k, "failure", v, "test", "ping")
		}
		return false
	}
	slog.Info("all nodes can ping each other")
	return true
}
