package models

import (
	"context"

	"github.com/lim-lq/dpm/core"
	"github.com/lim-lq/dpm/metadata"
)

var projectCol string = "project"

type ProjectModel struct {
	BaseModel `json:",inline" bson:",inline"`
	Name      string `json:"name" bson:"name"`
	Desc      string `json:"desc" bson:"desc"`
	Owner     string `json:"owner" bson:"owner"`
	Status    string `json:"status" bson:"status"`
}

func ProjectManager() *ProjectModel {
	return &ProjectModel{}
}

func (p *ProjectModel) GetList(ctx context.Context, cond *metadata.Condition) ([]ProjectModel, error) {
	projects := []ProjectModel{}
	mongocli := core.GetMongoClient()
	err := mongocli.FindAll(ctx, projectCol, cond, &projects)

	return projects, err
}

func (p *ProjectModel) GetDetail(ctx context.Context, cond *metadata.Condition) error {
	db := core.GetMongoClient()
	return db.FindOne(ctx, projectCol, cond, p)
}

func (p *ProjectModel) Update(ctx context.Context, cond *metadata.Condition, data *ProjectModel) error {
	db := core.GetMongoClient()
	return db.Update(ctx, projectCol, cond, data)
}
