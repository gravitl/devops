/*
Copyright © 2023 Netmaker Team <info@netmaker.io>

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

	"github.com/gravitl/devops/do"
	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
)

// getapiCmd represents the getapi command
var getapiCmd = &cobra.Command{
	Use:   "getapi -t <do tag>",
	Short: "get api endpoint of netmaker server",
	Long:  `get api endpoint of netmaker server with given tag`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getAPI(&config))
	},
}

func init() {
	rootCmd.AddCommand(getapiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getAPI(config *netmaker.Config) string {
	client, err := do.Name("server", config.Tag, config.DigitalOcean_Token)
	cobra.CheckErr(err)
	serverip, err := client.PublicIPv4()
	cobra.CheckErr(err)
	out, err := ssh.Run([]byte(config.Key), serverip, "grep NM_DOMAIN netmaker.env")
	cobra.CheckErr(err)
	if out == "" {
		cobra.CheckErr("api is blank")
	}
	parts := strings.Split(out, "=")
	temp := strings.ReplaceAll(parts[1], "\"", "")
	result := strings.TrimSpace(temp)
	return fmt.Sprintf("https://api.%s", result)
}
