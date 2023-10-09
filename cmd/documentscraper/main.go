package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/renatoaraujo/go-zenrows"
	"github.com/sashabaranov/go-openai"
)

func main() {
	key := getZenrowsAPIKey()
	client := createZenrowsClient(key)

	content := scrapeContent(client)
	doc := parseContent(content)

	links := extractLinks(doc)

	openaiClient := createOpenAIClient()
	var wg sync.WaitGroup
	summaries := make(chan string, len(links))

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			scrapedContent := scrapeContentForLink(client, link)
			summary := getSummary(openaiClient, scrapedContent)
			summaries <- summary
		}(link)
	}

	go func() {
		wg.Wait()
		close(summaries)
	}()

	var allSummaries []string
	for summary := range summaries {
		allSummaries = append(allSummaries, summary)
	}

	finalSummary := strings.Join(allSummaries, "\n")
	companyHealth := askCompanyHealth(openaiClient, finalSummary)
	fmt.Println(companyHealth)
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

func scrapeContentForLink(client *zenrows.Client, link string) string {
	content, err := client.Scrape(link)
	if err != nil {
		log.Println("error scraping link:", err)
		return ""
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

func createOpenAIClient() *openai.Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return client
}

func getSummary(client *openai.Client, content string) string {
	prompt := fmt.Sprintf("Provide a summary of the following content: %s", content)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Println("error getting summary:", err)
		return ""
	}
	return response.Choices[0].Message.Content
}

func askCompanyHealth(client *openai.Client, summaries string) string {
	prompt := fmt.Sprintf("Based on the following summaries, what can be inferred about the health of the company?\n%s", summaries)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Println("error asking about company health:", err)
		return ""
	}
	return response.Choices[0].Message.Content
}
