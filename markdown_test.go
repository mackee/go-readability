package readability

import (
	"strings"
	"testing"

	"github.com/mackee/go-readability/internal/parser"
)

func TestToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "basic HTML to Markdown",
			html: `
				<h1>Title</h1>
				<p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
				<p>Another paragraph with a <a href="http://example.com">link</a>.</p>
			`,
			expected: `# Title

This is a paragraph with **bold** and *italic* text.

Another paragraph with a [link](http://example.com).`,
		},
		{
			name: "headings",
			html: `
				<h1>H1</h1>
				<h2>H2</h2>
				<h3>H3</h3>
				<h4>H4</h4>
				<h5>H5</h5>
				<h6>H6</h6>
			`,
			expected: `# H1

## H2

### H3

#### H4

##### H5

###### H6`,
		},
		{
			name: "unordered lists",
			html: `
				<ul>
					<li>Item 1</li>
					<li>Item 2</li>
					<li>Item 3</li>
				</ul>
			`,
			expected: `- Item 1
- Item 2
- Item 3`,
		},
		{
			name: "ordered lists",
			html: `
				<ol>
					<li>First</li>
					<li>Second</li>
					<li>Third</li>
				</ol>
			`,
			expected: `1. First
1. Second
1. Third`,
		},
		{
			name:     "inline code",
			html:     `<p>Use <code>const</code> for constants.</p>`,
			expected: "Use `const` for constants.",
		},
		{
			name: "code blocks",
			html: `
				<pre><code>function greet() {
  console.log("Hello");
}</code></pre>
			`,
			expected: "```\nfunction greet() {\n  console.log(\"Hello\");\n}\n```",
		},
		{
			name: "code blocks with language class",
			html: `
				<pre><code class="language-javascript">function greet() {
  console.log("Hello");
}</code></pre>
			`,
			expected: "```javascript\nfunction greet() {\n  console.log(\"Hello\");\n}\n```",
		},
		{
			name:     "blockquotes",
			html:     `<blockquote>This is a quote.</blockquote>`,
			expected: `> This is a quote.`,
		},
		{
			name:     "images",
			html:     `<img src="image.png" alt="Alt text">`,
			expected: `![Alt text](image.png)`,
		},
		{
			name:     "horizontal rules",
			html:     `<hr>`,
			expected: `---`,
		},
		{
			name: "ignore script and style tags",
			html: `
				<p>Content</p>
				<script>alert('ignored');</script>
				<style>.ignored { color: red; }</style>
				<p>More Content</p>
			`,
			expected: `Content

More Content`,
		},
		{
			name: "nested lists (ul)",
			html: `
				<ul>
					<li>Item 1</li>
					<li>
						Item 2
						<ul>
							<li>Nested 2.1</li>
							<li>Nested 2.2</li>
						</ul>
					</li>
					<li>Item 3</li>
				</ul>
			`,
			expected: `- Item 1
- Item 2
  - Nested 2.1
  - Nested 2.2
- Item 3`,
		},
		{
			name: "nested lists (ol)",
			html: `
				<ol>
					<li>First</li>
					<li>
						Second
						<ol>
							<li>Nested 2.1</li>
							<li>Nested 2.2</li>
						</ol>
					</li>
					<li>Third</li>
				</ol>
			`,
			expected: `1. First
1. Second
  1. Nested 2.1
  1. Nested 2.2
1. Third`,
		},
		{
			name:     "image links",
			html:     `<a href="http://example.com"><img src="image.png" alt="Alt text"></a>`,
			expected: `[Alt text](http://example.com)`,
		},
		{
			name: "simple table",
			html: `
				<table>
					<thead>
						<tr>
							<th>Header 1</th>
							<th>Header 2</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>Data 1</td>
							<td>Data 2</td>
						</tr>
						<tr>
							<td>Data 3</td>
							<td>Data 4 <strong>bold</strong></td>
						</tr>
					</tbody>
				</table>
			`,
			expected: `| Header 1 | Header 2 |
| --- | --- |
| Data 1 | Data 2 |
| Data 3 | Data 4 **bold** |`,
		},
		{
			name: "table without thead",
			html: `
				<table>
					<tbody>
						<tr>
							<td>Row 1, Cell 1</td>
							<td>Row 1, Cell 2</td>
						</tr>
						<tr>
							<td>Row 2, Cell 1</td>
							<td>Row 2, Cell 2</td>
						</tr>
					</tbody>
				</table>
			`,
			expected: `| --- | --- |
| Row 1, Cell 1 | Row 1, Cell 2 |
| Row 2, Cell 1 | Row 2, Cell 2 |`,
		},
		{
			name: "table with varying columns (padded)",
			html: `
				<table>
					<thead>
						<tr><th>A</th><th>B</th><th>C</th></tr>
					</thead>
					<tbody>
						<tr><td>1</td><td>2</td></tr>
						<tr><td>3</td><td>4</td><td>5</td></tr>
					</tbody>
				</table>
			`,
			expected: `| A | B | C |
| --- | --- | --- |
| 1 | 2 |  |
| 3 | 4 | 5 |`,
		},
		{
			name: "nested blockquotes",
			html: `
				<blockquote>
					<p>Outer quote.</p>
					<blockquote>
						<p>Inner quote.</p>
					</blockquote>
					<p>Outer quote continued.</p>
				</blockquote>
			`,
			expected: `> Outer quote.
>
> > Inner quote.
>
> Outer quote continued.`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.ParseHTML(tt.html, "")
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}

			elementToConvert := doc.Body
			result := ToMarkdown(elementToConvert)

			// Normalize whitespace for comparison
			normalizedResult := normalizeWhitespace(result)
			normalizedExpected := normalizeWhitespace(tt.expected)

			if normalizedResult != normalizedExpected {
				t.Errorf("ToMarkdown() =\n%s\n\nwant:\n%s", normalizedResult, normalizedExpected)
			}
		})
	}
}

// normalizeWhitespace normalizes whitespace for comparison
func normalizeWhitespace(s string) string {
	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)
	// Normalize newlines
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return s
}

func TestEscapeMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escape asterisks",
			input:    "This *is* important",
			expected: `This \*is\* important`,
		},
		{
			name:     "escape underscores",
			input:    "This _is_ important",
			expected: `This \_is\_ important`,
		},
		{
			name:     "escape backticks",
			input:    "Use `code` here",
			expected: "Use \\`code\\` here",
		},
		{
			name:     "escape brackets",
			input:    "This [is] a link",
			expected: `This \[is\] a link`,
		},
		{
			name:     "escape backslashes",
			input:    `This \ is a backslash`,
			expected: `This \\ is a backslash`,
		},
		{
			name:     "decode HTML entities",
			input:    "This &amp; that &lt; this &gt; that",
			expected: `This & that < this > that`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("escapeMarkdown() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestJoinMarkdownParts(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "empty parts",
			parts:    []string{},
			expected: "",
		},
		{
			name:     "single part",
			parts:    []string{"Hello"},
			expected: "Hello",
		},
		{
			name:     "multiple parts",
			parts:    []string{"Hello", "world"},
			expected: "Hello world",
		},
		{
			name:     "parts with whitespace",
			parts:    []string{"Hello ", " world"},
			expected: "Hello  world",
		},
		{
			name:     "parts with punctuation",
			parts:    []string{"Hello", ".", "How", "are", "you", "?"},
			expected: "Hello. How are you?",
		},
		{
			name:     "skip empty parts",
			parts:    []string{"Hello", "", "world", "   "},
			expected: "Hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinMarkdownParts(tt.parts)
			if result != tt.expected {
				t.Errorf("joinMarkdownParts() = %v, want %v", result, tt.expected)
			}
		})
	}
}
