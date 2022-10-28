package models

import (
	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/metadata"
)

type Account struct {
	Username string `json:"username" bson:"username"`
	Nickname string `json:"nickname" bson:"nickname"`
	Password string `json:"" bson:"password"`
	Enabled  bool   `json:"enable" bson:"enable"`
}

func (a *Account) SearchByName(username string) error {
	mongocli := core.GetMongoClient()
	cond := metadata.Condition{
		Filters: map[string]interface{}{
			"username": map[string]string{
				"$eq": username,
			},
		},
	}
	return mongocli.FindOne("account", cond, a)
}
