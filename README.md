# CNU_Discod_Bot 프로젝트

[discordgo](https://github.com/bwmarrin/discordgo)와 [goquery](https://github.com/PuerkitoBio/goquery)를 활용했어요.

## 봇은 이렇게 동작해요😀

봇은 다음과 같이 동작해요.

1. 먼저 충남대학교 컴퓨터융합학부의 게시판들이 업데이트 되는 것을 주기적으로 확인합니다.
2. 해당 봇을 추가한 서버들의 특정 채팅채널에서 볼 수 있도록 업데이트 된 게시글의 정보를 웹에서 긁어옵니다.
3. 긁어온 데이터를 채팅채널에 메시지 형태로 보내줍니다.

## 이거 때문에 너무 힘들었어요😥

시간이 날 때마다 정리해서 올릴 예정이에요.

- 디스코드 봇 관련
  - 계속 꺼지는 프로그램
  - 비어있는 session.Guilds.Channels

* 웹 크롤링 관련
  - 날 죽이려는 whitespace
  - 이상한 날짜 변환

## 디스코드 봇 관련

디스코드 봇을 만들면서 겪었던 문제들과 해결 과정이에요.

### 계속 꺼지는 프로그램

이전에 Python이나 Node.js를 가지고 개발을 했었다.\
Python의 discord.py 라이브러리에는 run이라는 함수가 있었고, 이를 실행시키면 따로 프로그램을 끄기 전까지는 계속해서 돌아갔었다. Node.js의 discord.js도 마찬가지로 login이라는 함수가 같은 기능을 해줬다.\

근데 discord.go의 Open 함수는 그렇게 동작하지 않고, 봇에 로그인한 다음 바로 프로그램이 종료되었다.\

#### 해결 방법

os.Signal 타입 채널을 생성하고, Interrupt가 일어나면 바로 종료하도록 만들면 됐다.

```go
sc := make(chan os.Signal, 1)
signal.Notify(sc, os.Interrupt)
<-sc
```
