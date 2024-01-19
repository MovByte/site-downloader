package getLinksFrom

func SiteMap(origin string, links string[]) string {
	res, err := http.Get(origin + "/sitemap.xml")

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var sitemap Urlset
	if err := xml.Unmarshal(bodyBytes, &sitemap); err != nil {
		log.Fatal(err)
	}

	for _, url := range sitemap.Url {
		links = append(links, url.Loc)
	}
}

func WebAppManifest(origin string, links string[]) string {
	var manifest WebAppManifest
	
	err = json.Unmarshal(body, &manifest)
	if err != nil {
	 log.Fatal(err)
	}
	
	if manifest.Scope != "" || manifest.StartURL != "" {
	   links = append(links, manifest.Scope, manifest.StartURL)
	}
	
	if len(manifest.Icons) > 0 {
	   for _, icon := range manifest.Icons {
		   links = append(links, icon.Src)
	   }
	}
	
	if len(manifest.RelatedApplications) > 0 {
	   for _, app := range manifest.RelatedApplications {
		   links = append(links, app.Platform.URL)
	   }
	}
}

func cacheManifestAddToArr(link string, links string[]) {
	if !strings.Contains(link, "*")
		links = append(links, link)
}

func CacheManifest(origin string, pathToManifest string, links string[]) string {
	res, err := http.Get(origin + "/sitemap.xml")

	if (res.Headers.get("content-type") == "application/xml") {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(body, "\n")

		var currentDirective string
		for i, line := range lines {
			if strings.HasPrefix(line, "#") {
				// Ignore comments
			} else if strings.Contains(line, "CACHE:") || strings.Contains(line, "NETWORK:") {
				currentDirective = line
			} else if currentDirective == "CACHE:" || currentDirective == "NETWORK:" {
				parts := strings.Split(line, " ")
				cacheManifestAddToArr(parts[0], links)
			} else if currentDirective == "FALLBACK:" {
				parts := strings.Split(line, " ")
				fmt.Sprintf("%s %s", rewritePath(parts[0]), rewritePath(parts[1],))
			}
		}
	}
}

// TODO: RSS Feeds
// TODO: Atom Feeds
// TODO: Robots.txt