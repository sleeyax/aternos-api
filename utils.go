package aternos_api

import (
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"strings"
)

// findList finds and extracts a list of items from a HTML document.
//
// Example:
//
// <ul>
//   <li>foo</li>
//   <li>bar</li>
// </ul>
//
// findList(document, "li") -> ["foo", "bar"]
func findList(document *goquery.Document, selector string) []string {
	items := make([]string, 0)

	document.Find(selector).Each(func(i int, selection *goquery.Selection) {
		items = append(items, strings.TrimSpace(selection.Text()))
	})

	return items
}

// randomString generates a random lowercase string.
// E.g. mdlc2c9chx9, mbywjir33mm
func randomString(length int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	s := make([]rune, length)

	for i := range s {
		s[i] = charset[rand.Intn(len(charset))]
	}

	return string(s)
}

// GetStringInBetween finds a string between two substrings and returns an empty string if no string has been found.
func getStringInBetween(str string, left string, right string) string {
	s := strings.Index(str, left)
	if s == -1 {
		return ""
	}

	s += len(left)

	e := strings.Index(str, right)
	if e == -1 {
		return ""
	}

	return str[s:e]
}
