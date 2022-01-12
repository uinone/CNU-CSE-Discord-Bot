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

type scrappedData struct {
	contentId	int
	title 		string
	link 		string
	uploadedAt 	string
}

type infoData struct {
	idx 	int
	data 	[]scrappedData
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
// Get Info data parsed from scrapped data
func getInfoData(ds *discordgo.Session) [][]string {
	now := time.Now()
	
	results := make(chan infoData)
	defer close(results)

	lastIndexData := getLastIndexData()

	fmt.Println("Reciving Data...")
	for i:=0; i<len(urls); i++ {
		contentId, _ := strconv.Atoi(lastIndexData[i])
		go getScrappedData(i, contentId, results)
	}
	
	info := []infoData{}
	for i:=0; i<len(urls); i++ {
		info = append(info, <-results)
	}

	done := time.Since(now).Seconds()
	fmt.Println("Reciving Data done.", done)

	return formatScrappedData(info, lastIndexData)
}

func formatScrappedData(infoSet []infoData, lastIndexData []string) [][]string {
	formatedDataSet := [][]string{}

	for _, info := range infoSet {
		formatedData := []string{}
		if len(info.data) == 0 {
			continue
		}
		formatedData = append(formatedData, boardName[info.idx])

		var tmpMsg string
		for i, content := range info.data {
			if i == 0 {
				lastIndexData[info.idx] = strconv.Itoa(content.contentId)
			}
			tmpMsg = ""
			tmpMsg = fmt.Sprint(tmpMsg, contentPropertyName[0])
			tmpMsg = fmt.Sprintln(tmpMsg, content.title)
			formatedData = append(formatedData, tmpMsg)

			tmpMsg = ""
			tmpMsg = fmt.Sprint(tmpMsg, contentPropertyName[1]) 
			tmpMsg = fmt.Sprintln(tmpMsg, content.link)
			formatedData = append(formatedData, tmpMsg)

			tmpMsg = ""
			tmpMsg = fmt.Sprint(tmpMsg, contentPropertyName[2]) 
			tmpMsg = fmt.Sprintln(tmpMsg, content.uploadedAt)
			tmpMsg = fmt.Sprintln(tmpMsg, "+")

			formatedData = append(formatedData, tmpMsg)
		}

		formatedData = append(formatedData, "---")

		formatedDataSet = append(formatedDataSet, formatedData)
	}

	updateLastIndexData(lastIndexData)

	return formatedDataSet
}

func getScrappedData(idx int, lastContentId int, results chan<- infoData) {
	req, err := http.NewRequest("GET", urls[idx], nil)
	checkErr(err)
	req.Close = true

	client := &http.Client{}
	res, err := client.Do(req)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	scrapped := []scrappedData{}
	
	doc.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		num, _ := strconv.Atoi(cleanString(s.Find(".b-num-box").Text()))
		if num > lastContentId {
			title := cleanString(s.Find(".b-title-box>a").Text())
			link, _ := s.Find(".b-title-box>a").Attr("href")
			link = urls[idx] + link
			
			var uploadedAt string
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				if (i == 4) {
					uploadedAt = getDayCountFromNow(changeTimeToDate(cleanString(s.Text())))
				}
			})
			
			scrapped = append(scrapped, scrappedData{
				contentId: num,
				title: title,
				link: link,
				uploadedAt: uploadedAt,
			})
		}
	})

	results <- infoData{idx: idx, data: scrapped}
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

func changeTimeToDate(str string) time.Time {
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

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}