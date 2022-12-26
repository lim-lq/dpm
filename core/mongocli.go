package core

import (
	"context"
	"fmt"
	"time"

	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/core/log"
	"github.com/lim-lq/dpm/metadata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	cli    *mongo.Client
	Host   string
	Port   int
	User   string
	Passwd string
	DB     string
	// ctx    context.Context
	// opts   *options.FindOptions
}

type IdGen struct {
	Id         string `bson:"_id"`
	SequenceID int64  `bson:"SequenceID"`
}

var mongocli *mongoClient

func GetMongoClient() *mongoClient {
	return mongocli
}

func TransSetUpdate(data metadata.MapStr) *bson.D {
	update := bson.D{}
	for filed, value := range data {
		update = append(update, bson.E{Key: filed, Value: value})
	}
	return &update
}

func InitMongo() {
	host := config.GetString("mongo.host")
	port := config.GetInt("mongo.port")
	user := config.GetString("mongo.user")
	pass := config.GetString("mongo.pass")
	db := config.GetString("mongo.db")
	var url string
	if pass == "" {
		log.Logger.Warn("[WARN] You will connect mongodb without authentication")
		url = fmt.Sprintf("mongodb://%s:%d", host, port)
	} else {
		url = fmt.Sprintf("mongodb://%s:%s@%s:%d", user, pass, host, port)
	}
	ctx, cancle := context.WithTimeout(context.Background(), 20*time.Second)
	dbcli, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Logger.Fatalf("connect mongodb error - %v", err)
	}
	mongocli = &mongoClient{
		cli:    dbcli,
		Host:   host,
		Port:   port,
		User:   user,
		Passwd: pass,
		DB:     db,
		// ctx:    context.Background(),
	}
	defer func() {
		cancle()
	}()
}

func (m *mongoClient) NextSequence(ctx context.Context, seqName string) (int64, error) {
	colobj := m.cli.Database(m.DB).Collection("seqgenertor")
	update := bson.M{
		"$inc":         bson.M{"SequenceID": int64(1)},
		"$setOnInsert": bson.M{"createTime": time.Now()},
		"$set":         bson.M{"updateTime": time.Now()},
	}
	filter := bson.M{"_id": seqName}
	upsert := true
	returnChange := options.After
	opt := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnChange,
	}
	doc := IdGen{}
	err := colobj.FindOneAndUpdate(ctx, filter, update, opt).Decode(&doc)
	if err != nil {
		return 0, err
	}
	return doc.SequenceID, nil
}

func (m *mongoClient) Count(ctx context.Context, col string, cond *metadata.Condition) (int64, error) {
	colobj := m.cli.Database(m.DB).Collection(col)
	return colobj.CountDocuments(ctx, cond.Filters)
}

func (m *mongoClient) FindAll(ctx context.Context, col string, cond *metadata.Condition, result interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	opts := options.Find()
	opts.SetLimit(cond.Limit)
	opts.SetSkip(cond.Offset)
	cur, err := colobj.Find(ctx, cond.Filters, opts)
	if err != nil {
		return err
	}
	return cur.All(ctx, result)
}

func (m *mongoClient) FindOne(ctx context.Context, col string, cond *metadata.Condition, result interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	return colobj.FindOne(ctx, cond.Filters).Decode(result)
}

func (m *mongoClient) InsertOne(ctx context.Context, col string, data interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	_, err := colobj.InsertOne(ctx, data)
	return err
}

func (m *mongoClient) InsertMany(ctx context.Context, col string, data []interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	_, err := colobj.InsertMany(ctx, data)
	return err
}

func (m *mongoClient) Delete(ctx context.Context, col string, cond *metadata.Condition) error {
	colObj := m.cli.Database(m.DB).Collection(col)
	_, err := colObj.DeleteMany(ctx, cond.Filters)
	return err
}

func (m *mongoClient) Update(ctx context.Context, col string, cond *metadata.Condition, data metadata.MapStr) error {
	colObj := m.cli.Database(m.DB).Collection(col)

	_, err := colObj.UpdateMany(ctx, cond.Filters, bson.D{{Key: "$set", Value: *TransSetUpdate(data)}})
	return err
}
