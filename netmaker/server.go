package netmaker

import (
	"net/http"

	"github.com/gravitl/netmaker/models"
)

func ServerConfig() *models.ServerConfig {
	return callapi[models.ServerConfig](http.MethodGet, "/api/server/getserverinfo", nil)
}
