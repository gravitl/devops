package cmd

import (
	"fmt"

	"github.com/gravitl/devops/netmaker"
	"github.com/spf13/cobra"
)

var internetGatewayCmd = &cobra.Command{
	Use:   "internetGateway",
	Short: "run a internet gateway test",
	Long:  "create a internet gateway and report result",
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("internetGateway")
		fmt.Println(internetGateway(&config))
	},
}

func init() {
	rootCmd.AddCommand(internetGatewayCmd)
}

func internetGateway(config *netmaker.Config) bool {
	fmt.Println("\nTesting InternetGateway\n")
	fmt.Println(config)
	//TODO: setup a internet gateway
	//TODO: do a ping test
	//TODO: ping the internet
	//TODO: delete the gateway
	return false
}
