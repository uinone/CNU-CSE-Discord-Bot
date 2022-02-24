# CNU_CSE_Discod_Bot 프로젝트

[discordgo](https://github.com/bwmarrin/discordgo), [goquery](https://github.com/PuerkitoBio/goquery) 그리고 [gin](https://github.com/gin-gonic/gin)을 활용했어요.

서버에 추가하고 싶으시면 [링크](https://discordapp.com/api/oauth2/authorize?client_id=658973905435754526&permissions=100352&scope=bot)를 클릭해주세요.

## 봇은 이렇게 동작해요😀

봇은 다음과 같이 동작해요.

1. 충남대학교 컴퓨터융합학부의 게시판들이 업데이트 되는 것을 30분마다 한번씩 24시간동안 확인합니다.
2. 해당 봇을 추가한 서버들의 특정 채팅채널(이름은 "컴공과-공지사항"여야해요)에서 볼 수 있도록 업데이트 된 게시글의 정보를 웹에서 긁어옵니다.
3. 긁어온 데이터를 채팅채널에 메시지 형태로 보내줍니다.

봇은 아무리 느려도 1시간 30분안에는 정보를 제공해요!

## 업데이트 목록

변경되는 부분을 정리했어요.

### 22-01-13 업로드 날짜 표시방법 변경

[@dayo2n](https://github.com/dayo2n)님이 피드백 해주셨어요.

봇이 변경사항을 올리는 시점을 기준으로 몇일 지났는지를 표시하면 좋겠다고 생각했는데\
보는 사람 입장에서는 그냥 날짜를 보여주는 게 더 나은 것 같아서 업로드 날짜 그대로 표시하기로 했다.

#### 변경 사항 정리

> 변경 전\
> [업로드 날짜]: N일전\
> \
> 변경 후\
> [업로드 날짜]: 22.01.13

## 이거 때문에 너무 힘들었어요😥

- [디스코드 봇 관련](#디스코드-봇-관련)
  - 계속 꺼지는 프로그램
  - 비어있는 session.State.Guilds.Channels
  - error closing resp body
  - 60초 안에 해보세요

* [웹 크롤링 관련](#웹-크롤링-관련)
  - 날 죽이려는 whitespace
  - 이상한 날짜 변환

- [Go 관련](#Go-관련)
  - Get "url" EOF

* [Heroku 관련](#Heroku-관련)
  - 서버가 죽고나면 벌어지는 일들

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

### error closing resp body

heroku에서 봇을 돌리는데 자꾸 짜증나게 하는 에러가 있었다.

> 2022/01/26 17:00:26 error closing resp body

response body를 닫는 것에 대해 진짜 모든 방법을 다 동원했었다.

defer를 사용해서 함수가 끝나기 전에 response body를 닫기도 했고,\
GET요청을 보내는 Client객체에 TCP통신을 계속 잡고있지 않도록 하는 조치도 했었다.

하지만 error closing resp body라니.. 나랑 싸우자는 줄 알았다.

그래서 열심히 찾아보던 중 [다음과 같은 이슈](https://github.com/bwmarrin/discordgo/issues/1028)를 발견할 수 있었다.

알고보니 go 1.17.3 버전의 오류였던 것이다.

#### 해결 방법

이슈에 나온 방법처럼 heroku의 go 버전을 1.17.2 혹은 1.17.4 이상의 버전을 사용하면 해결되는 문제였다.

아쉽게도 heroku의 최신 go 버전이 1.17.3이었기 때문에 1.17.2로 내리는 결정을 해서 오류를 해결했다.

### 60초 안에 해보세요

갑자기 봇이 죽었다.

몇주 잘 돌다가 죽으니까 뭐가 문제인지 모르겠어서 눈물이 났다.

그래서 한참을 콘솔에 찍으면서 알아봤지만, 오류는 잘 보이지 않았다.

커멘드 창에 `heroku logs -a 봇이름 --tail`을 치면 로그가 나온다.

이걸 살펴보니까 다음과 같은 로그가 찍혀있었다.

> 2022-02-17T13:18:44.622932+00:00 heroku[web.1]: Error R10 (Boot timeout) -> Web process failed to bind to $PORT within 60 seconds of launch

저 로그를 왜 지나쳤을까.

60초 안에 크롤링과 디스코드에 메시지 보내기를 같이 해야하는 상황이 발생한 것이다.

디스코드에 메시지를 한번에 많이 보내게 되면\
디스코드가 몇개 보내다가 잠시 쉬고 다시 보내는 과정을 거친다.

이때 잠시 쉬는 시간이 얼마나 될지를 모르는데, 체감상 최소 30초는 되는 것 같다.

한마디로 웹 서버를 키기 전에 서버에 정보를 보내려고 시도하는 것은 도박에 가까웠다.

그동안 잘 돌아갔던 이유는 매일 올라오는 게시글의 양이 적어서 였기 떄문이었다.

#### 해결 방법

문제가 발생한건 main.go에서 다음 한 줄 때문이었다.

```go
utility.SendInfoToChannel(ds) // <--
utility.RunAlarm(ds, time.Minute * 30)
```

나는 웹서버를 키기 전에 utility.RunAlarm을 실행시켜야했다.

즉 웹서버를 키고 30분 뒤에 크롤링 결과를 알 수 있었는데, 난 이게 너무 답답했다.

그래서 켜지자 마자 한번 크롤링하고 시작하도록 했는데, 이게 악수였다.

따라서 저 부분을 지움으로써 해결할 수 있었다.

---

## 웹 크롤링 관련

Golang을 가지고 웹 크롤링을 하면서 겪었던 문제들과 해결 과정이에요.

### 날 죽이려는 whitespace

웹 크롤링에서 특정 태그 안에 있는 값을 생각없이 바로 문자열 형태로 만들어보면,\
생각과는 다르게 굉장히 많은 탭과 줄바꿈을 동반한 문자열이 반환된다는 것을 알 수 있다.

Node.js에서는 따로 이런 공백(/n)과 탭(/t)을 지울 수 있는 방법을 제공하지 않았으므로\
정규표현식(RegExp)을 활용하여 이를 해결해야했지만, Go에서는 이스케이프 문자들을 제거하고\
오직 그 안에 들어있는 문자열만 반환하도록 하는 메서드가 있었다.

#### 해결 방법

strings.TrimSpace 메서드는 따로 정규표현식을 쓰지 않아도 알아서 이런 whitespace들을 제거해줬다.

```go
func cleanString(str string) string {
	return strings.TrimSpace(str)
}
```

### 이상한 날짜 변환

어느 언어를 사용하던, 문자열을 날짜로 변환하고 이를 계산하는 형태가 많이 있다.

봇은 오늘 날짜와 게시물을 날짜(혹은 게시글 번호)를 확인해서 업데이트된 게시물을 확인해야했다.\
봇을 맨 처음 실행하는 상태에서는 마지막으로 업데이트 된 게시물을 확인할 수 없기 때문이다.

따라서 날짜 변환이 꼭 필요하게 되었지만, 컴공과 홈페이지의 날짜 형식은 YY.MM.DD 형식이었다.

Go에서 문자열을 날짜 형식으로 변환하는 방법을 찾아보면 알겠지만, 시간까지 같이 명시해줘야하는 경우가 많다.\
그런 경우들을 많이 접한 나는 겁을 먹었지만, 다행히 저런 형식을 날짜 형식으로 변환시킬 방법이 있었다.

#### 해결 방법

우선 '.'를 통해 구분되어있는 날짜 형식을 '-'을 통해 구분되도록 바꾼다.\
이후 time.parse 메서드룰 통해 원하는 형식으로 바꿀 수 있었다.

```go
func changeTimeToDate(str string) time.Time {
	strDate := strings.Join(strings.Split(str, "."), "-")
	t, _ := time.Parse("06-01-02", strDate)
	return t
}
```

웃기게도 2006년 1월 2일이 아니면 '-'를 통해 구분된 형식으로 바꿀 수 없다.\
왜인지 살펴보고 있지만, 결국 찾지 못했다..

이렇게 날짜와 관련된 문제는 쉽게 끝나나 했지만..\
time.parse 메서드를 통해 변환된 날짜는 UTC 형식이었다.

게시판에 표시된 것은 시간을 같이 표시하지 않아서 상관없었지만,\
UTC는 KST보다 9시간이 느리기 때문에 자칫하면 새로운 게시글을 찾지 못하는 경우가 생길 수 있었다.

하지만 Add 메서드를 사용하면 UTC형식의 시간에서 9시간을 더해서 KST로 만들 수 있었다.

```go
now := time.Now().UTC().Add(time.Hour * 9)
```

---

## Go 관련

Golang을 잘 몰라서 일어난 일이에요.

### Get "url" EOF

웹 크롤링을 시도할 때, 명시한 url들에게 GET 요청을 보내서 페이지를 가져온다.\
이를 url은 총 4개고 각각에 대해 동시성을 부여하기 위해 다음과 같이 시도했다.

```go
results := make(chan msgData)
for i:=0; i<len(urls); i++ {
	contentId, _ := strconv.Atoi(envData[i])
	go getScrappedData(i, contentId, results)
}
```

동시에 4개의 url에게 GET요청을 보내고 가져오므로 기존에 4개를 순서대로 요청했을 때보다 속도를 향상시킬 수 있었다.

문제는 갑자기 일어나게 되었는데, 그 전에는 문제가 없다가 go 모듈을 좀 건들이고 난 후에 다음과 같은 에러가 발생했다.

> Get "URL": EOF

처음보는 에러였기 때문에 뭔가싶어 구글링을 해보니 [다음과 같은 질문](https://stackoverflow.com/questions/28046100/golang-http-concurrent-requests-post-eof)을 찾을 수 있었다.

4개의 url로의 GET 요청이 따로 연결을 시도하는 것이 아니라, 먼저 사용하고 있던 연결을 유지하면서 다음 요청을 시도한다는 것이었다.

나는 GET요청을 시도하고 반환받은 response의 body부분을 사용한 후, 함수가 종료되는 시점에 response.body만 Close를 시켜줬었다.

```go
res, err := http.Get(urls[idx])
CheckErr(err)
checkCode(res)

defer res.Body.Close()
```

물론 저 부분도 필요하지만, 요청이 끝나면 연결 자체도 같이 닫아야 했다.

즉, 맨 처음 GET요청을 보낼 때 사용했던 연결을 다음 GET요청이 그대로 사용하려고 하기 때문에 일어난 일이다.

첫 요청때 사용되던 함수가 종료되는 시점에 연결을 끊기 때문에, 다음에 시도하려는 GET 요청들은 기존 연결을 사용해 GET 요청을 하려하고, 요청 도중 연결이 끊겨 EOF 에러를 뱉는 것이었다.

```go
res.Request.Close = true
```

위 코드를 집어넣어서 다음 요청들이 기존 연결을 계속 사용하지 않도록 하면 해결될 줄 알았다.

request.Close를 뜯어보니 다음과 같은 주석이 있었다.

```go
// Close indicates whether to close the connection after
// replying to this request (for servers) or after sending this
// request and reading its response (for clients).
//
// For server requests, the HTTP server handles this automatically
// and this field is not needed by Handlers.
//
// For client requests, setting this field prevents re-use of
// TCP connections between requests to the same hosts, as if
// Transport.DisableKeepAlives were set.
Close bool
```

Close는 서버에게 요청에 대한 답변을 주거나\
클라이언트에게 요청을 보내고 응답을 읽고난 후에\
연결을 닫아야하는지 아닌지를 나타낸다.

클라이언트 요청에 대해, 이 필드를 설정하면\
마치 Transport.DisableKeepAlives 가 설정된 것처럼\
같은 호스트에 대한 요청들 사이 TCP 연결을 재사용하지 않도록 한다.

아쉽게도 저 코드 한 줄이 문제를 해결해주지 못했다.\
되는줄 알았는데, 운이 좋아서 가끔 동작했을 뿐 여전히 TCP 연결을 재사용하고 있었다.

그래서 [golang-korea](https://www.facebook.com/groups/golangko/about)에 질문을 올렸고 좋은 답변을 얻을 수 있었다.

해결하지 못한 이유는 request.Close라는 설정이 되지 않았기 때문이었다.

#### 해결 방법

물론 request.Close라는 설정이 필요한것도 맞았다.

하지만 request.Close를 true로 설정한 클라이언트를 만들지 않아서 제대로 동작하지 않았던 것이었다.\
따라서 GET 요청 자체를 새로 만들어서 request.Close 설정을 만들어준 후, 클라이언트 객체를 가지고 요청을 보내주면 해결됐다.

```go
req, err := http.NewRequest("GET", urls[idx], nil)
CheckErr(err)
req.Close = true

client := &http.Client{}
res, err := client.Do(req)
CheckErr(err)
checkCode(res)

defer res.Body.Close()
```

너무 안일하게 GET요청을 보냈던 것 같다..

---

## Heroku 관련

### 서버가 죽고나면 벌어지는 일들

Heroku는 주기적으로 서버를 죽인다. [Heroku Kaffeine](https://kaffeine.herokuapp.com/)에 서버를 등록하고\
30분마다 한 번씩 서버를 흔들어놔도 자꾸 죽였다.

죽이는 건 아무래도 좋았는데, 문제는 업데이트 시 벌어지는 일들이었다.

같은 글을 다시 올리게 하고 싶지 않아서 lastIndex.txt파일에 최근 게시글 번호를 저장했었다.

서버가 죽더라도 txt파일은 살아있으니까 상관없을 줄 알았는데.. 알고보니 서버가 죽고 다시 켜지면\
배포할 때 넣어놓았던 lastIndex.txt파일로 다시 돌아가서 같은 글을 다시 올려버리는 대참사가 벌어졌다.

게시글의 번호 4개만 저장하는데, DB서버를 사용하기도 firebase를 이용하기에도 너무 부끄러워서 다른 방법을 모색했다.

그리고 떠올린 아이디어는 디스코드 서버에 이미 올린 글들에서 힌트를 얻는 것이었다.\
즉, 매 업데이트 시 마지막으로 업데이트한 게시글 번호 4개를 같이 올려주고 서버는 매번 그 글을 보고 업데이트 하는 방법을 떠올렸다.\
뭔소리냐면, 디스코드 서버의 채팅채널 자체를 DB로 사용하겠다는 아이디어다.

#### 해결 방법

현재 로그인된 서버들의 채팅서버에서 메세지들을 모두 가져온다.

이 중 마지막에서 두 번째 메세지에 마지막으로 업데이트된 게시물 번호를 적어놓았다.

따라서 이를 가지고 매번 업데이트를 하게 되면 헤로쿠가 서버를 다시 실행시킨다해도 문제가 없다.
