package competition

import (
	. "instarate/libs/competition/config"
)

func InitCompetitorsStorage() *CompetitorsStorage {
	competitors, err := NewCompetitorsStorage(Config.MongoCompetitors)
	if err != nil {
		panic(err)
	}
	return competitors
}

func InitProfilesStorage() *ProfilesStorage {
	profiles, err := NewProfilesStorage(Config.MongoProfiles)
	if err != nil {
		panic(err)
	}
	return profiles
}

func InitGoogleFilesStorage() *GoogleFilesStorage {
	filesStorage, err := NewGoogleStorage(Config.GoogleStorage.BucketName)
	if err != nil {
		panic(err)
	}
	return filesStorage
}

func InitVotersStorage() *VotersStorage {
	voters, err := NewVotersStorage(Config.MongoVoters)
	if err != nil {
		panic(err)
	}
	return voters
}

func InitCompetition() *Competition {
	competitors := InitCompetitorsStorage()
	profiles := InitProfilesStorage()
	filesStorage := InitGoogleFilesStorage()
	voters := InitVotersStorage()
	return NewCompetition(competitors, profiles, filesStorage, voters)
}
