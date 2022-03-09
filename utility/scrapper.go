package utility

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type scrappedData struct {
	articleNo	int
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
func getInfoData() [][]string {
	now := time.Now()
	
	results := make(chan infoData)
	defer close(results)

	//lastIndexData := getLastIndexData(ds)
	lastIndexData := []string{
		"291363",
		"291673",
		"292026",
		"291727",
	}

	fmt.Println("Reciving Data...")
	for i:=0; i<len(urls); i++ {
		articleNo, _ := strconv.Atoi(lastIndexData[i])
		go getScrappedData(i, articleNo, results)
	}
	
	info := []infoData{}
	for i:=0; i<len(urls); i++ {
		info = append(info, <-results)
	}

	done := time.Since(now).Seconds()
	fmt.Println("Reciving Data done.", done)

	return formatScrappedData(info, lastIndexData)
}

// Format scrapped data(infoData) to strings
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

// Get scrapped data by using web scrapping concurrently
func getScrappedData(idx int, articleNo int, results chan<- infoData) {
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
		link, _ := s.Find(".b-title-box>a").Attr("href")
		contentArticleNo := getArticleNo(link)

		if contentArticleNo != -1 && contentArticleNo > articleNo {
			title := cleanString(s.Find(".b-title-box>a").Text())
			link = urls[idx] + link
			
			var uploadedAt string
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				if (i == 4) {
					uploadedAt = cleanString(s.Text())
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

// Clean string by using strings.TrimSpace
func cleanString(str string) string {
	return strings.TrimSpace(str)
}

func getArticleNo(url string) int {
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

// Check that response's status code is 200
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

// Check err is nil
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}