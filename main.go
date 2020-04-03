package main

import (
	"./lib"
	"encoding/json"
	"fmt"
	"github.com/snabb/sitemap"
	"io/ioutil"
	"os"
	"path/filepath"
	str "strings"
	"time"
)

type Config struct {
	BaseUrl string   `json: "baseUrl"`
	PathXml string   `json: "pathXml"`
	Add     []string `json: "add"`
}

func main() {

	dir, e := filepath.Abs(filepath.Dir(os.Args[0]))

	if e != nil {
		panic(e)
	}

	jsonFile, e := os.Open(dir + "/config.json")

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
	now := time.Now().Local()

	for _, config := range configs {

		t := time.Now()

		formated = fmt.Sprintf("[%02d:%02d:%02d]", t.Hour(), t.Minute(), t.Second())

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

		var xmlPath string

		if str.HasPrefix(config.PathXml, ".") {
			xmlPath = dir + "/"
		}

		xmlPath = xmlPath + config.PathXml
		xmlPath, e := filepath.Abs(xmlPath)

		if e != nil {
			panic(e)
		}

		// check if dir exist
		if _, e := os.Stat(xmlPath); os.IsNotExist(e) {
			os.MkdirAll(xmlPath, 0775)
		}

		// write sitemap.xml
		file, e := os.Create(xmlPath + "/sitemap.xml")

		if e != nil {
			panic(e)
		}

		defer file.Close()

		t = time.Now()
		formated = fmt.Sprintf("[%02d:%02d:%02d]", t.Hour(), t.Minute(), t.Second())
		fmt.Println(formated, "writing:", xmlPath + "/sitemap.xml")

		sm.WriteTo(file)

	}

	end := time.Now()
	formated = fmt.Sprintf("[%02d:%02d:%02d]", end.Hour(), end.Minute(), end.Second())
	fmt.Println(formated, "Done")

}
