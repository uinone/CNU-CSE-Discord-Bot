package utility

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Discord bot initialization
func BotInit() *discordgo.Session {
	var err error
	var discordSession *discordgo.Session

	//_ = godotenv.Load()
	token := os.Getenv("TOKEN")
	discordSession, err = discordgo.New("Bot " + token)
	checkErr(err)
	

	err = discordSession.Open()
	checkErr(err)

	fmt.Printf("%s (%s)에 로그인 되었습니다.\n", discordSession.State.User.String(), discordSession.State.User.Username)

	return discordSession
}

// Run alarm regularly
func RunAlarm(ds *discordgo.Session, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for t := range ticker.C {
			fmt.Println(t)
			SendInfoToChannel(ds)
		}
	}()
}

// Send information to specified channel
func SendInfoToChannel(ds *discordgo.Session) {
	infoSet := getInfoData(ds)

	if len(infoSet) > 1 {
		sendMessageToChannel(ds, "모두 주목! 컴공과 공지 알림을 시작할게요🐧")

		for _, info := range infoSet {
			for _, msg := range info {
				sendMessageToChannel(ds, msg)
			}
		}

		sendMessageToChannel(ds, "업데이트가 완료됐어요!😀")
	}

	fmt.Println("Task Complete!")
}

// Send message to targeted channel
func sendMessageToChannel(ds *discordgo.Session, msg string) {
	channelIds := getChannelIds(ds)

	for _, channelId := range channelIds {
		ds.ChannelMessageSend(channelId, msg)
	}
}

// Get channel ids from guilds
func getChannelIds(ds *discordgo.Session) []string {
	targetedChannelName := "컴공과-공지사항"
	var channelIds []string

	guildIds := getGuildIds(ds)

	for _, guildId := range guildIds {
		channels, _ := ds.GuildChannels(guildId)
		for _, channel := range channels {
			if (channel.Name == targetedChannelName) {
				channelIds = append(channelIds, channel.ID)
			}
		}
	}

	return channelIds
}

// Get guild ids from guilds
func getGuildIds(ds *discordgo.Session) []string {
	var guildIds []string

	for _, guild := range ds.State.Guilds {
		guildIds = append(guildIds, guild.ID)
	}

	return guildIds
}

func getLastIndexData(ds *discordgo.Session) []string {
	var lastIndexData []string
	flag := false

	channelIds := getChannelIds(ds)
	for _, channelId := range channelIds {
		msgs, _ := ds.ChannelMessages(channelId, 3, "", "", "")
		for _, msg := range msgs {
			if msg.Content[0] == '$' {
				lastIndex := msg.Content[1:]
				lastIndexData = strings.Split(lastIndex, " ")
				flag = true
			}
		}

		if flag {
			break
		}
	}

	return lastIndexData
}