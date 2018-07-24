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

type Competitor struct {
	Username        string `bson:"username"`
	CompetitionCode string `bson:"competition"`
	Rating          int    `bson:"rating"`
	Matches         int    `bson:"matches"`
	Wins            int    `bson:"wins"`
	Loses           int    `bson:"loses"`
}

func CreateCompetitor(username, competitionCode string) *Competitor {
	return &Competitor{Username: username, CompetitionCode: competitionCode, Rating: initialRating}
}

type CompetitorsStorage struct {
	client *mgo.Collection
}

func NewCompetitorsStorage(mongoSettings *utils.MongoDBSettings) (*CompetitorsStorage, error) {
	collection, err := mongo.ConnectCollection(mongoSettings)
	if err != nil {
		return nil, err
	}
	return &CompetitorsStorage{collection}, nil
}

func (self *CompetitorsStorage) GetTop(ctx context.Context, competitionCode string, number, offset int) ([]*Competitor, error) {
	var result []*Competitor
	err := self.client.Find(bson.M{"competition": competitionCode}).
		Sort("-rating").Skip(offset).Limit(number).All(&result)
	return result, errors.Wrap(err, "get top from mongo")
}

func (self *CompetitorsStorage) Delete(ctx context.Context, usernames []string) error {
	_, err := self.client.RemoveAll(bson.M{"username": bson.M{"$in": usernames}})
	return errors.Wrap(err, "delete all by usernames")
}

func (self *CompetitorsStorage) GetCompetitorsNumber(ctx context.Context, competitionCode string) (int, error) {
	num, err := self.client.Find(bson.M{"competition": competitionCode}).Count()
	return num, errors.Wrap(err, "count all competitors documents")
}

func (self *CompetitorsStorage) Get(ctx context.Context, competitionCode, username string) (*Competitor, error) {
	result := &Competitor{}
	err := self.client.Find(bson.M{"competition": competitionCode, "username": username}).One(result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, &CompetitorNotFound{Username: username}
		}
		return nil, errors.Wrap(err, "get single Competitor document")
	}
	return result, nil
}

func (self *CompetitorsStorage) GetNumberWithHigherRating(ctx context.Context, competitionCode string, rating int) (int, error) {
	var result []*Competitor
	fmt.Println(competitionCode, rating)
	err := self.client.Find(bson.M{
		"competition": competitionCode, "rating": bson.M{"$gt": rating},
	}).All(&result)
	if err != nil {
		return 0, errors.Wrap(err, "get Competitor documents with higher rating")
	}
	ratings := map[int]bool{}
	for _, compettr := range result {
		ratings[compettr.Rating] = true
	}
	return len(ratings), nil
}

func (self *CompetitorsStorage) Update(ctx context.Context, model *Competitor) error {
	err := self.client.Update(
		bson.M{"competition": model.CompetitionCode, "username": model.Username},
		bson.M{"$set": bson.M{
			"rating":  model.Rating,
			"wins":    model.Wins,
			"loses":   model.Loses,
			"matches": model.Matches,
		}},
	)
	return errors.Wrap(err, "update Competitor document")
}

func (self *CompetitorsStorage) Create(ctx context.Context, model *Competitor) error {
	err := self.client.Insert(model)
	return errors.Wrap(err, "insert new Competitor document")
}

func (self *CompetitorsStorage) GetRandomPair(ctx context.Context, competitionCode string) (*Competitor, *Competitor, error) {
	for i := 0; i < randomPairGetAttempts; i++ {
		var result []*Competitor
		err := self.client.Pipe([]bson.M{
			{"$match": bson.M{"competition": competitionCode}},
			{"$sample": bson.M{"size": 2}},
		}).All(&result)
		if err != nil {
			return nil, nil, errors.Wrap(err, "get two random competitors documents")
		}
		if len(result) < 2 {
			return nil, nil, errors.Errorf("not enough competitors in %s%", competitionCode)
		}
		competitorOne, competitorTwo := result[0], result[1]
		if competitorOne.Username != competitorTwo.Username {
			return competitorOne, competitorTwo, nil
		}
	}
	return nil, nil, errors.Errorf("can't get two distinct Competitor in %s", competitionCode)
}

func (self *CompetitorsStorage) CreateIndexes() error {
	var err error

	err = self.client.EnsureIndex(
		mgo.Index{Key: []string{"competition", "username"}, Unique: true})
	if err != nil {
		return errors.Wrap(err, "unique key: competition,username")
	}

	err = self.client.EnsureIndex(mgo.Index{Key: []string{"competition", "-rating"}})
	if err != nil {
		return errors.Wrap(err, "key: competition,-rating")
	}

	return nil
}
