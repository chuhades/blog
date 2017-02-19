package main

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

var cache = struct {
	mutex       sync.RWMutex
	allArticles ArticleSlice
}{sync.RWMutex{}, ArticleSlice{}}

type Article struct {
	Title     string
	Link      string
	Published string
	Updated   string
}

type ArticleSlice []*Article

func (p ArticleSlice) Len() int {
	return len(p)
}

func (p ArticleSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// 时间倒序
func (p ArticleSlice) Less(i, j int) bool {
	return p[i].Published > p[j].Published
}

func getAllArticles() ArticleSlice {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	return cache.allArticles
}

func walkArchives() {
	wg := sync.WaitGroup{}
	cache.allArticles = ArticleSlice{}

	filepath.Walk("archives", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				articleInfo := parseArticleInfo(path)
				cache.mutex.Lock()
				defer cache.mutex.Unlock()
				cache.allArticles = append(cache.allArticles, articleInfo)
			}()

		}
		return err
	})
	wg.Wait()

	sort.Sort(cache.allArticles)
}

func parseArticleInfo(path string) *Article {
	// FIXME 不符合规则的文档
	line := []byte{}
	article := &Article{}

	fd, err := os.Open(path)
	if err != nil {
		return article
	}
	defer fd.Close()

	buf := bufio.NewReader(fd)

	line, _, _ = buf.ReadLine()
	article.Title = strings.Split(string(line), ": ")[1]
	line, _, _ = buf.ReadLine()
	article.Link = strings.Split(string(line), ": ")[1]
	line, _, _ = buf.ReadLine()
	article.Published = strings.Split(string(line), ": ")[1]
	line, _, _ = buf.ReadLine()
	article.Updated = strings.Split(string(line), ": ")[1]

	return article
}
