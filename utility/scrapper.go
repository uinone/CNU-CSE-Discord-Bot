package utility

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

type ScrappedData struct {
	contentId	int
	title 		string
	link 		string
	uploadedAt 	string
}

type msgData struct {
	idx 	int
	data 	[]ScrappedData
}

var (
	urls = [4]string{
		"https://computer.cnu.ac.kr/computer/notice/bachelor.do",
		"https://computer.cnu.ac.kr/computer/notice/notice.do",
		"https://computer.cnu.ac.kr/computer/notice/project.do",
		"https://computer.cnu.ac.kr/computer/notice/job.do",
	}

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

func SendScrappedData(ds *discordgo.Session, envData []string) {
	results := make(chan msgData)
	
	fmt.Println("Reciving Data...")
	for i:=0; i<len(urls); i++ {
		contentId, _ := strconv.Atoi(envData[i])
		go getScrappedData(i, contentId, results)
	}
	
	msgs := []msgData{}
	for i:=0; i<len(urls); i++ {
		msg := <-results
		fmt.Println(msg)
		msgs = append(msgs, msg)
	}
	
	fmt.Println("Reciving Data done.")

	SendMessageToChannel(ds, "모두 주목! 컴공과 공지 알림을 시작할게요🐧")

	for _, content := range msgs {
		SendMessageToChannel(ds, boardName[content.idx])
		if len(content.data) == 0 {
			SendMessageToChannel(ds, "새로 올라온 게시글이 없습니다.\n---")
		} else {
			var msg string
			for i, data := range content.data {
				if i == 0 {
					envData[content.idx] = strconv.Itoa(data.contentId)
				}
				msg = ""
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

	UpdateEnvData(envData)

	fmt.Println("스크랩이 끝났습니다.")
}

func getScrappedData(idx int, lastContentId int, results chan<- msgData) {
	scrapped := []ScrappedData{}

	res, err := http.Get(urls[idx])
	CheckErr(err)
	checkCode(res)
	
	res.Request.Close = true
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)

	doc.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		num, _ := strconv.Atoi(cleanString(s.Find(".b-num-box").Text()))
		if num > lastContentId {
			title := cleanString(s.Find(".b-title-box>a").Text())
			link, _ := s.Find(".b-title-box>a").Attr("href")
			link = urls[idx] + link
			
			var uploadedAt string
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				if (i == 4) {
					uploadedAt = getDayCountFromNow(ChangeTimeToDate(cleanString(s.Text())))
				}
			})
			
			scrapped = append(scrapped, ScrappedData{
				contentId: num,
				title: title,
				link: link,
				uploadedAt: uploadedAt,
			})
		}
	})

	results <- msgData{idx: idx, data: scrapped}
}

func getDayCountFromNow(t time.Time) string {
	now := time.Now().Add(time.Hour * 9)
	days := int(now.Sub(t).Hours() / 24)
	var dayCount string
	if days == 0 {
		dayCount = "오늘"
	} else {
		dayCount = strconv.Itoa(days) + "일전"
	}
	return dayCount
}

func ChangeTimeToDate(str string) time.Time {
	strDate := strings.Join(strings.Split(str, "."), "-")
	t, _ := time.Parse("06-01-02", strDate)
	return t
}

func cleanString(str string) string {
	return strings.TrimSpace(str)
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}