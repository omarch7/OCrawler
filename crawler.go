package main

import (
	"net/http"
	"sync"
	"log"
)

type Crawler struct {
	Domain    string
	MaxDepth  int
	wg        sync.WaitGroup
	DocParser *Parser
	Site      *SiteMap
}

// Creates a New Crawler Object which contains a DocParser and a New Empty SiteMap Object
func NewCrawler(domain string, maxDepth int) *Crawler {
	crawler := new(Crawler)
	crawler.Domain = domain
	crawler.MaxDepth = maxDepth
	crawler.DocParser = NewParser()
	crawler.Site = NewSiteMap(domain)
	return crawler
}

// Initialize the Crawler starting from / that is the root of the domain
// We wait until all threads are over to end the execution of this block
func (crawler *Crawler) Begin()  {
	crawler.wg.Add(1)
	go crawler.Crawl("/", 0)
	crawler.wg.Wait()
}

// When crawls, it parses the data and extracts the URL address from Links and Assets
// This function is called recursively to extract from the child pages
// The depth is important to stop the execution when that depth is reached
func (crawler *Crawler) Crawl(uri string, depth int) {
	defer crawler.wg.Done()
	domain := crawler.Domain
	log.Printf("Depth: %v - %s%s Crawling...", depth, domain, uri)
	resp, err := http.Get("http://" + domain + uri)
	if err != nil {
		log.Printf("Couldn't get %s%s", domain, uri)
		return
	}
	log.Printf("Depth: %v - %s%s Parsing...", depth, domain, uri)
	links, assets := crawler.DocParser.ParseBody(crawler.Domain, resp.Body)
	defer resp.Body.Close()
	log.Printf("Depth: %v - %s%s Extracted (%v links) and (%v assets)", depth, domain, uri, len(links), len(assets))
	crawler.Site.AddDocument(uri, depth)
	crawler.Site.Visit(uri)
	for asset := range assets {
		crawler.Site.AddAsset(uri, asset)
	}
	depth += 1
	for _, parsedLink := range links {
		if parsedLink.URI != "" {
			crawler.Site.AddLink(uri, parsedLink.URI, depth)
			if !crawler.Site.HasBeenVisited(parsedLink.URI) && depth <= crawler.MaxDepth {
				crawler.wg.Add(1)
				go crawler.Crawl(parsedLink.URI, depth)
			}
		}
	}
}

// Prints the Crawled Site Map
func (crawler *Crawler) PrintSiteMap() {
	crawler.Site.PrintSiteMap()
}