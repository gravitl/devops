package cmd

import (
	"fmt"

	"github.com/gravitl/devops/netmaker"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
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
	pass := true
	netclient := netmaker.GetNetclient(config.Network)
	internetGateway := netmaker.GetHost("node-gateway", netclient)
	if internetGateway == nil {
		slog.Error("did not find node-gateway")
		return false
	}

	ingressNode := netmaker.GetHost("node-ingress", netclient)
	if ingressNode == nil {
		slog.Error("did not find node-ingress")
		return false
	}

	slog.Info("found both nodes")

	netmaker.CreateInternetGateway(*internetGateway, *ingressNode)
	slog.Info("internet gateway was created")

	//TODO: do a ping test
	//TODO: ping the internet

	netmaker.DeleteInternetGateway(*internetGateway)

	slog.Info("internet gateway was deleted")
	return pass
}
