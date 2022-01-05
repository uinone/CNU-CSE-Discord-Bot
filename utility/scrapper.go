package utility

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ScrappedData struct {
	contentId 	int
	title 		string
	link 		string
	uploadedAt 	time.Time
}

func GetScrappedData(url string) []ScrappedData {
	scrapped := []ScrappedData{}

	res, err := http.Get(url)
	CheckErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	CheckErr(err)

	doc.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		num, _ := strconv.Atoi(cleanString(s.Find(".b-num-box").Text()))
		title := cleanString(s.Find(".b-title-box>a").Text())
		link, _ := s.Find(".b-title-box>a").Attr("href")
		link = url + link
		
		var uploadedAt time.Time
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if (i == 4) {
				uploadedAt = changeTimeToDate(cleanString(s.Text()))
			}
		})
		
		scrapped = append(scrapped, ScrappedData{
			contentId: num,
			title: title,
			link: link,
			uploadedAt: uploadedAt,
		})
	})

	return scrapped
}

func isDateBeforeToday(d time.Time) bool {
	now := time.Now().UTC().Add(time.Hour * 9)
	return now.After(d)
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

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}