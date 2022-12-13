package models

import (
	"context"
	"fmt"
	"time"

	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/metadata"
	"github.com/lim-lq/dpm/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Account struct {
	BaseModel   `json:",inline" bson:",inline"`
	Username    string `json:"username" bson:"username"`
	Nickname    string `json:"nickname" bson:"nickname"`
	Password    string `json:"password" bson:"password"`
	Email       string `json:"email" bson:"email"`
	Enabled     bool   `json:"enable" bson:"enable"`
	IsSuperuser bool   `json:"is_superuser" bson:"is_supersor"`
	IsAdmin     bool   `json:"is_admin" bson:"is_admin"`
}

type RoleModel struct {
	BaseModel `json:",inline" bson:",inline"`
	Name      string `json:"name" bson:"name"`
	Desc      string `json:"desc" bson:"desc"`
}

type RolesAccountsModel struct {
	BaseModel `json:",inline" bson:",inline"`
	RoleId    int64 `json:"roleid" bson:"roleid"`
	AccountId int64 `json:"accountid" bson:"accountid"`
}

// 权限记录模型
type PrivilegeRecordModel struct {
	ItemId    int64 `json:"itemid" bson:"itemid"` // 记录用户或者角色的id
	ProjectId int64 `json:"projectid" bson:"projectid"`
	ActionId  int64 `json:"actionid" bson:"actionid"`
}

// 权限动作模型
type PrivilegeActionModel struct {
	BaseModel `json:",inline" bson:",inline"`
	Name      string `json:"name" bson:"name"`
	Desc      string `json:"desc" bson:"desc"`
	CateEn    string `json:"cate_en" bson:"cate_en"` // 类别英文描述
	CateCh    string `json:"cate_ch" bson:"cate_ch"` // 类别中文描述
}

var colName = "account"

func AccountManager() *Account {
	return &Account{}
}

func (a *Account) Create(ctx context.Context) error {
	if a.Username == "" {
		return utils.DpmError("用户名不能空")
	}

	if a.Password == "" {
		return utils.DpmError("密码不能为空")
	}
	existAccount := Account{}
	err := existAccount.SearchByName(ctx, a.Username)
	if err == nil {
		return utils.DpmError("用户名已存在")
	}
	mongocli := core.GetMongoClient()
	seqId, err := mongocli.NextSequence(ctx, colName)
	if err != nil {
		return utils.DpmError(fmt.Sprintf("创建用户失败: %v", err))
	}
	a.Id = seqId
	a.Enabled = true
	a.CreateTime = time.Now()
	a.UpdateTime = time.Now()
	return mongocli.InsertOne(ctx, colName, a)
}

func (a *Account) Delete(ctx context.Context, cond *metadata.Condition) error {
	mongocli := core.GetMongoClient()
	return mongocli.Delete(ctx, colName, cond)
}

func (a *Account) GetList(ctx context.Context, cond *metadata.Condition) ([]Account, error) {
	mongocli := core.GetMongoClient()
	accList := []Account{}
	err := mongocli.FindAll(ctx, colName, cond, &accList)
	return accList, err
}

func (a *Account) Update(ctx context.Context, data metadata.MapStr) error {
	_, ok := data["password"]
	if ok {
		delete(data, "password")
	}
	mongocli := core.GetMongoClient()
	// 查询用户是否重复
	cond := metadata.Condition{}
	if username, ok := data["username"]; ok {
		cond.Filters = metadata.Filters{"username": username}
		tmpAcc := Account{}
		err := mongocli.FindOne(ctx, colName, &cond, &tmpAcc)
		if err == nil {
			return utils.DpmError(fmt.Sprintf("用户名重复: %v", err))
		}
	}
	cond.Filters = metadata.Filters{"id": a.Id}

	return mongocli.Update(ctx, colName, &cond, bson.D{{Key: "$set", Value: *TransSetUpdate(data)}})
}

func (a *Account) ChangePassword(ctx context.Context, cipher string) error {
	mongocli := core.GetMongoClient()
	cond := metadata.Condition{Filters: metadata.Filters{"id": a.Id}}
	return mongocli.Update(ctx, colName, &cond, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: cipher},
		}}})
}

func (a *Account) SearchByName(ctx context.Context, username string) error {
	mongocli := core.GetMongoClient()
	cond := &metadata.Condition{
		Filters: map[string]interface{}{
			"username": map[string]string{
				"$eq": username,
			},
		},
	}
	return mongocli.FindOne(ctx, colName, cond, a)
}
