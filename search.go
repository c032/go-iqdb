package iqdb

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const homeURL = "https://iqdb.org/"

func init() {
	_, err := url.Parse(homeURL)
	if err != nil {
		panic(err)
	}
}

func SearchURL(imageURL string) ([]Match, error) {
	var (
		err error

		searchURL *url.URL
		doc       *goquery.Document
	)

	searchURL, err = url.Parse(homeURL)
	if err != nil {
		return nil, err
	}

	q := searchURL.Query()
	q.Set("url", imageURL)

	searchURL.RawQuery = q.Encode()

	var resp *http.Response

	resp, err = http.Get(searchURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Url = searchURL

	matches := matchesFromDocument(doc)

	return matches, nil
}

func matchesFromDocument(doc *goquery.Document) []Match {
	var matches []Match

	doc.Find("#pages > div > table").Each(func(i int, table *goquery.Selection) {
		if i == 0 {
			return
		}

		var (
			err error

			sauceURL   *url.URL
			similarity float64
		)

		a := table.Find(".image a")
		href := a.AttrOr("href", "")
		if href == "" {
			return
		}

		sauceURL, err = doc.Url.Parse(href)
		if err != nil {
			return
		}

		tr := table.Find("tr:last-of-type").First()
		if tr.Length() != 1 {
			return
		}

		similarityText := strings.TrimSpace(tr.Text())
		parts := strings.Split(similarityText, " ")
		if len(parts) == 0 || !strings.HasSuffix(parts[0], "%") {
			return
		}

		similarityStr := strings.TrimSuffix(parts[0], "%")

		similarity, err = strconv.ParseFloat(similarityStr, 64)
		if err != nil {
			return
		}

		matches = append(matches, Match{
			URL:        sauceURL.String(),
			Similarity: similarity / 100,
		})
	})

	return matches
}
