package utility

import "github.com/bwmarrin/discordgo"

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