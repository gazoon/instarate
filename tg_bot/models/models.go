package models

import (
	"github.com/gazoon/go-utils"
	"instarate/libs/competition"
	"time"
)

var (
	DefaultTimeoutSeconds = 5
	DefaultTimeout        = time.Duration(DefaultTimeoutSeconds) * time.Second
	CancelKeyboardAfter   = time.Duration(20) * time.Second

	RuLanguage  = "ru"
	EnLanguage  = "en"
	DefaultLang = RuLanguage
)

type Match struct {
	MessageId         int       `bson:"message_id"`
	LeftGirlUsername  string    `bson:"left_girl"`
	RightGirlUsername string    `bson:"right_girl"`
	ShownAt           time.Time `bson:"shown_at"`
}

func NewMatch(messageId int, leftGirlUsername, rightGirlUsername string) *Match {
	return &Match{
		MessageId:         messageId,
		LeftGirlUsername:  leftGirlUsername,
		RightGirlUsername: rightGirlUsername,
		ShownAt:           utils.UTCNow(),
	}
}

//TopRecord model is needed to implement show top with timeouts,
// like the voting process.
//type TopRecord struct {
//	MessageId int       `bson:"message_id"`
//	ShownAt   time.Time `bson:"shown_at"`
//}
//
//func NewTopRecord(messageId int) *TopRecord {
//	return &TopRecord{
//		MessageId: messageId,
//		ShownAt:   utils.UTCNow(),
//	}
//}

type Chat struct {
	Id          int    `bson:"chat_id"`
	MembersNum  int    `bson:"members_number"`
	IsGroupChat bool   `bson:"is_group_chat"`
	LastMatch   *Match `bson:"last_match"`
	//LastTopRecord         *TopRecord `bson:"top_record"`
	CurrentTopOffset      int       `bson:"current_top_offset"`
	CreatedAt             time.Time `bson:"created_at"`
	CompetitionCode       string    `bson:"competition_code"`
	SelfActivationAllowed bool      `bson:"self_activation_allowed"`
	VotingTimeout         int       `bson:"voting_timeout"`
	ShowTopTimeout        int       `bson:"top_timeout"`
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
		VotingTimeout:         DefaultTimeoutSeconds,
		ShowTopTimeout:        DefaultTimeoutSeconds,
		Language:              DefaultLang,
	}
}

func (self *Chat) GetVotingTimeout() int {
	if self.IsGroupChat {
		return self.VotingTimeout
	}
	return 0
}

func (self *Chat) GetShowTopTimeout() int {
	if self.IsGroupChat {
		return self.ShowTopTimeout
	}
	return 0
}

func (self *Chat) ResetTopOffset() {
	self.CurrentTopOffset = 0
}

func (self *Chat) ResetLastMatch() {
	self.LastMatch = nil
}
