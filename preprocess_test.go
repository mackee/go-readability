package readability

import (
	"testing"

	"github.com/mackee/go-readability/internal/dom"
	"github.com/mackee/go-readability/internal/parser"
)

func TestPreprocessDocument(t *testing.T) {
	t.Run("should remove script tags", func(t *testing.T) {
		html := `
			<html>
				<body>
					<h1>Title</h1>
					<p>Some content.</p>
					<script>alert('Hello');</script>
					<p>More content.</p>
					<script src="script.js"></script>
				</body>
			</html>
		`
		doc, err := parser.ParseHTML(html, "")
		if err != nil {
			t.Fatalf("Failed to parse HTML: %v", err)
		}

		PreprocessDocument(doc)

		scriptElements := dom.GetElementsByTagName(doc.Body, "script")
		if len(scriptElements) != 0 {
			t.Errorf("Expected 0 script elements, got %d", len(scriptElements))
		}

		// Ensure content paragraphs are still present
		pElements := dom.GetElementsByTagName(doc.Body, "p")
		if len(pElements) != 2 {
			t.Errorf("Expected 2 paragraph elements, got %d", len(pElements))
		}
	})

	t.Run("should remove style tags", func(t *testing.T) {
		html := `
			<html>
				<head>
					<style>body { background: red; }</style>
				</head>
				<body>
					<h1>Title</h1>
					<style>.content { color: blue; }</style>
					<p>Some content.</p>
				</body>
			</html>
		`
		doc, err := parser.ParseHTML(html, "")
		if err != nil {
			t.Fatalf("Failed to parse HTML: %v", err)
		}

		PreprocessDocument(doc)

		// Check head and body separately
		styleElementsInDoc := dom.GetElementsByTagName(doc.DocumentElement, "style")
		if len(styleElementsInDoc) != 0 {
			t.Errorf("Expected 0 style elements in document, got %d", len(styleElementsInDoc))
		}

		// Double check body specifically
		styleElementsInBody := dom.GetElementsByTagName(doc.Body, "style")
		if len(styleElementsInBody) != 0 {
			t.Errorf("Expected 0 style elements in body, got %d", len(styleElementsInBody))
		}

		// Ensure content is still present
		pElements := dom.GetElementsByTagName(doc.Body, "p")
		if len(pElements) != 1 {
			t.Errorf("Expected 1 paragraph element, got %d", len(pElements))
		}

		h1Elements := dom.GetElementsByTagName(doc.Body, "h1")
		if len(h1Elements) != 1 {
			t.Errorf("Expected 1 h1 element, got %d", len(h1Elements))
		}
	})

	t.Run("should remove both script and style tags", func(t *testing.T) {
		html := `
			<html>
				<body>
					<style>h1 { font-size: 2em; }</style>
					<h1>Title</h1>
					<script>console.log('Logging');</script>
					<p>Content between tags.</p>
					<script src="another.js"></script>
					<style>.footer { text-align: center; }</style>
				</body>
			</html>
		`
		doc, err := parser.ParseHTML(html, "")
		if err != nil {
			t.Fatalf("Failed to parse HTML: %v", err)
		}

		PreprocessDocument(doc)

		scriptElements := dom.GetElementsByTagName(doc.Body, "script")
		if len(scriptElements) != 0 {
			t.Errorf("Expected 0 script elements, got %d", len(scriptElements))
		}

		styleElements := dom.GetElementsByTagName(doc.Body, "style")
		if len(styleElements) != 0 {
			t.Errorf("Expected 0 style elements, got %d", len(styleElements))
		}

		// Ensure content is preserved
		h1Elements := dom.GetElementsByTagName(doc.Body, "h1")
		if len(h1Elements) != 1 {
			t.Errorf("Expected 1 h1 element, got %d", len(h1Elements))
		}

		pElements := dom.GetElementsByTagName(doc.Body, "p")
		if len(pElements) != 1 {
			t.Errorf("Expected 1 paragraph element, got %d", len(pElements))
		}
	})

	t.Run("should not remove content when no script or style tags are present", func(t *testing.T) {
		html := `
			<html>
				<body>
					<h1>Main Title</h1>
					<p>This is the first paragraph.</p>
					<div><p>Nested paragraph.</p></div>
				</body>
			</html>
		`
		doc, err := parser.ParseHTML(html, "")
		if err != nil {
			t.Fatalf("Failed to parse HTML: %v", err)
		}

		// Get a count of elements before preprocessing
		h1BeforeCount := len(dom.GetElementsByTagName(doc.Body, "h1"))
		pBeforeCount := len(dom.GetElementsByTagName(doc.DocumentElement, "p"))
		divBeforeCount := len(dom.GetElementsByTagName(doc.Body, "div"))

		PreprocessDocument(doc)

		// Check if the body structure remains the same
		h1AfterCount := len(dom.GetElementsByTagName(doc.Body, "h1"))
		if h1AfterCount != h1BeforeCount {
			t.Errorf("Expected %d h1 elements, got %d", h1BeforeCount, h1AfterCount)
		}

		pAfterCount := len(dom.GetElementsByTagName(doc.DocumentElement, "p"))
		if pAfterCount != pBeforeCount {
			t.Errorf("Expected %d paragraph elements, got %d", pBeforeCount, pAfterCount)
		}

		divAfterCount := len(dom.GetElementsByTagName(doc.Body, "div"))
		if divAfterCount != divBeforeCount {
			t.Errorf("Expected %d div elements, got %d", divBeforeCount, divAfterCount)
		}
	})

	t.Run("should remove ad elements", func(t *testing.T) {
		html := `
			<html>
				<body>
					<h1>Main Title</h1>
					<div class="ad-container">This is an ad</div>
					<p>This is content.</p>
					<div id="banner-ad">Another ad</div>
					<div data-ad="true">Yet another ad</div>
				</body>
			</html>
		`
		doc, err := parser.ParseHTML(html, "")
		if err != nil {
			t.Fatalf("Failed to parse HTML: %v", err)
		}

		PreprocessDocument(doc)

		// Check if ad elements were removed
		adElements := dom.GetElementsByTagName(doc.Body, "div")
		if len(adElements) != 0 {
			t.Errorf("Expected 0 div elements (ads), got %d", len(adElements))
		}

		// Ensure content is preserved
		h1Elements := dom.GetElementsByTagName(doc.Body, "h1")
		if len(h1Elements) != 1 {
			t.Errorf("Expected 1 h1 element, got %d", len(h1Elements))
		}

		pElements := dom.GetElementsByTagName(doc.Body, "p")
		if len(pElements) != 1 {
			t.Errorf("Expected 1 paragraph element, got %d", len(pElements))
		}
	})
}
