package main

import (
	"fmt"
	"regexp"
	s "strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("College Data Gatherer - Admissions")

	startsWithHttpsRegExp, _ := regexp.Compile("^https")

	c := colly.NewCollector(
		colly.AllowedDomains("drexel.edu"),
		colly.CacheDir("cache/"),
		colly.URLFilters(startsWithHttpsRegExp),
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

	earlyDecisionTerm := "early decision"
	earlyDecisionUrls := map[string]bool{}
	c.OnHTML("h1, h2, h3, h4, h5, h6", func(e *colly.HTMLElement) {
		if s.Contains(s.ToLower(e.Text), s.ToLower(earlyDecisionTerm)) {
			fmt.Printf("Term Match: %q -> <%s> %s\n", earlyDecisionTerm, e.Name, e.Text)
			earlyDecisionUrls[e.Request.URL.String()] = true
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://drexel.edu/admissions")

	fmt.Println("URLs with content matching Early Decision term")
	for u, _ := range earlyDecisionUrls {
		fmt.Println(u)
	}
}
