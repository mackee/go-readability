package readability

import (
	"testing"

	"github.com/mackee/go-readability/internal/dom"
)

func TestGetArticleTitle(t *testing.T) {
	testCases := []struct {
		name     string
		setupDoc func() *dom.VDocument
		expected string
	}{
		{
			name: "simple title",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Simple Title"))
				head.AppendChild(title)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Simple Title",
		},
		{
			name: "title with separator",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Main Title | Site Name"))
				head.AppendChild(title)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Site Name",
		},
		{
			name: "title with colon",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Site Name: Article Title"))
				head.AppendChild(title)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Article Title",
		},
		{
			name: "long title with h1",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("This is a very long title that exceeds the 150 character limit and should be replaced with the h1 content. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."))
				head.AppendChild(title)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				h1 := dom.NewVElement("h1")
				h1.AppendChild(dom.NewVText("H1 Title"))
				body.AppendChild(h1)
				
				return dom.NewVDocument(html, body)
			},
			expected: "H1 Title",
		},
		{
			name: "title with matching h1",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Exact Match: This is the title"))
				head.AppendChild(title)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				h1 := dom.NewVElement("h1")
				h1.AppendChild(dom.NewVText("Exact Match: This is the title"))
				body.AppendChild(h1)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Exact Match: This is the title",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := tc.setupDoc()
			title := GetArticleTitle(doc)
			
			if title != tc.expected {
				t.Errorf("Expected title '%s', got '%s'", tc.expected, title)
			}
		})
	}
}

func TestGetArticleByline(t *testing.T) {
	testCases := []struct {
		name     string
		setupDoc func() *dom.VDocument
		expected string
	}{
		{
			name: "meta author tag",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				meta := dom.NewVElement("meta")
				meta.SetAttribute("name", "author")
				meta.SetAttribute("content", "John Doe")
				head.AppendChild(meta)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "John Doe",
		},
		{
			name: "dc:creator meta tag",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				meta := dom.NewVElement("meta")
				meta.SetAttribute("property", "dc:creator")
				meta.SetAttribute("content", "Jane Smith")
				head.AppendChild(meta)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Jane Smith",
		},
		{
			name: "article:author meta tag (not URL)",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				meta := dom.NewVElement("meta")
				meta.SetAttribute("property", "article:author")
				meta.SetAttribute("content", "Alice Johnson")
				head.AppendChild(meta)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Alice Johnson",
		},
		{
			name: "article:author meta tag (URL)",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				meta := dom.NewVElement("meta")
				meta.SetAttribute("property", "article:author")
				meta.SetAttribute("content", "https://example.com/author/bob")
				head.AppendChild(meta)
				
				// Add another meta tag that should be used instead
				meta2 := dom.NewVElement("meta")
				meta2.SetAttribute("name", "author")
				meta2.SetAttribute("content", "Bob Williams")
				head.AppendChild(meta2)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Bob Williams",
		},
		{
			name: "HTML entity in author name",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				head := dom.NewVElement("head")
				html.AppendChild(head)
				
				meta := dom.NewVElement("meta")
				meta.SetAttribute("name", "author")
				meta.SetAttribute("content", "Charlie &amp; Dave")
				head.AppendChild(meta)
				
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				return dom.NewVDocument(html, body)
			},
			expected: "Charlie & Dave",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := tc.setupDoc()
			byline := GetArticleByline(doc)
			
			if byline != tc.expected {
				t.Errorf("Expected byline '%s', got '%s'", tc.expected, byline)
			}
		})
	}
}

func TestUnescapeHTMLEntities(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no entities",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "named entities",
			input:    "This &amp; that &lt; other &gt; things &quot;quoted&quot; and &apos;single&apos;",
			expected: "This & that < other > things \"quoted\" and 'single'",
		},
		{
			name:     "decimal entities",
			input:    "&#65;&#66;&#67;",
			expected: "ABC",
		},
		{
			name:     "hex entities",
			input:    "&#x41;&#x42;&#x43;",
			expected: "ABC",
		},
		{
			name:     "mixed entities",
			input:    "&#65;&amp;&#x42;",
			expected: "A&B",
		},
		{
			name:     "invalid entities",
			input:    "&#xFFFFF;&#x110000;&#xD800;",
			expected: "\uFFFD\uFFFD\uFFFD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := UnescapeHTMLEntities(tc.input)
			
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestTextSimilarity(t *testing.T) {
	testCases := []struct {
		name     string
		textA    string
		textB    string
		expected float64
		delta    float64 // Allowed difference for floating point comparison
	}{
		{
			name:     "identical texts",
			textA:    "This is a test",
			textB:    "This is a test",
			expected: 1.0,
			delta:    0.001,
		},
		{
			name:     "completely different texts",
			textA:    "This is a test",
			textB:    "Something else entirely",
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "partially similar texts",
			textA:    "This is a test sentence",
			textB:    "This is a different sentence",
			expected: 0.6, // Approximate value
			delta:    0.1,
		},
		{
			name:     "case insensitive",
			textA:    "This Is A Test",
			textB:    "this is a test",
			expected: 1.0,
			delta:    0.001,
		},
		{
			name:     "empty texts",
			textA:    "",
			textB:    "",
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "one empty text",
			textA:    "This is a test",
			textB:    "",
			expected: 0.0,
			delta:    0.001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			similarity := TextSimilarity(tc.textA, tc.textB)
			
			if similarity < tc.expected-tc.delta || similarity > tc.expected+tc.delta {
				t.Errorf("Expected similarity around %.2f, got %.2f", tc.expected, similarity)
			}
		})
	}
}
