package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vashe25/queue"
	"go-crawl/logger"
	"net/http"
	str "strings"
	"sync"
)

type WebCrawler struct {
	base        string
	vlt         *vault
	client      *http.Client
	filterRules []rule
	loadQueue   *queue.Queue
	stop        chan struct{}
}

func NewWebCrawler(base string, rules []string) (*WebCrawler, error) {
	compiledRules, err := compileFilterRules(rules)
	if err != nil {
		return nil, err
	}
	return &WebCrawler{
		base:        base,
		filterRules: compiledRules,
		vlt:         newVault(),
		client:      &http.Client{},
		loadQueue:   queue.NewQueue(),
		stop:        make(chan struct{}),
	}, nil
}

func compileFilterRules(rules []string) ([]rule, error) {
	length := len(rules)
	compiledRules := make([]rule, length)
	for i, pattern := range rules {
		r, err := compileRule(pattern)
		if err != nil {
			return nil, err
		}
		compiledRules[i] = *r
	}
	return compiledRules, nil
}

// Run starts crawling
// base must be in format "https://host"
func (_this *WebCrawler) Run() {
	var (
		loadPageCh    chan string
		processUrlsCh chan map[string]string
		bufferCh      chan string
		closeBuffer   chan struct{}
	)

	workers := 4
	loadPageCh = make(chan string, 1)
	processUrlsCh = make(chan map[string]string, workers)
	bufferCh = make(chan string, 1)
	closeBuffer = make(chan struct{})

	go _this.getFromBufferTo(loadPageCh, closeBuffer)
	go _this.putToBufferFrom(bufferCh, closeBuffer)
	go _this.loadPagesPool(workers, loadPageCh, processUrlsCh)
	go _this.processUrls(processUrlsCh, bufferCh)
	_this.vlt.addToQueue("/")
	loadPageCh <- "/"

	// lock until the end
	<-_this.stop
}

func (_this *WebCrawler) loadPagesPool(workers int, ch1 <-chan string, ch2 chan<- map[string]string) {
	defer close(_this.stop)
	defer close(ch2)

	var wg sync.WaitGroup
	wg.Add(workers)
	for ; workers > 0; workers-- {
		go _this.loadPageWorker(ch1, ch2, &wg)
	}
	wg.Wait()
}

func (_this *WebCrawler) loadPageWorker(ch1 <-chan string, ch2 chan<- map[string]string, wg *sync.WaitGroup) {
	for url := range ch1 {
		_this.vlt.addVisited(url)
		response := _this.loadPage(url)
		if response == nil {
			continue
		}
		_this.vlt.collect(url)

		urls := _this.parseResponse(response)
		response.Body.Close()

		if urls != nil {
			ch2 <- urls
		}
	}
	wg.Done()
}

func (_this *WebCrawler) putToBufferFrom(ch1 <-chan string, closeBuffer chan<- struct{}) {
	defer close(closeBuffer)
	for url := range ch1 {
		_this.loadQueue.Push(queue.NewTask(url))
	}
	closeBuffer <- struct{}{}
}

func (_this *WebCrawler) getFromBufferTo(ch1 chan<- string, closeBuffer <-chan struct{}) {
	defer close(ch1)
	for {
		select {
		case <-closeBuffer:
			for _this.loadQueue.Length() != 0 {
				url := _this.loadQueue.Pop().Value().(string)
				ch1 <- url
			}
			return
		default:
			if _this.loadQueue.Length() != 0 {
				url := _this.loadQueue.Pop().Value().(string)
				ch1 <- url
			}
		}
	}
}

func (_this *WebCrawler) processUrls(ch1 <-chan map[string]string, bufferCh chan<- string) {
	defer close(bufferCh)
	for urls := range ch1 {
		for url := range urls {
			url = _this.sanitizeURL(url)
			if _this.isSkipURL(url) {
				continue
			}
			if _this.vlt.isVisited(url) {
				continue
			}
			if _this.vlt.isInQueue(url) {
				continue
			}
			_this.vlt.addToQueue(url)
			bufferCh <- url
		}

		if _this.vlt.isFull() {
			break
		}
	}
}

// Print to stdout
func (_this *WebCrawler) Print() {
	for item := range _this.vlt.collected().Items() {
		fmt.Println(item)
	}
}

func (_this *WebCrawler) GetLinks() map[string]bool {
	return _this.vlt.collected().Items()
}

// GetChunked returns result in chunks
func (_this *WebCrawler) GetChunked(size int) map[int][]string {
	if size == 0 {
		size = 5000
	}

	chunks := make(map[int][]string)

	i := 0
	j := 0
	for item := range _this.vlt.collected().Items() {
		chunks[i] = append(chunks[i], item)

		if j == size {
			j = 0
			i++
		} else {
			j++
		}
	}

	return chunks
}

// CheckAdditionalURL adds url into loadQueue
func (_this *WebCrawler) CheckAdditionalURL(url string) {
	url = _this.sanitizeURL(url)
	if _this.isSkipURL(url) {
		return
	}
	_this.loadQueue.Push(queue.NewTask(url))
}

func (_this *WebCrawler) loadPage(url string) *http.Response {
	request := _this.makeRequest(url)
	response, err := _this.client.Do(request)
	if err != nil {
		logger.Log("[c] [loadPage] get url fail: '%s' reason: '%s'", url, err.Error())
		return nil
	} else {
		if response.StatusCode >= 400 {
			response.Body.Close()
			logger.Log("[c] [loadPage] %v '%s'", response.StatusCode, url)
			return nil
		}
		return response
	}
}

func (_this *WebCrawler) parseResponse(response *http.Response) map[string]string {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil
	}

	var urls = make(map[string]string)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			url := _this.sanitizeURL(href)
			if !_this.isSkipURL(url) && urls[url] == "" {
				urls[url] = url
			}
		}
	})

	if len(urls) == 0 {
		return nil
	}

	return urls
}

func (_this *WebCrawler) makeRequest(url string) *http.Request {
	request, err := http.NewRequest("GET", _this.base+url, nil)
	if err != nil {
		return nil
	}
	request.Header.Add("User-Agent", "GO-Crawl")
	return request
}

func (_this *WebCrawler) isSkipURL(url string) bool {
	if url == "" {
		return true
	}

	for _, r := range _this.filterRules {
		if r.match(url) {
			return true
		}
	}

	return false
}

func (_this *WebCrawler) sanitizeURL(url string) string {
	url = str.TrimSpace(url)

	// cut base
	pos := str.Index(url, _this.base)
	if pos == 0 {
		pos = len(_this.base)
		url = url[pos:]
	}

	// drop anchor
	pos = str.Index(url, "#")

	// check utm params
	if str.Contains(url, "utm_") {
		pos = str.Index(url, "?")
	}

	if pos == 0 {
		return ""
	}

	if pos != -1 {
		url = url[0:pos]
	}

	if !str.HasPrefix(url, "/") {
		return ""
	}

	return url
}
