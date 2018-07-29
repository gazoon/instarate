package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gazoon/go-utils/logging"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/telegram-bot-api.v4"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	httpTimeout = time.Second * 10
	apiURL      = "https://api.telegram.org/bot"
)

const (
	sendPhotoMethod = "sendPhoto"
)

type messageResponse struct {
	Ok          bool              `json:"ok"`
	Result      *tgbotapi.Message `json:"result"`
	ErrorCode   int               `json:"error_code"`
	Description string            `json:"description"`
}

type Telegram struct {
	*logging.LoggerMixin
	botAPI *tgbotapi.BotAPI
	token  string
	client *http.Client
}

func NewTelegram(token string) (*Telegram, error) {
	httpClient := &http.Client{Timeout: httpTimeout}
	bot, err := tgbotapi.NewBotAPIWithClient(token, httpClient)
	if err != nil {
		return nil, errors.Wrap(err, "Telegram API initialization")
	}
	return &Telegram{
		botAPI:      bot,
		token:       token,
		client:      httpClient,
		LoggerMixin: logging.NewLoggerMixin("messenger", nil),
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
	num, err := self.botAPI.GetChatMembersCount(c)
	return num, errors.Wrap(err, "get chat members tg API call")
}

func (self *Telegram) SendNotification(ctx context.Context, callbackId, text string) error {
	_, err := self.botAPI.AnswerCallbackQuery(
		tgbotapi.CallbackConfig{CallbackQueryID: callbackId, Text: text})
	return errors.Wrap(err, "answer callback tg API call, with text")
}

func (self *Telegram) AnswerCallback(ctx context.Context, callbackId string) error {
	_, err := self.botAPI.AnswerCallbackQuery(
		tgbotapi.CallbackConfig{CallbackQueryID: callbackId})
	return errors.Wrap(err, "answer callback tg API call")
}

func (self *Telegram) SendBinaryPhoto(ctx context.Context, chatId int,
	fileData io.Reader, opts ...func(settings *tgbotapi.PhotoConfig)) (int, string, error) {
	c := tgbotapi.PhotoConfig{}
	for _, o := range opts {
		o(&c)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", uuid.NewV4().String())
	if err != nil {
		return 0, "", errors.Wrap(err, "create form for file uploading")
	}
	_, err = io.Copy(part, fileData)
	writer.WriteField("chat_id", strconv.Itoa(chatId))
	if c.Caption != "" {
		writer.WriteField("caption", c.Caption)
	}
	if c.ReplyMarkup != nil {
		replayMarkupData, err := json.Marshal(c.ReplyMarkup)
		if err != nil {
			return 0, "", errors.Wrap(err, "serialize replay markup into json")
		}

		writer.WriteField("reply_markup", string(replayMarkupData))
	}
	err = writer.Close()
	if err != nil {
		return 0, "", errors.Wrap(err, "close multipart writer")
	}

	resp, err := self.client.Post(self.buildUrl("sendPhoto"), writer.FormDataContentType(), body)
	return handleSendPhotoResponse(resp, err)
}

func (self *Telegram) SendPhoto(ctx context.Context, chatId int,
	photoUri string, opts ...func(settings *tgbotapi.PhotoConfig)) (int, string, error) {

	c := tgbotapi.NewPhotoUpload(int64(chatId), photoUri)
	for _, o := range opts {
		o(&c)
	}
	data := url.Values{}
	data.Set("chat_id", strconv.Itoa(chatId))
	data.Set("photo", photoUri)
	if c.Caption != "" {
		data.Set("caption", c.Caption)
	}
	if c.ReplyMarkup != nil {
		replayMarkupData, err := json.Marshal(c.ReplyMarkup)
		if err != nil {
			return 0, "", errors.Wrap(err, "serialize replay markup into json")
		}
		data.Set("reply_markup", string(replayMarkupData))
	}

	self.GetLogger(ctx).WithField("photo_uri", photoUri).Info("Upload photo")
	resp, err := self.client.PostForm(self.buildUrl(sendPhotoMethod), data)
	return handleSendPhotoResponse(resp, err)
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
		return 0, errors.Wrapf(err, "send %T tg API call", c)
	}
	return resp.MessageID, nil
}

func (self *Telegram) buildUrl(methodName string) string {
	return apiURL + self.token + "/" + methodName
}

func handleSendPhotoResponse(resp *http.Response, err error) (int, string, error) {
	if err != nil {
		return 0, "", errors.Wrap(err, "send multipart post request")
	}
	defer resp.Body.Close()
	apiResponse := &messageResponse{}
	err = json.NewDecoder(resp.Body).Decode(apiResponse)
	if err != nil {
		return 0, "", errors.Wrap(err, "parse sendPhoto response into json")
	}
	if !apiResponse.Ok || apiResponse.Result == nil {
		return 0, "", errors.Errorf("invalid sendPhoto response: %s", apiResponse.Description)
	}
	return retrievePhotoInfo(apiResponse.Result)

}

func retrievePhotoInfo(resp *tgbotapi.Message) (int, string, error) {
	if resp.Photo == nil || len(*resp.Photo) == 0 {
		return 0, "", errors.New("send photo response without photos info")
	}
	fileId := (*resp.Photo)[len(*resp.Photo)-1].FileID
	return resp.MessageID, fileId, nil

}
