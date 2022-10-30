package wp

import (
	"context"
	"github.com/myOmikron/echotools/worker"

	"github.com/myOmikron/bbb-rec-controller/models"
)

type w struct {
	queue chan worker.Task
	quit  chan bool
}

func (w *w) SetQueue(c chan worker.Task) {
	w.queue = c
}

func (w *w) Start() {
	ctx := context.WithValue(context.Background(), "conn", "todo!")
	for {
		select {
		case <-w.quit:
			return
		case t := <-w.queue:
			t.ExecuteWithContext(ctx)
		}
	}
}

func (w *w) Stop() {
	go func() {
		w.quit <- true
	}()
}

func CreateSeleniumWorker(conf *models.Config) func() (worker.Worker, error) {
	return func() (worker.Worker, error) {
		return &w{
			quit: make(chan bool),
		}, nil
	}
}
