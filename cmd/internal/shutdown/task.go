// Package shutdown defines a graceful shutdown behavior
package shutdown

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/financial_advisor/app/config"
	"github.com/sirupsen/logrus"
)

var (
	isShuttingDown atomic.Bool
	singleton      *ShutdownTask
	once           sync.Once
	ErrAborted     = errors.New("aborted due to shutdown")
)

type ShutdownTask struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	tasks     []Task
}

type Task interface {
	Shutdown() error
	Name() string
}

func Init(ctx context.Context, cacel context.CancelFunc) {
	once.Do(func() {
		singleton = &ShutdownTask{
			tasks:     []Task{},
			ctx:       ctx,
			ctxCancel: cacel,
		}
	})
}

func (t *ShutdownTask) Add(task ...Task) {
	t.tasks = append(t.tasks, task...)
}

func (t *ShutdownTask) Shutdown() {
	t.ctxCancel()
}

func (t *ShutdownTask) ShutdownTasks() {
	// helps to prevent unnecessary resources loading
	// and setting the readiness probe to false
	isShuttingDown.Store(true)

	logrus.Infoln("waiting for readiness drain delay...")

	// wait for readiness drain delay to let load balancer know this instance is going down
	time.Sleep(config.ReadinessDrainDelay)

	for i := range t.tasks {
		if t.tasks[i] == nil {
			continue
		}

		logrus.WithField("task", t.tasks[i].Name()).Infoln("shutting down...")

		if err := t.tasks[i].Shutdown(); err != nil {
			logrus.WithField("task", t.tasks[i].Name()).WithError(err).Errorln("shutdown")
		}
	}
}

func Get() *ShutdownTask {
	return singleton
}

func (t *ShutdownTask) WaitForSignals() {
	<-t.ctx.Done()

	t.ctxCancel()
}

func (t *ShutdownTask) IsShuttingDown() bool {
	return isShuttingDown.Load()
}
