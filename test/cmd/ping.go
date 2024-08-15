package cmd

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "run a ping test",
	Long:  `ping all nodes on network and reports result`,
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("ping")
		fmt.Println(runPingTest(&config))
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}

func runPingTest(config *netmaker.Config) bool {
	netclients := netmaker.GetNetclient(config.Network)
	hostMap := createHostMap(netclients)
	destinations, err := netmaker.GetWireGuardIPs(config.Network)
	if err != nil {
		slog.Error(
			"unable to get wireguard IP for network",
			"network",
			config.Network,
			"test",
			"ping",
			"err",
			err,
		)
		return false
	}

	results := pingAllHosts(config, netclients, hostMap, destinations)
	return analyzeResults(results)
}

func createHostMap(netclients []netmaker.Netclient) map[string]string {
	hostMap := make(map[string]string)
	for _, client := range netclients {
		addToHostMap(hostMap, client.Node.Address, client.Host.Name)
		addToHostMap(hostMap, client.Node.Address6, client.Host.Name+"6")
	}
	return hostMap
}

func addToHostMap(hostMap map[string]string, address, name string) {
	ip, _, err := net.ParseCIDR(address)
	if err == nil {
		hostMap[ip.String()] = name
	}
}

func pingAllHosts(
	config *netmaker.Config,
	netclients []netmaker.Netclient,
	hostMap map[string]string,
	destinations []net.IP,
) map[string]map[string]bool {
	results := make(map[string]map[string]bool)
	var wg sync.WaitGroup

	for _, client := range netclients {
		wg.Add(1)
		go func(client netmaker.Netclient) {
			defer wg.Done()
			pingFromHost(config, client, hostMap, destinations, results)
		}(client)
	}

	wg.Wait()
	return results
}

func pingFromHost(
	config *netmaker.Config,
	client netmaker.Netclient,
	hostMap map[string]string,
	destinations []net.IP,
	results map[string]map[string]bool,
) {
	sourceName := client.Host.Name
	sourceIP := client.Host.EndpointIP
	slog.Info("pinging from", "host", sourceName, "ip", sourceIP)

	results[sourceName] = make(map[string]bool)

	for _, destIP := range destinations {
		destName := hostMap[destIP.String()]
		if destName == sourceName {
			continue // Skip self
		}

		success := pingWithRetry([]byte(config.Key), sourceIP, destIP.String(), 3)
		results[sourceName][destName] = success

		if !success {
			slog.Error("failed to ping", "source", sourceName, "destination", destName)
		}
	}
}

func pingWithRetry(sshKey []byte, sourceIP, destIP string, maxRetries int) bool {
	for i := 0; i < maxRetries; i++ {
		out, err := ssh.Run(
			sshKey,
			sourceIP,
			fmt.Sprintf("ping -c 10 -W 2 %s | grep packet", destIP),
		)
		if err != nil {
			slog.Error("error connecting to host", "ip", sourceIP, "err", err)
			time.Sleep(time.Second) // Wait before retrying
			continue
		}

		if strings.Contains(out, ", 0% packet loss") {
			slog.Info(
				"ping success",
				"source",
				sourceIP,
				"destination",
				destIP,
				"output",
				strings.TrimSpace(out),
			)
			return true
		}

		slog.Warn(
			"packet loss",
			"source",
			sourceIP,
			"destination",
			destIP,
			"output",
			strings.TrimSpace(out),
		)
		time.Sleep(time.Second) // Wait before retrying
	}

	return false
}

func analyzeResults(results map[string]map[string]bool) bool {
	allSuccessful := true
	for source, destinations := range results {
		for dest, success := range destinations {
			if !success {
				slog.Error("ping failure", "source", source, "destination", dest)
				allSuccessful = false
			}
		}
	}

	if allSuccessful {
		slog.Info("all nodes can ping each other")
	}
	return allSuccessful
}
