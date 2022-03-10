package controller

import (
	"GO/nomad/app/model"
	"time"
)

type controller struct {
	Bot *model.Bot
	Web *model.Web
}

func NewController(noticeDuration time.Duration) *controller {
	c := new(controller)

	c.Bot = model.NewBot(noticeDuration)
	c.Web = model.NewWeb(c.Bot.GetDiscordSession())

	return c
}

func (c *controller) BotRun() {
	c.Bot.RunAlarm()
}

func (c *controller) WebRun() {
	c.Web.Run()
}