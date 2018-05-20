package models

import (
	"github.com/gazoon/go-utils"
	"instarate/libs/competition"
)

var (
	DefaultVotingTimeout = 5
	DefaultLang          = "ru"
)

type Chat struct {
	Id                    int    `bson:"chat_id"`
	MembersNum            int    `bson:"members_number"`
	IsGroupChat           bool   `bson:"is_group_chat"`
	LastMatch             *Match `bson:"last_match"`
	CurrentTopOffset      int    `bson:"current_top_offset"`
	CreatedAt             int    `bson:"created_at"`
	CompetitionCode       string `bson:"competition_code"`
	SelfActivationAllowed bool   `bson:"self_activation_allowed"`
	VotingTimeout         int    `bson:"voting_timeout"`
	Language              string `bson:"language"`
}

func NewChat(chatId, membersNum int, isGroupChat bool) *Chat {
	return &Chat{
		Id:                    chatId,
		MembersNum:            membersNum,
		IsGroupChat:           isGroupChat,
		CreatedAt:             utils.TimestampMilliseconds(),
		CompetitionCode:       competition.GlobalCompetition,
		SelfActivationAllowed: true,
		VotingTimeout:         DefaultVotingTimeout,
		Language:              DefaultLang,
	}
}
