package netmaker

import (
	"net/http"
	"os"

	"github.com/gravitl/netmaker/models"
	"golang.org/x/exp/slog"
)

func CreateIngress(name Netclient) {
	slog.Info("creating ingress on node", "node", name.Node.ID)
	callapi[models.ApiNode](http.MethodPost, "/api/nodes/"+name.Node.Network+"/"+name.Node.ID+"/createingress", nil)
}

func GetExtClient(m Netclient, ext string) *models.ExtClient {
	return callapi[models.ExtClient](http.MethodGet, "/api/extclients/"+m.Node.Network+"/"+ext, nil)
}

func ChangeClient(m Netclient, key string, value int) {
	slog.Info("changing", key, "to", value)
	url := "/api/hosts/" + m.Host.ID
	data := struct {
		mtu int
	}{
		mtu: value,
	}
	callapi[models.ApiHost](http.MethodPut, url, data)
}

// func changeNode(m Netclient, key string, value string) {
// 	slog.Info("changing", key, "to", value)
// 	url := "api/nodes/" + m.Node.Network + "/" + m.Node.ID
// 	data := struct {
// 		key string
// 	}{
// 		key: value,
// 	}
// 	callapi[models.ApiNode](http.MethodPut, url, data)
// }

func CreateExtClient(client Netclient, network string) string {
	slog.Info("creating ingress on node", "node", client.Node.ID)
	clients := map[string]string{
		"devops":   "road-warrior",
		"devops4":  "road-warrior2",
		"devopsv6": "road-warrior3",
		"netmaker": "road-warrior4",
	}
	slog.Info("creating ext client", clients[network], network)
	clientID, exists := clients[network]
	if !exists {
		slog.Error("No client ID found for network", network)
		return ""
	}

	data := struct {
		Clientid string `json:"clientid"`
	}{
		Clientid: clientID,
	}

	callapi[models.ApiNode](http.MethodPost, "/api/extclients/"+client.Node.Network+"/"+client.Node.ID, data)

	slog.Info("Successfully created client '%s'\n", clientID)
	return clientID
}

func DownloadExtClientConfig(client Netclient, ext string) error {
	slog.Info("downloading config for", client.Node.Network, ext)
	file := download(http.MethodGet, "/api/extclients/"+client.Node.Network+"/"+ext+"/file", nil)
	slog.Debug("received file", "file", string(file))
	save, err := os.Create("/tmp/netmaker.conf")
	if err != nil {
		return err
	}
	if _, err := save.Write(file); err != nil {
		return err
	}
	return nil
}

func CreateEgress(client Netclient, ranges []string) *models.ApiNode {
	slog.Info("creatting egress", "node", client.Node.ID)
	data := models.EgressGatewayRequest{
		NodeID:     client.Node.ID,
		NetID:      client.Node.Network,
		NatEnabled: "yes",
		Ranges:     ranges,
	}
	return callapi[models.ApiNode](http.MethodPost, "/api/nodes/"+data.NetID+"/"+data.NodeID+"/creategateway", data)
}

func CreateRelay(relay, relayed *Netclient) {
	data := models.RelayRequest{
		NodeID:       relay.Node.ID,
		NetID:        relay.Node.Network,
		RelayedNodes: []string{relayed.Node.ID},
	}
	slog.Debug("debuging", "data", data)
	callapi[models.ApiHost](http.MethodPost, "/api/nodes/"+relay.Node.Network+"/"+relay.Node.ID+"/createrelay", data)
}
