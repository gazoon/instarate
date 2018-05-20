package competition

import (
	"context"
	"fmt"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

const (
	initialRating         = 1500
	randomPairGetAttempts = 5
)

var (
	CompetitorNotFoundErr = errors.New("competitor doesn't exist")
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
	collection, err := mongo.ConnectCollection(mongoSettings)
	if err != nil {
		return nil, err
	}
	return &competitorsStorage{collection}, nil
}

func (self *competitorsStorage) getTop(ctx context.Context, competitionCode string, number, offset int) ([]*competitor, error) {
	var result []*competitor
	err := self.client.Find(bson.M{"competition": competitionCode}).
		Sort("-rating").Skip(offset).Limit(number).All(&result)
	return result, err
}

func (self *competitorsStorage) delete(ctx context.Context, usernames []string) error {
	_, err := self.client.RemoveAll(bson.M{"username": bson.M{"$in": usernames}})
	return err
}

func (self *competitorsStorage) getCompetitorsNumber(ctx context.Context, competitionCode string) (int, error) {
	return self.client.Find(bson.M{"competition": competitionCode}).Count()
}

func (self *competitorsStorage) get(ctx context.Context, competitionCode, username string) (*competitor, error) {
	result := &competitor{}
	err := self.client.Find(bson.M{"competition": competitionCode, "username": username}).One(result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, CompetitorNotFoundErr
		}
		return nil, err
	}
	return result, nil
}

func (self *competitorsStorage) getNumberWithHigherRating(ctx context.Context, competitionCode string, rating int) (int, error) {
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

func (self *competitorsStorage) update(ctx context.Context, model *competitor) error {
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

func (self *competitorsStorage) create(ctx context.Context, model *competitor) error {
	err := self.client.Insert(model)
	return err
}

func (self *competitorsStorage) getRandomPair(ctx context.Context, competitionCode string) (*competitor, *competitor, error) {
	for i := 0; i < randomPairGetAttempts; i++ {
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
