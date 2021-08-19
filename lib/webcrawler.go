/*
WebCrawler - package is used for recursive walking by host
and collecting links.
*/
package webcrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	str "strings"
	"regexp"
)

/*
base - is host
links - collection of good urls
visited - collection of all visited urls, a wayout from recursion.
*/
type WebCrawler struct {
	base    string
	links   map[string]string
	visited map[string]string
	client  *http.Client
	filterRules []string
}

/*
Loads rules to skip or filter urls
*/
func (this *WebCrawler) LoadFilterRules(filterRules []string) {
	this.filterRules = filterRules
}

/*
Main method to start crawling.
base - must be in a "http://host" format.
*/
func (this *WebCrawler) Run(base string) {
	this.base = base
	this.links = make(map[string]string)
	this.visited = make(map[string]string)
	this.client = &http.Client{}

	this.processUrl("/")
}

/*
Simply prints links map into stdout.
*/
func (this *WebCrawler) Print() {
	for _, v := range this.links {
		if v != "" {
			fmt.Println(v)
		}
	}
}

/*
Returns a map of links.
*/
func (this *WebCrawler) GetLinks() map[string]string {
	return this.links
}

/*
Slice links map into chunks.
Mainly used for creating index sitemap.xml and sitemap-%d.xml
*/
func (this *WebCrawler) GetChunked(size int) map[int][]string {
	if size == 0 {
		size = 5000
	}

	chunks := make(map[int][]string)

	i := 0
	j := 0
	for link := range this.links {
		chunks[i] = append(chunks[i], link)

		if j == size {
			j = 0
			i++
		} else {
			j++
		}
	}

	return chunks
}

/*
Adds url into links map and starts crawling.
*/
func (this *WebCrawler) AddLink(url string) {
	this.processUrl(url)
}

/*
Starts crawling from url, but not adds url into links map.
*/
func (this *WebCrawler) GetLinksFromUrl(url string) {
	if this.visited[url] != "" {
		return
	}

	this.visited[url] = url

	items, e := this.getLinks(this.base + url)
	if e != nil {
		fmt.Println(e)
		return
	}

	for k, v := range items {
		if this.visited[k] == "" {
			this.processUrl(v)
		}
	}
}

/*
Core of Webcrawler.
Processing url.
*/
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

/*
Get document. Checks status code. Parse document body. Returns found links.
*/
func (this *WebCrawler) getLinks(url string) (items map[string]string, e error) {
	// throttler
	// time.Sleep(2 * time.Millisecond)

	items = make(map[string]string)

	request, e := http.NewRequest("GET", url, nil)

	if e != nil {
		return items, e
	}

	request.Header.Add("User-Agent", "GO-Crawl")

	response, e := this.client.Do(request)

	if e != nil {
		return items, e
	} else {
		defer response.Body.Close()

		if response.StatusCode >= 400 {
			e = fmt.Errorf("%d | %s", response.StatusCode, url)
			return
		}

		doc, e := goquery.NewDocumentFromReader(response.Body)

		if e != nil {
			return items, e
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if href, ok := s.Attr("href"); ok {
				if url, ok := this.filterLink(href); ok {
					if items[url] == "" {
						items[url] = url
					}
				}
			}
		})

		return items, e
	}
}

/*
Filter.
A place for rules to skip link.
*/
func (this *WebCrawler) filterLink(url string) (string, bool) {
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

	// check utm params
	if str.Contains(url, "utm") {
		pos = str.Index(url, "?")
	}

	if pos == 0 {
		ok = false
		return url, ok
	}

	if pos != -1 {
		url = url[0:pos]
	}

	if url == "" || !str.HasPrefix(url, "/") {
		ok = false
		return url, ok
	}

	for _, r := range this.filterRules {
		matched, e := regexp.MatchString(r, url)
		if e != nil {
			panic(e)
		}
		if matched {
			ok = false
			return url, ok
		}
	}

	return url, ok
}
