package main

import (
	"fmt"
	"github.com/snabb/sitemap"
	"go-crawl/config"
	webcrawler "go-crawl/crawler"
	"go-crawl/logger"
	"go-crawl/utility"
	"os"
	"time"
)

func main() {
	currentDir, err := utility.CurrentDir()
	if err != nil {
		utility.Exit(err)
	}

	configs, err := config.LoadConfig(currentDir + "/config.json")
	if err != nil {
		utility.Exit(err)
	}

	now := time.Now().Local()
	for _, conf := range configs {
		startTime := time.Now()
		crawler, err := webcrawler.NewWebCrawler(conf.GetBaseUrl(), conf.GetFilterRules())
		if err != nil {
			utility.Exit(err)
		}
		logger.Log("[main] crawling '%s'", conf.GetBaseUrl())
		for _, url := range conf.GetGetFrom() {
			crawler.CheckAdditionalURL(url)
		}
		crawler.Run()

		var sitemapPath string
		sitemapPath, err = utility.MakeDir(conf.GetSitemapPath())
		if err != nil {
			utility.Exit(err)
		}

		chunkedData := crawler.GetChunked(5000)
		countChunks := len(chunkedData)

		logger.Log("[main] collected %d chunks", countChunks)
		logger.Log("[main] writting sitemap '%s'", conf.GetBaseUrl())
		if countChunks > 1 {
			// create sitemap index
			smi := sitemap.NewSitemapIndex()
			for i := 0; i < countChunks; i++ {
				smi.Add(&sitemap.URL{
					Loc:     fmt.Sprintf("%v%vsitemap-%v.xml", conf.GetBaseUrl(), conf.GetSitemapUrlPath(), i),
					LastMod: &now,
				})
			}

			file, err := os.Create(sitemapPath + "/sitemap.xml")
			if err != nil {
				utility.Exit(err)
			}

			smi.WriteTo(file)
			file.Close()

			// create sitemaps
			for i := 0; i < countChunks; i++ {
				sm := sitemap.New()

				for _, link := range chunkedData[i] {
					sm.Add(&sitemap.URL{
						Loc:        conf.GetBaseUrl() + link,
						LastMod:    &now,
						ChangeFreq: sitemap.Daily,
					})

				}

				file, err := os.Create(fmt.Sprintf("%vsitemap-%v.xml", sitemapPath, i))
				if err != nil {
					panic(err)
				}

				sm.WriteTo(file)
				file.Close()
			}
		} else {
			// create sitemap
			sm := sitemap.New()

			for _, link := range chunkedData[0] {
				sm.Add(&sitemap.URL{
					Loc:        conf.GetBaseUrl() + link,
					LastMod:    &now,
					ChangeFreq: sitemap.Daily,
				})

			}

			file, err := os.Create(sitemapPath + "/sitemap.xml")
			if err != nil {
				utility.Exit(err)
			}

			sm.WriteTo(file)
			file.Close()
		}

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		logger.Log("[main] done in %s with '%s'", duration, conf.GetBaseUrl())
	}

	select {
	case <-time.After(1 * time.Second):
		utility.Finish("[main] exit")
	}
}
