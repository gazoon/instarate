package models

import (
	"github.com/gazoon/go-utils"
	"instarate/libs/competition"
	"time"
)

var (
	MinVotingTimeoutSeconds = 5
	MinVotingTimeout        = time.Duration(MinVotingTimeoutSeconds) * time.Second
	RuLanguage              = "ru"
	EnLanguage              = "en"
	DefaultLang             = RuLanguage
)

type Chat struct {
	Id                    int       `bson:"chat_id"`
	MembersNum            int       `bson:"members_number"`
	IsGroupChat           bool      `bson:"is_group_chat"`
	LastMatch             *Match    `bson:"last_match"`
	CurrentTopOffset      int       `bson:"current_top_offset"`
	CreatedAt             time.Time `bson:"created_at"`
	CompetitionCode       string    `bson:"competition_code"`
	SelfActivationAllowed bool      `bson:"self_activation_allowed"`
	VotingTimeout         int       `bson:"voting_timeout"`
	Language              string    `bson:"language"`
}

func NewChat(chatId, membersNum int, isGroupChat bool) *Chat {
	return &Chat{
		Id:                    chatId,
		MembersNum:            membersNum,
		IsGroupChat:           isGroupChat,
		CreatedAt:             utils.UTCNow(),
		CompetitionCode:       competition.GlobalCompetition,
		SelfActivationAllowed: true,
		VotingTimeout:         MinVotingTimeoutSeconds,
		Language:              DefaultLang,
	}
}
