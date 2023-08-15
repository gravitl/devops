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
	"golang.org/x/exp/slog"
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
	Server             string
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

func DeleteRelay(id, network string) {
	api := "/api/nodes/" + network + "/" + id + "/deleterelay"
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
	slog.Info("deleted ingress from node ", "node", id)
}

func DeleteEgress(id, network string) {
	api := "/api/nodes/" + network + "/" + id + "/deletegateway"
	callapi[models.ApiNode](http.MethodDelete, api, nil)
	slog.Info("deleted egress from node ", "node", id)
}

func StartExtClient(config *Config) error {
	client, err := do.Name("extclient", config.Tag, config.DigitalOcean_Token)
	if err != nil {
		return fmt.Errorf("failed to get extclient %w", err)
	}
	clientip, err := client.PublicIPv4()
	if err != nil {
		return fmt.Errorf("error retrieving extclient ip address %w", err)
	}
	if err := ssh.CopyTo([]byte(config.Key), clientip, "/tmp/netmaker.conf", "/etc/wireguard"); err != nil {
		return fmt.Errorf("error copying config to extclient %w", err)
	}
	if _, err := ssh.Run([]byte(config.Key), clientip, "wg-quick up netmaker"); err != nil {
		return fmt.Errorf("error starting wireguard on extclient %w", err)
	}
	return nil
}

func RestoreExtClient(config *Config) error {
	client, err := do.Name("extclient", config.Tag, config.DigitalOcean_Token)
	if err != nil {
		return fmt.Errorf("error getting DO client %w", err)
	}
	clientip, err := client.PublicIPv4()
	if err != nil {
		return fmt.Errorf("error retrieving extclient ip address %w", err)
	}
	if out, err := ssh.Run([]byte(config.Key), clientip, "wg-quick down netmaker"); err != nil {
		slog.Warn("error stopping wireguard on extclient", out, err)
	}
	if out, err := ssh.Run([]byte(config.Key), clientip, "rm /etc/wireguard/netmaker.conf"); err != nil {
		slog.Warn("error removing wireguard conf on extclient", out, err)
	}
	return nil
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
	for _, node := range *nodes {
		slog.Info("checking node ", "node", node.ID)
		if node.Network != network {
			continue
		}
		slog.Info("node", "Address", node.Address)
		if node.Address != "" {
			ip, _, err := net.ParseCIDR(node.Address)
			if err != nil {
				slog.Error("error parsing cidr ", "node", node.Address, "err", err)
			} else {
				ips = append(ips, ip)
			}
		}
		slog.Info("node", "Address6", node.Address6)
		if node.Address6 != "" {
			ip, _, err := net.ParseCIDR(node.Address6)
			if err != nil {
				slog.Error("error parsing cidr ", "node", node.Address, "err", err)
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

// func GetToken(net, url string, uses int) (string, error) {
// 	key := models.AccessKey{}
// 	key.Uses = uses
// 	response, err := Api(key, http.MethodPost, url+"/api/networks/"+net+"/keys", "secretkey")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer response.Body.Close()
// 	if response.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("error retriving token: status %s", response.Status)
// 	}
// 	if err := json.NewDecoder(response.Body).Decode(&key); err != nil {
// 		return "", err
// 	}
// 	return key.AccessString, nil
// }

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
	slog.Debug("calling api", slog.String("endpoint", ctx.Endpoint+route))
	if payload == nil {
		req, err = http.NewRequest(method, ctx.Endpoint+route, nil)
		if err != nil {
			slog.Error("Client could not create request:", "err", err)
			return nil
		}
	} else {
		slog.Debug("debugging", "payload", payload)
		payloadBytes, jsonErr := json.Marshal(payload)
		if jsonErr != nil {
			slog.Error("Error in request JSON marshalling:", "err", err)
			return nil
		}
		req, err = http.NewRequest(method, ctx.Endpoint+route, bytes.NewReader(payloadBytes))
		if err != nil {
			slog.Error("Client could not create request:", "err", err)
			return nil
		}
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+ctx.MasterKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Client error making http request:", "err", err)
		return nil
	}
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Client could not read response body:", "err", err)
		return nil
	}
	if res.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Error Status: %d Response: %s", res.StatusCode, string(resBodyBytes)))
		return nil
	}
	body := new(T)
	if len(resBodyBytes) > 0 {
		if err := json.Unmarshal(resBodyBytes, body); err != nil {
			slog.Error("Error unmarshalling JSON:", "err", err)
			return nil
		}
	}
	return body
}

func download(method, route string, payload any) []byte {
	var (
		req *http.Request
		err error
	)
	slog.Debug("calling api", slog.String("endpoint", ctx.Endpoint+route))
	if payload == nil {
		req, err = http.NewRequest(method, ctx.Endpoint+route, nil)
		if err != nil {
			slog.Error("Client could not create request:", "err", err)
			return []byte{}
		}
	} else {
		payloadBytes, jsonErr := json.Marshal(payload)
		if jsonErr != nil {
			slog.Error("Error in request JSON marshalling:", "err", err)
			return []byte{}
		}
		req, err = http.NewRequest(method, ctx.Endpoint+route, bytes.NewReader(payloadBytes))
		if err != nil {
			slog.Error("Client could not create request:", "err", err)
			return []byte{}
		}
		req.Header.Set("Content-Type", "application/json")
	}
	slog.Debug("debuging", "request", req)
	req.Header.Set("Authorization", "Bearer "+ctx.MasterKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Client error making http request:", "err", err)
		return []byte{}
	}
	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Client could not read response body:", "err", err)
		return []byte{}
	}
	if res.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("Error Status: %d Response: %s", res.StatusCode, string(resBodyBytes)))
		return []byte{}
	}
	return resBodyBytes
}

func SetVerbosity(value int) {
	hosts := callapi[[]models.ApiHost](http.MethodGet, "/api/hosts", nil)
	for _, host := range *hosts {
		host.Verbosity = value
		callapi[models.ApiHost](http.MethodPut, "/api/hosts/"+host.ID, host)
	}
}
