package core

import (
	"context"
	"fmt"
	"time"

	"github.com/lim-lq/dpm/core/config"
	"github.com/lim-lq/dpm/core/log"
	"github.com/lim-lq/dpm/metadata"
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
	ctx    context.Context
	// opts   *options.FindOptions
}

var mongocli *mongoClient

func GetMongoClient() *mongoClient {
	return mongocli
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
		ctx:    context.Background(),
	}
	defer func() {
		cancle()
	}()
}

func (m *mongoClient) Count(col string, cond metadata.Condition) (int64, error) {
	colobj := m.cli.Database(m.DB).Collection(col)
	return colobj.CountDocuments(m.ctx, cond.Filters)
}

func (m *mongoClient) FindAll(col string, cond metadata.Condition, result interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	opts := options.Find()
	opts.SetLimit(cond.Limit)
	opts.SetSkip(cond.Offset)
	cur, err := colobj.Find(m.ctx, cond.Filters, opts)
	if err != nil {
		return err
	}
	return cur.All(m.ctx, result)
}

func (m *mongoClient) FindOne(col string, cond metadata.Condition, result interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	return colobj.FindOne(m.ctx, cond.Filters).Decode(result)
}

func (m *mongoClient) InsertOne(col string, data interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	_, err := colobj.InsertOne(m.ctx, data)
	return err
}

func (m *mongoClient) InsertMany(col string, data []interface{}) error {
	colobj := m.cli.Database(m.DB).Collection(col)
	_, err := colobj.InsertMany(m.ctx, data)
	return err
}
