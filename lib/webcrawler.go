package webcrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	str "strings"
)

type WebCrawler struct {
	base    string
	links   map[string]string
	visited map[string]string
}

func (this *WebCrawler) Run(base string) {
	this.base = base
	this.links = make(map[string]string)
	this.visited = make(map[string]string)

	this.processUrl("/")
}

func (this *WebCrawler) Print() {
	for _, v := range this.links {
		if v != "" {
			fmt.Println(v)
		}
	}
}

func (this *WebCrawler) GetLinks() map[string]string {
	return this.links
}

func (this *WebCrawler) AddLink(url string) bool {
	if this.visited[url] == "" {
		this.processUrl(url)
		return true
	}

	return false
}

func (this *WebCrawler) processUrl(url string) {

	if this.visited[url] != "" {
		return
	}

	this.visited[url] = url

	items, e := this.getLinks(this.base + url)
	if e != nil {
		fmt.Println(e)
		return
	}

	this.links[url] = url

	for k, v := range items {
		if this.visited[k] == "" {
			this.processUrl(v)
		}
	}
}

func (this *WebCrawler) getLinks(url string) (items map[string]string, e error) {
	// time.Sleep(1000 * time.Millisecond)

	items = make(map[string]string)

	response, e := http.Get(url)

	if e != nil {
		return items, e
	}

	if response.StatusCode >= 400 {
		e = fmt.Errorf("%d | %s", response.StatusCode, url)
		return
	}

	defer response.Body.Close()

	doc, e := goquery.NewDocumentFromReader(response.Body)

	if e != nil {
		return
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			if url, ok := this.filterLinks(href); ok {
				if items[url] == "" {
					items[url] = url
				}
			}
		}
	})

	return
}

func (this *WebCrawler) filterLinks(url string) (string, bool) {
	ok := true

	url = str.TrimSpace(url)

	// cut base
	pos := str.Index(url, this.base)
	if pos == 0 {
		pos = len(this.base)
		url = url[pos:]
	}

	// drop anchor
	pos = str.Index(url, "#")
	if pos == 0 {
		ok = false
		return url, ok
	}

	if pos != -1 {
		url = url[0:pos]
	}

	if url == "" ||
		!str.HasPrefix(url, "/") ||
		// str.HasPrefix(url, "http") ||
		// str.HasPrefix(url, "javascript") ||
		// str.HasPrefix(url, "tel") ||
		// str.HasPrefix(url, "callto") ||
		// str.HasPrefix(url, "mailto") ||
		str.HasPrefix(url, "/upload") ||
		str.HasPrefix(url, "/document.php") ||
		str.HasSuffix(url, ".pdf") {

		ok = false
		return url, ok
	}

	return url, ok
}
