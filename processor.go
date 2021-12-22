package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oglinuk/goccer"
	"golang.org/x/net/html"
)

// processor orchestrates the crawling and creation of document entries
// that are stored in an sqlite3 database
type processor struct {
	DomainCheck map[string]struct{}
	URLCheck map[string]struct{}
	Seeds []string

	Doc *html.Tokenizer
	Ttype html.TokenType
	T html.Token

	C *http.Client
	Wp *goccer.Workerpool
	Dbi *DBInstance
}

// NewProcessor constructor
func NewProcessor() *processor {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 15,
	}

	return &processor{
		DomainCheck: make(map[string]struct{}),
		URLCheck: make(map[string]struct{}),
		Seeds: []string{
			"https://fourohfournotfound.com",
			"https://fourohfournotfound.com",
			"https://en.wikipedia.org/wiki/Search_engine",
			"https://creativecommons.org",
		},
		Doc: nil,
		Ttype: 0,
		T: html.Token{},
		C: c,
		Wp: goccer.NewWorkerpool(),
		Dbi: NewDBInstance(),
	}
}

// iterate p.Doc TokenType and Token
func (p *processor) iterate() {
	p.Ttype = p.Doc.Next()
	p.T = p.Doc.Token()
}

// getMetadata extracts a title and description from an HTML page
func (p *processor) getMetadata(body io.ReadCloser) (string, string) {
	title := ""
	desc := ""

	p.Doc = html.NewTokenizer(body)
	for {
		p.iterate()
		if p.Ttype == html.ErrorToken {
			break
		}

		// Get title from <title> tag
		if p.Ttype == html.StartTagToken && p.T.Data == "title" {
			p.iterate()
			title = p.T.Data
		}

		// TODO: Add check for schema.org format

		// Check for open graph <meta property="og:..." content="..."> tags
		if p.Ttype == html.StartTagToken && p.T.Data == "meta" {
			prop := ""
			for _, a := range p.T.Attr {
				if a.Key == "content" && prop == "og:title" {
					title = a.Val
				}

				if a.Key == "content" && prop == "og:description" {
					desc = a.Val
				}

				prop = a.Val
			}
		}

		// If all else fails use <h1> for title
		if p.Ttype == html.StartTagToken && p.T.Data == "h1" && title == "" {
			p.iterate()
			title = p.T.Data
		}

		// If all else fails use the first <p> tag for desc
		if p.Ttype == html.StartTagToken && p.T.Data == "p" && desc == "" {
			for {
				p.iterate()
				if p.Ttype == html.TextToken {
					if desc == "" {
						desc = p.T.Data
					} else {
						desc = fmt.Sprintf("%s%s", desc, p.T.Data)
					}
				}

				if p.Ttype == html.EndTagToken && p.T.Data == "p" {
					desc = strings.Join(strings.Fields(desc), " ")
					break
				}

				if p.Ttype == html.ErrorToken {
					break
				}
			}
		}
	}

	return title, desc
}

// createDoc from given URL and insert into database
func (p *processor) createDoc(URL string) {
	resp, err := p.C.Get(URL)
	if err != nil {
		log.Printf("p.C.Get: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	title, desc := p.getMetadata(resp.Body)
	if title == "" && desc == "" {
		return
	}
	log.Printf("Processed %s | Title: %s | Desc: %s\n", URL, title, desc)
	doc := NewDocument(title, URL, desc)
	p.Dbi.Insert(doc)
}

// filter already visited URLs
func (p *processor) filter() []string {
	var filtered []string
	for _, URL := range p.Seeds {
		if _, exists := p.URLCheck[URL]; !exists {
			p.URLCheck[URL] = struct{}{}
			filtered = append(filtered, URL)
		}
	}

	return filtered
}

func (p *processor) process() {
	collected := p.Wp.Queue(p.filter())

	for _, c := range collected {
		parsed, err := url.Parse(c)
		if err != nil {
			log.Printf("url.Parse: %s", err.Error())
		}

		URL := fmt.Sprintf("http://%s", parsed.Hostname())

		if _, exists := p.DomainCheck[URL]; !exists {
			p.DomainCheck[URL] = struct{}{}
			p.createDoc(URL)
		}
	}

	p.Seeds = collected
}
