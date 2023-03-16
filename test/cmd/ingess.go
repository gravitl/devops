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
	"os"
	"strings"
	"time"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

// ingessCmd represents the ingess command
var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "run ingress test",
	Long: `create an ingress gateway and extclient;
	verify all nodes received update`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ingess called")
		ingresstest(&config)
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

func ingresstest(config *netmaker.Config) {
	netclient := netmaker.GetNetclient(config.Network)
	ingress := netmaker.GetHost("ingress", netclient)
	if ingress == nil {
		log.Fatal("did not find ingress host/node")
	}
	pretty.Println(ingress)
	//create ingress
	log.Println("creating ingress node")
	netmaker.CreateIngress(*ingress)
	//create extclient
	log.Println("creating extclient")
	netmaker.CreateExtClient(*ingress)
	log.Println("downloading client config")
	if err := netmaker.DownloadExtClientConfig(*ingress); err != nil {
		log.Fatal(err)
	}
	log.Println("copying file to extclient")
	netmaker.StartExtClient(config)
	//verify
	failedmachines := []string{}
	extclient := netmaker.GetExtClient(*ingress)
	ip := extclient.Address
	log.Println("waiting for update to propogate")
	time.Sleep(time.Second * 30)
	for _, machine := range netclient {

		log.Printf("checking that %s @ %s received the update", machine.Host.Name, machine.Host.EndpointIP)
		out, err := ssh.Run([]byte(config.Key), machine.Host.EndpointIP, "wg show netmaker allowed-ips | grep "+ip)
		if err != nil {
			log.Printf("err connecting to %s\n", machine.Host.Name)
			log.Println(out, err)
			failedmachines = append(failedmachines, machine.Host.Name)
			continue
		}
		if !strings.Contains(out, ip) {
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
}
