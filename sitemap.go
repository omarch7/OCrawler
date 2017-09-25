package main

import (
	"sync"
	"fmt"
	"strings"
	"github.com/beego/bee/logger/colors"
)

type Document struct {
	URI       string
	Depth     int
	Links []*Document
	Assets []string
	Visited bool
}

// Because only one domain is being crawled, the URI (Uniform Resource Identifier) is used instead
func NewDocument(uri string, depth int) *Document{
	document := new(Document)
	document.URI = uri
	document.Depth = depth
	document.Visited = false
	return document
}

// This is basically a Graph Data structure, instead of using an adjacency list,
// an array of references of documents is used
type SiteMap struct {
	RootURL string
	Documents map[string]*Document
	mux sync.Mutex
}

func NewSiteMap(rootURL string) *SiteMap {
	siteMap := new(SiteMap)
	siteMap.Documents = make(map[string]*Document)
	siteMap.RootURL = rootURL
	return siteMap
}

// This data structure implements the sync/mux object to control the locking and unlocking of processes
// this allow us to ensure the concurrent execution of the algorithm while using the data structure
func (siteMap *SiteMap) AddDocument(uri string, depth int) {
	if !siteMap.documentExists(uri) {
		siteMap.mux.Lock()
		siteMap.Documents[uri] = NewDocument(uri, depth)
		siteMap.mux.Unlock()
	}
}

func (siteMap *SiteMap) AddAsset(uri string, asset string)  {
	if siteMap.documentExists(uri) {
		siteMap.mux.Lock()
		siteMap.Documents[uri].Assets = append(siteMap.Documents[uri].Assets, asset)
		siteMap.mux.Unlock()
	}
}

func (siteMap *SiteMap) AddLink(uri string, link string, depth int)  {
	if siteMap.documentExists(uri) {
		siteMap.AddDocument(link, depth)
		siteMap.mux.Lock()
		siteMap.Documents[uri].Links = append(siteMap.Documents[uri].Links, siteMap.Documents[link])
		siteMap.mux.Unlock()
	}
}

func (siteMap *SiteMap) documentExists(uri string) bool {
	siteMap.mux.Lock()
	var exists bool
	if _, ok := siteMap.Documents[uri]; ok {
		exists = ok
	}
	defer siteMap.mux.Unlock()
	return exists
}

func (siteMap *SiteMap) Visit(uri string)  {
	siteMap.mux.Lock()
	siteMap.Documents[uri].Visited = true
	siteMap.mux.Unlock()
}

func (siteMap *SiteMap) HasBeenVisited(uri string) bool  {
	siteMap.mux.Lock()
	defer siteMap.mux.Unlock()
	return siteMap.Documents[uri].Visited
}

const edge string = "│   "
const space string = "\t"
const node string = "├── "

func (siteMap *SiteMap) PrintSiteMap()  {
	fmt.Println(colors.CyanBold(fmt.Sprintf("\n%s Site Map %s", strings.Repeat("*", 45),
		strings.Repeat("*", 45))))
	fmt.Println("\n"+colors.Magenta(siteMap.RootURL+"/"))
	PrintDocuments(siteMap.Documents["/"])
}

// To print, a tree is implemented, lower levels are called recursively
func PrintDocuments(document *Document) {
	for _, link := range document.Links {
		PrintSpaces(document.Depth)
		fmt.Print(colors.White(node))
		if link.Depth > document.Depth{
			fmt.Println(colors.Blue(link.URI))
			PrintDocuments(link)
		}else{
			fmt.Println(colors.Magenta(link.URI))
		}
	}
	for _, asset := range document.Assets {
		PrintSpaces(document.Depth)
		fmt.Printf("%s%s\n", colors.White(node), colors.Green(asset))
	}
}

func PrintSpaces(depth int)  {
	if depth > 0 {
		spaces := strings.Repeat(space, depth)
		fmt.Printf("%s%s", colors.White(edge), spaces)
	}
}