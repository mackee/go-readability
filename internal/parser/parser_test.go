package parser

import (
	"strings"
	"testing"

	"github.com/mackee/go-readability/internal/dom"
)

func TestParseHTML(t *testing.T) {
	// Test case 1: Basic HTML parsing
	html := `<!DOCTYPE html>
<html>
<head>
  <title>Test Page</title>
</head>
<body>
  <h1>Hello World!</h1>
  <p>This is a <strong>test</strong> paragraph.</p>
</body>
</html>`

	doc, err := ParseHTML(html, "https://example.com")
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}

	// Check document structure
	if doc.DocumentElement == nil {
		t.Fatal("DocumentElement is nil")
	}

	if doc.Body == nil {
		t.Fatal("Body is nil")
	}

	if doc.BaseURI != "https://example.com" {
		t.Errorf("Expected BaseURI to be %q, got %q", "https://example.com", doc.BaseURI)
	}

	// Debug: Print the document structure
	t.Logf("Document structure: %s", SerializeDocumentToHTML(doc))
	
	// Find h1 element - search recursively through the body
	var findH1 func(*dom.VElement) *dom.VElement
	findH1 = func(el *dom.VElement) *dom.VElement {
		for _, child := range el.Children {
			if element, ok := dom.AsVElement(child); ok {
				t.Logf("Found element: %s", element.TagName)
				if element.TagName == "h1" {
					return element
				}
				if found := findH1(element); found != nil {
					return found
				}
			}
		}
		return nil
	}
	
	h1 := findH1(doc.Body)

	if h1 == nil {
		t.Fatal("h1 element not found")
	}

	// Check h1 content
	if len(h1.Children) != 1 {
		t.Fatalf("Expected h1 to have 1 child, got %d", len(h1.Children))
	}

	h1Text, ok := dom.AsVText(h1.Children[0])
	if !ok {
		t.Fatal("h1 child is not a text node")
	}

	if h1Text.TextContent != "Hello World!" {
		t.Errorf("Expected h1 text to be %q, got %q", "Hello World!", h1Text.TextContent)
	}

	// Test case 2: HTML with attributes
	html = `<div id="content" class="main">
  <a href="https://example.com" target="_blank">Link</a>
</div>`

	doc, err = ParseHTML(html, "https://example.com")
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}

	// Find div element
	var div *dom.VElement
	for _, child := range doc.Body.Children {
		if element, ok := dom.AsVElement(child); ok && element.TagName == "div" {
			div = element
			break
		}
	}

	if div == nil {
		t.Fatal("div element not found")
	}

	// Check div attributes
	if div.ID() != "content" {
		t.Errorf("Expected div ID to be %q, got %q", "content", div.ID())
	}

	if div.ClassName() != "main" {
		t.Errorf("Expected div class to be %q, got %q", "main", div.ClassName())
	}

	// Find a element
	var a *dom.VElement
	for _, child := range div.Children {
		if element, ok := dom.AsVElement(child); ok && element.TagName == "a" {
			a = element
			break
		}
	}

	if a == nil {
		t.Fatal("a element not found")
	}

	// Check a attributes
	if a.GetAttribute("href") != "https://example.com" {
		t.Errorf("Expected a href to be %q, got %q", "https://example.com", a.GetAttribute("href"))
	}

	if a.GetAttribute("target") != "_blank" {
		t.Errorf("Expected a target to be %q, got %q", "_blank", a.GetAttribute("target"))
	}
}

func TestSerializeToHTML(t *testing.T) {
	// Create a test document
	div := dom.NewVElement("div")
	div.SetAttribute("id", "content")
	div.SetAttribute("class", "main")

	p := dom.NewVElement("p")
	p.AppendChild(dom.NewVText("Hello, "))

	strong := dom.NewVElement("strong")
	strong.AppendChild(dom.NewVText("World"))
	p.AppendChild(strong)
	p.AppendChild(dom.NewVText("!"))

	div.AppendChild(p)

	// Test serialization
	html := SerializeToHTML(div)
	
	// 属性の順序に依存しないテスト
	if !strings.Contains(html, `<div`) ||
	   !strings.Contains(html, ` id="content"`) ||
	   !strings.Contains(html, ` class="main"`) ||
	   !strings.Contains(html, `<p>Hello, <strong>World</strong>!</p>`) {
		t.Errorf("SerializeToHTML produced unexpected output.\nGot: %q", html)
	}

	// Test with self-closing tag
	img := dom.NewVElement("img")
	img.SetAttribute("src", "image.jpg")
	img.SetAttribute("alt", "Test Image")

	html = SerializeToHTML(img)
	
	// 属性の順序に依存しないテスト
	if !strings.Contains(html, `<img`) ||
	   !strings.Contains(html, ` src="image.jpg"`) ||
	   !strings.Contains(html, ` alt="Test Image"`) ||
	   !strings.Contains(html, `/>`) {
		t.Errorf("SerializeToHTML produced unexpected output for self-closing tag.\nGot: %q", html)
	}

	// Test with special characters
	text := dom.NewVText("This & that < > \" '")
	html = SerializeToHTML(text)
	expected := "This &amp; that &lt; &gt; &#34; &#39;"

	if html != expected {
		t.Errorf("SerializeToHTML produced unexpected output for text with special characters.\nExpected: %q\nGot: %q", expected, html)
	}
}

func TestSerializeDocumentToHTML(t *testing.T) {
	// Create a test document
	html := dom.NewVElement("html")
	head := dom.NewVElement("head")
	title := dom.NewVElement("title")
	title.AppendChild(dom.NewVText("Test Page"))
	head.AppendChild(title)
	html.AppendChild(head)

	body := dom.NewVElement("body")
	h1 := dom.NewVElement("h1")
	h1.AppendChild(dom.NewVText("Hello, World!"))
	body.AppendChild(h1)
	html.AppendChild(body)

	doc := dom.NewVDocument(html, body)

	// Test serialization
	output := SerializeDocumentToHTML(doc)
	expected := "<!DOCTYPE html>\n<html><head><title>Test Page</title></head><body><h1>Hello, World!</h1></body></html>"

	if output != expected {
		t.Errorf("SerializeDocumentToHTML produced unexpected output.\nExpected: %q\nGot: %q", expected, output)
	}
}

func TestRoundTrip(t *testing.T) {
	// Test round-trip conversion (HTML -> VDocument -> HTML)
	originalHTML := `<!DOCTYPE html>
<html>
<head>
  <title>Test Page</title>
</head>
<body>
  <h1>Hello, World!</h1>
  <p>This is a <strong>test</strong> paragraph.</p>
  <img src="image.jpg" alt="Test Image"/>
</body>
</html>`

	// Normalize the HTML by removing whitespace between tags
	normalizedOriginal := strings.ReplaceAll(originalHTML, "\n", "")
	normalizedOriginal = strings.ReplaceAll(normalizedOriginal, "  ", "")
	
	// Use normalizedOriginal for comparison later if needed
	_ = normalizedOriginal

	// Parse the HTML
	doc, err := ParseHTML(originalHTML, "https://example.com")
	if err != nil {
		t.Fatalf("ParseHTML failed: %v", err)
	}

	// Serialize back to HTML
	outputHTML := SerializeDocumentToHTML(doc)

	// Normalize the output HTML
	normalizedOutput := strings.ReplaceAll(outputHTML, "\n", "")

	// Compare the normalized HTML strings
	// Note: We're not expecting an exact match due to potential differences in formatting,
	// but the structure and content should be preserved.
	if !strings.Contains(normalizedOutput, "<h1>Hello, World!</h1>") {
		t.Errorf("Round-trip conversion failed to preserve h1 content")
	}

	if !strings.Contains(normalizedOutput, "<p>This is a <strong>test</strong> paragraph.</p>") {
		t.Errorf("Round-trip conversion failed to preserve p content")
	}

	// 属性の順序に依存しないテスト
	if !strings.Contains(normalizedOutput, "<img") &&
	   !strings.Contains(normalizedOutput, " src=\"image.jpg\"") &&
	   !strings.Contains(normalizedOutput, " alt=\"Test Image\"") &&
	   !strings.Contains(normalizedOutput, "/>") {
		t.Errorf("Round-trip conversion failed to preserve img element")
	}
}
