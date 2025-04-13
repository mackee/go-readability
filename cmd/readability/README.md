# readability CLI

A command-line tool for extracting the main content from web pages using the go-readability library.

## Installation

```bash
go install github.com/mackee/go-readability/cmd/readability@latest
```

## Usage

```bash
readability [options] <url>
```

### Options

- `--format <format>`: Output format (html or markdown, default: html)
- `--metadata`: Output metadata as JSON instead of content
- `--help`: Show help message

### Examples

Extract content from a URL and output as HTML:
```bash
readability https://example.com/article
```

Extract content and output as Markdown:
```bash
readability --format markdown https://example.com/article
```

Extract metadata only:
```bash
readability --metadata https://example.com/article
```

Save the extracted content to a file:
```bash
readability https://example.com/article > article.html
```

## Output

By default, the tool outputs the extracted HTML content to stdout. You can redirect this to a file if needed.

When using the `--metadata` flag, the tool outputs a JSON object containing metadata such as title, byline, and other information about the article.
