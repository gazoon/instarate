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
	CallbackType            = "callback"
	NextPairTaskType        = "next_pair_task"
	DailyActivationTaskType = "daily_activation_task"
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

func UserFromData(data interface{}) (*User, error) {
	u := &User{}
	err := mapstructure.Decode(data, u)
	if err != nil {
		return nil, errors.Wrap(err, "can't create user from data")
	}
	return u, nil
}

func (self User) String() string {
	return utils.ObjToString(&self)
}

type userMessage struct {
	ChatId      int   `mapstructure:"chat_id"`
	IsGroupChat bool  `mapstructure:"is_group_chat"`
	User        *User `mapstructure:"-"`
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
	callbackData, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("callback data must be a map, got: %v", callbackData)
	}
	user, err := UserFromData(callbackData["user"])
	if err != nil {
		return nil, errors.Wrap(err, "callback")
	}
	callback := &Callback{}
	err = mapstructure.Decode(callbackData, callback)
	if err != nil {
		return nil, errors.Wrap(err, "can't create callback from data")
	}
	callback.User = user
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
	Text         string       `mapstructure:"text"`
	MessageId    int          `mapstructure:"message_id"`
	ReplyTo      *TextMessage `mapstructure:"-"`
	wordsLowered mapset.Set
}

func TextMessageFromData(data interface{}) (*TextMessage, error) {
	messageData, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("text message data must be a map, got: %v", data)
	}
	user, err := UserFromData(messageData["user"])
	if err != nil {
		return nil, errors.Wrap(err, "text message")
	}

	replyToData := messageData["reply_to"]
	var replyTo *TextMessage
	if replyToData != nil {
		replyTo, err = TextMessageFromData(replyToData)
		if err != nil {
			return nil, errors.Wrap(err, "reply to")
		}
	}
	message := &TextMessage{}
	err = mapstructure.Decode(messageData, message)
	if err != nil {
		return nil, errors.Wrap(err, "can't create text message from data")
	}
	message.User = user
	message.ReplyTo = replyTo
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
	if self.ReplyTo == nil {
		return false
	}
	return self.ReplyTo.User.Username == botUsername
}

func (self *TextMessage) IsAppealToBot(botName string) bool {
	return strings.Contains(strings.ToLower(self.Text), strings.ToLower(botName))
}

func (self *TextMessage) GetCommandArg() string {
	args := self.GetCommandArgs()
	if len(args) == 0 {
		return ""
	}
	return args[len(args)-1]
}

func (self *TextMessage) GetCommandArgs() []string {
	var tokens []string
	for _, token := range strings.Split(self.Text, " ") {
		if token != "" {
			tokens = append(tokens, token)
		}
	}
	if len(tokens) == 0 {
		return []string{}
	}
	return tokens[1:]
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
