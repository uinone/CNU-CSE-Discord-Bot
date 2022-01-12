package utility

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	boardName = [4]string {
		"🎨 학사공지 🎨",
		"📜 일반소식 📜",
		"🔆 사업단소식 🔆",
		"🎈 취업정보 🎈 ※취업정보는 로그인해야 볼 수 있어요!😅",
	}

	contentPropertyName = [3]string {
		"[제목] ",
		"[링크] ",
		"[업로드 날짜] ",
	}
)

// Send information to specified channel
func SendInfoToChannel(ds *discordgo.Session, lastIndexData []string) {
	info := GetInfoData(ds, lastIndexData)

	SendMessageToChannel(ds, "모두 주목! 컴공과 공지 알림을 시작할게요🐧")

	var msg string
	for _, content := range info {
		SendMessageToChannel(ds, boardName[content.idx])

		if len(content.data) == 0 {
			SendMessageToChannel(ds, "새로 올라온 게시글이 없습니다.\n---")
		} else {
			for i, data := range content.data {
				if i == 0 {
					lastIndexData[content.idx] = strconv.Itoa(data.contentId)
				}
				msg = fmt.Sprint(msg, contentPropertyName[0])
				msg = fmt.Sprintln(msg, data.title)
				msg = fmt.Sprint(msg, contentPropertyName[1]) 
				msg = fmt.Sprintln(msg, data.link)
				msg = fmt.Sprint(msg, contentPropertyName[2]) 
				msg = fmt.Sprintln(msg, data.uploadedAt)
				msg = fmt.Sprintln(msg, "+")
				SendMessageToChannel(ds, msg)
			}
			SendMessageToChannel(ds, "---")
		}
	}

	SendMessageToChannel(ds, "업데이트가 완료됐어요!😀")

	UpdateLastIndexData(lastIndexData)

	fmt.Println("Task Complete!")
}

// Send message to targeted channel
func SendMessageToChannel(ds *discordgo.Session, msg string) {
	channelIds := GetChannelIds(ds)

	for _, channelId := range channelIds {
		ds.ChannelMessageSend(channelId, msg)
	}
}

// Get channel ids from guilds
func GetChannelIds(ds *discordgo.Session) []string {
	targetedChannelName := "bot-test"
	var channelIds []string

	guildIds := GetGuildIds(ds)

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
func GetGuildIds(ds *discordgo.Session) []string {
	var guildIds []string

	for _, guild := range ds.State.Guilds {
		guildIds = append(guildIds, guild.ID)
	}

	return guildIds
}