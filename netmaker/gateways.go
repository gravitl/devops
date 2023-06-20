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

func GetExtClient(m Netclient) *models.ExtClient {
	return callapi[models.ExtClient](http.MethodGet, "/api/extclients/"+m.Node.Network+"/road-warrior", nil)
}

func CreateExtClient(client Netclient) {
	slog.Info("creating ingress on node", "node", client.Node.ID)
	data := struct {
		Clientid string
	}{
		Clientid: "road-warrior",
	}
	callapi[models.ApiNode](http.MethodPost, "/api/extclients/"+client.Node.Network+"/"+client.Node.ID, data)
}

func DownloadExtClientConfig(client Netclient) error {
	file := download(http.MethodGet, "/api/extclients/"+client.Node.Network+"/road-warrior/file", nil)
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
