// DPM internal login plugin
package dpm

import (
	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/middleware/login/plugins"
	"github.com/lim-lq/dpm/models"
)

func init() {
	plugins.RegisterPlugin(&plugins.LoginPlugin{
		Version: "dpm",
		Handler: &user{},
	})
}

type user struct {
}

func (u *user) IsAuthed(c *gin.Context) bool {
	return true
}

func (u *user) GetLoginUrl(c *gin.Context) string {
	return "/login"
}
func (u *user) GetUser(g *gin.Context) (*models.Account, error) {
	return nil, nil
}
func (u *user) GetUserList(c *gin.Context) ([]*models.Account, error) {
	return nil, nil
}
