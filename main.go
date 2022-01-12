package main

import (
	"GO/nomad/utility"
	"os"
	"os/signal"
)

func main() {
	ds := utility.BotInit()
	defer ds.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}