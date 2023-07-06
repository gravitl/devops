package netmaker

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gravitl/netmaker/logic"
	"github.com/gravitl/netmaker/models"
	"github.com/kr/pretty"
)

type Netclient struct {
	Host models.ApiHost
	Node models.ApiNode
}

func GetNetclient(network string) []Netclient {
	response := []Netclient{}
	nodes := callapi[[]models.ApiNode](http.MethodGet, "/api/nodes", nil)
	hosts := callapi[[]models.ApiHost](http.MethodGet, "/api/hosts", nil)
	for _, node := range *nodes {
		net := Netclient{}
		if node.Network == network {
			net.Node = node
			host := GetHostByID(node.HostID, hosts)
			net.Host = *host
			response = append(response, net)
		}
	}
	return response
}

func GetHost(hostname string, nets []Netclient) *Netclient {
	for _, net := range nets {
		if net.Host.Name == hostname {
			return &net
		}
	}
	return nil
}

func FindNode(name string) *models.ApiNode {
	nodes := callapi[[]models.ApiNode](http.MethodGet, "/api/nodes", nil)
	for _, node := range *nodes {
		hosts := callapi[models.ApiHost](http.MethodGet, "/api/hosts/"+node.HostID, nil)
		if hosts.Name == name {
			return &node
		}
	}
	return nil
}

func GetNode(network string, ids []string) *models.ApiNode {
	for _, id := range ids {
		log.Println("checking node ", id)
		node := callapi[models.ApiNode](http.MethodGet, "/api/nodes/"+network+"/"+id, nil)
		pretty.Println(node)
		log.Println("checking network ", node.Network, network)
		if node.Network == network {
			log.Println("found node")
			return node
		}
		log.Println("did not find node - wrong network ", node.Network)
	}
	log.Println("did not find node -- out of options")
	return nil
}

func UpdateNode(node *models.ApiNode) {
	callapi[models.ApiNode](http.MethodPut, fmt.Sprintf("/api/nodes/%s/%s", node.Network, node.ID), node)
}

func UpdateNodeWGAddress(name, network, ip string) error {
	found := false
	node := models.ApiNode{}
	nodes := GetNetworkNodes(network)
	if nodes != nil {
		return errors.New("no nodes found")
	}
	for _, node = range *nodes {
		host, err := logic.GetHost(node.HostID)
		if err != nil {
			return err
		}
		if name == host.Name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("node %s was not found in network %s", name, network)
	}
	node.Address = ip
	fmt.Println("updating WG IP", node.Address, network, node.ID)
	callapi[models.Node](http.MethodPut, fmt.Sprintf("/api/nodes/%s/%s", network, node.ID), node)
	return nil
}

func GetAllNodes() *[]models.ApiNode {
	return callapi[[]models.ApiNode](http.MethodGet, "/api/nodes", nil)
}
