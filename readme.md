###### GO-Crawl
It is a crawler used for generating sitemap.xml.
That's it and that's all!


###### How to build a binary
Run the following commands to build a binary for your OS:

Linux: `$ env GOOS=linux GOARCH=amd64 go build -o go-crawl`

MacOS: `$ env GOOS=darwin GOARCH=amd64 go build -o go-crawl`

Windows: `$ env GOOS=windows GOARCH=amd64 go build -o go-crawl.exe`


###### Configuration
Near go-crawl binary must be a **config.json** file.
There are array of objects into it.

Example:

```
[
    {
        "baseUrl": "https://perm.medsi.ru",
        "pathXml": "../../../subdomains/perm/",
        "add": [
            "/any/additional/url/",
            "/any/other/additional/url/"
        ]
    }
]
```

**baseUrl** - protocol + host.

**pathXml** - absolute or relative path to where sitemap.xml will be written.

**add** - array of strings for additional urls.

###### Run
`$ go-crawl`

And it will be done.
