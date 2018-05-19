package tasks

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Task struct {
	ChatId int                    `bson:"chat_id" mapstructure:"chat_id"`
	Name   string                 `bson:"name" mapstructure:"name"`
	Args   map[string]interface{} `bson:"args" mapstructure:"args"`
	DoAt   int                    `bson:"do_at" mapstructure:"do_at"`
}

func TaskFromData(data interface{}) (*Task, error) {
	t := &Task{}
	err := mapstructure.Decode(data, t)
	if err != nil {
		return nil, errors.Wrap(err, "can't create task from data")
	}
	return t, nil
}

func (self Task) String() string {
	return utils.ObjToString(&self)
}

var (
	TaskAlreadyExistsErr = errors.New(
		"Task with that name already exists for the chat")
)

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

func (self *Storage) GetTask() *Task {
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

func (self *Storage) CreateTask(task *Task) error {
	err := self.client.Insert(task)
	if mgo.IsDup(err) {
		return TaskAlreadyExistsErr
	}
	return err
}

func (self *Storage) CreateOrReplaceTask(task *Task) error {
	_, err := self.client.Upsert(bson.M{"chat_id": task.ChatId, "name": task.Name}, task)
	return err
}

func (self *Storage) DeleteTask(chatId int, name string) error {
	return self.client.Remove(bson.M{"chat_id": chatId, "name": name})
}
