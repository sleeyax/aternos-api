package aternos_api

import (
	"github.com/PuerkitoBio/goquery"
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
