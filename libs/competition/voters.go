package competition

import (
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/mongo"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"strings"
)

type votersStorage struct {
	client *mgo.Collection
}

func newVotersStorage(mongoSettings *utils.MongoDBSettings) (*votersStorage, error) {
	db, err := mongo.Connect(mongoSettings)
	if err != nil {
		return nil, err
	}
	collection := db.C(mongoSettings.Collection)
	return &votersStorage{collection}, nil
}

func (self *votersStorage) tryVote(competitionCode, votersGroup, voter,
	competitorOne, competitorTwo string) (bool, error) {
	unitedCompetitorsId := buildUnitedId(competitorOne, competitorTwo)
	err := self.client.Insert(bson.M{
		"competition":    competitionCode,
		"voters_group":   votersGroup,
		"voter":          voter,
		"competitors_id": unitedCompetitorsId,
	})
	if err != nil {
		if mgo.IsDup(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (self *votersStorage) haveSeenPair(competitionCode, votersGroup, competitorOne, competitorTwo string) (bool, error) {
	unitedCompetitorsId := buildUnitedId(competitorOne, competitorTwo)
	rows, err := self.client.Find(bson.M{
		"competition":    competitionCode,
		"voters_group":   votersGroup,
		"competitors_id": unitedCompetitorsId,
	}).Count()
	if err != nil {
		return false, err
	}
	return rows == 1, err
}

func buildUnitedId(ids ...string) string {
	sort.Strings(ids)
	return strings.Join(ids, " | ")
}
