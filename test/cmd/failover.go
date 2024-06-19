package cmd

import (
	"fmt"
	"log/slog"

	"github.com/gravitl/devops/netmaker"
	"github.com/spf13/cobra"
)

var failovercmd = &cobra.Command{
	Use:   "failover",
	Short: "run failover test",
	Long: `create a failover;
	verify all nodes received update`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("failover")
		fmt.Println(failovertest(&config))
	},
}

func init() {
	rootCmd.AddCommand(failovercmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ingessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ingessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func failovertest(config *netmaker.Config) bool {
	pass := true
	netclient := netmaker.GetNetclient(config.Network)
	failover := netmaker.GetHost("server", netclient)
	if failover == nil {
		slog.Error("did not find server host/node")
		return false
	}
	slog.Debug("debugging", "failover", failover)
	slog.Info("Creating failover")
	netmaker.CreateFailover(*failover)
	//TODO test if failover addition successful.
	return pass
}
