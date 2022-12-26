package models

import (
	"encoding/json"

	"github.com/lim-lq/dpm/metadata"
	"go.mongodb.org/mongo-driver/bson"
)

// type JSONTime time.Time

// func (t *JSONTime) MarshalJSON() ([]byte, error) {
// 	stamp := time.Time(*t).Format("2006-01-02 15:04:05")
// 	return []byte(stamp), nil
// }

type BaseModel struct {
	Id         int64            `json:"id" bson:"id"`
	CreateTime metadata.DpmTime `json:"createTime" bson:"createTime"`
	UpdateTime metadata.DpmTime `json:"updateTime" bson:"updateTime"`
}

func (b *BaseModel) ToMap() metadata.MapStr {
	result := metadata.MapStr{}
	jsonBytes, _ := json.Marshal(b)
	json.Unmarshal(jsonBytes, &result)
	return result
}

func TransSetUpdate(data metadata.MapStr) *bson.D {
	update := bson.D{}
	for filed, value := range data {
		update = append(update, bson.E{Key: filed, Value: value})
	}
	return &update
}
