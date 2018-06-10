package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/telegram-bot-api.v4"
)

type ProcessUpdate func(ctx context.Context, queueName string, update *tgbotapi.Update) error

type Webhook struct {
	*logging.LoggerMixin
	httpServer      *http.Server
	botTokenToQueue map[string]string
	processUpdate   ProcessUpdate
}

func New(port int, botTokenToQueue map[string]string, processUpdate ProcessUpdate) *Webhook {
	logger := logging.NewLoggerMixin("webhook", nil)
	webhook := &Webhook{processUpdate: processUpdate, botTokenToQueue: botTokenToQueue, LoggerMixin: logger}
	r := httprouter.New()
	r.POST("/update/:bot_token", webhook.updateHandler)
	webhook.httpServer = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: utils.RecoveryHandler(r)}
	return webhook
}

func (self *Webhook) Run() {
	log.Infof("Run webhook on %s", self.httpServer.Addr)
	go func() {
		if err := self.httpServer.ListenAndServe(); err != nil {
			log.Panicf("Webhook server failed: %s", err)
		}
	}()
}

func (self *Webhook) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return self.httpServer.Shutdown(ctx)
}

func (self *Webhook) updateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := utils.FillContext(r.Context())
	defer func() {
		if r := recover(); r != nil {
			self.handleError(ctx, w, r)
		}
	}()
	logger := self.GetLogger(ctx)
	botToken := ps.ByName("bot_token")
	queueName, ok := self.botTokenToQueue[botToken]
	if !ok {
		logger.WithField("bot_token", botToken).Error("Unknown bot token")
		http.NotFound(w, r)
		return
	}
	update := &tgbotapi.Update{}
	err := json.NewDecoder(r.Body).Decode(update)
	if err != nil {
		logger.Errorf("Cannot parse http request into Update: %s", err)
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}
	err = self.processUpdate(ctx, queueName, update)
	if err != nil {
		self.handleError(ctx, w, err)
		return
	}
}

func (self *Webhook) handleError(ctx context.Context, w http.ResponseWriter, err interface{}) {
	self.LogError(ctx, err)
	http.Error(w, "Internal error", http.StatusInternalServerError)
}
