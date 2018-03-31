package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"instarate/tg_gateway/worker"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/go-utils"
	"github.com/gazoon/go-utils/logging"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/telegram-bot-api.v4"
)

type Webhook struct {
	*logging.LoggerMixin
	httpServer      *http.Server
	botTokenToQueue map[string]string
	worker          *worker.Worker
}

func New(port int, botTokenToQueue map[string]string, worker *worker.Worker) *Webhook {
	logger := logging.NewLoggerMixin("webhook", nil)
	webhook := &Webhook{worker: worker, botTokenToQueue: botTokenToQueue, LoggerMixin: logger}
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
	botToken := ps.ByName("bot_token")
	queueName, ok := self.botTokenToQueue[botToken]
	if !ok {
		http.NotFound(w, r)
		return
	}
	ctx := utils.CreateContext()
	logger := self.GetLogger(ctx)
	update := &tgbotapi.Update{}
	err := json.NewDecoder(r.Body).Decode(update)
	if err != nil {
		logger.Errorf("Cannot parse http request into Update: %s", err)
		http.Error(w, "Bad request data", http.StatusBadRequest)
		return
	}
	err = self.worker.ProcessUpdate(ctx, queueName, update)
	if err != nil {
		logger.Errorf("Process update: %s", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
