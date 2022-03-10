package main

import (
	"GO/nomad/app/controller"
	"time"
)

func main() {
	controller := controller.NewController(time.Minute * 30)
	defer controller.Bot.GetDiscordSession().Close()

	controller.BotRun()
	controller.WebRun()
}