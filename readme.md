# GO-Crawl
WebCrawler + Sitemap.xml writer.

## How to build a binary
Run the following commands to build a binary for your OS:

Linux:
`$ env GOOS=linux GOARCH=amd64 go build -o go-crawl`

MacOS:
`$ env GOOS=darwin GOARCH=amd64 go build -o go-crawl`

Windows:
`$ env GOOS=windows GOARCH=amd64 go build -o go-crawl.exe`


## Configuration
Near go-crawl binary must be a **config.json** file.
There are array of objects into it.

Example:

```
[
  {
    "baseUrl": "http://localhost.local",
    "sitemapUrlPath": "/sitemaps/",
    "sitemapPath": "~/web/go-crawl/sitemap/",
    "getFrom": [
    ],
    "filterRules": [
      "^/about/press-centr/press-relize",
      "special_version",
      "^/upload",
      "^/document.php",
      "^/review/\\?",
      "\\.pdf",
      "\\.PDF"
    ]
  }
]
```

**baseUrl** - protocol + host.

**sitemapUrlPath** - used for generating sitemap index, when total links more then 5000

**sitemapPath** - absolute or relative path to where sitemap.xml will be written.

**getFrom** - array of urls, where to crawl for new urls.

**filterRules** - array of regexp rules to skip urls. (case sensetive)


Add `User-Agent: GO-Crawl` into the table `b_stat_searcher` to skip bitrix throttler:
http://localhost.local/bitrix/admin/perfmon_table.php?lang=ru&table_name=b_stat_searcher

### Run
`$ go-crawl`

And it will be done.
