package main

import (
	"fmt"
	"log"
	"net/http"
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

var baseUrl string = "https://www.jobkorea.co.kr/Search/?stext=python"

func main() {
	var jobs []extractedJob
	totalPages := getPages()
	for i := 0; i <= totalPages; i++ {
		extractedJob := getPage(i)
		jobs = append(jobs, extractedJob...)
	}
	fmt.Println(jobs)
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageUrl := baseUrl + "&tabType=recruit&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".list-post")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs

}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-gino")
	company := cleanString(card.Find(".post-list-corp>a").Text())
	title := cleanString(card.Find(".post-list-info>a").Text())
	exp := cleanString(card.Find(".option .exp").Text())
	edu := cleanString(card.Find(".option .edu").Text())
	return extractedJob{
		id:      id,
		company: company,
		title:   title,
		exp:     exp,
		edu:     edu,
	}

}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
func getPages() int {
	pages := 0
	res, err := http.Get(baseUrl)
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
