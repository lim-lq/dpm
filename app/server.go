package app

import (
	"fmt"

	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/service"
)

func StartServer() error {
	ws := service.WebService()
	host := config.GetString("http.host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := config.GetInt("http.port")
	if port == 0 {
		port = 5002
	}
	return ws.Run(fmt.Sprintf("%s:%d", host, port))
}
