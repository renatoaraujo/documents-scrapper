package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/renatoaraujo/go-zenrows"
)

func main() {
	key := getZenrowsAPIKey()
	client := createZenrowsClient(key)

	content := scrapeContent(client)
	doc := parseContent(content)

	links := extractLinks(doc)
	for _, link := range links {
		fmt.Println(link)
	}
}

func getZenrowsAPIKey() string {
	key := os.Getenv("ZENROWS_API_KEY")
	if key == "" {
		log.Fatal("ZENROWS_API_KEY environment variable is not set")
	}
	return key
}

func createZenrowsClient(key string) *zenrows.Client {
	hc := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}
	client := zenrows.NewClient(hc).WithApiKey(key)

	return client
}

func scrapeContent(client *zenrows.Client) string {
	jsInstructions := `[
  {"evaluate": "if ($.fn.dataTable.isDataTable('#diversified_gas_and_oil_rns') && $('#diversified_gas_and_oil_rns').DataTable().ajax.json()) { document.body.setAttribute('data-table-loaded', 'true'); }"},
  {"wait_for": "body[data-table-loaded='true']"}
]
`
	content, err := client.Scrape("https://polaris.brighterir.com/public/diversified_gas_and_oil/news/rns", zenrows.WithJSInstructions(jsInstructions))
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func parseContent(content string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func extractLinks(doc *goquery.Document) []string {
	var links []string
	doc.Find(".dataTables_wrapper a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && !strings.HasSuffix(href, "/export") {
			links = append(links, href)
		}
	})
	return links
}
