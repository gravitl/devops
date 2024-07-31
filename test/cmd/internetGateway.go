package cmd

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/gravitl/devops/netmaker"
	"github.com/gravitl/devops/ssh"
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

	out, err := ssh.Run(
		[]byte(config.Key),
		ingressNode.Host.EndpointIP,
		"curl -s icanhazip.com && hostname -I | awk '{print $1}'",
	)

	if err != nil {
		slog.Error("ssh failed", ingressNode.Host.Name)
		pass = false
	}

	ips := strings.Split(string(out), "\n")
	parsedIP1 := net.ParseIP(ips[0])
	parsedIP2 := net.ParseIP(ips[1])
	if parsedIP1 == nil || parsedIP2 == nil {
		slog.Error("invalid IP address")
		pass = false
	}

	if parsedIP1.Equal(parsedIP2) {
		slog.Error("internet gateway was not used")
		pass = false
	}
	slog.Info("internet gateway was used")

	out, err = ssh.Run(
		[]byte(config.Key),
		ingressNode.Host.EndpointIP,
		"ping -c 10 1.1.1.1 | grep packet",
	)

	if err != nil {
		slog.Error("ssh failed", ingressNode.Host.Name)
		pass = false
	}

	if !strings.Contains(out, ", 0% packet loss") {
		slog.Error("error connecting to the internet")
		pass = false
	}
	slog.Info("host can reach the internet")

	peers, err := netmaker.GetWireGuardIPs(config.Network)
	if err != nil {
		slog.Error("failed to get peers", ingressNode.Host.Name)
		pass = false
	}
	slog.Info("pinging all the peers")

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 4)

	for _, peer := range peers {
		wg.Add(1)
		go func(peer net.IP) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			pass = pingPeers(ingressNode.Host.EndpointIP, peer)
		}(peer)
	}
	wg.Wait()

	netmaker.DeleteInternetGateway(*internetGateway)
	slog.Info("internet gateway was deleted")

	return pass
}

func pingPeers(source string, ip net.IP) bool {
	out, err := ssh.Run([]byte(config.Key), source, "ping -c 10 "+ip.String()+" | grep packet")
	if err != nil {
		slog.Error("error connecting to peer", "peer", ip.String(), "test", "ping", "err", err)
	}

	if strings.Contains(out, ", 10% packet loss") || strings.Contains(out, ", 20% packet loss") {
		slog.Warn(
			"ping success",
			"peer",
			ip.String(),
			"output",
			strings.TrimSuffix(out, "\n"),
		)
		return true
	}

	if strings.Contains(out, ", 0% packet loss") {
		slog.Info(
			"ping success",
			"peer",
			ip.String(),
			"output",
			strings.TrimSuffix(out, "\n"),
		)
		return true
	}

	return false
}
