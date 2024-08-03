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

// ingessCmd represents the ingess command
var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "run ingress test",
	Long: `create an ingress gateway and extclient;
	verify all nodes received update`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("ingress")
		fmt.Println(ingresstest(&config))
	},
}

func init() {
	rootCmd.AddCommand(ingressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ingessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ingessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ingresstest(config *netmaker.Config) bool {
	pass := true
	netclient := netmaker.GetNetclient(config.Network)
	ingress := netmaker.GetHost("ingress", netclient)
	if ingress == nil {
		slog.Error("did not find ingress host/node")
		return false
	}
	slog.Debug("debuging", "ingess", ingress)
	//create ingress
	slog.Info("creating ingress node")
	netmaker.CreateIngress(*ingress)
	//create extclient
	slog.Info("creating extclient")
	ext := netmaker.CreateExtClient(*ingress)
	slog.Info("downloading client config")
	if err := netmaker.DownloadExtClientConfig(*ingress); err != nil {
		slog.Error("failed to download extclient config", "test", "ingress", "err", err)
		return false
	}
	slog.Info("copying file to extclient")
	if err := netmaker.StartExtClient(config); err != nil {
		slog.Error("failed to start extclient", "test", "ingress", "err", err)
		return false
	}
	//verify
	failedmachines := []string{}
	extclient := netmaker.GetExtClient(*ingress, ext)
	ip := extclient.Address
	// wait for update to be propoated
	time.Sleep(time.Second * 30)
	for _, machine := range netclient {
		slog.Info(fmt.Sprintf("checking that %s @ %s received the update", machine.Host.Name, machine.Host.EndpointIP))
		out, err := ssh.Run([]byte(config.Key), machine.Host.EndpointIP, "wg show netmaker allowed-ips ")
		if err != nil {
			slog.Error("err connecting", "machine", machine.Host.Name, "err", err)
			failedmachines = append(failedmachines, machine.Host.Name)
			pass = false
			continue
		}
		if !strings.Contains(out, ip) {
			slog.Error("update not received", "host", machine.Host.Name, "output", out)
			failedmachines = append(failedmachines, machine.Host.Name)
			pass = false
			continue
		}
	}
	if len(failedmachines) > 0 {
		slog.Error("not all machines were updated")
		for _, machine := range failedmachines {
			slog.Error("failures", "machine", machine)
		}
		return false
	}
	slog.Info("all nodes received the ingress ips")
	return pass
}
