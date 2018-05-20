package cache

import (
	"context"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Mongo struct {
	client *mgo.Collection
}

func NewMongo(mongoSettings *utils.MongoDBSettings) (*Mongo, error) {
	collection, err := mongo.ConnectCollection(mongoSettings)
	if err != nil {
		return nil, err
	}
	return &Mongo{collection}, nil
}

func (self *Mongo) Get(ctx context.Context, key string) (string, error) {
	var result map[string]string
	err := self.client.Find(bson.M{"key": key}).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return result["value"], nil
}

func (self *Mongo) Set(ctx context.Context, key, value string) error {
	_, err := self.client.Upsert(
		bson.M{"key": key}, bson.M{"key": key, "value": value},
	)
	return err
}
