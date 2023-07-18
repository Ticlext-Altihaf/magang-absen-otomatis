package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// return array of string or error if fail to parse
func parse_html_text(html string, selector string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var result []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		result = append(result, text)
	})
	return result, nil
}
