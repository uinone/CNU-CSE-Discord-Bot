package controller

import (
	"GO/nomad/app/model"
	"time"
)

type controller struct {
	Bot *model.Bot
	Web *model.Web
}

// Create controller object
func NewController() *controller {
	c := new(controller)

	c.Bot = model.NewBot()
	c.Web = model.NewWeb(c.Bot.GetDiscordSession())

	return c
}
 
// Run alarm of bot
func (c *controller) BotRun(alarmDuration time.Duration) {
	c.Bot.RunAlarm(alarmDuration)
}

// Run web router
func (c *controller) WebRun() {
	c.Web.Run()
}