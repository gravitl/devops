package netmaker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gravitl/devops/do"
	"github.com/gravitl/devops/ssh"
	"github.com/gravitl/netmaker/models"
	"github.com/kr/pretty"
)

var Debug bool
var baseURL string = "https://api.clustercat.com/"
var ctx Endpoint

func SetBaseURL(url string) {
	baseURL = url
}

type Endpoint struct {
	Endpoint  string
	MasterKey string
}
type EnvMap map[string]string

type Config struct {
	Network            string
	CleanUp            bool
	CleanUpTimeout     int
	DigitalOcean_Token string
	Tag                string
	Api                string
	Masterkey          string
	Ranges             []string
	Key                string
	Debug              bool
}

func SetCxt(endpoint, masterkey string) {
	ctx.Endpoint = endpoint
	ctx.MasterKey = masterkey
}

func DeleteRelay(id string) {
	api := "/api/hosts/" + id + "/relay"
	callapi[models.ApiHost](http.MethodDelete, api, nil)
}

//	_, err := Api("", http.MethodDelete, baseURL+api, "secretkey")
//	if err != nil {
//		log.Println("err deleting relay on node ", host.ID, err)
//	}
//	log.Println("deleted relay from node ", host.ID)
//}

func DeleteIngress(id, network string) {
	api := "/api/nodes/" + network + "/" + id + "/deleteingress"
	callapi[models.ApiNode](http.MethodDelete, api, nil)
	log.Println("deleted ingress from node ", id)
}

func DeleteEgress(id, network string) {
	api := "/api/nodes/" + network + "/" + id + "/deletegateway"
	callapi[models.ApiNode](http.MethodDelete, api, nil)
	log.Println("deleted egress from node ", id)
}

func StartExtClient(config *Config) {
	client, err := do.Name("extclient", config.DigitalOcean_Token)
	if err != nil {
		log.Fatal(err)
	}
	clientip, err := client.PublicIPv4()
	if err != nil {
		log.Fatal("error retrieving extclient ip address")
	}
	if err := ssh.CopyTo([]byte(config.Key), clientip, "/tmp/netmaker.conf", "/etc/wireguard"); err != nil {
		log.Fatal("error copying config to extclient", err)
	}
	if _, err := ssh.Run([]byte(config.Key), clientip, "wg-quick up netmaker"); err != nil {
		log.Fatal("error starting wireguard on extclient", err)
	}
}

func RestoreExtClient(config *Config) {
	client, err := do.Name("extclient", config.DigitalOcean_Token)
	if err != nil {
		log.Fatal(err)
	}
	clientip, err := client.PublicIPv4()
	if err != nil {
		log.Fatal("error retrieving extclient ip address")
	}
	if out, err := ssh.Run([]byte(config.Key), clientip, "wg-quick down netmaker"); err != nil {
		log.Println("error stopping wireguard on extclient", out, err)
	}
	if out, err := ssh.Run([]byte(config.Key), clientip, "rm /etc/wireguard/netmaker.conf"); err != nil {
		log.Println("error removing wireguard conf on extclient", out, err)
	}

}

func GetNetworkNodes(network string) *[]models.ApiNode {
	return callapi[[]models.ApiNode](http.MethodGet, "/api/nodes/"+network, "")
}

func GetWireGuardIPs(network string) ([]net.IP, error) {
	ips := []net.IP{}
	nodes := GetNetworkNodes(network)
	if nodes == nil {
		return ips, errors.New("failled to retrieve network nodes")
	}
	for _, apinode := range *nodes {
		log.Println("network nodes", apinode.ID)
	}
	for _, node := range *nodes {
		log.Println("checking node ", node.ID)
		if node.Network != network {
			continue
		}
		log.Println(node.Address)
		if node.Address != "" {
			ip, _, err := net.ParseCIDR(node.Address)
			if err != nil {
				log.Println("error parsing cidr ", node.Address, err)
			} else {
				ips = append(ips, ip)
			}
		}
		log.Println(node.Address6)
		if node.Address6 != "" {
			ip, _, err := net.ParseCIDR(node.Address6)
			if err != nil {
				log.Println("error parsing cidr ", node.Address, err)
			} else {
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}

func AddAdmin(url string) error {
	admin := models.UserAuthParams{}
	admin.UserName = "admin"
	admin.Password = "password"
	resp, err := Api(admin, http.MethodPost, url+"/api/users/adm/createadmin", "")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadGateway {
			log.Println("Bad Gateway error ... waiting and retrying")
			time.Sleep(time.Second * 10)
			if err := AddAdmin(url); err != nil {
				return errors.New("second error after Gateway error" + err.Error())
			}
		}
		return errors.New(resp.Status)
	}
	return nil
}

func GetToken(net, url string, uses int) (string, error) {
	key := models.AccessKey{}
	key.Uses = uses
	response, err := Api(key, http.MethodPost, url+"/api/networks/"+net+"/keys", "secretkey")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error retriving token: status %s", response.Status)
	}
	if err := json.NewDecoder(response.Body).Decode(&key); err != nil {
		return "", err
	}
	return key.AccessString, nil
}

// func CreateRelay(name, network, relayrange string) error {
// 	var response *http.Response
// 	node, err := FindNode(name, network)
// 	if err != nil {
// 		return fmt.Errorf("failed to get nodes: %w", err)
// 	}
// 	relay := models.RelayRequest{}
// 	relay.NetID = network
// 	relay.NodeID = node.ID
// 	relay.RelayAddrs = append(relay.RelayAddrs, relayrange)
// 	if response, err = Api(relay, http.MethodPost, baseURL+"api/nodes/"+network+"/"+node.ID+"/createrelay", "secretkey"); err != nil {
// 		return fmt.Errorf("api call failed %w", err)
// 	}
// 	if response.StatusCode != http.StatusOK {
// 		return fmt.Errorf("recieved unexpected status code from api: %s", response.Status)
// 	}
// 	return nil
// }

//func CreateEgress(name, network, ranges string) error {
//	egress := models.EgressGatewayRequest{}
//	var response *http.Response
//	node, err := FindNode(name, network)
//	if err != nil {
//		return err
//	}
//	egress.Interface = "eth1"
//	egress.Ranges = append(egress.Ranges, ranges)
//	if response, err = Api(egress, http.MethodPost, baseURL+"api/nodes/"+network+"/"+node.ID+"/creategateway", "secretkey"); err != nil {
//		return fmt.Errorf("api call failed %w", err)
//	}
//	if response.StatusCode != http.StatusOK {
//		return fmt.Errorf("recieved unexpected status code from api: %s", response.Status)
//	}
//	return nil
//}

func FindExtClient(id, network string) (*models.ExtClient, error) {
	extclient := &models.ExtClient{}
	response, err := Api("", http.MethodGet, baseURL+"api/extclients/"+network, "secretkey")
	if err != nil {
		return &models.ExtClient{}, fmt.Errorf("api err: %w", err)
	}
	var extclients []models.ExtClient
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&extclients)
	if err != nil {
		return &models.ExtClient{}, fmt.Errorf("json decoding err in FindExtClient: %w", err)
	}
	found := false
	for _, *extclient = range extclients {
		log.Println("checking extclient ", extclient.ClientID)
		if extclient.IngressGatewayID == id {
			found = true
			log.Println("found ", extclient.ClientID)
			break
		}
	}
	if !found {
		return extclient, fmt.Errorf("extclient not found in network %s", network)
	}
	log.Println("returning ", extclient.ClientID, extclient.IngressGatewayID)
	return extclient, nil
}

func Api(data interface{}, method, url, authorization string) (*http.Response, error) {
	var request *http.Request
	var response *http.Response
	var err error
	if data != "" {
		payload, err := json.Marshal(data)
		if err != nil {
			return response, err
		}
		request, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
		if err != nil {
			return response, err
		}
		request.Header.Set("Content-Type", "application/json")
	} else {
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			return response, err
		}
	}
	if authorization != "" {
		request.Header.Set("authorization", "Bearer "+authorization)
	}
	client := http.Client{}
	client.Timeout = time.Second * 10
	return client.Do(request)
}

func callapi[T any](method, route string, payload any) *T {
	var (
		req *http.Request
		err error
	)
	log.Println("calling api", ctx.Endpoint+route)
	if payload == nil {
		req, err = http.NewRequest(method, ctx.Endpoint+route, nil)
		if err != nil {
			log.Fatalf("Client could not create request: %s", err)
		}
	} else {
		if Debug {
			pretty.Println(payload)
		}
		payloadBytes, jsonErr := json.Marshal(payload)
		if jsonErr != nil {
			log.Fatalf("Error in request JSON marshalling: %s", err)
		}
		req, err = http.NewRequest(method, ctx.Endpoint+route, bytes.NewReader(payloadBytes))
		if err != nil {
			log.Fatalf("Client could not create request: %s", err)
		}
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+ctx.MasterKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Client error making http request: %s", err)
	}
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Client could not read response body: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Error Status: %d Response: %s", res.StatusCode, string(resBodyBytes))
	}
	body := new(T)
	if len(resBodyBytes) > 0 {
		if err := json.Unmarshal(resBodyBytes, body); err != nil {
			log.Fatalf("Error unmarshalling JSON: %s", err)
		}
	}
	return body
}

func download(method, route string, payload any) []byte {
	var (
		req *http.Request
		err error
	)
	log.Println("calling api", ctx.Endpoint+route)
	if payload == nil {
		req, err = http.NewRequest(method, ctx.Endpoint+route, nil)
		if err != nil {
			log.Fatalf("Client could not create request: %s", err)
		}
	} else {
		payloadBytes, jsonErr := json.Marshal(payload)
		if jsonErr != nil {
			log.Fatalf("Error in request JSON marshalling: %s", err)
		}
		req, err = http.NewRequest(method, ctx.Endpoint+route, bytes.NewReader(payloadBytes))
		if err != nil {
			log.Fatalf("Client could not create request: %s", err)
		}
		req.Header.Set("Content-Type", "application/json")
	}
	if Debug {
		pretty.Println(req)
	}
	req.Header.Set("Authorization", "Bearer "+ctx.MasterKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Client error making http request: %s", err)
	}
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Client could not read response body: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Error Status: %d Response: %s", res.StatusCode, string(resBodyBytes))
	}
	return resBodyBytes
}
