package main

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	// "log/slog"
	"os"
	"path/filepath"

	"github.com/BishopFox/jsluice"
	"github.com/gocolly/colly"
	"github.com/tdewolff/parse/css"
	// TODO: Import htmlCollection.go
	// TODO: Import getConfig.go
)

// TODO: For dynamic sites like BGS, combine colly with ferret for headless scraping
func main() {

	// Setup logging
	// JSON error logging
	if (config.ErrorLogFile != "") {
		// Create the file, if it doesn't already exist
		jsonFile, err := os.OpenFile(config.ErrorLogFile, os.O_CREATE|os.O_WRONLY, 0644)
		if (err != nil) {
			panic(err)
		}
		defer jsonFile.Close()
	}

	// TODO: Use
	/*
	var handler slog.Handler
	if (config.ErrorLogFile != "") {
		options := &slog.HandlerOptions{}
		handler = slog.NewJSONHandler(jsonFile, options)
	}
	*/

	parsedWebsiteURL, err := url.Parse(config.Website)
	if (err != nil) {
		panic(err)
	}

	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.IgnoreRobotsTxt(),
		colly.AllowedDomains(parsedWebsiteURL.Hostname()),
	)


	
	// Where to find external references in XML
	// attribute, selectors for it
	/*
	xmlSelectors := map[string]string{
		"url": "[src],[href],[action],[background],[cite],[classid],[codebase],[data],[longdesc],[profile],[usemap]",
	}
	*/
	
	for attr, selectors := range htmlSelectors {
		for _, selector := range selectors {
			c.OnHTML(selector, func(e *colly.HTMLElement) {
				e.Request.Visit(e.Attr(attr))
				c.Visit(e.Attr(attr))
			})
		}
	}
	// TODO: Parse the inline CSS
	// TODO: Parse the inline JS


	/*
	for attr, selector := range xmlSelectors {
		c.OnXML(selector, func(e *colly.XMLElement) {
			e.Request.Visit(e.Attr(attr))
		})
	}
	*/

	// TODO: Crawl SVG (xlink:href)

	crawlStartTime := time.Now().Unix()
	c.OnResponse(func(r *colly.Response) {
		if (r.Headers.Get("content-type") == "text/javascript" || r.Headers.Get("content-type") == "application/javascript") {
			ext := filepath.Ext(r.Request.URL.Path)
			if (ext == ".js" || ext == ".mjs") {
				analyzer := jsluice.NewAnalyzer(r.Body)

				for _, url := range analyzer.GetURLs() {
					r.Request.Visit(url.URL)
				}
			}
		}
		if (r.Headers.Get("content-type") == "text/css") {
			p := css.NewParser(bytes.NewReader(r.Body), false)
			for {
				gt, tt, data := p.Next()
				if gt == css.GrammarType(css.ErrorToken) || tt == css.ErrorToken {
					// TODO: Log that the CSS is broken
					break
				}
		 
				// Check if the token is a URL token
				if tt == css.URLToken {
					r.Request.Visit(string(data))
				}
			}
		}

		fmt.Println(r.Request.URL.String());

		path := strings.Trim(r.Request.URL.Path, "/")
		segments := strings.Split(path, "/")
		dir := filepath.Join(config.OutDir, strconv.FormatInt(crawlStartTime, 10), r.Request.URL.Hostname())
		fmt.Println(dir);
		for _, segment := range segments {
		   if segment != "" {
			   dir = filepath.Join(dir, segment)
			   if err := os.MkdirAll(dir, 0755); err != nil {
				   panic(err)
			   }
		   }
		}
		// Still create the base directory then
		if (len(segments) == 1) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				panic(err)
			}
		}

		networkFile := filepath.Base(r.Request.URL.Path)
		if (networkFile == "/") {
            networkFile = "/index.html"
        }
		fmt.Println(filepath.Join(dir, networkFile))
		dlFile, err := os.OpenFile(filepath.Join(dir, networkFile), os.O_CREATE|os.O_WRONLY, 0644)
		if (err != nil) {
			panic(err)
		}
		defer dlFile.Close()

		// Write the file to the disk
		_, err = io.Copy(dlFile, bytes.NewBuffer(r.Body))
		if (err != nil) {
			panic(err)
		}
	})

	// The intitial visit
	c.Visit(config.Website)
}
