package main

import (
	"os"
	"log"
	"regexp"
	"strconv"
	"runtime"
	"time"
)

func isValidUrlFormat(url string) bool {
	r, _ := regexp.Compile(`^([\w]+)\.?([\w])+(\.[a-zA-Z]{2,11})+$`)
	log.Printf("Checking %s", url)
	return r.MatchString(url)
}

func main()  {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatal("Didn't provide enough arguments")
	}
	if !isValidUrlFormat(args[0]) {
		log.Fatalln("Invalid Domain, Do not include the http or slash at the end of the domain")
	}
	depth, error := strconv.Atoi(args[1])
	if error != nil {
		log.Fatalln("You should specify the depth with an integer")
	}
	n, error := strconv.Atoi(args[2])
	if error != nil {
		log.Fatalln("You should specify the max of processes with an integer")
	}
	runtime.GOMAXPROCS(n)
	crawler := NewCrawler(args[0], depth)

	crawler.Begin()
	time.Sleep(1 * time.Second)
	crawler.PrintSiteMap()
}
