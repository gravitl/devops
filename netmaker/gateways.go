package netmaker

import (
	"log"
	"net/http"
	"os"

	"github.com/gravitl/netmaker/models"
	"github.com/kr/pretty"
)

func CreateIngress(name Netclient) {
	log.Println("creating ingress on node", name.Node.ID)
	callapi[models.ApiNode](http.MethodPost, "/api/nodes/"+name.Node.Network+"/"+name.Node.ID+"/createingress", nil)
}

func GetExtClient(m Netclient) *models.ExtClient {
	return callapi[models.ExtClient](http.MethodGet, "/api/extclients/"+m.Node.Network+"/road-warrior", nil)
}

func CreateExtClient(client Netclient) {
	log.Println("creating ingress on node", client.Node.ID)
	data := struct {
		Clientid string
	}{
		Clientid: "road-warrior",
	}
	callapi[models.ApiNode](http.MethodPost, "/api/extclients/"+client.Node.Network+"/"+client.Node.ID, data)
}

func DownloadExtClientConfig(client Netclient) error {
	file := download(http.MethodGet, "/api/extclients/"+client.Node.Network+"/road-warrior/file", nil)
	if Debug {
		log.Println("received file \n", string(file))
	}
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
	log.Println("creatting egress on node", client.Node.ID)
	data := models.EgressGatewayRequest{
		NodeID:     client.Node.ID,
		NetID:      client.Node.Network,
		NatEnabled: "yes",
		Ranges:     ranges,
	}
	return callapi[models.ApiNode](http.MethodPost, "/api/nodes/"+data.NetID+"/"+data.NodeID+"/creategateway", data)
}

func CreateRelay(relay, relayed *Netclient) {
	data := models.HostRelayRequest{
		HostID:       relay.Host.ID,
		RelayedHosts: []string{relayed.Host.ID},
	}
	if Debug {
		pretty.Println(data)
	}
	callapi[models.ApiHost](http.MethodPost, "/api/hosts/"+relay.Host.ID+"/relay", data)
}
