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
	"time"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// egressCmd represents the egress command
var egressCmd = &cobra.Command{
	Use:   "egress",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("egress")
		fmt.Println(egresstest(&config))
	},
}

func init() {
	rootCmd.AddCommand(egressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// egressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// egressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func egresstest(config *netmaker.Config) bool {
	pass := true
	netmaker.SetCxt(config.Api, config.Masterkey)
	netclient := netmaker.GetNetclient(config.Network)
	egress := netmaker.GetHost("egress", netclient)
	if egress == nil {
		slog.Error("did not find egress host/node")
		return false
	}
	slog.Debug("debuging", "egress", egress)
	//create egress
	slog.Info("creating egress gateway")
	netmaker.CreateEgress(*egress, config.Ranges)
	//verify
	failedmachines := []string{}
	ip := config.Ranges[0]
	// wait for update to be propogated
	time.Sleep(time.Second * 30)
	for _, machine := range netclient {
		if machine.Host.Name == "egress" {
			continue
		}

		slog.Info(fmt.Sprintf("checking that %s @ %s received the update", machine.Host.Name, machine.Host.EndpointIP))
		out, err := ssh.Run([]byte(config.Key), machine.Host.EndpointIP, "wg show netmaker allowed-ips | grep "+ip)
		if err != nil {
			slog.Error("err connecting to", "host", machine.Host.Name, "err", err)
			failedmachines = append(failedmachines, machine.Host.Name)
			pass = false
			continue
		}
		if !strings.Contains(out, ip) {
			slog.Error("update not received", "machine", machine.Host.Name, "output", out)
			failedmachines = append(failedmachines, machine.Host.Name)
			pass = false
			continue
		}
	}
	if len(failedmachines) > 0 {
		slog.Error("not all machines were updated")
		for _, machine := range failedmachines {
			slog.Error("Failures", "machine", machine)
		}
		return pass
	}
	slog.Info("all nodes received the egress range")
	return pass
}
