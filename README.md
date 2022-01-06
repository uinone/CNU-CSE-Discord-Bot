# CNU_Discod_Bot 프로젝트

[discordgo](https://github.com/bwmarrin/discordgo)와 [goquery](https://github.com/PuerkitoBio/goquery)를 활용했어요.

## 봇은 이렇게 동작해요😀

봇은 다음과 같이 동작해요.

1. 먼저 충남대학교 컴퓨터융합학부의 게시판들이 업데이트 되는 것을 주기적으로 확인합니다.
2. 해당 봇을 추가한 서버들의 특정 채팅채널에서 볼 수 있도록 업데이트 된 게시글의 정보를 웹에서 긁어옵니다.
3. 긁어온 데이터를 채팅채널에 메시지 형태로 보내줍니다.

## 이거 때문에 너무 힘들었어요😥

시간이 날 때마다 정리해서 올릴 예정이에요.

- [디스코드 봇 관련](#디스코드-봇-관련)
  - 계속 꺼지는 프로그램
  - 비어있는 session.State.Guilds.Channels

* [웹 크롤링 관련](#웹-크롤링-관련)
  - 날 죽이려는 whitespace
  - 이상한 날짜 변환

## 디스코드 봇 관련

디스코드 봇을 만들면서 겪었던 문제들과 해결 과정이에요.

### 계속 꺼지는 프로그램

이전에 Python이나 Node.js를 가지고 개발을 했었다.\
Python의 discord.py 라이브러리에는 run이라는 함수가 있었고, 이를 실행시키면 따로 프로그램을 끄기 전까지는 계속해서 돌아갔었다. Node.js의 discord.js도 마찬가지로 login이라는 함수가 같은 기능을 해줬다.

근데 discord.go의 Open 함수는 그렇게 동작하지 않고, 봇에 로그인한 다음 바로 프로그램이 종료되었다.

#### 해결 방법

os.Signal 타입 채널을 생성하고, Interrupt가 일어나면 바로 종료하도록 만들면 됐다.

```go
sc := make(chan os.Signal, 1)
signal.Notify(sc, os.Interrupt)
<-sc
```

### 비어있는 session.State.Guilds.Channels

봇이 정상적으로 동작하려면, 봇을 추가한 서버들의 리스트와 서버 각각이 가지고 있는 채널들의 정보가 필요했다.\
디스코드에서는 서버를 Guild라고 표현한다. 따라서 로그인한 세션에서 Guilds 변수를 찾아냈다.

Guilds는 봇이 추가된 서버들의 리스트를 담고있었고, 그 안에는 Channels라는 채널 포인터들이 들어있는 리스트가 있었다.

따라서 이를 통해 한 번에 정보를 얻을 수 있다고 생각했다.

```go
for _, guild := range session.State.Guilds {
  fmt.Println(guild.Channels)
}
```

근데 콘솔에는 빈 리스트만 출력됐다.. 그래서 코드를 뜯어보니 struct.go파일에 Channels와 관련된 코드를 찾을 수 있었다.

```go
// A list of channels in the guild.
// This field is only present in GUILD_CREATE events and websocket
// update events, and thus is only present in state-cached guilds.
Channels []*Channel `json:"channels"`
```

처음에는 저게 뭔소리인지 모르겠어서 discordgo 레포에서 많은 이슈들을 돌아봤다. 그리고 한 이슈를 발견하고 안 사실은 [Intent를 설정해야만 데이터가 보인다는 것이었다.](https://github.com/bwmarrin/discordgo/issues/812) 디스코드 개발자 페이지 + 코드작성 시 Intent 설정을 따로 해줘야하는 것을 알아서 뛸듯이 기뻤다.

하지만 디스코드 개발자 페이지에서 Intent 설정을 하고 코드에서 따로 Intent를 주어도 의미가 없었다. [0.23 버전](https://github.com/bwmarrin/discordgo/releases/tag/v0.23.0)에서는 Intent 설정이 강제여서 Intent 설정을 따로 해줄 필요는 없었기 때문이다.

이 모든 원인은 저 주석을 해석하지 못한 내 잘못이었다.

계속 생각해보니, 따로 저 데이터를 요청해야 받을 수 있다는 말 같았다. [이 이슈](https://github.com/bwmarrin/discordgo/issues/999)를 읽어보니 성능의 이유로 직접 API를 호출하지 않으면 데이터를 주지 않도록 짠것같았다. 따라서 직접 API를 호출하는 방향으로 코드를 짜서 해결했다.

#### 해결 방법

서버 아이디는 session.State.Guilds에서 가져온다.

```go
// Get guild ids from guilds
func GetGuildIds(ds *discordgo.Session) []string {
	var guildIds []string

	for _, guild := range ds.State.Guilds {
		guildIds = append(guildIds, guild.ID)
	}

	return guildIds
}
```

이후 session.GuildChannels 메서드를 사용해 특정 길드의 채널 데이터를 직접 요청하면 됐다.

```go
// Get channel ids from guilds
func GetChannelIds(ds *discordgo.Session) []string {
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
```
