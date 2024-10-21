package timekeeper

import (
	"time"

	"github.com/kshyr/pact/config"
	"github.com/kshyr/pact/invoker"
	"github.com/kshyr/pact/script"
	"github.com/robfig/cron/v3"
)

type Timekeeper struct {
	cron    *cron.Cron
	scripts []script.Script
	cfg     config.Config
}

func New(scripts []script.Script, cfg config.Config) *Timekeeper {
	return &Timekeeper{
		cron:    cron.New(cron.WithSeconds()),
		scripts: scripts,
		cfg:     cfg,
	}
}

func (t *Timekeeper) Start() {
	for _, script := range t.scripts {
		if script.Active {
			t.scheduleScript(script)
		}
	}
	t.cron.Start()
}

func (t *Timekeeper) Stop() {
	t.cron.Stop()
}

func (t *Timekeeper) scheduleScript(script script.Script) {
	if script.Schedule != "" {
		t.cron.AddFunc(script.Schedule, func() {
			invoker.InvokeScript(script, t.cfg)
		})
	} else if !script.At.IsZero() {
		delay := time.Until(script.At)
		time.AfterFunc(delay, func() {
			invoker.InvokeScript(script, t.cfg)
		})
	}
}
