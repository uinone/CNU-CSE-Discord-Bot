# CNU_Discod_Bot 프로젝트

[discordgo](https://github.com/bwmarrin/discordgo)와 [goquery](https://github.com/PuerkitoBio/goquery)를 활용했어요.

## 봇은 이렇게 동작해요😀

봇은 다음과 같이 동작해요.

1. 먼저 충남대학교 컴퓨터융합학부의 게시판들이 업데이트 되는 것을 확인합니다.
2. 해당 봇을 추가한 서버들의 특정 채팅채널에서 볼 수 있도록 업데이트 된 게시글의 정보를 웹에서 긁어옵니다.
3. 긁어온 데이터를 채팅채널에 메시지 형태로 보내줍니다.

## 이거 때문에 너무 힘들었어요😥

시간이 날 때마다 정리해서 올릴 예정이에요.

- 디스코드 봇 관련
  - 계속 꺼지는 서버
  - 비어있는 session.Guilds.Channels

* 웹 크롤링 관련
  - 날 죽이려는 whitespace
  - 이상한 날짜 변환
