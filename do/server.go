package do

import (
	"log"
	"strings"

	"golang.org/x/exp/slog"
)

type Server struct {
	FQDN      string
	Broker    string
	API       string
	Dashboard string
	PublicIP  string
	Subdomain string
	Branch    string
	UIBranch  string
}

func (request *Request) ResetServer(server *Server) {
	ssh.Server = server.FQDN
	log.Println("shutting down docker")
	_, err := ssh.Run("docker-compose down")
	if err != nil {
		log.Println("err running docker-compose down ", err)
	}
	log.Println("removing docker artifacts")
	_, err = ssh.Run("docker system prune -f --all --volumes")
	if err != nil {
		log.Println("error running docker prune", err)
	}
	log.Println("copying files")
	if err := request.CopyServerFiles(server); err != nil {
		log.Println("error copying files to server", err)
	}
	log.Println("starting docker")
	if err := request.StartDocker(server); err != nil {
		log.Println("error starting docker", err)
	}
}

func (request *Request) SoftResetServer(server *Server) {
	ssh.Server = server.FQDN
	slog.Info("retrieving new docker images docker")
	_, err := ssh.Run("docker-compose pull")
	if err != nil {
		log.Println("err running docker-compose pull ", err)
	}
	slog.Info("starting docker")
	if err := request.StartDocker(server); err != nil {
		log.Println("error starting docker", err)
	}
}

//func (request *Request) CreateAdmin(server *Server) {
//	log.Println("creating admin on server ", server.FQDN)
//	url := "https://" + server.API
//	if err := netmaker.AddAdmin(url); err != nil {
//		log.Println("error creating admin", err)
//	}
//
//}

// CopyServerFiles - copies file to server
func (request *Request) CopyServerFiles(server *Server) error {
	ssh.Server = server.FQDN
	log.Println("copying files to server", server.FQDN)
	if err := ssh.Scp("files/certs-traefik.yml", "~/"); err != nil {
		return err
	}
	source := "files/server.docker-compose.yml"
	cmd1 := "cp /etc/letsencrypt/live/dashboard.clustercat.com/fullchain.pem /root/certs/"
	cmd2 := "cp /etc/letsencrypt/live/dashboard.clustercat.com/privkey.pem /root/certs/"
	if strings.Contains(server.FQDN, "2") {
		source = "files/server2.docker-compose.yml"
		cmd1 = "cp /etc/letsencrypt/live/dashboard2.sandbox.clustercat.com/fullchain.pem /root/certs/"
		cmd2 = "cp /etc/letsencrypt/live/dashboard2.sandbox.clustercat.com/privkey.pem /root/certs/"
	} else if strings.Contains(server.FQDN, "1") {
		source = "files/server1.docker-compose.yml"
		cmd1 = "cp /etc/letsencrypt/live/dashboard1.sandbox.clustercat.com/fullchain.pem /root/certs/"
		cmd2 = "cp /etc/letsencrypt/live/dashboard1.sandbox.clustercat.com/privkey.pem /root/certs/"
	}
	if err := ssh.Scp(source, "~/docker-compose.yml"); err != nil {
		return err
	}
	cmd := cmd1 + "; " + cmd2
	if _, err := ssh.Run(cmd); err != nil {
		return err
	}
	return nil
}

//func (server *Server) AddNetwork(name, addressRange string, udpholepunch bool) {
//	log.Printf("creating network %s on server %s\n", name, server.FQDN)
//	network := models.Network{}
//	network.NetID = name
//	network.AddressRange = addressRange
//	if udpholepunch {
//		network.DefaultUDPHolePunch = "yes"
//	}
//	baseURL := "https://" + server.API
//	resp, err := netmaker.Api(network, http.MethodPost, baseURL+"/api/networks", "secretkey")
//	if err != nil {
//		log.Println("error creating network ", err)
//	} else if resp.StatusCode != http.StatusOK {
//		log.Println("error creating newtork ", resp.Status)
//	}
//}
