package config

import (
	"encoding/json"
	"errors"
	"os"
)

type config struct {
	BaseUrl        string   `json: "baseUrl"`
	SitemapUrlPath string   `json: "sitemapUrlPath"`
	SitemapPath    string   `json: "sitemapPath"`
	GetFrom        []string `json: "getFrom"`
	FilterRules    []string `json: "filterRules"`
}

func (_this *config) GetBaseUrl() string {
	return _this.BaseUrl
}

func (_this *config) GetSitemapUrlPath() string {
	return _this.SitemapUrlPath
}

func (_this *config) GetSitemapPath() string {
	return _this.SitemapPath
}

func (_this *config) GetGetFrom() []string {
	return _this.GetFrom
}

func (_this *config) GetFilterRules() []string {
	return _this.FilterRules
}

func LoadConfig(path string) ([]config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configs []config
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, errors.New("bad config")
	}

	return configs, nil
}
