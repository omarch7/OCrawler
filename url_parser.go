package main

import (
	"regexp"
	"fmt"
)

type ParsedURL struct {
	RootUrl string
	Domain string
	URI string
	Vars string
	URLRegExp *regexp.Regexp
}

func UnpackURL(matches []string) (string, string, string) {
	return matches[1], matches[2], matches[3]
}

// It is easier to check the parsed URL using a custom data structure, to check if belongs to the same domain
// and to verify if the page has the correct format
func NewParsedURL(rootUrl string, url string, urlRegExp *regexp.Regexp) *ParsedURL {
	parsedURL := new(ParsedURL)
	parsedURL.RootUrl = rootUrl
	parsedURL.Domain,
	parsedURL.URI,
	parsedURL.Vars = UnpackURL(urlRegExp.FindStringSubmatch(url))
	return parsedURL
}

func (parsedURL *ParsedURL) IsSameDomain() bool  {
	return parsedURL.Domain == parsedURL.RootUrl || parsedURL.Domain == ""
}

func (parsedURL *ParsedURL) GetURL() string  {
	var domain string
	if parsedURL.IsSameDomain() {
		domain = parsedURL.RootUrl
	}else{
		domain = parsedURL.Domain
	}
	return fmt.Sprintf("%s%s", domain, parsedURL.URI)
}

func (parsedURL *ParsedURL) URIHasExtension(extRegEx *regexp.Regexp) (bool, string)  {
	var ext string
	hasExt := false
	if extRegEx.MatchString(parsedURL.URI) {
		ext = extRegEx.FindStringSubmatch(parsedURL.URI)[1]
		if ext != "" {
			hasExt = true
		}
	}
	return hasExt, ext
}