package competition

import (
	"context"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"strings"
)

type VotersStorage struct {
	client *mgo.Collection
}

func newVotersStorage(mongoSettings *utils.MongoDBSettings) (*VotersStorage, error) {
	collection, err := mongo.ConnectCollection(mongoSettings)
	if err != nil {
		return nil, err
	}
	return &VotersStorage{collection}, nil
}

func (self *VotersStorage) tryVote(ctx context.Context, competitionCode, votersGroup, voter,
	competitorOne, competitorTwo string) (bool, error) {
	unitedCompetitorsId := buildUnitedId(competitorOne, competitorTwo)
	err := self.client.Insert(bson.M{
		"voters_group":   votersGroup,
		"competition":    competitionCode,
		"competitors_id": unitedCompetitorsId,
		"voter":          voter,
	})
	if err != nil {
		if mgo.IsDup(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "insert new vote-record document")
	}
	return true, nil
}

func (self *VotersStorage) haveSeenPair(ctx context.Context, competitionCode, votersGroup, competitorOne, competitorTwo string) (bool, error) {
	unitedCompetitorsId := buildUnitedId(competitorOne, competitorTwo)
	rows, err := self.client.Find(bson.M{
		"voters_group":   votersGroup,
		"competition":    competitionCode,
		"competitors_id": unitedCompetitorsId,
	}).Count()
	if err != nil {
		return false, errors.Wrap(err, "count vote-record documents")
	}
	return rows == 1, err
}

func (self *VotersStorage) CreateIndexes() error {
	var err error

	// This index is useful in case we need the history for a specific user.
	err = self.client.EnsureIndex(mgo.Index{Key: []string{"voter"}})
	if err != nil {
		return errors.Wrap(err, "key: voter")
	}

	// voters_group is the first because we may want to retrieve the history for a group.
	err = self.client.EnsureIndex(mgo.Index{Key: []string{
		"voters_group", "competition", "competitors_id", "voter"}, Unique: true})
	if err != nil {
		return errors.Wrap(err, "unique key: voters_group,competition,competitors_id,voter")
	}

	return nil
}

func buildUnitedId(ids ...string) string {
	sort.Strings(ids)
	return strings.Join(ids, " | ")
}
