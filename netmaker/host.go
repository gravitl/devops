package netmaker

import (
	"net/http"

	"github.com/gravitl/netmaker/models"
	"golang.org/x/exp/slog"
)

func GetHostByName(name string) *models.ApiHost {
	hosts := GetHosts()
	for _, host := range *hosts {
		if host.Name == name {
			return &host
		}
	}
	slog.Error("failed to find host", "func", "GetHostByName")
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
	slog.Error("failed to find host", "func", "GetHostByID")
	return nil
}
