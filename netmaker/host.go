package netmaker

import (
	"log"
	"net/http"

	"github.com/gravitl/netmaker/models"
)

func GetHostByName(name string) *models.ApiHost {
	hosts := GetHosts()
	for _, host := range *hosts {
		if host.Name == name {
			return &host
		}
	}
	log.Println("failed to find host")
	return nil
}

func GetHosts() *[]models.ApiHost {
	return callapi[[]models.ApiHost](http.MethodGet, "/api/hosts", nil)
}

func GetHostByID(id string, hosts *[]models.ApiHost) *models.ApiHost {
	for _, host := range *hosts {
		if host.ID == id {
			return &host
		}
	}
	return nil
}
