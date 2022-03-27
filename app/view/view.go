package view

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Viewer struct {
	discordSession *discordgo.Session
}

// Create new Viewer object
func NewViewer() *Viewer {
	v := new(Viewer)

	return v
}

// Set discordSession of Viewer instance
func (v *Viewer) SetDiscordSession(ds *discordgo.Session) {
	v.discordSession = ds
}

// Send information to specified channel
func (v *Viewer) SendInfoToChannels(channelIds *[]string, infoSet *[][]string) {
	if len(*infoSet) > 1 {
		v.SendMessageToChannels(channelIds, "모두 주목! 컴공과 공지 알림을 시작할게요🐧")

		for _, info := range *infoSet {
			for _, msg := range info {
				v.SendMessageToChannels(channelIds, msg)
			}
		}

		v.SendMessageToChannels(channelIds, "업데이트가 완료됐어요!😀")
	}
}

// Send message to targeted channel
func (v *Viewer) SendMessageToChannels(channelIds *[]string, msg string) {
	for _, channelId := range *channelIds {
		v.discordSession.ChannelMessageSend(channelId, msg)
	}
}

// Print msg to console with line change
func (v *Viewer) PrintlnMsgToConsole(msg string) {
	fmt.Println(msg)
}

// Print msg to console with line change
func (v *Viewer) PrintlnTimeToConsole(t time.Time) {
	fmt.Println(t)
}

// Print error to console
func (v *Viewer) PrintlnErrorToConsole(err error) {
	fmt.Println(err)
}

// Print error type log to console and exit
func (v *Viewer) FatallnErrorToConsole(err error) {
	log.Fatalln(err)
}

// Print string type log to console and exit
func (v *Viewer) FatallnMsgToConsole(msg string) {
	log.Fatalln(msg)
}