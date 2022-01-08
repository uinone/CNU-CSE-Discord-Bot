package main

import (
	"GO/nomad/utility"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	envData []string = utility.GetEnvData()
	token string = envData[4]
	discordSession *discordgo.Session
)

func init() {
	var err error
	discordSession, err = discordgo.New("Bot " + token)
	utility.CheckErr(err)

	err = discordSession.Open()
	utility.CheckErr(err)

	fmt.Printf("%s (%s)에 로그인 되었습니다.\n", discordSession.State.User.String(), discordSession.State.User.Username)
}

func main() {
	defer discordSession.Close()
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}