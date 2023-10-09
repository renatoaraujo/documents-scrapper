package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/renatoaraujo/go-zenrows"
)

func main() {
	key := os.Getenv("ZENROWS_API_KEY")
	hc := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}
	client := zenrows.NewClient(hc).WithApiKey(key)

	woodboisPage := "https://www.woodbois.com/investors/"
	content, err := client.Scrape(woodboisPage)
	if err != nil {
		panic(err)
	}

	fmt.Println(content)
}
