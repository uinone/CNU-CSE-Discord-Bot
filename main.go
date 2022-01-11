package main

import (
	"GO/nomad/utility"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	lastIndexData []string = utility.GetLastIndexData()
	discordSession *discordgo.Session
)

func init() {
	var err error

	err = godotenv.Load()
	utility.CheckErr(err)
	
	token := os.Getenv("TOKEN")
	discordSession, err = discordgo.New("Bot " + token)
	utility.CheckErr(err)
	

	err = discordSession.Open()
	utility.CheckErr(err)

	fmt.Printf("%s (%s)에 로그인 되었습니다.\n", discordSession.State.User.String(), discordSession.State.User.Username)
}

func main() {
	defer discordSession.Close()

	utility.SendScrappedData(discordSession, lastIndexData)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}