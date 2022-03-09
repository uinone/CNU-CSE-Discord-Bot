package controller

import (
	"GO/nomad/app/model"
	"time"
)

type controller struct {
	Bot *model.Bot
}

func NewController(noticeDuration time.Duration) *controller {
	c := new(controller)

	c.Bot = model.NewBot(noticeDuration)

	return c
}

func (c *controller) Run() {
	c.Bot.RunAlarm()
}