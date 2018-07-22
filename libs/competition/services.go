package competition

import (
	. "instarate/libs/competition/config"
)

func InitCompetitorsStorage() *CompetitorsStorage {
	competitors, err := newCompetitorsStorage(Config.MongoCompetitors)
	if err != nil {
		panic(err)
	}
	return competitors
}

func InitProfilesStorage() *ProfilesStorage {
	profiles, err := newProfilesStorage(Config.MongoProfiles)
	if err != nil {
		panic(err)
	}
	return profiles
}

func InitGoogleFilesStorage() *GoogleFilesStorage {
	filesStorage, err := newGoogleStorage(Config.GoogleStorage.BucketName)
	if err != nil {
		panic(err)
	}
	return filesStorage
}

func InitVotersStorage() *VotersStorage {
	voters, err := newVotersStorage(Config.MongoVoters)
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
