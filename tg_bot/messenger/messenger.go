package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/telegram-bot-api.v4"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

var (
	httpTimeout = time.Second * 3
	apiURL      = "https://api.telegram.org/bot"
)

type messageResponse struct {
	Ok          bool              `json:"ok"`
	Result      *tgbotapi.Message `json:"result"`
	ErrorCode   int               `json:"error_code"`
	Description string            `json:"description"`
}

type Telegram struct {
	botAPI *tgbotapi.BotAPI
	token  string
	client *http.Client
}

func NewTelegram(token string) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Telegram{
		botAPI: bot,
		token:  token,
		client: &http.Client{Timeout: httpTimeout},
	}, nil
}

func (self *Telegram) SendText(ctx context.Context, chatId int, text string,
	opts ...func(msg *tgbotapi.MessageConfig)) (int, error) {
	msg := tgbotapi.NewMessage(int64(chatId), text)
	for _, o := range opts {
		o(&msg)
	}
	return self.send(msg)
}

func (self *Telegram) SendMarkdown(ctx context.Context, chatId int, text string,
	opts ...func(msg *tgbotapi.MessageConfig)) (int, error) {
	opts = append(opts, func(msg *tgbotapi.MessageConfig) {
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
	})
	return self.SendText(ctx, chatId, text, opts...)
}

func (self *Telegram) GetMembersNum(ctx context.Context, chatId int) (int, error) {
	c := tgbotapi.ChatConfig{ChatID: int64(chatId)}
	return self.botAPI.GetChatMembersCount(c)
}

func (self *Telegram) SendNotification(ctx context.Context, callbackId, text string) error {
	_, err := self.botAPI.AnswerCallbackQuery(
		tgbotapi.CallbackConfig{CallbackQueryID: callbackId, Text: text})
	return err
}

func (self *Telegram) AnswerCallback(ctx context.Context, callbackId string) error {
	_, err := self.botAPI.AnswerCallbackQuery(
		tgbotapi.CallbackConfig{CallbackQueryID: callbackId})
	return err
}

func (self *Telegram) SendBinaryPhoto(ctx context.Context, chatId int,
	fileData io.Reader, opts ...func(config *tgbotapi.PhotoConfig)) (int, string, error) {
	c := tgbotapi.PhotoConfig{}
	for _, o := range opts {
		o(&c)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", uuid.NewV4().String())
	if err != nil {
		return 0, "", err
	}
	_, err = io.Copy(part, fileData)
	writer.WriteField("chat_id", strconv.Itoa(chatId))
	if c.Caption != "" {
		writer.WriteField("caption", c.Caption)
	}
	if c.ReplyMarkup != nil {
		replayMarkupData, err := json.Marshal(c.ReplyMarkup)
		if err != nil {
			return 0, "", err
		}

		writer.WriteField("reply_markup", string(replayMarkupData))
	}
	err = writer.Close()
	if err != nil {
		return 0, "", err
	}

	resp, err := self.client.Post(self.buildUrl("sendPhoto"), writer.FormDataContentType(), body)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	apiResponse := &messageResponse{}
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return 0, "", err
	}
	if !apiResponse.Ok || apiResponse.Result == nil {
		return 0, "", errors.Errorf("invalid sendPhoto response: %s", apiResponse.Description)
	}
	return retrievePhotoInfo(apiResponse.Result)
}

func (self *Telegram) SendPhoto(ctx context.Context, chatId int,
	photoUri string, opts ...func(config *tgbotapi.PhotoConfig)) (int, string, error) {

	c := tgbotapi.NewPhotoUpload(int64(chatId), photoUri)
	for _, o := range opts {
		o(&c)
	}
	resp, err := self.botAPI.Send(c)
	if err != nil {
		return 0, "", err
	}
	return retrievePhotoInfo(&resp)
}

func (self *Telegram) DeleteAttachedKeyboard(ctx context.Context, chatId,
	messageId int) (int, error) {
	c := tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{ChatID: int64(chatId), MessageID: messageId},
	}
	return self.send(c)
}

func (self *Telegram) send(c tgbotapi.Chattable) (int, error) {
	resp, err := self.botAPI.Send(c)
	if err != nil {
		return 0, err
	}
	return resp.MessageID, nil
}

func (self *Telegram) buildUrl(methodName string) string {
	return apiURL + self.token + "/" + methodName
}

func retrievePhotoInfo(resp *tgbotapi.Message) (int, string, error) {
	if resp.Photo == nil || len(*resp.Photo) == 0 {
		return 0, "", errors.New("send photo response without photos info")
	}
	fileId := (*resp.Photo)[len(*resp.Photo)-1].FileID
	return resp.MessageID, fileId, nil

}
