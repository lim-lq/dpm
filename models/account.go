package models

import (
	"context"
	"fmt"
	"time"

	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/metadata"
	"github.com/lim-lq/dpm/utils"
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

type AccountInfo struct {
	*Account `json:",inline"`
	Roles    []string               `json:"roles"`
	Actions  []PrivilegeActionModel `json:"actions"`
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
	BaseModel  `json:",inline" bson:",inline"`
	ItemId     int64 `json:"itemid" bson:"itemid"`         // 记录用户或者角色的id
	ResourceId int64 `json:"resourceid" bson:"resourceid"` // 资源ID，0为所有资源
	ActionId   int64 `json:"actionid" bson:"actionid"`
}

// 权限动作模型
type PrivilegeActionModel struct {
	BaseModel    `json:",inline" bson:",inline"`
	Name         string `json:"name" bson:"name"`
	Desc         string `json:"desc" bson:"desc"`
	CateEn       string `json:"cate_en" bson:"cate_en"`             // 类别英文描述
	CateCh       string `json:"cate_ch" bson:"cate_ch"`             // 类别中文描述
	NeedResource bool   `json:"need_resource" bson:"need_resource"` // 是否需要资源绑定
	Method       string `json:"method" bson:"method"`               // http 方法
	Url          string `json:"url" bson:"url"`                     // 接口url
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
	a.CreateTime.Time = time.Now()
	a.UpdateTime.Time = time.Now()
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

func (a *Account) GetInfoByUsername(ctx context.Context, username string) (*AccountInfo, error) {
	cond := metadata.Condition{
		Filters: metadata.Filters{"username": username},
	}
	return a.GetInfo(ctx, &cond)
}

func (a *Account) GetInfoById(ctx context.Context, id int64) (*AccountInfo, error) {
	cond := metadata.Condition{
		Filters: metadata.Filters{"id": id},
	}
	return a.GetInfo(ctx, &cond)
}

func (a *Account) GetInfo(ctx context.Context, cond *metadata.Condition) (*AccountInfo, error) {
	mongocli := core.GetMongoClient()
	err := mongocli.FindOne(ctx, colName, cond, a)
	if err != nil {
		return nil, err
	}
	actions := []PrivilegeActionModel{}
	if !a.IsAdmin {
		// 获取角色相关actions
		cond := metadata.Condition{
			Filters: metadata.Filters{"accountid": a.Id},
		}
		rolesAccounts := []RolesAccountsModel{}
		err = mongocli.FindAll(ctx, "roles_accounts", &cond, &rolesAccounts)
		if err != nil {
			return nil, err
		}
		for _, rc := range rolesAccounts {
			// 获取角色对应的actions
			cond = metadata.Condition{
				Filters: metadata.Filters{"itemid": rc.RoleId},
			}
			records := []PrivilegeRecordModel{}
			err = mongocli.FindAll(ctx, "privilege_record", &cond, &records)
			if err != nil {
				return nil, err
			}
			for _, pr := range records {
				// 获取action名字
				cond = metadata.Condition{
					Filters: metadata.Filters{"id": pr.ActionId},
				}
				action := PrivilegeActionModel{}
				err = mongocli.FindOne(ctx, "privilege_action", &cond, action)
				if err != nil {
					return nil, err
				}
				actions = append(actions, action)
			}
		}
	}
	return &AccountInfo{Account: a, Actions: actions}, nil
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

	return mongocli.Update(ctx, colName, &cond, data)
}

func (a *Account) ChangePassword(ctx context.Context, cipher string) error {
	mongocli := core.GetMongoClient()
	cond := metadata.Condition{Filters: metadata.Filters{"id": a.Id}}
	return mongocli.Update(ctx, colName, &cond, metadata.MapStr{"password": cipher})
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

func (pa *PrivilegeActionModel) Update(ctx context.Context, actions []metadata.MapStr) error {
	var err error
	var newId int64
	mongocli := core.GetMongoClient()
	colString := "privilege_action"
	actionNames := []string{}
	for _, action := range actions {
		cond := metadata.Condition{Filters: metadata.Filters{"name": action["name"]}}
		actionObj := PrivilegeActionModel{}
		err = mongocli.FindOne(ctx, colString, &cond, &actionObj)
		action["updateTime"] = time.Now()
		// 存在
		if err == nil {
			err = mongocli.Update(ctx, colString, &cond, action)
			if err != nil {
				return err
			}
		} else {
			// 不存在
			newId, err = mongocli.NextSequence(ctx, colString)
			if err != nil {
				return err
			}
			action["id"] = newId
			action["createTime"] = time.Now()
			err = mongocli.InsertOne(ctx, colString, action)
			if err != nil {
				return err
			}
		}
		actionNames = append(actionNames, action["name"].(string))
	}
	// 清除多余的
	cond := metadata.Condition{
		Filters: metadata.Filters{"name": map[string][]string{"$nin": actionNames}},
	}
	return mongocli.Delete(ctx, colString, &cond)
}
