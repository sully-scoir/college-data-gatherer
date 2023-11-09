package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	s "strings"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("College Data Gatherer - Admissions")

	collegesCsv := "college-admissions.csv"
	colleges := readCollegesCsv(collegesCsv)

	// fmt.Println("colleges:", colleges)

	for _, c := range colleges {
		crawlCollege(c.Domain)
	}
}

type College struct {
	Name   string
	Domain string
}

func readCollegesCsv(collegesCsv string) []College {
	f, err := os.Open(collegesCsv)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	colleges := []College{}
	headerSkipped := false
	trimHttpWww := func(domain string) string {
		domain = s.TrimPrefix(domain, "http://")
		domain = s.TrimPrefix(domain, "https://")
		domain = s.TrimPrefix(domain, "www.")
		return domain
	}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if !headerSkipped {
			headerSkipped = true
			continue
		}

		college := College{}
		college.Name = row[0]
		college.Domain = trimHttpWww(row[2])
		colleges = append(colleges, college)
	}

	return colleges
}

func crawlCollege(collegeDomain string) {
	startsWithHttpsRegExp, _ := regexp.Compile("^https")
	c := colly.NewCollector(
		colly.AllowedDomains(collegeDomain, "www."+collegeDomain),
		colly.CacheDir("cache/"),
		colly.URLFilters(startsWithHttpsRegExp),
		// colly.Debugger(&debug.LogDebugger{}),
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

		// fmt.Printf("Link : %q -> %s\n", text, link)

		c.Visit(e.Request.AbsoluteURL(link))
	})

	deadlineTerm := "deadline"
	deadlineUrls := map[string]bool{}
	c.OnHTML("p, h1, h2, h3, h4, h5, h6", func(e *colly.HTMLElement) {
		if s.Contains(s.ToLower(e.Text), s.ToLower(deadlineTerm)) {
			fmt.Printf("Term Match: %q -> <%s> %s\n", deadlineTerm, e.Name, e.Text)
			deadlineUrls[e.Request.URL.String()] = true
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error:", e)
	})

	c.Visit("https://" + collegeDomain)

	fmt.Println("URLs with content matching Deadline term")
	for u, _ := range deadlineUrls {
		fmt.Println(u)
	}
}
