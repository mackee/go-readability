# go-readability

A Go implementation of Mozilla's Readability library, inspired by [@mizchi/readability](https://github.com/mizchi/readability). This library extracts the main content from web pages, removing clutter like navigation, ads, and unnecessary elements to provide a clean reading experience.

## Installation

```bash
go get github.com/mackee/go-readability
```

## Usage

### As a Library

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mackee/go-readability"
)

func main() {
	// Fetch a web page
	resp, err := http.Get("https://example.com/article")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse and extract the main content
	options := readability.DefaultOptions()
	article, err := readability.Extract(string(body), options)
	if err != nil {
		log.Fatal(err)
	}

	// Access the extracted content
	fmt.Println("Title:", article.Title)
	fmt.Println("Byline:", article.Byline)

	// Get content as HTML
	if article.Root != nil {
		html := readability.ToHTML(article.Root)
		fmt.Println("HTML Content:", html)
	}

	// Convert to Markdown
	if article.Root != nil {
		markdown := readability.ToMarkdown(article.Root)
		fmt.Println("Markdown Content:", markdown)
	}
}
```

### Using the CLI Tool

The package includes a command-line tool that can extract content from a URL:

```bash
# Install the CLI tool
go install github.com/mackee/go-readability/cmd/readability@latest

# Extract content from a URL
readability https://example.com/article

# Save the extracted content to a file
readability https://example.com/article > article.html

# Output as markdown
readability --format markdown https://example.com/article > article.md

# Output metadata as JSON
readability --metadata https://example.com/article
```

## Features

- Extracts the main content from web pages
- Removes clutter like navigation, ads, and unnecessary elements
- Preserves important images and formatting
- Extracts metadata (title, byline, excerpt, etc.)
- Supports output in HTML or Markdown format
- Command-line interface for easy content extraction

## Testing

This library uses test fixtures based on [Mozilla's Readability](https://github.com/mozilla/readability) test suite. Currently, we have implemented a subset of the test cases, with the source HTML files being identical to the original Mozilla implementation.

### Test Fixtures

The test fixtures in `testdata/fixtures/` are sourced from Mozilla's Readability test suite, with some differences:

- The source HTML files (`source.html`) are identical to Mozilla's Readability
- The expected output HTML (`expected.html`) may differ due to implementation differences between JavaScript and Go
- The expected metadata extraction results are aligned with Mozilla's implementation where possible

While not all test cases from Mozilla's Readability are currently implemented, using the same source HTML helps ensure that:

1. The Go implementation handles the same input as the JavaScript implementation
2. Regressions can be easily detected
3. Users can trust the library to process the same types of content as Mozilla's Readability

### Fixture Licensing

- `testdata/fixtures/001`: © Nicolas Perriault, [CC BY-SA 3.0](http://creativecommons.org/licenses/by-sa/3.0/)

These fixtures are identical to those used in Mozilla's Readability implementation.

## License

[Apache License 2.0](LICENSE)
