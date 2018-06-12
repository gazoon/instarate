package models

import (
	"github.com/gazoon/go-utils"
	"time"
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
