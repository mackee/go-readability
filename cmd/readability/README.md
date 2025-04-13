# readability CLI

A command-line tool for extracting the main content from web pages using the go-readability library.

## Installation

```bash
go install github.com/mackee/go-readability/cmd/readability@latest
```

## Usage

```bash
readability [options] <url|file_path>
```

or

```bash
cat <file_path> | readability [options]
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

Extract content from a local file and output as HTML:
```bash
readability ./article.html
```

Extract content and output as Markdown:
```bash
readability --format markdown https://example.com/article
```

Extract metadata only:
```bash
readability --metadata https://example.com/article
```

Extract content from stdin and output as Markdown:
```bash
cat ./article.html | readability --format markdown
```

Save the extracted content to a file:
```bash
readability https://example.com/article > article.html
```

## Output

By default, the tool outputs the extracted HTML content to stdout. You can redirect this to a file if needed.

When using the `--metadata` flag, the tool outputs a JSON object containing metadata such as title, byline, and other information about the article.
