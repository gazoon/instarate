package chats

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
		return nil, err
	}
	return chat, nil
}

func (self *MongoStorage) Save(ctx context.Context, chat *models.Chat) error {
	_, err := self.client.Upsert(bson.M{"chat_id": chat.Id}, chat)
	return err
}
