package models

import (
	"context"
	"time"

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

func (p *ProjectModel) Create(ctx context.Context) error {
	mongocli := core.GetMongoClient()
	id, err := mongocli.NextSequence(ctx, projectCol)
	if err != nil {
		return err
	}
	p.Id = id
	p.CreateTime.Time = time.Now()
	p.UpdateTime.Time = time.Now()
	return mongocli.InsertOne(ctx, projectCol, p)
}

func (p *ProjectModel) GetList(ctx context.Context, cond *metadata.Condition) ([]ProjectModel, error) {
	projects := []ProjectModel{}
	mongocli := core.GetMongoClient()
	err := mongocli.FindAll(ctx, projectCol, cond, &projects)

	return projects, err
}

func (p *ProjectModel) GetPageList(ctx context.Context, cond *metadata.Condition) (*metadata.PageData, error) {
	projects, err := p.GetList(ctx, cond)
	if err != nil {
		return nil, err
	}
	mongocli := core.GetMongoClient()
	totalCount, err := mongocli.Count(ctx, projectCol, cond)
	if err != nil {
		return nil, err
	}
	return &metadata.PageData{
		TotalCount: totalCount,
		Data:       projects,
		PageNo:     cond.Offset/cond.Limit + 1,
		PageSize:   cond.Limit,
	}, err
}

func (p *ProjectModel) GetDetail(ctx context.Context, cond *metadata.Condition) error {
	db := core.GetMongoClient()
	return db.FindOne(ctx, projectCol, cond, p)
}

func (p *ProjectModel) Update(ctx context.Context, cond *metadata.Condition, data metadata.MapStr) error {
	db := core.GetMongoClient()
	return db.Update(ctx, projectCol, cond, data)
}
