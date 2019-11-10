package search

import (
	"github.com/gocolly/colly"
	"log"
	neturl "net/url"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type CrawlData struct {
	Url   string
	Title string
	Words []string
}

func Crawl(url string, toIndexWriter chan *CrawlData) (int64, int64) {

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.Async(true))

	_ = colly.LimitRule{Parallelism: 30}

	urlCount := int64(0)
	wordCount := int64(0)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	c.OnHTML("*:not(script)", func(e *colly.HTMLElement) {

		title := strings.TrimSpace(e.DOM.Find("title").Text())
		e.ForEach("body", func(i int, eIn *colly.HTMLElement) {
			if err != nil {
				log.Fatal(err)
			}

			words := strings.Split(strings.ToLower(strings.TrimSpace(reg.ReplaceAllString(e.Text, " "))), " ")
			wordCount += int64(len(words))
			msg := CrawlData{
				eIn.Request.URL.String(),
				title,
				words}

			if len(msg.Words) > 0 {
				toIndexWriter <- &msg
			}
		})
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link, err := neturl.Parse(e.Request.AbsoluteURL(e.Attr("href")))

		if err != nil {
			return
		}
		host := link.Host
		scheme := link.Scheme

		if scheme != "http" && scheme != "https" {
			return
		}

		next := scheme + "://" + host + "/"
		urlCount += int64(1)
		_ = c.Visit(next)
	})

	c.OnRequest(func(request *colly.Request) {
		//log.Println("Found URL:", request.URL)
	})

	e := c.Visit(url)
	check(e)
	c.Wait()
	return urlCount, wordCount
}
