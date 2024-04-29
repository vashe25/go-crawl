package main

import (
	"encoding/json"
	"fmt"
	"github.com/snabb/sitemap"
	webcrawler "go-crawl/lib"
	"io/ioutil"
	"os"
	"path/filepath"
	str "strings"
	"time"
)

type Config struct {
	BaseUrl        string   `json: "baseUrl"`
	SitemapUrlPath string   `json: "sitemapUrlPath"`
	SitemapPath    string   `json: "sitemapPath"`
	GetFrom        []string `json: "getFrom"`
	FilterRules    []string `json: "filterRules"`
}

func createPath(path string) string {
	var result string

	if str.HasPrefix(path, ".") {
		currentDir, e := filepath.Abs(filepath.Dir(os.Args[0]))

		if e != nil {
			panic(e)
		}

		result = currentDir + "/"
	}

	result = result + path
	result, e := filepath.Abs(result)

	if e != nil {
		panic(e)
	}

	// check if dir exist
	if _, e := os.Stat(result); os.IsNotExist(e) {
		os.MkdirAll(result, 0775)
	}

	return result
}

func loadConfig(path string) []Config {
	jsonFile, e := os.Open(path)

	if e != nil {
		panic(e)
	}

	defer jsonFile.Close()

	byteValue, e := ioutil.ReadAll(jsonFile)

	if e != nil {
		panic(e)
	}

	var configs []Config
	json.Unmarshal([]byte(byteValue), &configs)

	if len(configs) == 0 {
		panic("Bad configs")
	}

	return configs
}

func main() {
	currentDir, e := filepath.Abs(filepath.Dir(os.Args[0]))
	if e != nil {
		panic(e)
	}

	now := time.Now().Local()

	configs := loadConfig(currentDir + "/config.json")

	for _, config := range configs {
		crawler := new(webcrawler.WebCrawler)
		crawler.LoadFilterRules(config.FilterRules)
		crawler.Run(config.BaseUrl)

		for _, url := range config.GetFrom {
			crawler.GetLinksFromUrl(url)
		}

		sitemapPath := createPath(config.SitemapPath)

		chunkedData := crawler.GetChunked(5000)
		countChunks := len(chunkedData)

		if countChunks > 1 {
			// create sitemap index
			smi := sitemap.NewSitemapIndex()
			for i := 0; i < countChunks; i++ {
				smi.Add(&sitemap.URL{
					Loc:     fmt.Sprintf("%v%vsitemap-%v.xml", config.BaseUrl, config.SitemapUrlPath, i),
					LastMod: &now,
				})
			}

			file, e := os.Create(sitemapPath + "/sitemap.xml")
			if e != nil {
				panic(e)
			}

			defer file.Close()
			smi.WriteTo(file)

			// create sitemaps
			for i := 0; i < countChunks; i++ {
				sm := sitemap.New()

				for _, link := range chunkedData[i] {
					sm.Add(&sitemap.URL{
						Loc:        config.BaseUrl + link,
						LastMod:    &now,
						ChangeFreq: sitemap.Daily,
					})

				}

				file, e := os.Create(fmt.Sprintf("%v/sitemap-%v.xml", sitemapPath, i))
				if e != nil {
					panic(e)
				}

				defer file.Close()
				sm.WriteTo(file)
			}
		} else {
			// create sitemap
			sm := sitemap.New()

			for _, link := range chunkedData[0] {
				sm.Add(&sitemap.URL{
					Loc:        config.BaseUrl + link,
					LastMod:    &now,
					ChangeFreq: sitemap.Daily,
				})

			}

			file, e := os.Create(sitemapPath + "/sitemap.xml")
			if e != nil {
				panic(e)
			}

			defer file.Close()
			sm.WriteTo(file)
		}
	}

}
