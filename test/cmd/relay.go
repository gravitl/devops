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
	"log"
	"net"
	"strings"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
)

// relayCmd represents the relay command
var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "run a relay test",
	Long: `creates a relay and
	verifies all other nodes have received the update`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("relay called")
		relaytest(&config)
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

func relaytest(config *netmaker.Config) {
	netmaker.SetCxt(config.Api, config.Masterkey)
	netclient := netmaker.GetNetclient(config.Network)
	relay := netmaker.GetHost("relay", netclient)
	if relay == nil {
		log.Fatal("did not find relay netclient")
	}
	relayed := netmaker.GetHost("relayed", netclient)
	if relayed == nil {
		log.Fatal("did not find relayed netclient")
	}
	log.Println("creating relay")
	netmaker.CreateRelay(relay, relayed)
	log.Println("setting firewall rules")
	egress := netmaker.GetHost("egress", netclient)
	if egress == nil {
		log.Fatal("did not find egress netclient")
	}
	_, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "iptables -A OUTPUT -d "+egress.Host.EndpointIP+" -j DROP")
	if err != nil {
		log.Fatal("failed to set firewall rule on relayed")
	}
	_, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "iptables -A OUTPUT -d "+relayed.Host.EndpointIP+" -j DROP")
	if err != nil {
		log.Fatal("failed to set firewall rule on relayed")
	}
	defer resetFirewall(relayed, egress)
	log.Println("ping egress from relayed")
	ip, _, err := net.ParseCIDR(egress.Node.Address)
	if err != nil {
		log.Fatal("failed to parse egress address", egress.Node.Address)
	}
	out, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "ping -c 3 "+ip.String())
	if err != nil {
		log.Println("error connecting to relayed")
		log.Println(out, err)
	}
	if !strings.Contains(out, "3 received") {
		log.Fatal("failed to ping egress from relayed")
	}
	log.Println("ping relayed from egress")
	ip, _, err = net.ParseCIDR(relayed.Node.Address)
	if err != nil {
		log.Fatal("failed to parse relayed address", relayed.Node.Address)
	}
	out, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "ping -c 3 "+ip.String())
	if err != nil {
		log.Println("error connecting to egress")
		log.Println(out, err)
	}
	if !strings.Contains(out, "3 received") {
		log.Fatal("failed to ping relayed from egress")
	}
}

func resetFirewall(relayed, egress *netmaker.Netclient) {
	log.Println("reseting firewall on relayed/egress")
	_, err := ssh.Run([]byte(config.Key), relayed.Host.EndpointIP, "iptables -D OUTPUT -d "+egress.Host.EndpointIP+" -j DROP")
	if err != nil {
		log.Fatal("failed to set firewall rule on relayed")
	}
	_, err = ssh.Run([]byte(config.Key), egress.Host.EndpointIP, "iptables -D OUTPUT -d "+relayed.Host.EndpointIP+" -j DROP")
	if err != nil {
		log.Fatal("failed to set firewall rule on relayed")
	}
}
