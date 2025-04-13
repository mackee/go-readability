package readability

import (
	"strings"
	"testing"

	"github.com/mackee/go-readability/internal/dom"
)

func TestToHTML(t *testing.T) {
	t.Run("should remove span tags but keep their content", func(t *testing.T) {
		element := dom.NewVElement("div")
		element.AppendChild(dom.NewVText("Hello "))
		
		span1 := dom.NewVElement("span")
		span1.AppendChild(dom.NewVText("world"))
		element.AppendChild(span1)
		
		element.AppendChild(dom.NewVText("!"))

		expectedHTML := "<div>Hello world!</div>"
		if html := ToHTML(element); html != expectedHTML {
			t.Errorf("Expected HTML: %s, got: %s", expectedHTML, html)
		}
	})

	t.Run("should remove class attributes from all elements", func(t *testing.T) {
		element := dom.NewVElement("div")
		element.SetAttribute("class", "container") // Add class to div
		element.SetAttribute("id", "main")         // Keep other attributes like id

		p1 := dom.NewVElement("p")
		p1.SetAttribute("class", "intro") // Add class to p
		p1.AppendChild(dom.NewVText("This is a paragraph."))
		element.AppendChild(p1)

		spanInP := dom.NewVElement("span")
		spanInP.SetAttribute("class", "highlight") // Add class to span
		spanInP.AppendChild(dom.NewVText(" Important text."))
		p1.AppendChild(spanInP) // Add span inside p

		expectedHTML := "<div id=\"main\"><p>This is a paragraph. Important text.</p></div>"
		if html := ToHTML(element); html != expectedHTML {
			t.Errorf("Expected HTML: %s, got: %s", expectedHTML, html)
		}
	})

	t.Run("should handle nested spans and other elements correctly", func(t *testing.T) {
		element := dom.NewVElement("article")
		element.SetAttribute("class", "post")

		h1 := dom.NewVElement("h1")
		h1.SetAttribute("class", "title")
		h1.AppendChild(dom.NewVText("Test Title"))
		element.AppendChild(h1)

		p1 := dom.NewVElement("p")
		p1.SetAttribute("class", "content")
		p1.AppendChild(dom.NewVText("Some text "))

		outerSpan := dom.NewVElement("span")
		outerSpan.SetAttribute("class", "outer")
		outerSpan.AppendChild(dom.NewVText("with an "))

		innerSpan := dom.NewVElement("span")
		innerSpan.SetAttribute("class", "inner important") // Multiple classes
		innerSpan.AppendChild(dom.NewVText("inner span"))
		outerSpan.AppendChild(innerSpan)

		outerSpan.AppendChild(dom.NewVText(" inside."))
		p1.AppendChild(outerSpan)
		element.AppendChild(p1)

		img := dom.NewVElement("img")
		img.SetAttribute("src", "image.jpg")
		img.SetAttribute("class", "featured") // Class on self-closing tag
		element.AppendChild(img)

		expectedHTML := "<article><h1>Test Title</h1><p>Some text with an inner span inside.</p><img src=\"image.jpg\"/></article>"
		if html := ToHTML(element); html != expectedHTML {
			t.Errorf("Expected HTML: %s, got: %s", expectedHTML, html)
		}
	})

	t.Run("should handle self-closing tags correctly, removing class", func(t *testing.T) {
		element := dom.NewVElement("div")
		
		br := dom.NewVElement("br")
		br.SetAttribute("class", "break") // Class on br
		element.AppendChild(br)
		
		hr := dom.NewVElement("hr")
		hr.SetAttribute("class", "divider") // Class on hr
		element.AppendChild(hr)
		
		img := dom.NewVElement("img")
		img.SetAttribute("src", "test.png")
		img.SetAttribute("class", "icon") // Class on img
		img.SetAttribute("alt", "test")   // Keep alt attribute
		element.AppendChild(img)

		// 属性の順序は保証されないため、両方の順序を許容する
		html := ToHTML(element)
		validHTML1 := "<div><br/><hr/><img src=\"test.png\" alt=\"test\"/></div>"
		validHTML2 := "<div><br/><hr/><img alt=\"test\" src=\"test.png\"/></div>"
		
		if html != validHTML1 && html != validHTML2 {
			t.Errorf("Expected HTML to be either %s or %s, got: %s", validHTML1, validHTML2, html)
		}
	})

	t.Run("should return empty string for nil input", func(t *testing.T) {
		if html := ToHTML(nil); html != "" {
			t.Errorf("Expected empty string for nil input, got: %s", html)
		}
	})

	t.Run("should preserve whitespace including nbsp when removing span tags", func(t *testing.T) {
		// Create a structure similar to the one in the TypeScript test
		p1 := dom.NewVElement("p")
		p1.AppendChild(dom.NewVText("Some text "))
		span1 := dom.NewVElement("span")
		span1.AppendChild(dom.NewVText("with a span"))
		p1.AppendChild(span1)
		p1.AppendChild(dom.NewVText(" inside."))

		p2 := dom.NewVElement("p")
		p2.AppendChild(dom.NewVText("Another text\u00a0")) // Use \u00a0 for nbsp
		span2 := dom.NewVElement("span")
		span2.AppendChild(dom.NewVText("with nbsp"))
		p2.AppendChild(span2)
		p2.AppendChild(dom.NewVText("\u00a0around.")) // Use \u00a0 for nbsp

		p3 := dom.NewVElement("p")
		p3.AppendChild(dom.NewVText("Text"))
		span3 := dom.NewVElement("span")
		span3.AppendChild(dom.NewVText("without space"))
		p3.AppendChild(span3)
		p3.AppendChild(dom.NewVText("around."))

		p4 := dom.NewVElement("p")
		p4.AppendChild(dom.NewVText(" Text with leading space"))
		span4 := dom.NewVElement("span")
		span4.AppendChild(dom.NewVText(" and span"))
		p4.AppendChild(span4)
		p4.AppendChild(dom.NewVText("."))

		p5 := dom.NewVElement("p")
		span5 := dom.NewVElement("span")
		span5.AppendChild(dom.NewVText("Span at start"))
		p5.AppendChild(span5)
		p5.AppendChild(dom.NewVText(" and text."))

		container := dom.NewVElement("div")
		container.AppendChild(p1)
		container.AppendChild(p2)
		container.AppendChild(p3)
		container.AppendChild(p4)
		container.AppendChild(p5)

		// Expected HTML, note &nbsp; for non-breaking spaces
		expectedHTML := "<div><p>Some text with a span inside.</p><p>Another text&nbsp;with nbsp&nbsp;around.</p><p>Textwithout spacearound.</p><p> Text with leading space and span.</p><p>Span at start and text.</p></div>"
		
		html := ToHTML(container)
		if html != expectedHTML {
			t.Errorf("Expected HTML: %s, got: %s", expectedHTML, html)
		}
	})
}

func TestStringify(t *testing.T) {
	t.Run("should convert element to readable string format", func(t *testing.T) {
		article := dom.NewVElement("article")
		
		h1 := dom.NewVElement("h1")
		h1.AppendChild(dom.NewVText("Article Title"))
		article.AppendChild(h1)
		
		p1 := dom.NewVElement("p")
		p1.AppendChild(dom.NewVText("This is the first paragraph."))
		article.AppendChild(p1)
		
		p2 := dom.NewVElement("p")
		p2.AppendChild(dom.NewVText("This is the second paragraph with "))
		em := dom.NewVElement("em")
		em.AppendChild(dom.NewVText("emphasized"))
		p2.AppendChild(em)
		p2.AppendChild(dom.NewVText(" text."))
		article.AppendChild(p2)

		result := Stringify(article)
		
		// The result should contain the text content with appropriate line breaks
		// but without HTML tags
		if result == "" {
			t.Error("Stringify returned empty string")
		}
		
		// Check that the result contains the text content
		if !formatContains(result, "Article Title") || 
		   !formatContains(result, "This is the first paragraph") ||
		   !formatContains(result, "This is the second paragraph with emphasized text") {
			t.Errorf("Stringify did not preserve text content: %s", result)
		}
	})

	t.Run("should handle special tags like br and hr", func(t *testing.T) {
		div := dom.NewVElement("div")
		
		p1 := dom.NewVElement("p")
		p1.AppendChild(dom.NewVText("Line 1"))
		div.AppendChild(p1)
		
		br := dom.NewVElement("br")
		div.AppendChild(br)
		
		p2 := dom.NewVElement("p")
		p2.AppendChild(dom.NewVText("Line 2"))
		div.AppendChild(p2)
		
		hr := dom.NewVElement("hr")
		div.AppendChild(hr)
		
		p3 := dom.NewVElement("p")
		p3.AppendChild(dom.NewVText("Line 3"))
		div.AppendChild(p3)

		result := Stringify(div)
		
		// Check that br is converted to a line break
		if !formatContains(result, "Line 1") || !formatContains(result, "Line 2") || !formatContains(result, "Line 3") {
			t.Errorf("Stringify did not preserve text content: %s", result)
		}
		
		// Check that hr is converted to a horizontal rule
		if !formatContains(result, "----------") {
			t.Errorf("Stringify did not convert hr to horizontal rule: %s", result)
		}
	})

	t.Run("should return empty string for nil input", func(t *testing.T) {
		if result := Stringify(nil); result != "" {
			t.Errorf("Expected empty string for nil input, got: %s", result)
		}
	})
}

func TestFormatDocument(t *testing.T) {
	t.Run("should merge consecutive line breaks", func(t *testing.T) {
		input := "Line 1\n\n\nLine 2\n\nLine 3"
		expected := "Line 1\nLine 2\nLine 3"
		
		if result := FormatDocument(input); result != expected {
			t.Errorf("Expected: %s, got: %s", expected, result)
		}
	})

	t.Run("should remove leading and trailing line breaks", func(t *testing.T) {
		input := "\n\nContent\n\n"
		expected := "Content"
		
		if result := FormatDocument(input); result != expected {
			t.Errorf("Expected: %s, got: %s", expected, result)
		}
	})

	t.Run("should trim whitespace", func(t *testing.T) {
		input := "  \t  Content  \t  "
		expected := "Content"
		
		if result := FormatDocument(input); result != expected {
			t.Errorf("Expected: %s, got: %s", expected, result)
		}
	})
}

func TestExtractTextContent(t *testing.T) {
	t.Run("should extract text content from element", func(t *testing.T) {
		div := dom.NewVElement("div")
		
		p1 := dom.NewVElement("p")
		p1.AppendChild(dom.NewVText("Paragraph 1"))
		div.AppendChild(p1)
		
		p2 := dom.NewVElement("p")
		p2.AppendChild(dom.NewVText("Paragraph "))
		strong := dom.NewVElement("strong")
		strong.AppendChild(dom.NewVText("2"))
		p2.AppendChild(strong)
		div.AppendChild(p2)

		expected := "Paragraph 1Paragraph 2"
		if result := ExtractTextContent(div); result != expected {
			t.Errorf("Expected: %s, got: %s", expected, result)
		}
	})

	t.Run("should return empty string for nil input", func(t *testing.T) {
		if result := ExtractTextContent(nil); result != "" {
			t.Errorf("Expected empty string for nil input, got: %s", result)
		}
	})
}

func TestCountNodes(t *testing.T) {
	t.Run("should count nodes correctly", func(t *testing.T) {
		div := dom.NewVElement("div")
		
		p1 := dom.NewVElement("p")
		p1.AppendChild(dom.NewVText("Text 1"))
		div.AppendChild(p1)
		
		p2 := dom.NewVElement("p")
		p2.AppendChild(dom.NewVText("Text 2"))
		strong := dom.NewVElement("strong")
		strong.AppendChild(dom.NewVText("Bold"))
		p2.AppendChild(strong)
		div.AppendChild(p2)

		// Count: div(1) + p1(1) + "Text 1"(1) + p2(1) + "Text 2"(1) + strong(1) + "Bold"(1) = 7
		expected := 7
		if result := CountNodes(div); result != expected {
			t.Errorf("Expected: %d, got: %d", expected, result)
		}
	})

	t.Run("should return 0 for nil input", func(t *testing.T) {
		if result := CountNodes(nil); result != 0 {
			t.Errorf("Expected 0 for nil input, got: %d", result)
		}
	})
}

// Helper function to check if a string contains a substring
func formatContains(s, substr string) bool {
	return strings.Contains(s, substr)
}
