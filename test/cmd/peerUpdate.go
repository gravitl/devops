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
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

// peerUpdateCmd represents the peerUpdate command
var peerUpdateCmd = &cobra.Command{
	Use:   "peerUpdate",
	Short: "run peerupdate test",
	Long: `updates wg address of a node and
	verifies that all other nodes received the update
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("peerUpdate called")
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
			log.Fatal("invalid configuration, ssh key not set", err)
		}
		config.Key = string(key)
	}
	fmt.Printf("key is of type %T\n", config.Key)
	netmaker.SetCxt(config.Api, config.Masterkey)
	netclient := netmaker.GetNetclient(config.Network)
	server := netmaker.GetHost("server", netclient)
	if server == nil {
		log.Fatal("did not find server")
	}
	log.Println("updating wg ip on server node")
	taken := make(map[string]bool)
	for _, machine := range netclient {
		taken[machine.Node.Address] = true
	}
	if debug {
		pretty.Println("exclued ips ", taken)
	}
	newip := getNextIP(server.Node.Address, taken)
	server.Node.Address = newip
	log.Println("updating wg address of ", server.Host.Name, " to ", newip)
	netmaker.UpdateNode(&server.Node)
	// check node received update
	//check if other nodes received update
	failedmachines := []string{}
	ip, _, err := net.ParseCIDR(newip)
	if err != nil {
		log.Fatal("could not parse newip", ip, err)
	}
	for _, machine := range netclient {
		time.Sleep(time.Second)
		if machine.Host.Name == "server" {
			continue
		}
		if machine.Host.IsRelayed {
			continue
		}
		log.Printf("checking that %s @ %s received the update", machine.Host.Name, machine.Host.EndpointIP)
		out, err := ssh.Run([]byte(config.Key), machine.Host.EndpointIP, "wg show netmaker allowed-ips | grep "+ip.String())
		if err != nil {
			log.Printf("err connecting to %s\n", machine.Host.Name)
			log.Println(out, err)
			failedmachines = append(failedmachines, machine.Host.Name)
			continue
		}
		if !strings.Contains(out, ip.String()) {
			log.Printf("%s did not receive the update %s\n", machine.Host.Name, out)
			failedmachines = append(failedmachines, machine.Host.Name)
			continue
		}
	}
	if len(failedmachines) > 0 {
		log.Println("not all machines were updated")
		for _, machine := range failedmachines {
			log.Printf("%s ", machine)
		}
		os.Exit(1)
	}
	log.Println("all netclients received the update")
}

func getNextIP(current string, taken map[string]bool) string {
	var newip net.IP
	if len(taken) > 253 {
		log.Fatal("no free ips")
	}
	ip, cidr, err := net.ParseCIDR(current)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("getting free ip")
	net4 := iplib.Net4FromStr(current)
	for {
		newip, err = net4.NextIP(ip)
		if errors.Is(err, iplib.ErrBroadcastAddress) {
			newip, err = net4.NextIP(net4.FirstAddress())
		}
		if err != nil {
			log.Fatal("NextIP", err)
		}
		if !taken[newip.String()] {
			break
		}
	}
	cidr.IP = newip
	return cidr.String()
}
