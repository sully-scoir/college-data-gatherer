package main

import (
	"fmt"
	s "strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("College Data Gatherer - Admissions")

	c := colly.NewCollector(
		colly.AllowedDomains("drexel.edu"),
	)

	admissionTextTerms := []string{
		"admission",
		"apply",
		"deadline",
	}
	matchesAdmissionTextTerms := func(text string) bool {
		for i := 0; i < len(admissionTextTerms); i++ {
			if s.Contains(s.ToLower(text), s.ToLower(admissionTextTerms[i])) {
				return true
			}
		}

		return false
	}
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		text := e.Text
		link := e.Attr("href")

		if !matchesAdmissionTextTerms(text) {
			return
		}

		fmt.Printf("Link : %q -> %s\n", text, link)

		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://drexel.edu/admissions")
}