package tasks

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Task struct {
	ChatId int                    `bson:"chat_id"`
	Name   string                 `bson:"name"`
	Args   map[string]interface{} `bson:"args"`
	DoAt   int                    `bson:"do_at"`
}

func (self Task) String() string {
	return utils.ObjToString(&self)
}

type Storage struct {
	client *mgo.Collection
}

func NewStorage(mongoSettings *utils.MongoDBSettings) (*Storage, error) {
	db, err := mongo.Connect(mongoSettings)
	if err != nil {
		return nil, err
	}
	collection := db.C(mongoSettings.Collection)
	return &Storage{collection}, nil
}

func (self *Storage) GetTask() interface{} {
	currentTime := utils.TimestampMilliseconds()
	task := &Task{}
	_, err := self.client.Find(bson.M{"do_at": bson.M{"$lte": currentTime}}).
		Apply(mgo.Change{Remove: true}, task)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf("Cannot get task from mongo: %s", err)
		}
		return nil
	}
	return task
}
