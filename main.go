package main

import (
	"./lib"
	"encoding/json"
	"fmt"
	"github.com/snabb/sitemap"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	BaseUrl string   `json: "baseUrl"`
	PathXml string   `json: "pathXml"`
	Add     []string `json: "add"`
}

func main() {

	jsonFile, e := os.Open("config.json")

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

	var formated string
	now := time.Unix(0, 0).UTC()

	for _, config := range configs {

		t := time.Now()

		formated = fmt.Sprintf("[%02d:%02d]", t.Minute(), t.Second())

		fmt.Println(formated, "crawling:", config.BaseUrl)

		crawler := new(webcrawler.WebCrawler)

		crawler.Run(config.BaseUrl)

		for _, url := range config.Add {
			crawler.AddLink(url)
		}

		// create sitemap
		sm := sitemap.New()

		for _, link := range crawler.GetLinks() {
			sm.Add(&sitemap.URL{
				Loc:        config.BaseUrl + link,
				LastMod:    &now,
				ChangeFreq: sitemap.Daily,
			})

		}

		// check if dir exist
		if _, e := os.Stat(config.PathXml); os.IsNotExist(e) {
			os.MkdirAll(config.PathXml, 0755)
		}

		// write sitemap.xml
		file, e := os.Create(config.PathXml + "sitemap.xml")

		if e != nil {
			panic(e)
		}

		defer file.Close()

		t = time.Now()
		formated = fmt.Sprintf("[%02d:%02d]", t.Minute(), t.Second())
		fmt.Println(formated, "writing:", config.PathXml+"sitemap.xml")

		sm.WriteTo(file)

	}

	end := time.Now()
	formated = fmt.Sprintf("[%02d:%02d]", end.Minute(), end.Second())
	fmt.Println(formated, "Done")

}
