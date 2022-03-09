package model

import (
	"GO/nomad/app/view"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	discordSession *discordgo.Session
	noticeDuration time.Duration
	scrapper *scrapper
	viewer *view.Viewer
}

// Create new bot object
func NewBot(noticeDuration time.Duration) *bot {
	var err error

	b := new(bot)

	b.noticeDuration = noticeDuration

	b.viewer = view.NewViewer(b.discordSession)

	token := os.Getenv("TOKEN")
	b.discordSession, err = discordgo.New("Bot " + token)
	if err != nil {
		b.viewer.FatallnErrorToConsole(err)
	}

	err = b.discordSession.Open()
	if err != nil {
		b.viewer.FatallnErrorToConsole(err)
	}

	b.scrapper = NewScrapper(b.discordSession)

	loginMsg := b.discordSession.State.User.String() + " (" + b.discordSession.State.User.Username + ")에 로그인 되었습니다.\n"
	b.viewer.PrintlnMsgToConsole(loginMsg)

	return b
}

// Send information to specified channel
func (b *bot) RunAlarm() {
	ticker := time.NewTicker(b.noticeDuration)
	go func() {
		for t := range ticker.C {
			b.viewer.PrintlnTimeToConsole(t)

			channelIds := b.getChannelIds()
			lastIndexedData := b.getLastArticleNoOfData()
			infoSet := b.scrapper.getInfoData(lastIndexedData)

			b.viewer.SendInfoToChannels(&channelIds, &infoSet)
		}
	}()
}

// Get last articleNo of article
func (b *bot) getLastArticleNoOfData() []string {
	var lastIndexData []string
	flag := false

	channelIds := b.getChannelIds()

	for _, channelId := range channelIds {
		msgs, _ := b.discordSession.ChannelMessages(channelId, 3, "", "", "")
		for _, msg := range msgs {
			if msg.Content[0] == '$' {
				lastIndex := msg.Content[1:]
				lastIndexData = strings.Split(lastIndex, " ")
				flag = true
				break
			}
		}

		if flag {
			break
		}
	}

	return lastIndexData[:4]
}

// Get channel ids from guilds
func (b *bot) getChannelIds() []string {
	targetedChannelName := "컴공과-공지사항"
	var channelIds []string

	guildIds := b.getGuildIds()

	for _, guildId := range guildIds {
		channels, _ := b.discordSession.GuildChannels(guildId)
		for _, channel := range channels {
			if (channel.Name == targetedChannelName) {
				channelIds = append(channelIds, channel.ID)
			}
		}
	}

	return channelIds
}

// Get guild ids from guilds
func (b *bot) getGuildIds() []string {
	var guildIds []string

	for _, guild := range b.discordSession.State.Guilds {
		guildIds = append(guildIds, guild.ID)
	}

	return guildIds
}