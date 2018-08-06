package messages

import (
	"github.com/deckarep/golang-set"
	"github.com/gazoon/go-utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var (
	TextType                = "text"
	NewChatUsersType        = "new_chat_users"
	CallbackType            = "callback"
	NextPairTaskType        = "next_pair_task"
	DailyActivationTaskType = "daily_activation_task"
	CancelKeyboardTaskType  = "cancel_keyboard_task"
)

type Message interface {
	GetChatId() int
}

type UserMessage interface {
	Message
	GetIsGroupChat() bool
	GetUser() *User
}

type User struct {
	Id       int    `mapstructure:"id"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
}

func (self User) String() string {
	return utils.ObjToString(&self)
}

type userMessage struct {
	ChatId      int   `mapstructure:"chat_id"`
	IsGroupChat bool  `mapstructure:"is_group_chat"`
	User        *User `mapstructure:"user"`
}

func (self userMessage) GetChatId() int       { return self.ChatId }
func (self userMessage) GetIsGroupChat() bool { return self.IsGroupChat }
func (self userMessage) GetUser() *User       { return self.User }

var callbackPayloadSeparator = ":"

type Callback struct {
	userMessage `mapstructure:",squash"`
	CallbackId  string `mapstructure:"callback_id"`
	ParentMsgId int    `mapstructure:"parent_msg_id"`
	Payload     string `mapstructure:"payload"`
}

func CallbackFromData(data interface{}) (*Callback, error) {
	callback := &Callback{}
	err := mapstructure.Decode(data, callback)
	if err != nil {
		return nil, errors.Wrap(err, "can't create callback from data")
	}
	return callback, nil
}

func (self Callback) String() string {
	return utils.ObjToString(&self)
}

func (self *Callback) GetName() string {
	return self.splitPayload()[0]
}

func (self *Callback) GetArgs() string {
	tokens := self.splitPayload()
	return tokens[len(tokens)-1]
}

func (self *Callback) splitPayload() []string {
	return strings.Split(self.Payload, callbackPayloadSeparator)
}

func BuildCallbackPayload(callbackName, args string) string {
	return callbackName + callbackPayloadSeparator + args
}

type TextMessage struct {
	userMessage  `mapstructure:",squash"`
	Text         string     `mapstructure:"text"`
	MessageId    int        `mapstructure:"message_id"`
	ReplyToUser  *User      `mapstructure:"reply_to_user"`
	wordsLowered mapset.Set `mapstructure:"-"`
}

func TextMessageFromData(data interface{}) (*TextMessage, error) {
	message := &TextMessage{}
	err := mapstructure.Decode(data, message)
	if err != nil {
		return nil, errors.Wrap(err, "can't create text message from data")
	}
	return message, nil
}

func (self TextMessage) String() string {
	return utils.ObjToString(&self)
}

func (self *TextMessage) TextContains(phrase string) bool {
	phraseWords := utils.SplitWordsLowered(phrase)
	messageWords := self.GetWordsLowered()
	for _, w := range phraseWords {
		if !messageWords.Contains(w) {
			return false
		}
	}
	return true
}

func (self *TextMessage) GetWordsLowered() mapset.Set {
	if self.wordsLowered != nil {
		return self.wordsLowered
	}
	self.wordsLowered = mapset.NewSet()
	words := utils.SplitWordsLowered(self.Text)
	for _, w := range words {
		self.wordsLowered.Add(w)
	}

	return self.wordsLowered
}

func (self *TextMessage) IsReplyToBot(botUsername string) bool {
	if self.ReplyToUser == nil {
		return false
	}
	return self.ReplyToUser.Username == botUsername
}

func (self *TextMessage) IsAppealToBot(botName string) bool {
	return strings.Contains(strings.ToLower(self.Text), strings.ToLower(botName))
}

func (self *TextMessage) GetCommandArg() string {
	tokens := strings.Fields(self.Text)
	if len(tokens) < 2 {
		return ""
	}
	return tokens[len(tokens)-1]
}

func (self *TextMessage) GetLastWord() string {
	tokens := strings.Fields(self.Text)
	if len(tokens) == 0 {
		return ""
	}
	return tokens[len(tokens)-1]
}

type NewChatUsers struct {
	userMessage `mapstructure:",squash"`
	MessageId   int     `mapstructure:"message_id"`
	NewUsers    []*User `mapstructure:"new_users"`
}

func NewChatUsersFromData(data interface{}) (*NewChatUsers, error) {
	message := &NewChatUsers{}
	err := mapstructure.Decode(data, message)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new chat users message from data")
	}
	return message, nil
}

func (self *NewChatUsers) IsBotAdded(botUsername string) bool {
	for _, u := range self.NewUsers {
		if u.Username == botUsername {
			return true
		}
	}
	return false
}

type task struct {
	ChatId int       `mapstructure:"chat_id"`
	DoAt   time.Time `mapstructure:"do_at"`
}

func (self *task) GetChatId() int {
	return self.ChatId
}

type NextPairTask struct {
	task               `mapstructure:",squash"`
	LastMatchMessageId int `mapstructure:"last_match_message_id"`
}

func NextPairTaskFromData(data interface{}) (*NextPairTask, error) {
	t := &NextPairTask{}
	err := mapstructure.Decode(data, t)
	if err != nil {
		return nil, errors.Wrap(err, "can't create next pair task from data")
	}
	return t, nil
}

func (self NextPairTask) String() string {
	return utils.ObjToString(&self)
}

type DailyActivationTask struct {
	task `mapstructure:",squash"`
}

func DailyActivationTaskFromData(data interface{}) (*DailyActivationTask, error) {
	t := &DailyActivationTask{}
	err := mapstructure.Decode(data, t)
	if err != nil {
		return nil, errors.Wrap(err, "can't create daily activation task from data")
	}
	return t, nil
}

func (self DailyActivationTask) String() string {
	return utils.ObjToString(&self)
}

type CancelKeyboardTask struct {
	task `mapstructure:",squash"`
}

func CancelKeyboardTaskFromData(data interface{}) (*CancelKeyboardTask, error) {
	t := &CancelKeyboardTask{}
	err := mapstructure.Decode(data, t)
	if err != nil {
		return nil, errors.Wrap(err, "can't create cancel keyboard task from data")
	}
	return t, nil
}

func (self CancelKeyboardTask) String() string {
	return utils.ObjToString(&self)
}
