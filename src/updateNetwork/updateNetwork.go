package main

import (
	"os"

	"github.com/gravitl/devops/do"
	"github.com/gravitl/devops/logging"
	"github.com/gravitl/devops/netmaker"
	"golang.org/x/exp/slog"
)

func main() {
	//SetUp Logging
	logging.SetupLoging("updateNetwork")
	do_token := os.Getenv("DIGITALOCEAN_TOKEN")
	masterkey, ok := os.LookupEnv("MASTERKEY")
	if !ok {
		masterkey = "secretkey"
	}
	request, _, _, _ := do.Default()
	request.Token = do_token
	//SetUp Servers
	ip, err := request.GetPublicIP("server", "devops")
	if err != nil {
		slog.Error("get public ip for server ", "ERROR", err)
	}
	server := &do.Server{
		FQDN:      "server.clustercat.com",
		Broker:    "broker.clustercat.com",
		API:       "api.clustercat.com",
		Dashboard: "dashboard.clustercat.com",
		PublicIP:  ip,
		Subdomain: "clustercat.com",
		Branch:    "testing",
		UIBranch:  "testing",
	}
	request.SoftResetServer(server)
	netmaker.SetCxt("https://"+server.API, masterkey)
	nodes := netmaker.GetAllNodes()
	slog.Info("checking nodes for gateways")
	for _, node := range *nodes {
		if node.IsEgressGateway {
			netmaker.DeleteEgress(node.ID, node.Network)
		}
		if node.IsRelay {
			netmaker.DeleteRelay(node.ID, node.Network)
		}
		if node.IsIngressGateway {
			netmaker.DeleteIngress(node.ID, node.Network)
		}
	}
	//Update Nodes
	if request.DropletsExist("devops") {
		request.UpdateNodes("devops", "testing")
		//request.StopDocker("docker", "devops-docker")
	} else {
		CreateNodes(request)
	}
	slog.Info("success")
}

func CreateNodes(request *do.Request) {
	//create nomal nodes
	request.Names = append(request.Names, "node1", "relay", "relayed", "egress", "ingress")
	request.Tags = append(request.Tags, "normal")
	if err := request.CreateNodes(); err != nil {
		slog.Error("creating nodes ", "ERROR", err)
		os.Exit(1)
	}
	//create special nodes
	request.Names = []string{"docker", "egressrange", "extclient"}
	request.Tags = []string{"testing", "special"}
	if err := request.CreateNodes(); err != nil {
		slog.Error("creating nodes ", "ERROR", err)
		os.Exit(1)
	}
	slog.Info("wait for nodes to be fully available")
	request.WaitForCloudInit("testing")
	request.VerifyDNS("testing")
	slog.Info("copying netclient to new nodes")
	if err := request.CopyNodeFiles("normal", "develop"); err != nil {
		slog.Error("copying netclient to nodes ", "ERRROR", err)
	}
}
