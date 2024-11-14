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
	"sync"

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
	hostmap := getMap(netclient)
	destinations, err := netmaker.GetWireGuardIPs(config.Network)
	if err != nil {
		slog.Error("unable to get wireguard IP for network", "network", config.Network, "test", "ping", "err", err)
	}

	var resultsMutex sync.Mutex
	var failuresMutex sync.Mutex
	failures := make(map[string]string)
	results := make(map[string]map[string]bool)

	wg := sync.WaitGroup{}
	for _, hosts := range netclient {
		hosts := hosts
		wg.Add(1)
		go func() {
			defer wg.Done()
			source := hosts.Host.EndpointIP
			slog.Info("ping from", "host", hosts.Host.Name, "ip", source)

			localResults := make(map[string]bool)
			var localFailures string

			for _, destination := range destinations {
				if hostmap[destination.String()] == hosts.Host.Name {
					//skip self
					continue
				}
				out, err := ssh.Run([]byte(config.Key), source, "ping -c 10 "+destination.String()+" | grep packet")
				if err != nil {
					slog.Error("error connecting to host", "host", hosts.Host.Name, "ip", source, "test", "ping", "err", err)
					localFailures = "unable to connect"
					break
				}
				localResults[hostmap[destination.String()]] = true
				if strings.Contains(out, ", 10% packet loss") || strings.Contains(out, ", 20% packet loss") || strings.Contains(out, ", 30% packet loss") {
					slog.Warn("packet loss", "host", hosts.Host.Name, "destination", destination, "output", strings.TrimSuffix(out, "\n"))
					continue
				}
				if strings.Contains(out, ", 0% packet loss") {
					slog.Info("ping success", "host", hosts.Host.Name, "destination", destination, "output", strings.TrimSuffix(out, "\n"))
					continue
				}
				slog.Error("failed to ping", "host", hosts.Host.Name, "destination", destination, "output", out)
				if localFailures == "" {
					localFailures = hostmap[destination.String()]
				} else {
					localFailures += " " + hostmap[destination.String()]
				}
				localResults[hostmap[destination.String()]] = false
			}

			resultsMutex.Lock()
			results[hosts.Host.Name] = localResults
			resultsMutex.Unlock()

			if localFailures != "" {
				failuresMutex.Lock()
				failures[hosts.Host.Name] = localFailures
				failuresMutex.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(failures) > 0 {
		for k, v := range failures {
			slog.Error("ping failures", "host", k, "failure", v)
		}
		return false
	}
	slog.Info("all nodes can ping each other")
	return true
}

func getMap(netclient []netmaker.Netclient) map[string]string {
	hosts := make(map[string]string)
	for _, client := range netclient {
		ip, _, err := net.ParseCIDR(client.Node.Address)
		if err != nil {
			continue
		}
		hosts[ip.String()] = client.Host.Name
		ip, _, err = net.ParseCIDR(client.Node.Address6)
		if err != nil {
			continue
		}
		hosts[ip.String()] = client.Host.Name + "6"
	}
	return hosts
}
