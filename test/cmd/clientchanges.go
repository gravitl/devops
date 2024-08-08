package cmd

import (
	"fmt"
	"time"

	"github.com/gravitl/devops/netmaker"
	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/slog"
)

// ingessCmd represents the ingess command
var clientChangesCmd = &cobra.Command{
	Use:   "clientChanges",
	Short: "run client changes test",
	Long: `make changes to a client;
	verify all nodes received update`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("clientChanges")
		fmt.Println(clientchangestest(&config))
	},
}

func init() {
	rootCmd.AddCommand(clientChangesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ingessCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ingessCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clientchangestest(config *netmaker.Config) bool {
	pass := true
	netclient := netmaker.GetNetclient(config.Network)
	rand.Seed(uint64(time.Now().UnixNano()))
	changer := netclient[rand.Intn(len(netclient))]
	slog.Info("making changes to", changer.Host.Name)
	mtu := rand.Intn(221) + 1280 // Generates a random MTU between 1280 and 1500
	slog.Info("changing mtu to", mtu, "on", changer.Host.Name)
	changer.Host.MTU = mtu
	netmaker.UpdateHost(&changer.Host)
	return pass
}
