package main

import (
	"testing"
	"os"
	"bufio"
	"strings"
	"strconv"
)

func check(e error)  {
	if e != nil {
		panic(e)
	}
}

// Unit Test to verify the URL REGEX is valid
func TestIsValidURL(t *testing.T) {
	file, err := os.Open("valid_urls.txt")
	defer file.Close()
	check(err)
	scanner := bufio.NewScanner(file)
	parser := NewParser()
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ", ")
		isValid, _ := strconv.ParseBool(line[1])
		if isV := parser.IsValidURL(line[0]); isValid != isV {
			t.Errorf("Expected %v, got %v: %s", isValid, isV, line[0])
		}
	}
}