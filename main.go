package main

import (
	"context"
	"log"
	"time"

	config "github.com/MovByte/site-downloader/getConfig"
	getLinksFrom "github.com/MovByte/site-downloader/getLinksFrom"
	resourceAttrsMap "github.com/MovByte/site-downloader/resourceAttrsMap"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// TODO: Use slog

// These are only partially implemented, for what is needed
type SitemapIndex struct {
	Loc string `xml:"loc"`
}
type Urlset struct {
	Url []SitemapIndex `xml:"url"`
}

type WebAppManifest struct {
	Scope string `json:"scope"`
	StartURL string `json:"start_url"`
	Icons []struct {
		Src string `json:"src"`
	} `json:"icons"`
	RelatedApplications []struct {
		Platform struct {
			URL string `json:"url"`
		} `json:"platform"`
	} `json:"related_applications"`
}

var links []string

func main() {
	config := config.GetConfig()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// listen for network events
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			log.Printf("request %s %s", ev.Request.Method, ev.Request.URL)
		}
	})

	for attr, selectors := range resourceAttrsMap.HTMLResourceSelectors {
		for _, selector := range selectors {
            err := chromedp.Run(ctx, chromedp.Nodes(selector, chromedp.Attr(attr)))
            if err!= nil {
                log.Fatal(err)
            }
        }
	}

	parsedUrl, err := url.parse(config.Website)

	if err != nil {
		log.Fatal(err)
	}

	getLinksFrom.SiteMap(parsedUrl.origin, links)

	cacheManifestPath := ""
	manifestPath := "manifest.json"

	err = chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(config.Website),
		chromedp.WaitReady(`body`),
		// TODO: Meta Refresh Tag
		// TODO: https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls
		// TODO: Open Graph Tags
		chromedp.Evaluate(`document.getElementsByClassname("body")?.[0].getAttribute("manifest")`, manifestPath),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, links),
		// TODO: Get every form element and look at the action and method attributes
	)

	if (manifestPath != "") {
		getLinksFrom.WebAppManifest(parsedUrl.origin, manifestPath, links)
	}
	getLinksFrom.CacheManifest(parsedUrl.origin, manifestPath, links)

	if err != nil {
		log.Fatal(err)
	}

	for _, link := range links {
		err := chromedp.Run(ctx,
			chromedp.Click(link, chromedp.NodeVisible),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}