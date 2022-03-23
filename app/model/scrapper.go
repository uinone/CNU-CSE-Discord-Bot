package model

import (
	"GO/nomad/app/view"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/PuerkitoBio/goquery"
)

var (
	urls = [4]string{
		"https://computer.cnu.ac.kr/computer/notice/bachelor.do",
		"https://computer.cnu.ac.kr/computer/notice/notice.do",
		"https://computer.cnu.ac.kr/computer/notice/project.do",
		"https://computer.cnu.ac.kr/computer/notice/job.do",
	}

	boardName = [4]string{
		"🎨 학사공지 🎨",
		"📜 일반소식 📜",
		"🔆 사업단소식 🔆",
		"🎈 취업정보 🎈 ※취업정보는 로그인해야 볼 수 있어요!😅",
	}

	contentPropertyName = [3]string{
		"[제목] ",
		"[링크] ",
		"[업로드 날짜] ",
	}
)

type scrapper struct {
	viewer *view.Viewer
}

type scrappedData struct {
	articleNo  int
	title      string
	link       string
	uploadedAt string
}

type infoData struct {
	idx  int
	data []scrappedData
}

// Create scrapper object
func NewScrapper(ds *discordgo.Session) *scrapper {
	s := new(scrapper)

	s.viewer = view.NewViewer()
	s.viewer.SetDiscordSession(ds)

	return s
}

// Get Info data parsed from scrapped data
func (s *scrapper) getInfoData(lastIndexData []string) [][]string {
	now := time.Now()

	results := make(chan infoData)
	defer close(results)

	s.viewer.PrintlnMsgToConsole("Reciving Data...")
	for i := 0; i < len(urls); i++ {
		articleNo, _ := strconv.Atoi(lastIndexData[i])
		go s.getScrappedData(i, articleNo, results)
	}

	info := []infoData{}
	for i := 0; i < len(urls); i++ {
		result := <-results
		if result.data != nil {
			info = append(info, result)
		}
	}

	done := strconv.Itoa(int(time.Since(now).Seconds()))

	recivingEndMsg := "Reciving Data done in " + done + "sec."
	s.viewer.PrintlnMsgToConsole(recivingEndMsg)

	return s.formatScrappedData(info, lastIndexData)
}

// Get scrapped data by using web scrapping concurrently
func (s *scrapper) getScrappedData(idx int, articleNo int, results chan<- infoData) {
	req, err := http.NewRequest("GET", urls[idx], nil)
	if err != nil {
		s.viewer.FatallnErrorToConsole(err)
	}
	req.Close = true

	client := &http.Client{}
	res, err := client.Do(req)
	
	if err != nil {
		s.viewer.FatallnErrorToConsole(err)
		results <- infoData{idx: idx, data: nil} // When "Get EOF" error occur
		res.Body.Close()
		return
	}

	defer res.Body.Close()
	
	if res.StatusCode != 200 {
		statusErrorMsg := "Request failed with Status:" + strconv.Itoa(res.StatusCode)
		s.viewer.FatallnMsgToConsole(statusErrorMsg)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		s.viewer.FatallnErrorToConsole(err)
	}

	scrapped := []scrappedData{}
	
	doc.Find("tbody").Find("tr").Each(func(i int, gs *goquery.Selection) {
		link, _ := gs.Find(".b-title-box>a").Attr("href")
		contentArticleNo := s.getArticleNo(link)

		if contentArticleNo != -1 && contentArticleNo > articleNo {
			title := strings.TrimSpace(gs.Find(".b-title-box>a").Text())
			link = urls[idx] + link
			
			var uploadedAt string
			gs.Find("td").Each(func(i int, s *goquery.Selection) {
				if (i == 4) {
					uploadedAt = strings.TrimSpace(s.Text())
				}
			})
			
			scrapped = append(scrapped, scrappedData{
				articleNo: contentArticleNo,
				title: title,
				link: link,
				uploadedAt: uploadedAt,
			})
		}
	})

	results <- infoData{idx: idx, data: scrapped}
}

// Get articleNo processing url string
func (s *scrapper) getArticleNo(url string) int {
	queryString := strings.Split(url[1:], "&")

	for _, querySubString := range queryString {
		query := strings.Split(querySubString, "=")
		if query[0] == "articleNo" {
			atricleNo, _ := strconv.Atoi(query[1])
			return atricleNo
		}
	}

	return -1
}

// Format scrapped data(infoData) to strings
func (s *scrapper) formatScrappedData(infoSet []infoData, lastIndexData []string) [][]string {
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
				lastIndexData[info.idx] = strconv.Itoa(content.articleNo)
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

	lastIndexes := "$"
	for _, lastIndex := range lastIndexData {
		lastIndexes = fmt.Sprint(lastIndexes, lastIndex + " ")
	}

	formatedDataSet = append(formatedDataSet, []string{lastIndexes})

	return formatedDataSet
}