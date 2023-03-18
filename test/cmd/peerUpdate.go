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
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// peerUpdateCmd represents the peerUpdate command
var peerUpdateCmd = &cobra.Command{
	Use:   "peerUpdate",
	Short: "run peerupdate test",
	Long: `updates wg address of a node and
	verifies that all other nodes received the update
`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("peerupdate")
		peerupdatetest(&config)
	},
}

func init() {
	rootCmd.AddCommand(peerUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// peerUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// peerUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func peerupdatetest(config *netmaker.Config) {
	if config.Key == "" {
		key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_ed25519")
		if err != nil {
			slog.Error("invalid configuration, ssh key not set", "test", "peerupdate", "err", err)
			return
		}
		config.Key = string(key)
	}
	netmaker.SetCxt(config.Api, config.Masterkey)
	netclient := netmaker.GetNetclient(config.Network)
	server := netmaker.GetHost("server", netclient)
	if server == nil {
		slog.Error("did not find server", "test", "peerupdate")
		return
	}
	slog.Info("updating wg ip on server node")
	taken := make(map[string]bool)
	for _, machine := range netclient {
		taken[machine.Node.Address] = true
	}
	slog.Debug("debugging", "exclued ips ", taken)

	newip := getNextIP(server.Node.Address, taken)
	if newip == "" {
		return
	}
	server.Node.Address = newip
	slog.Info(fmt.Sprintf("updating wg address of %s to %s", server.Host.Name, newip))

	netmaker.UpdateNode(&server.Node)
	// check node received update
	//check if other nodes received update
	failedmachines := []string{}
	ip, _, err := net.ParseCIDR(newip)
	if err != nil {
		slog.Error("could not parse newip", "test", "peerupdate", "ip", ip, "err", err)
		return
	}
	for _, machine := range netclient {
		time.Sleep(time.Second)
		if machine.Host.Name == "server" {
			continue
		}
		if machine.Host.IsRelayed {
			continue
		}
		slog.Info(fmt.Sprintf("checking that %s @ %s received the update", machine.Host.Name, machine.Host.EndpointIP))
		out, err := ssh.Run([]byte(config.Key), machine.Host.EndpointIP, "wg show netmaker allowed-ips | grep "+ip.String())
		if err != nil {
			slog.Error("err connecting", "machine", machine.Host.Name, "test", "peerupdate", "err", err)
			failedmachines = append(failedmachines, machine.Host.Name)
			continue
		}
		if !strings.Contains(out, ip.String()) {
			slog.Error("node did not receive the update", "machine", machine.Host.Name, "test", "peerupdate", "ouput", out)
			failedmachines = append(failedmachines, machine.Host.Name)
			continue
		}
	}
	if len(failedmachines) > 0 {
		slog.Error("not all machines were updated", "test", "peerupdate")
		for _, machine := range failedmachines {
			slog.Error(machine, "test", "peerupdate")
		}
		return
	}
	slog.Info("all netclients received the update")
}

func getNextIP(current string, taken map[string]bool) string {
	var newip net.IP
	if len(taken) > 253 {
		slog.Error("no free ips", "test", "peerupdate")
		return ""
	}
	ip, cidr, err := net.ParseCIDR(current)
	if err != nil {
		slog.Error("failed to parse cidr", "test", "peerupdate", "err", err)
		return ""
	}
	slog.Info("getting free ip")
	net4 := iplib.Net4FromStr(current)
	for {
		newip, err = net4.NextIP(ip)
		if errors.Is(err, iplib.ErrBroadcastAddress) {
			newip, err = net4.NextIP(net4.FirstAddress())
		}
		if err != nil {
			slog.Error("NextIP", "test", "peerupdate", "err", err)
			return ""
		}
		if !taken[newip.String()] {
			break
		}
	}
	cidr.IP = newip
	return cidr.String()
}
