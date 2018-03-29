package webhook

import (
	"context"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"instarate/tg_gateway/worker"
	"net/http"
	"time"
)

type Webhook struct {
	httpServer      *http.Server
	botTokenToQueue map[string]string
	worker          *worker.Worker
}

func New(port int, botTokenToQueue map[string]string, worker *worker.Worker) *Webhook {
	webhook := &Webhook{worker: worker, botTokenToQueue: botTokenToQueue}
	r := httprouter.New()
	r.POST("/update/:bot_token", webhook.updateHandler)
	webhook.httpServer = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r}
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
	log.Info(ps)
}
