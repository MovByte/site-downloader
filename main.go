package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/gocolly/colly"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Verbose bool
	ErrorLogFile string
	OutDir string
	Website string
}

func main() {
	// Setup the config
	config := &Config{}

	// Check if the config file exists.
	_, err := os.Stat("config.toml")
	if (err == nil) {
		tree, err := toml.LoadFile("config.toml")
		if err != nil {
			panic(err)
		}

		err = tree.Unmarshal(config)
		if err != nil {
			panic(err)
		}
	}

	// Flags
	verbosePtr := flag.Bool("v", false, "Enable verbose output")
	errorLogFilePtr := flag.String("elog", "", "Path to the error log file")
	websitePtr := flag.String("website", "", "Website to crawl")
	outDirPtr := flag.String("outdir", "", "Output directory")
	flag.Parse()

	// Map the flags onto the config.
	config.Verbose = *verbosePtr
	config.ErrorLogFile = *errorLogFilePtr
	config.Website = *websitePtr
	config.OutDir = *outDirPtr

	// Setup logging
	// JSON error logging
	jsonFile, err := os.Create(config.ErrorLogFile)
	if (err != nil) {
		panic(err)
	}
	defer jsonFile.Close()

	// TODO: Use
	var handler slog.Handler
	if (config.ErrorLogFile != "") {
		options := &slog.HandlerOptions{}
		handler = slog.NewJSONHandler(jsonFile, options)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(config.Website),
	)

	// Where to find external references in HTML
	// attribute, selectors for it
	selectors := map[string]string{
		"src": "script[src], img[src], video[src], audio[src], embed[src], iframe[src], source[src], track[src], picture[srcset], meta[content]",
		"href": "a[href], link[rel='stylesheet'][href]",
		"data": "object[data]",
	}
	// TODO: Also parse XML, SVGs, and CSS for external references.

	for attr, selector := range selectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			e.Request.Visit(e.Attr(attr))
		})
	}

	c.OnResponse(func(r *colly.Response) {
		if (r.Headers.Get("content-type") == "text/javascript" || r.Headers.Get("content-type") == "application/javascript") {
			ext := filepath.Ext(r.Request.URL.Path)
			if (ext == ".js" || ext == ".mjs") {
				// TODO: Use https://github.com/BishopFox/jsluice to get the file references out of the JS files
			}
		}

		dlFile, err := os.Create(fmt.Sprintf("archive/%d/", time.Now().Unix())+ r.Request.URL.String())
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
}