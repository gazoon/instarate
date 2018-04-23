package competition

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"instarate/libs/instagram"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

var (
	profileExistsErr = errors.New("profile already exists")
)

type InstProfile struct {
	Username      string
	ProfileLink   string `bson:"-"`
	PhotoPath     string `bson:"photo"`
	PhotoUrl      string `bson:"-"`
	PhotoInstCode string `bson:"photo_code"`
	Followers     int
	AddedAt       int `bson:"added_at"`
}

func createProfile(username, photoStoragePath, photoInstCode string, followers int) *InstProfile {
	addedAt := utils.TimestampSeconds()
	return &InstProfile{
		Username: username, PhotoPath: photoStoragePath,
		PhotoInstCode: photoInstCode, Followers: followers, AddedAt: addedAt,
	}
}

func (self *InstProfile) getProfileUrl() string {
	return instagram.BuildProfileUrl(self.Username)
}

type profilesStorage struct {
	client *mgo.Collection
}

func newProfilesStorage(mongoSettings *utils.MongoDBSettings) (*profilesStorage, error) {
	db, err := mongo.Connect(mongoSettings)
	if err != nil {
		return nil, err
	}
	collection := db.C(mongoSettings.Collection)
	return &profilesStorage{collection}, nil
}

func (self *profilesStorage) create(model *InstProfile) error {
	err := self.client.Insert(model)
	if mgo.IsDup(err) {
		return profileExistsErr
	}
	return err
}

func (self *profilesStorage) get(username string) (*InstProfile, error) {
	result := &InstProfile{}
	err := self.client.Find(bson.M{"username": username}).One(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (self *profilesStorage) getMultiple(usernames []string) ([]*InstProfile, error) {
	var result []*InstProfile
	err := self.client.Find(bson.M{"username": bson.M{"$in": usernames}}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (self *profilesStorage) delete(usernames []string) error {
	_, err := self.client.RemoveAll(bson.M{"username": bson.M{"$in": usernames}})
	return err
}
