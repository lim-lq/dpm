package login

import (
	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/middleware/login/plugins"
	_ "github.com/lim-lq/dpm/middleware/login/plugins/register"
)

func newSessionID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 32)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = alphabet[v%byte(len(alphabet))]
	}
	return string(bytes)
}

func IsAuthed() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionid, err := c.Cookie(config.SessionName)
		if err != nil {
			sessionid = newSessionID()
			c.SetCookie(config.SessionName, sessionid, 24*60*60, "/", "", false, true)
		}
		userManager := plugins.CurrentLoginPlugin(config.LoginVersion)
		if userManager == nil {
			c.Next()
			return
		}

		if c.Request.RequestURI == "/api/login" {
			c.Next()
			return
		}

		// 获取对应的登录信息
		rediscli := core.GetRedisClient()
		_, err = rediscli.Get(sessionid)
		if err != nil {
			c.Redirect(302, userManager.GetLoginUrl(c))
			c.Abort()
			return
		}
		c.Next()
	}
}
