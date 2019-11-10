package main

import (
	"context"
	"fmt"
	micro "github.com/micro/go-micro"
	"github.com/wangjia184/sortedset"
	"log"
	indexer "pinpoint/api"
	"pinpoint/search"
	"strings"
)

type Indexer struct {
	ch chan *search.CrawlData
	ix *search.Index
}

func (c *Indexer) Index(ctx context.Context, req *indexer.IndexRequest, resp *indexer.IndexResponse) error {

	log.Println("Received index request for", req.Url)
	urls, words := search.Crawl(req.Url, c.ch)
	log.Printf("Urls indexed: %d, words indexed: %d\n", urls, words)
	resp.Body = fmt.Sprintf("Urls indexed: %d, words indexed: %d", urls, words)

	return nil
}

func (c *Indexer) Search(ctx context.Context, req *indexer.SearchRequest, resp *indexer.SearchResponse) error {

	word := strings.ToLower(req.Word)
	ix := (map[string]*sortedset.SortedSet)(*c.ix)
	if v, ok := ix[word]; ok {

		var pages []*indexer.SearchResponse_Page
		for _, node := range v.GetByRankRange(-1, 1, false) {
			pages = append(
				pages,
				&indexer.SearchResponse_Page{
					Url:   node.Key(),
					Title: node.Value.(search.WordData).Title,
					Count: int64(node.Score()),
				})
		}
		resp.Pages = pages
	}

	return nil
}

func main() {

	// Start search writer
	index := search.Index{}
	toIndexWriter := make(chan *search.CrawlData)
	go search.IndexWriter(toIndexWriter, &index)

	// Start service
	service := micro.NewService(
		micro.Name("search"))
	service.Init()

	// Register handler
	indexerHandler := Indexer{toIndexWriter, &index}
	err := indexer.RegisterIndexerHandler(service.Server(), &indexerHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Start server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
