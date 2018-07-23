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

func BuildWebhookPath(botToken string) string {
	return "/update/" + botToken
}

type ProcessUpdate func(ctx context.Context, update *tgbotapi.Update) error

type Webhook struct {
	*logging.LoggerMixin
	httpServer    *http.Server
	processUpdate ProcessUpdate
}

func New(port int, botToken string, processUpdate ProcessUpdate) *Webhook {
	logger := logging.NewLoggerMixin("webhook", nil)
	webhook := &Webhook{processUpdate: processUpdate, LoggerMixin: logger}
	r := httprouter.New()
	r.POST(BuildWebhookPath(botToken), webhook.updateHandler)
	webhook.httpServer = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r}
	return webhook
}

func (self *Webhook) Run() {
	log.WithField("address", self.httpServer.Addr).Info("Run webhook")
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

func (self *Webhook) updateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := utils.FillContext(r.Context())
	defer func() {
		if r := recover(); r != nil {
			self.handleError(ctx, w, r)
		}
	}()
	logger := self.GetLogger(ctx)

	update := &tgbotapi.Update{}
	err := json.NewDecoder(r.Body).Decode(update)
	if err != nil {
		logger.WithError(err).Error("Cannot parse http request into Update")
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}
	err = self.processUpdate(ctx, update)
	if err != nil {
		self.handleError(ctx, w, err)
		return
	}
}

func (self *Webhook) handleError(ctx context.Context, w http.ResponseWriter, err interface{}) {
	self.LogError(ctx, err)
	http.Error(w, "Internal error", http.StatusInternalServerError)
}
