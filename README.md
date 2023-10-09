# ZenRows + OpenAI GPT Experiment: Document Scraper

This repository, `documents-scraper`, is an experiment combining the scraping capabilities of ZenRows with the natural language processing strengths of OpenAI's GPT. It aims to scrape content from a website, summarize it using GPT, and then provide insights on a company's health based on these summaries.

## Prerequisites

- Go (Golang) installed on your system.
- API keys for both ZenRows and OpenAI.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/renatoaraujo/documents-scraper.git
cd documents-scraper
```

2. Install the required Go packages:

```bash
go get -u github.com/PuerkitoBio/goquery
go get -u github.com/renatoaraujo/go-zenrows
go get -u github.com/sashabaranov/go-openai
```

## Usage

1. Set up the necessary environment variables:

```bash
export ZENROWS_API_KEY=your_zenrows_api_key
export OPENAI_API_KEY=your_openai_api_key
```

2. Navigate to the `cmd/documentscraper` directory:

```bash
cd cmd/documentscraper
```

3. Run the experiment:

```bash
go run main.go
```

This will scrape content from the specified website, generate summaries, and then provide an inference on the company's health based on the gathered summaries.

## Notes

This is an experimental project and is meant for demonstration and experimentation purposes only.
