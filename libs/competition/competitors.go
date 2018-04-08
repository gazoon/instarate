package competition

import (
	"fmt"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

const (
	initialRating        = 1500
	maxRandomGetAttempts = 5
)

var (
	competitorNotFoundErr = errors.New("competitor doesn't exist")
)

type competitor struct {
	Username        string
	CompetitionCode string `bson:"competition"`
	Rating          int
	Matches         int
	Wins            int
	Loses           int
}

func createCompetitor(username, competitionCode string) *competitor {
	return &competitor{Username: username, CompetitionCode: competitionCode, Rating: initialRating}
}

type competitorsStorage struct {
	client *mgo.Collection
}

func newCompetitorsStorage(mongoSettings *utils.MongoDBSettings) (*competitorsStorage, error) {
	db, err := mongo.Connect(mongoSettings)
	if err != nil {
		return nil, err
	}
	collection := db.C(mongoSettings.Collection)
	return &competitorsStorage{collection}, nil
}

func (self *competitorsStorage) getTop(competitionCode string, number, offset int) ([]*competitor, error) {
	var result []*competitor
	err := self.client.Find(bson.M{"competition": competitionCode}).
		Sort("-rating").Skip(offset).Limit(number).All(&result)
	return result, err
}

func (self *competitorsStorage) delete(usernames []string) error {
	_, err := self.client.RemoveAll(bson.M{"username": bson.M{"$in": usernames}})
	return err
}

func (self *competitorsStorage) getCompetitorsNumber(competitionCode string) (int, error) {
	return self.client.Find(bson.M{"competition": competitionCode}).Count()
}

func (self *competitorsStorage) get(competitionCode, username string) (*competitor, error) {
	result := &competitor{}
	err := self.client.Find(bson.M{"competition": competitionCode, "username": username}).One(result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, competitorNotFoundErr
		}
		return nil, err
	}
	return result, nil
}

func (self *competitorsStorage) getNumberWithHigherRating(competitionCode string, rating int) (int, error) {
	var result []*competitor
	fmt.Println(competitionCode, rating)
	err := self.client.Find(bson.M{
		"competition": competitionCode, "rating": bson.M{"$gt": rating},
	}).All(&result)
	if err != nil {
		return 0, err
	}
	ratings := map[int]bool{}
	for _, compettr := range result {
		ratings[compettr.Rating] = true
	}
	return len(ratings), nil
}

func (self *competitorsStorage) update(model *competitor) error {
	err := self.client.Update(
		bson.M{"competition": model.CompetitionCode, "username": model.Username},
		bson.M{"$set": bson.M{
			"rating":  model.Rating,
			"wins":    model.Wins,
			"loses":   model.Loses,
			"matches": model.Matches,
		}},
	)
	return err
}

func (self *competitorsStorage) create(model *competitor) error {
	err := self.client.Insert(model)
	return err
}

func (self *competitorsStorage) getRandomPair(competitionCode string) (*competitor, *competitor, error) {
	for i := 0; i < maxRandomGetAttempts; i++ {
		var result []*competitor
		err := self.client.Pipe([]bson.M{
			{"$match": bson.M{"competition": competitionCode}},
			{"$sample": bson.M{"size": 2}},
		}).All(&result)
		if err != nil {
			return nil, nil, err
		}
		if len(result) < 2 {
			return nil, nil, errors.Errorf("not enough competitors in %s%", competitionCode)
		}
		competitorOne, competitorTwo := result[0], result[1]
		if competitorOne.Username != competitorTwo.Username {
			return competitorOne, competitorTwo, nil
		}
	}
	return nil, nil, errors.Errorf("can't get two distinct competitor in %s", competitionCode)
}
