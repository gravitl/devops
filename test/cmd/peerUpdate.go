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
	Use:   "peerUpdate -s <server>",
	Short: "run peerupdate test",
	Long: `updates wg address of a node and
	verifies that all other nodes received the update
`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("peerupdate")
		config.Server = cmd.Flag("server").Value.String()
		fmt.Println(peerupdatetest(&config))
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
	peerUpdateCmd.Flags().StringP("server", "s", "server", "server name")
}

func peerupdatetest(config *netmaker.Config) bool {
	pass := true
	if config.Key == "" {
		key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_ed25519")
		if err != nil {
			slog.Error("invalid configuration, ssh key not set", "test", "peerupdate", "err", err)
			return false
		}
		config.Key = string(key)
	}
	netmaker.SetCxt(config.Api, config.Masterkey)
	netclient := netmaker.GetNetclient(config.Network)
	slog.Info("getting server", "server-name", config.Server)
	server := netmaker.GetHost(config.Server, netclient)
	if server == nil {
		slog.Error("did not find server", "test", "peerupdate")
		return false
	}
	slog.Info("updating wg ip on server")
	taken := make(map[string]bool)
	var addressToUse string
	for _, machine := range netclient {
		if config.Network == "devopsv6" {
			addressToUse = machine.Node.Address6
		} else {
			addressToUse = machine.Node.Address
		}
		ip, _, err := net.ParseCIDR(addressToUse)
		if err != nil {
			slog.Warn(fmt.Sprintf("%s is not a cidr", addressToUse))
		}
		taken[ip.String()] = true
	}
	slog.Debug("debugging", "exclued ips ", taken)
	var serverAddressToUse string
	if server.Node.Network == "devopsv6" {
		serverAddressToUse = server.Node.Address6
	} else {
		serverAddressToUse = server.Node.Address
	}

	newip := getNextIP(serverAddressToUse, taken)
	if newip == "" {
		return false
	}
	if server.Node.Network == "devopsv6" {
		server.Node.Address6 = newip
	} else {
		server.Node.Address = newip
	}
	slog.Info(fmt.Sprintf("updating wg address of %s to %s", server.Host.Name, newip))

	netmaker.UpdateNode(&server.Node)
	//verify that server received update
	time.Sleep(time.Second * 5)
	slog.Info("checking that server received the update")
	out, err := ssh.Run([]byte(config.Key), server.Host.EndpointIP, "ip a show netmaker")
	if err != nil {
		slog.Error("ssh connect err", "machine", server.Host.Name, "cmd", "ssh "+server.Host.EndpointIP+" ip -a show netmaker", "err", err)
		return false
	}
	if !strings.Contains(out, newip) {
		slog.Error("server did not receive the update", "machine", server.Host.Name, "ouput", out)
		return false
	}
	// check node received update
	//check if other nodes received update
	failedmachines := []string{}
	ip, _, err := net.ParseCIDR(newip)
	if err != nil {
		slog.Error("could not parse newip", "test", "peerupdate", "ip", ip, "err", err)
		return false
	}
	// wait for update to be propogated
	time.Sleep(time.Second * 30)
	for _, machine := range netclient {
		if machine.Host.Name == config.Server {
			continue
		}
		if machine.Node.IsRelayed {
			continue
		}
		slog.Info(fmt.Sprintf("checking that %s @ %s received the update", machine.Host.Name, machine.Host.EndpointIP))
		out, err := ssh.Run([]byte(config.Key), machine.Host.EndpointIP, "wg show netmaker allowed-ips")
		if err != nil {
			slog.Error("ssh connect err", "machine", machine.Host.Name, "cmd", "ssh "+machine.Host.EndpointIP+" wg show netmaker allowed-ips ", "err", err)
			failedmachines = append(failedmachines, machine.Host.Name)
			pass = false
			continue
		}
		if !strings.Contains(out, ip.String()) {
			slog.Error("node did not receive the update", "machine", machine.Host.Name, "ouput", out)
			failedmachines = append(failedmachines, machine.Host.Name)
			pass = false
			continue
		}
	}
	if len(failedmachines) > 0 {
		slog.Error("not all machines were updated")
		for _, machine := range failedmachines {
			slog.Error(machine)
		}
		return false
	}
	slog.Info("all netclients received the update")
	return pass
}

func getNextIP(current string, taken map[string]bool) string {
	var newip net.IP
	if len(taken) > 253 {
		slog.Error("no free ips")
		return ""
	}
	ip, cidr, err := net.ParseCIDR(current)
	if err != nil {
		slog.Error("failed to parse cidr", "err", err)
		return ""
	}
	slog.Info("getting free ip")
	net4 := iplib.Net4FromStr(current)
	newip, err = net4.NextIP(ip)
	for {
		if errors.Is(err, iplib.ErrBroadcastAddress) {
			newip, err = net4.NextIP(net4.FirstAddress())
		}
		if err != nil {
			slog.Error("NextIP", "err", err)
			return ""
		}
		if !taken[newip.String()] {
			break
		}
		newip, err = net4.NextIP(newip)
	}
	cidr.IP = newip
	return cidr.String()
}
