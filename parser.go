package main

import (
	"regexp"
	"io"
	"golang.org/x/net/html"
)

type Parser struct {
	ValidUrlRegEx *regexp.Regexp
	ExtRegEx *regexp.Regexp
	Tags map[string]string
	Extensions map[string]bool
}

// The regular expression checks if the URL is valid, and at the same time extracts the sub-domain, domain, uri and variables
func NewParser() *Parser  {
	parser := new(Parser)
	parser.ValidUrlRegEx, _ = regexp.Compile(`^(?:(?:https?:)?(?://)?(?:([\w]+\.))*((?:\w+)\.[a-zA-Z]{2,11}))?((?:/[\w-_.+]*)+)?(#[\w]+|\?(?:\w+=\w+&?)+)?$`)
	parser.ExtRegEx, _ = regexp.Compile(`\.([a-z]{2,4})$`)
	parser.Tags = map[string]string{
		"a": "href",
		"img": "src",
		"script": "src",
		"link": "href",
		"iframe": "src",
		"video": "src"}
	parser.Extensions = map[string]bool{
		"html": true,
		"htm": true,
		"asp": true,
		"aspx": true,
		"php": true,
		"php3": true,
		"php4": true,
		"xhtml": true,
		"jhtml": true,
		"cgi": true,
		"shtml": true}
	return parser
}

func (parser *Parser) IsValidURL(url string) bool  {
	return parser.ValidUrlRegEx.MatchString(url)
}

func (parser *Parser) GetAttr(attr string, attrs []html.Attribute) string {
	var val string
	for _, a := range attrs {
		if a.Key == attr {
			val = a.Val
			break
		}
	}
	return val
}

func (parser *Parser) ParseBody(rootURL string, body io.ReadCloser) (map[string]*ParsedURL, map[string]*ParsedURL) {
	links := make(map[string]*ParsedURL)
	assets := make(map[string]*ParsedURL)
	tokens := html.NewTokenizer(body)
	for {
		token := tokens.Next()
		switch token {
		case html.ErrorToken:
			return links, assets
		case html.StartTagToken:
			tag := tokens.Token()
			if attr, ok := parser.Tags[tag.Data]; ok {
				if url := parser.GetAttr(attr, tag.Attr); url != "" {
					if parser.IsValidURL(url){
						parsedURL := NewParsedURL(rootURL, url, parser.ValidUrlRegEx)
						if tag.Data == "a" {
							if hasExt, ext := parsedURL.URIHasExtension(parser.ExtRegEx); hasExt {
								if parser.Extensions[ext] {
									if parsedURL.IsSameDomain() && !parsedURL.IsSubDomain() {
										links[parsedURL.GetURL()] = parsedURL
									}
								}
							}else{
								if parsedURL.IsSameDomain() && !parsedURL.IsSubDomain() {
									links[parsedURL.GetURL()] = parsedURL
								}
							}
						}else{
							assets[parsedURL.GetURL()] = parsedURL
						}
					}
				}
			}
		}
	}
}