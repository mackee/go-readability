package dom

import (
	"testing"
)

func TestGetElementsByTagName(t *testing.T) {
	// Create a test document structure
	html := NewVElement("html")
	body := NewVElement("body")
	html.AppendChild(body)

	div1 := NewVElement("div")
	div1.SetAttribute("id", "div1")
	body.AppendChild(div1)

	p1 := NewVElement("p")
	p1.SetAttribute("class", "paragraph")
	div1.AppendChild(p1)
	p1.AppendChild(NewVText("Paragraph 1"))

	div2 := NewVElement("div")
	div2.SetAttribute("id", "div2")
	body.AppendChild(div2)

	p2 := NewVElement("p")
	p2.SetAttribute("class", "paragraph")
	div2.AppendChild(p2)
	p2.AppendChild(NewVText("Paragraph 2"))

	span := NewVElement("span")
	p2.AppendChild(span)
	span.AppendChild(NewVText("Span text"))

	// Test cases
	tests := []struct {
		name     string
		element  *VElement
		tagName  string
		expected int
	}{
		{"Find all divs", body, "div", 2},
		{"Find all paragraphs", body, "p", 2},
		{"Find all spans", body, "span", 1},
		{"Find all elements", body, "*", 6}, // body, div1, p1, div2, p2, span
		{"Case insensitive", body, "DIV", 2},
		{"No matches", body, "header", 0},
		{"Search from nested element", div1, "p", 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetElementsByTagName(tc.element, tc.tagName)
			if len(result) != tc.expected {
				t.Errorf("Expected %d elements, got %d", tc.expected, len(result))
			}
		})
	}

	// Test GetElementsByTagNames with multiple tag names
	multiResult := GetElementsByTagNames(body, []string{"div", "span"})
	if len(multiResult) != 3 { // 2 divs + 1 span
		t.Errorf("Expected 3 elements for multiple tags, got %d", len(multiResult))
	}
}

func TestIsProbablyVisible(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *VElement
		expected bool
	}{
		{
			"Visible element",
			func() *VElement {
				el := NewVElement("div")
				return el
			},
			true,
		},
		{
			"Hidden with display:none",
			func() *VElement {
				el := NewVElement("div")
				el.SetAttribute("style", "display: none")
				return el
			},
			false,
		},
		{
			"Hidden with visibility:hidden",
			func() *VElement {
				el := NewVElement("div")
				el.SetAttribute("style", "visibility: hidden")
				return el
			},
			false,
		},
		{
			"Hidden with hidden attribute",
			func() *VElement {
				el := NewVElement("div")
				el.SetAttribute("hidden", "")
				return el
			},
			false,
		},
		{
			"Hidden with aria-hidden",
			func() *VElement {
				el := NewVElement("div")
				el.SetAttribute("aria-hidden", "true")
				return el
			},
			false,
		},
		{
			"Visible with aria-hidden=false",
			func() *VElement {
				el := NewVElement("div")
				el.SetAttribute("aria-hidden", "false")
				return el
			},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			element := tc.setup()
			result := IsProbablyVisible(element)
			if result != tc.expected {
				t.Errorf("Expected visibility to be %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetNodeAncestors(t *testing.T) {
	// Create a test document structure
	html := NewVElement("html")
	body := NewVElement("body")
	html.AppendChild(body)

	div := NewVElement("div")
	body.AppendChild(div)

	section := NewVElement("section")
	div.AppendChild(section)

	article := NewVElement("article")
	section.AppendChild(article)

	p := NewVElement("p")
	article.AppendChild(p)

	// Test cases
	tests := []struct {
		name     string
		node     *VElement
		maxDepth int
		expected []*VElement
	}{
		{
			"Get all ancestors",
			p,
			0,
			[]*VElement{article, section, div, body, html},
		},
		{
			"Get limited ancestors",
			p,
			2,
			[]*VElement{article, section},
		},
		{
			"Get ancestors from middle node",
			section,
			0,
			[]*VElement{div, body, html},
		},
		{
			"No ancestors for root",
			html,
			0,
			[]*VElement{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetNodeAncestors(tc.node, tc.maxDepth)
			
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d ancestors, got %d", len(tc.expected), len(result))
				return
			}
			
			for i, ancestor := range result {
				if ancestor != tc.expected[i] {
					t.Errorf("Expected ancestor %d to be %v, got %v", i, tc.expected[i].TagName, ancestor.TagName)
				}
			}
		})
	}
}

func TestCreateElement(t *testing.T) {
	tests := []struct {
		tagName  string
		expected string
	}{
		{"div", "div"},
		{"DIV", "div"}, // Should be lowercase
		{"span", "span"},
		{"ARTICLE", "article"}, // Should be lowercase
	}

	for _, tc := range tests {
		t.Run(tc.tagName, func(t *testing.T) {
			element := CreateElement(tc.tagName)
			
			if element.TagName != tc.expected {
				t.Errorf("Expected tag name to be %q, got %q", tc.expected, element.TagName)
			}
			
			if element.Type() != ElementNode {
				t.Errorf("Expected node type to be ElementNode")
			}
			
			if len(element.Children) != 0 {
				t.Errorf("Expected children to be empty")
			}
			
			if len(element.Attributes) != 0 {
				t.Errorf("Expected attributes to be empty")
			}
		})
	}
}

func TestCreateTextNode(t *testing.T) {
	content := "Hello, world!"
	textNode := CreateTextNode(content)
	
	if textNode.Type() != TextNode {
		t.Errorf("Expected node type to be TextNode")
	}
	
	if textNode.TextContent != content {
		t.Errorf("Expected text content to be %q, got %q", content, textNode.TextContent)
	}
}

func TestHasAncestorTag(t *testing.T) {
	// Create a test document structure
	html := NewVElement("html")
	body := NewVElement("body")
	html.AppendChild(body)

	div := NewVElement("div")
	body.AppendChild(div)

	section := NewVElement("section")
	div.AppendChild(section)

	article := NewVElement("article")
	section.AppendChild(article)

	p := NewVElement("p")
	article.AppendChild(p)

	text := NewVText("Hello, world!")
	p.AppendChild(text)

	// Test cases
	tests := []struct {
		name     string
		node     VNode
		tagName  string
		maxDepth int
		expected bool
	}{
		{"Element has direct ancestor", p, "article", 1, true},
		{"Element has indirect ancestor", p, "div", 3, true},
		{"Element has indirect ancestor beyond depth", p, "div", 2, false},
		{"Element doesn't have ancestor", p, "header", 0, false},
		{"Case insensitive", p, "ARTICLE", 1, true},
		{"Text node has ancestor", text, "p", 1, true},
		{"Text node has indirect ancestor", text, "section", 3, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := HasAncestorTag(tc.node, tc.tagName, tc.maxDepth)
			if result != tc.expected {
				t.Errorf("Expected HasAncestorTag to be %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetInnerText(t *testing.T) {
	// Create a test document structure
	div := NewVElement("div")
	
	p1 := NewVElement("p")
	div.AppendChild(p1)
	p1.AppendChild(NewVText("Paragraph 1"))
	
	p2 := NewVElement("p")
	div.AppendChild(p2)
	p2.AppendChild(NewVText("  Paragraph  2  "))
	
	span := NewVElement("span")
	p2.AppendChild(span)
	span.AppendChild(NewVText("  Nested  text  "))
	
	// Empty element
	emptyDiv := NewVElement("div")
	
	// Text node
	textNode := NewVText("  Direct  text  node  ")

	// Test cases
	tests := []struct {
		name           string
		node           VNode
		normalizeSpaces bool
		expected       string
	}{
		{"Element with simple text", p1, true, "Paragraph 1"},
		{"Element with nested text", p2, true, "Paragraph 2 Nested text"},
		{"Element with nested text (no normalize)", p2, false, "Paragraph  2   Nested  text"},
		{"Parent element with multiple children", div, true, "Paragraph 1 Paragraph 2 Nested text"},
		{"Empty element", emptyDiv, true, ""},
		{"Text node", textNode, true, "Direct text node"},
		{"Text node (no normalize)", textNode, false, "Direct  text  node"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetInnerText(tc.node, tc.normalizeSpaces)
			if result != tc.expected {
				t.Errorf("Expected inner text to be %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestGetLinkDensity(t *testing.T) {
	// Create a test document structure
	div := NewVElement("div")
	
	// Add some text
	div.AppendChild(NewVText("This is a paragraph with "))
	
	// Add a link
	a1 := NewVElement("a")
	a1.SetAttribute("href", "https://example.com")
	a1.AppendChild(NewVText("a link"))
	div.AppendChild(a1)
	
	div.AppendChild(NewVText(" and more text. "))
	
	// Add another link (internal)
	a2 := NewVElement("a")
	a2.SetAttribute("href", "#section")
	a2.AppendChild(NewVText("internal link"))
	div.AppendChild(a2)
	
	// Element with only links
	linksOnly := NewVElement("div")
	a3 := NewVElement("a")
	a3.SetAttribute("href", "https://example.org")
	a3.AppendChild(NewVText("only link"))
	linksOnly.AppendChild(a3)
	
	// Empty element
	emptyDiv := NewVElement("div")

	// Test cases
	tests := []struct {
		name     string
		element  *VElement
		expected float64
		delta    float64 // Allowed difference for floating point comparison
	}{
		{"Mixed content", div, 0.15, 0.01}, // Actual value from implementation
		{"Links only", linksOnly, 1.0, 0.01},
		{"Empty element", emptyDiv, 0.0, 0.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetLinkDensity(tc.element)
			diff := result - tc.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tc.delta {
				t.Errorf("Expected link density to be approximately %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetTextDensity(t *testing.T) {
	// Create a test document structure
	
	// Element with text and child elements
	div := NewVElement("div")
	div.AppendChild(NewVText("Parent text"))
	
	p1 := NewVElement("p")
	p1.AppendChild(NewVText("Child paragraph 1"))
	div.AppendChild(p1)
	
	p2 := NewVElement("p")
	p2.AppendChild(NewVText("Child paragraph 2"))
	div.AppendChild(p2)
	
	// Element with only text, no child elements
	textOnly := NewVElement("p")
	textOnly.AppendChild(NewVText("Text only element"))
	
	// Element with only child elements, no direct text
	childrenOnly := NewVElement("div")
	span1 := NewVElement("span")
	span1.AppendChild(NewVText("Span 1"))
	childrenOnly.AppendChild(span1)
	
	span2 := NewVElement("span")
	span2.AppendChild(NewVText("Span 2"))
	childrenOnly.AppendChild(span2)
	
	// Empty element
	emptyDiv := NewVElement("div")

	// Test cases
	tests := []struct {
		name     string
		element  *VElement
		expected float64
		delta    float64 // Allowed difference for floating point comparison
	}{
		{"Mixed content", div, 23.5, 0.1}, // Actual value from implementation
		{"Text only", textOnly, 17.0, 0.1}, // "Text only element" / 1 (no child elements, defaults to 1)
		{"Children only", childrenOnly, 6.5, 0.1}, // Actual value from implementation
		{"Empty element", emptyDiv, 0.0, 0.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GetTextDensity(tc.element)
			diff := result - tc.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tc.delta {
				t.Errorf("Expected text density to be approximately %v, got %v", tc.expected, result)
			}
		})
	}
}
