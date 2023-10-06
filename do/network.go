// Package do provides an abstraction layer for the creation, deletion
//of droplets and dns records on Digital Ocean.  It also provides functions
//to configure servers and nodes for testing purposes.

package do

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/digitalocean/godo"
	"github.com/gravitl/netclient/ncutils"
	"golang.org/x/exp/slog"
)

var ctx context.Context
var opts *godo.ListOptions

// var ssh *easyssh.MakeConfig
var ssh *SshConf

// Request - contains data for building a set of nodes and dns records on digital ocean
type Request struct {
	Token        string
	Names        []string
	Region       string
	Distribution string
	SubDomain    string
	Image        string
	UserData     string
	Tags         []string
}

// SSH - contains data for ssh/scp connections
type SshConf struct {
	User    string
	Server  string
	Key     string
	Options string
}

// init -- initialize global structures and vars
func init() {
	ctx = context.TODO()
	opts = &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	user := os.Getenv("USER")
	ssh = &SshConf{
		User:    "root",
		Key:     "/home/" + user + "/.ssh/id_devops",
		Options: "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=15",
	}
}

// Default - returns a DO network creation request with sane values
func Default() (*Request, context.Context, *godo.ListOptions, *SshConf) {
	//upgrade takes too long so it is commented out
	//docker-compose is overkill for nodes but easier to have single
	//configuration
	userdata := `#!bin/bash
	apt-get update
	#apt-get upgrade -y
	apt-get install wireguard-tools docker-compose 
	`
	request := Request{
		Region:       "nyc3",
		Distribution: "ubuntu-20-04-x64",
		//ServerName:   "server",
		Tags: []string{
			"testing",
		},
		UserData: userdata,
	}
	return &request, ctx, opts, ssh
}

// CreateNodes- creates a set of droplets and dns records on digital ocean as specified in the request
// returns public ip of created server
func (request *Request) CreateNodes(tags ...string) error {
	request.Tags = append(request.Tags, tags...)
	client := godo.NewFromToken(request.Token)
	//get all ssh keys
	keys, _, err := client.Keys.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to retrieve ssh keys %w", err)
	}
	var keyIDs []godo.DropletCreateSSHKey
	var keyID godo.DropletCreateSSHKey
	for _, key := range keys {
		keyID.Fingerprint = key.Fingerprint
		keyIDs = append(keyIDs, keyID)
	}
	log.Println("creating droplets")
	createRequest := &godo.DropletMultiCreateRequest{
		Names:  request.Names,
		Region: request.Region,
		Size:   "s-1vpcu-1gb",
		Image: godo.DropletCreateImage{
			Slug: request.Distribution,
		},
		SSHKeys:    keyIDs,
		Monitoring: true,
		Tags:       request.Tags,
		UserData:   request.UserData,
	}
	droplets, response, err := client.Droplets.CreateMultiple(ctx, createRequest)
	if err != nil {
		return fmt.Errorf("failure creating droplet %w", err)
	}
	//wait for droplet creation to complete
	wg := sync.WaitGroup{}
	for _, action := range response.Links.Actions {
		wg.Add(1)
		go func(action godo.LinkAction) {
			defer wg.Done()
			for {
				status, _, err := client.Actions.Get(ctx, action.ID)
				if err != nil {
					continue
				}
				if status.Status == "completed" {
					break
				}
				log.Print("waiting for droplet creation to complete ... ", status.Status, " ", status.ID)
				time.Sleep(time.Second << 2)
			}
		}(action)
	}
	wg.Wait()
	//assign dns for nodes
	// for some reason the droplet returned from CreateMultiple doesn't contain networking info so need to retrieve data again
	log.Println("creating dns entries for nodes")
	domainRequest := &godo.DomainRecordEditRequest{
		Type: "A",
		TTL:  300,
	}
	for _, droplet_old := range droplets {
		droplet, _, err := client.Droplets.Get(ctx, droplet_old.ID)
		if err != nil {
			log.Println("failed to retrive droplet data ", err)
			continue
		}
		publicIP, err := droplet.PublicIPv4()
		if err != nil || publicIP == "" {
			log.Println("failed to retrieve droplet ip", err)
		}
		if request.SubDomain == "" {
			domainRequest.Name = droplet.Name
		} else {
			domainRequest.Name = droplet.Name + "." + request.SubDomain
		}
		domainRequest.Data = publicIP
		_, _, err = client.Domains.CreateRecord(ctx, "clustercat.com", domainRequest)
		if err != nil {
			log.Println("failed to create dns record for ", droplet.Name, err)
		}
	}
	log.Println("done, droplets & dns records created")
	return nil
}

func (request *Request) DeleteDroplets(tag string) error {
	client := godo.NewFromToken(request.Token)
	//delete droplets
	log.Println("deleting droplets")
	_, err := client.Droplets.DeleteByTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("error deleting droplets %w", err)
	}
	return nil
}

// WaitForCloudInit - queries droplets for status of cloud-init completion
// returns when cloud-init is complete
func (request *Request) WaitForCloudInit(tag string) {
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		log.Fatal("could not retrieve droplets ", err)
	}
	for _, droplet := range droplets {
		ssh.Server, err = droplet.PublicIPv4()
		if err != nil {
			log.Println("droplet ", droplet.Name, " does not have public ip ... skipping ", err)
			continue
		}
		fmt.Println("checking cloud-init", droplet.Name, ssh.Server)
		done := false
		for {
			if done {
				break
			}
			time.Sleep(time.Second * 10)
			out, err := ssh.Run("cloud-init status")
			if err != nil {
				log.Println("ssh error connecting to node ", droplet.Name, " ", err)
				continue
			}
			if strings.Contains(out, "status: done") {
				done = true
			}
			log.Print(out)
		}
	}
}

func (request *Request) DropletsExist(tag string) bool {
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		log.Println("error retrieving droplets", err)
		return false
	}
	return len(droplets) > 0
}

// StartDocker - starts dockers on server
func (request *Request) StartDocker(server *Server) error {
	log.Println("starting docker on server", server.FQDN)
	ssh.Server = server.FQDN
	_, err := ssh.Run("docker-compose up -d")
	if err != nil {
		return errors.New("error starting docker " + err.Error())
	}
	return nil
}

// CopyNodeFiles -- copies netclient binary from fileserver to nodes
func (request *Request) CopyNodeFiles(tag, branch string) error {
	//check that desired netclient version is available
	ssh.Server = "fileserver.clustercat.com"
	out, err := ssh.Run("ls /var/www/files")
	if !strings.Contains(out, branch) || err != nil {
		return errors.New(branch + " does not exist on fileserver")
	}
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		return err
	}
	for _, droplet := range droplets {
		ssh.Server, err = droplet.PublicIPv4()
		if err != nil {
			log.Println("unable to get public ip for droplet ", droplet.Name, err)
			continue
		}
		log.Println("copying netclient to ", droplet.Name)
		cmd := fmt.Sprintf("wget -O netclient https://fileserver.clustercat.com/%s/netclient", branch)
		cmd += "; chmod +x netclient"
		//needed (on ubuntu at least) to be able to view logs
		cmd += "; systemctl restart systemd-journald"
		cmd += "; ./netclient install"
		_, err = ssh.Run(cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// JoinNetwork -- cycles through nodes and joins network
func (request *Request) JoinNetwork(tag, token string) (bool, []string) {
	success := true
	failednodes := []string{}
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		//return error message in failed nodes array -- kludgy
		failednodes = append(failednodes, err.Error())
		return false, failednodes
	}
	for _, droplet := range droplets {
		ssh.Server, err = droplet.PublicIPv4()
		if err != nil {
			log.Println("unable to get public ip for droplet ", droplet.Name, err)
			success = false
			failednodes = append(failednodes, droplet.Name)
			continue
		}
		log.Println("starting netclient on ", droplet.Name)
		//ensure journalctl is working
		_, err = ssh.Run("systemctl restart systemd-journald")
		if err != nil {
			log.Println("failed to restart journald", err)
		}
		out, err := ssh.Run("./netclient join -t " + token)
		if err != nil {
			success = false
			failednodes = append(failednodes, droplet.Name)
		}
		log.Println("netclient started on ", droplet.Name, "\n", out)
	}
	return success, failednodes
}

// UpdateNodes -- updateNodes to latest version of netclient
func (request *Request) UpdateNodes(tag, branch string) (bool, []string) {
	slog.Info("updating droplets with tag" + tag + " to branch " + branch)
	success := true
	failednodes := []string{}
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		//return error message in failed nodes array -- kludgy
		failednodes = append(failednodes, err.Error())
		return false, failednodes
	}
	wg := sync.WaitGroup{}
	for _, droplet := range droplets {
		ip, _ := droplet.PublicIPv4()
		slog.Info("droplet to update " + droplet.Name + " " + ip)
	}
	for _, droplet := range droplets {
		if droplet.Name == "extclient" {
			slog.Info("skipping " + droplet.Name)
			continue
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, droplet godo.Droplet) {
			defer wg.Done()
			cmd := "apt-get update; apt-get upgrade -y; systemctl restart netclient"
			ssh := &SshConf{
				User:    "root",
				Key:     os.Getenv("HOME") + "/.ssh/id_devops",
				Options: "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=15",
			}
			ssh.Server, _ = droplet.PublicIPv4()
			slog.Info("updating droplet ", "droplet", droplet.Name, "ip", ssh.Server)
			if err != nil {
				slog.Error("unable to get public ip for droplet "+droplet.Name, "droplet", droplet.Name, "ERROR", err)
				failednodes = append(failednodes, droplet.Name)
				return
			}
			if droplet.Name == "docker" {
				cmd = "docker exec -it netclient ./netclient leave devops; /usr/bin/docker-compose pull; /usr/bin/docker-compose up -d"
			}
			_, err := ssh.Run(cmd)
			if err != nil {
				slog.Warn("netclient update failed", "droplet", droplet.Name, "ip", ssh.Server, "cmd", cmd)
				slog.Warn("droplet update failed", "droplet", droplet.Name, "ERROR", err)
				failednodes = append(failednodes, droplet.Name)
			}
			_, err = ssh.Run("journalctl --rotate; journalctl --vacuum-time=10s")
			if err != nil {
				slog.Warn("clearing logs failed", "droplet", droplet.Name, "ERROR", err)
				failednodes = append(failednodes, droplet.Name)
			}
			//clear logs

			//slog.Info("update result", "droplet", droplet.Name, "update", out, "clear log", out2)
			slog.Info("update finished", "droplet", droplet.Name)

		}(&wg, droplet)
	}
	slog.Info("waiting for updates to complete")
	wg.Wait()
	return success, failednodes
}

func (request *Request) InstallDocker(name, tag string) error {
	ip, err := request.GetPublicIP(name, tag)
	if err != nil {
		log.Printf("unable to get public ip for droplet %s: %s ", name, err.Error())
		return err
	}
	//docker needs to be installed
	ssh.Server = ip
	log.Println(" installing docker")
	_, err = ssh.Run("apt-get install docker.io -y")
	if err != nil {
		return err
	}
	return nil
}

func (request *Request) StopDocker(name, tag string) error {
	ip, err := request.GetPublicIP(name, tag)
	if err != nil {
		log.Printf("unable to get public ip for droplet %s: %s ", name, err.Error())
		return err
	}
	ssh.Server = ip
	log.Println("removing existing docker")
	_, err = ssh.Run("docker stop netclient; docker rm netclient")
	if err != nil {
		return err
	}
	return nil
}

func (request *Request) JoinDocker(name, tag, token string, first, useCommunity bool) error {
	ip, err := request.GetPublicIP(name, tag)
	if err != nil {
		log.Printf("unable to get public ip for droplet %s: %s ", name, err.Error())
		return err
	}
	ssh.Server = ip
	//allowcat tty
	ssh.Options = "-t " + ssh.Options
	log.Println("starting netclient docker")
	if first {
		_, err = ssh.Run("docker run -d --network host --privileged -e TOKEN=" + token + " -v /etc/netclient:/etc/netclient --name netclient gravitl/netclient:testing")
		if err != nil {
			return err
		}
	} else {
		_, err = ssh.Run("docker exec -it netclient ./netclient join -t " + token)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ssh *SshConf) Run(cmd string) (string, error) {
	str := fmt.Sprintf("ssh %s -i %s %s@%s %s", ssh.Options, ssh.Key, ssh.User, ssh.Server, cmd)
	out, err := ncutils.RunCmd(str, true)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (ssh *SshConf) Scp(source, destination string) error {
	str := fmt.Sprintf("scp %s -i %s %s %s@%s:%s", ssh.Options, ssh.Key, source, ssh.User, ssh.Server, destination)
	_, err := ncutils.RunCmd(str, true)
	if err != nil {
		return err
	}
	return nil
}

// GetPublicIP returns the public ipv4 address of a named droplet
func (request *Request) GetPublicIP(name, tag string) (string, error) {
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		log.Fatal("could not retrieve droplets ", err)
	}
	for _, droplet := range droplets {
		if droplet.Name != name {
			continue
		}
		return droplet.PublicIPv4()
	}
	return "", errors.New("droplet not found")
}

// GetPrivateIP returns the private ipv4 address of a named droplet
func (request *Request) GetPrivateIP(name, tag string) (string, error) {
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		log.Fatal("could not retrieve droplets ", err)
	}
	for _, droplet := range droplets {
		if droplet.Name != name {
			continue
		}
		return droplet.PrivateIPv4()
	}
	return "", errors.New("droplet not found")
}
