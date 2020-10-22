// main.go
package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

func extractHtmlVersionFromDoctype(docType string) string {
	res := ""

	if docType == "html" {
		res = "HTML 5"
	} else {
		// regular expression match
		re := regexp.MustCompile(`HTML [-]?\d[\d,]*[\.]?[\d{2}]*`)
		if re.MatchString(docType) {
			submatchall := re.FindAllString(docType, -1)
			res = submatchall[0]
		}
	}
	return res
}

func getDoctype(url string) string {
	var res string
	resp, _ := http.Get(url)
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.DoctypeToken:
			// tokenize the html file and extract the doctype
			res = strings.TrimSpace(token.Data)
			break
		}
	}
	return extractHtmlVersionFromDoctype(res)
}

func goQueryDocHandler(url string) *goquery.Document {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	return doc
}

func getPageTitle(url string) string {
	doc := goQueryDocHandler(url)
	if doc.Find("title").Text() == "" {
		return "No title is included in this webpage!"
	} else {
		return doc.Find("title").Text()
	}
}

func getHeadingsCount(url string) int {

	level := 0
	resp, _ := http.Get(url)
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()
		err := tokenizer.Err()
		if err == io.EOF {
			break
		}
		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken, html.SelfClosingTagToken:
			tag := token.Data
			// get the highest heading level e.g. h5 > h4 > h3 ...
			match, _ := regexp.MatchString(`h\d`, tag)
			if match {
				count, _ := strconv.Atoi(tag[1:])
				if count > level {
					level = count
				}
			}
		}
	}
	return level
}

func getAllLinks(url string) []string {

	doc := goQueryDocHandler(url)
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		band, ok := s.Attr("href")
		if ok {
			links = append(links, band)
		}
	})
	return links
}

func isinternallinks(link string) bool {
	// internal links start with hashtag
	if len(link) > 0 && link[0] == '#' {
		return true
	} else {
		return false
	}
}

func isexternallinks(link string, inputURL string) bool {
	u, _ := url.Parse(link)
	inputU, _ := url.Parse(inputURL)

	if u.Host == "" {
		// if host domain is empty, it is not an enteral link
		// e.g. href=#section1, href="/html/test.php"
		// href="href="mailto:mail@way2tutorial.com?cc=xyz@mail.com&bcc=abc@mail.com&subject=Feedback&body=Message"
		return false
	} else if u.Host == inputU.Host {
		// if the link has same host domain with the webpage, it is not an external link
		return false
	} else {
		return true
	}
}

func getNumberOfInternalAndExternalLinks(url string) (int, int) {
	links := getAllLinks(url)

	internalLinks := 0
	externalLinks := 0

	for _, link := range links {
		if isinternallinks(link) {
			internalLinks++
		} else if isexternallinks(link, url) {
			externalLinks++
		}
	}

	return internalLinks, externalLinks
}

func findElementID(url string, id string) bool {
	exist := false
	// TODO: not pass the test
	doc := goQueryDocHandler(url)
	doc.Find(id).Each(func(i int, s *goquery.Selection) {
		exist = true
	})
	return exist
}

func isExternalLinkAccessible(link string) bool {
	resp, _ := http.Get(link)
	// accessible external links
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func getNumberOfAccessiblelinks(url string) int {
	links := getAllLinks(url)

	accessibleLinks := 0
	for _, link := range links {
		if isinternallinks(link) && findElementID(url, link) {
			// interanl link #id exists
			accessibleLinks++
		} else if isexternallinks(link, url) && isExternalLinkAccessible(link) {
			// external link is accessible
			accessibleLinks++
		}
	}

	return accessibleLinks
}

func containsLoginForm(url string) bool {
	res := false
	// if a input element with type of password, it is very likely that
	// it is a login form
	doc := goQueryDocHandler(url)
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		inputType, ok := s.Attr("type")
		if ok {
			if inputType == "password" {
				res = true
			}
		}
	})

	return res
}

var router *gin.Engine

func main() {
	//url := "http://go-colly.org/"

	// Set the router as the default one provided by Gin
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"htmlversion":             "?",
			"urltitle":                "?",
			"headingslevel":           "?",
			"numberofinternallinks":   "?",
			"numberofexternallinks":   "?",
			"numberofaccessiblelinks": "?",
			"loginform":               "?",
		},
		)

	})

	router.POST("/test", func(c *gin.Context) {
		url := c.DefaultPostForm("url", "")
		internalLinks, externalLinks := getNumberOfInternalAndExternalLinks(url)
		LoginForm := ""
		if containsLoginForm(url) {
			LoginForm = "Yes!"
		} else {
			LoginForm = "No!"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"htmlversion":             getDoctype(url),
			"urltitle":                getPageTitle(url),
			"headingslevel":           getHeadingsCount(url),
			"numberofinternallinks":   internalLinks,
			"numberofexternallinks":   externalLinks,
			"numberofaccessiblelinks": getNumberOfAccessiblelinks(url),
			"loginform":               LoginForm,
		},
		)

	})

	// Start serving the application
	router.Run()

}
