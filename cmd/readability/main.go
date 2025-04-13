package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mackee/go-readability"
)

func main() {
	// Define command-line flags
	formatFlag := flag.String("format", "html", "Output format: html or markdown")
	metadataFlag := flag.Bool("metadata", false, "Output metadata as JSON instead of content")
	helpFlag := flag.Bool("help", false, "Show help")
	flag.Parse()

	// Show help if requested or no URL provided
	if *helpFlag || flag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	// Get the URL from command-line arguments
	url := flag.Arg(0)

	// Fetch the content from the URL
	article, err := fetchAndParse(url)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Output based on flags
	if *metadataFlag {
		// Output metadata as JSON
		metadata := map[string]string{
			"title":     article.Title,
			"byline":    article.Byline,
			"nodeCount": fmt.Sprintf("%d", article.NodeCount),
			"pageType":  string(article.PageType),
		}
		jsonData, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	} else {
		// Output content in the specified format
		switch strings.ToLower(*formatFlag) {
		case "html":
			if article.Root != nil {
				fmt.Println(readability.ToHTML(article.Root))
			} else {
				log.Fatalf("No content was extracted from the URL")
			}
		case "markdown":
			if article.Root != nil {
				fmt.Println(readability.ToMarkdown(article.Root))
			} else {
				log.Fatalf("No content was extracted from the URL")
			}
		default:
			log.Fatalf("Unknown format: %s", *formatFlag)
		}
	}
}

// fetchAndParse fetches content from a URL and parses it using readability
func fetchAndParse(url string) (*readability.ReadabilityArticle, error) {
	// Fetch the content
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the content
	options := readability.DefaultOptions()
	article, err := readability.Extract(string(body), options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}

	return &article, nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Usage: readability [options] <url>")
	fmt.Println("\nOptions:")
	fmt.Println("  --format <format>  Output format: html or markdown (default: html)")
	fmt.Println("  --metadata         Output metadata as JSON instead of content")
	fmt.Println("  --help             Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  readability https://example.com/article")
	fmt.Println("  readability --format markdown https://example.com/article")
	fmt.Println("  readability --metadata https://example.com/article")
}
