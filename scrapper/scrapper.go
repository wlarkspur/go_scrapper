package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id      string
	company string
	title   string
	exp     string
	edu     string
}

// Scrape JobKorea by a term
func Scrape(term string) {

	var baseUrl string = "https://www.jobkorea.co.kr/Search/?stext=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseUrl)

	for i := 0; i <= totalPages; i++ {
		go getPage(i, baseUrl, c)

	}

	for i := 0; i < totalPages; i++ {
		extractedJob := <-c
		jobs = append(jobs, extractedJob...)
	}
	writeJobs(jobs, term)
	/* done := make(chan bool)
	go func() {

		done <- true
	}()
	<-done */
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageUrl := url + "&tabType=recruit&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".list-post")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs

}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-gno")
	company := CleanString(card.Find(".post-list-corp>a").Text())
	title := CleanString(card.Find(".post-list-info>a").Text())
	exp := CleanString(card.Find(".option .exp").Text())
	edu := CleanString(card.Find(".option .edu").Text())

	c <- extractedJob{
		id:      id,
		company: company,
		title:   title,
		exp:     exp,
		edu:     edu,
	}
}

// CleanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	fmt.Println(res.StatusCode)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".tplPagination.newVer.wide").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func writeJobs(jobs []extractedJob, term string) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	uft8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(uft8bom)

	w := csv.NewWriter(file)
	defer w.Flush()
	defer file.Close()

	headers := []string{"Link", "Company", "Title", "Exp", "Edu"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{`https://www.jobkorea.co.kr/Recruit/GI_Read/` + job.id + `?Oem_Code=C1&logpath=1&stext=` + term + `&listno=`, job.company, job.title, job.exp, job.edu}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
