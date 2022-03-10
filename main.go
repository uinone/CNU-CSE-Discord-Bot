package main

import (
	"GO/nomad/app/controller"
	"time"
)

func main() {
	controller := controller.NewController()
	defer controller.Bot.GetDiscordSession().Close()

	controller.BotRun(time.Minute * 30)
	controller.WebRun()
}