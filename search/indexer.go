package search

import (
	"github.com/wangjia184/sortedset"
)

type WordData struct {
	Title string
}

type Index map[string]*sortedset.SortedSet

func IndexWriter(ch chan *CrawlData, index *Index) {
	for {
		m := <-ch
		title := m.Title
		words := m.Words
		url := m.Url

		for _, word := range words {
			if ss, ok := (*index)[word]; !ok {
				ss := sortedset.New()
				ss.AddOrUpdate(url, 1, WordData{title})
				(*index)[word] = ss
			} else {
				node := ss.GetByKey(url)
				if node != nil {
					ss.AddOrUpdate(url, node.Score()+1, node.Value)
				} else {
					ss.AddOrUpdate(url, 1, WordData{title})
				}
			}
		}
	}
}
