package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/models"
)

type UserInterface interface {
	IsAuthed(c *gin.Context) bool
	GetLoginUrl(c *gin.Context) string
	GetUser(g *gin.Context) (*models.Account, error)
	GetUserList(c *gin.Context) ([]*models.Account, error)
}

type LoginPlugin struct {
	Version string
	Handler UserInterface
}

var allLoginPlugins []*LoginPlugin

func RegisterPlugin(plugin *LoginPlugin) {
	allLoginPlugins = append(allLoginPlugins, plugin)
}

func CurrentLoginPlugin(version string) UserInterface {
	for _, plugin := range allLoginPlugins {
		if plugin.Version == version {
			return plugin.Handler
		}
	}
	return nil
}
