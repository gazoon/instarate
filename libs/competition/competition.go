package competition

import (
	"fmt"
	"github.com/gazoon/go-utils"
	"instarate/libs/competition/config"
	"path"
)

type Competition struct {
	competitors *competitorsStorage
}

func New() (*Competition, error) {

	conf := &config.Config{}
	configPath := path.Join(utils.GetCurrentFileDir(), "config")
	err := utils.ParseConfig(configPath, &conf)
	if err != nil {
		return nil, err
	}
	competitors, err := newCompetitorsStorage(conf.MongoCompetitors)
	if err != nil {
		return nil, err
	}
	return &Competition{competitors: competitors}, nil
}

func (self *Competition) Test() {
	x, err := self.competitors.getTop("global", 2, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(x[0], x[1])
	n, err := self.competitors.getCompetitorsNumber("global")
	if err != nil {
		panic(err)
	}
	fmt.Printf("number %d\n", n)
	m, err := self.competitors.get("global", "mirgaeva_galinka")
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
	higherRating, err := self.competitors.getNumberWithHigherRating("global", 1490)
	if err != nil {
		panic(err)
	}
	fmt.Printf("With higher rating: %d\n", higherRating)
	m.Wins = 228
	err = self.competitors.update(m)
	if err != nil {
		panic(err)
	}
	one, two, err := self.competitors.getRandomPair("global")
	fmt.Printf("one %v\n", one)
	fmt.Printf("two %v\n", two)
}
