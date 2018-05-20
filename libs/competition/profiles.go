package competition

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"instarate/libs/instagram"

	"context"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	profileExistsErr = errors.New("profile already exists")
)

type InstProfile struct {
	Username      string
	PhotoPath     string `bson:"photo"`
	PhotoInstCode string `bson:"photo_code"`
	Followers     int
	AddedAt       int `bson:"added_at"`
}

func createProfile(username, photoInstCode string, followers int) *InstProfile {
	addedAt := utils.TimestampSeconds()
	photoStoragePath := username + "-" + uuid.NewV4().String()
	return &InstProfile{
		Username: username, PhotoPath: photoStoragePath,
		PhotoInstCode: photoInstCode, Followers: followers, AddedAt: addedAt,
	}
}

func (self *InstProfile) getProfileLink() string {
	return instagram.BuildProfileUrl(self.Username)
}

func (self InstProfile) String() string {
	return utils.ObjToString(&self)
}

type profilesStorage struct {
	client *mgo.Collection
}

func newProfilesStorage(mongoSettings *utils.MongoDBSettings) (*profilesStorage, error) {
	collection, err := mongo.ConnectCollection(mongoSettings)
	if err != nil {
		return nil, err
	}
	return &profilesStorage{collection}, nil
}

func (self *profilesStorage) save(ctx context.Context, model *InstProfile) error {
	err := self.client.Insert(model)
	if mgo.IsDup(err) {
		return profileExistsErr
	}
	return err
}

func (self *profilesStorage) get(ctx context.Context, username string) (*InstProfile, error) {
	result := &InstProfile{}
	err := self.client.Find(bson.M{"username": username}).One(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (self *profilesStorage) getMultiple(ctx context.Context, usernames []string) ([]*InstProfile, error) {
	var result []*InstProfile
	err := self.client.Find(bson.M{"username": bson.M{"$in": usernames}}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (self *profilesStorage) delete(ctx context.Context, usernames []string) error {
	_, err := self.client.RemoveAll(bson.M{"username": bson.M{"$in": usernames}})
	return err
}
