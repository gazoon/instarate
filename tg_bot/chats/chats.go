package chats

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"instarate/tg_bot/models"
)

type MongoStorage struct {
	client *mgo.Collection
}

func NewMongoStorage(mongoSettings *utils.MongoDBSettings) (*MongoStorage, error) {
	collection, err := mongo.ConnectCollection(mongoSettings)
	if err != nil {
		return nil, err
	}
	return &MongoStorage{collection}, nil
}

func (self *MongoStorage) Get(ctx context.Context, chatId int) (*models.Chat, error) {
	chat := &models.Chat{}
	err := self.client.Find(bson.M{"chat_id": chatId}).One(chat)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "get chat document")
	}
	return chat, nil
}

func (self *MongoStorage) Save(ctx context.Context, chat *models.Chat) error {
	_, err := self.client.Upsert(bson.M{"chat_id": chat.Id}, chat)
	return errors.Wrap(err, "upsert chat document")
}

func (self *MongoStorage) ResetState(ctx context.Context, chatId int) error {
	err := self.client.Update(bson.M{"chat_id": chatId},
		bson.M{"$set": bson.M{"last_match": nil, "current_top_offset": 0}})
	return errors.Wrap(err, "reset chat state")
}

func (self *MongoStorage) CreateIndexes() error {
	var err error

	err = self.client.EnsureIndex(mgo.Index{Key: []string{"chat_id"}, Unique: true})
	if err != nil {
		return errors.Wrap(err, "unique key: chat_id")
	}

	return nil
}
