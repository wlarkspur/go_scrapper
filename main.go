package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseUrl string = "https://www.jobkorea.co.kr/Search/?stext=python"

func main() {
	totalPages := getPages()
	for i := 0; i <= totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageUrl := baseUrl + "&tabType=recruit&Page_No=" + strconv.Itoa(page)
	fmt.Println("Requesting", pageUrl)
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
