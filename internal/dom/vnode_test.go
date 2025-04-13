package dom

import (
	"testing"
)

func TestVText(t *testing.T) {
	text := NewVText("Hello, world!")
	
	if text.Type() != TextNode {
		t.Errorf("Expected node type to be %v, got %v", TextNode, text.Type())
	}
	
	if text.TextContent != "Hello, world!" {
		t.Errorf("Expected text content to be %q, got %q", "Hello, world!", text.TextContent)
	}
	
	if text.Parent() != nil {
		t.Errorf("Expected parent to be nil, got %v", text.Parent())
	}
	
	// Test ReadabilityData
	if text.GetReadabilityData() != nil {
		t.Errorf("Expected ReadabilityData to be nil initially")
	}
	
	data := &ReadabilityData{ContentScore: 42}
	text.SetReadabilityData(data)
	
	if text.GetReadabilityData() != data {
		t.Errorf("Expected ReadabilityData to be %v, got %v", data, text.GetReadabilityData())
	}
	
	if text.GetReadabilityData().ContentScore != 42 {
		t.Errorf("Expected ContentScore to be %v, got %v", 42, text.GetReadabilityData().ContentScore)
	}
}

func TestVElement(t *testing.T) {
	element := NewVElement("div")
	
	if element.Type() != ElementNode {
		t.Errorf("Expected node type to be %v, got %v", ElementNode, element.Type())
	}
	
	if element.TagName != "div" {
		t.Errorf("Expected tag name to be %q, got %q", "div", element.TagName)
	}
	
	if len(element.Children) != 0 {
		t.Errorf("Expected children count to be 0, got %d", len(element.Children))
	}
	
	if len(element.Attributes) != 0 {
		t.Errorf("Expected attributes count to be 0, got %d", len(element.Attributes))
	}
	
	// Test attributes
	element.SetAttribute("id", "test-id")
	element.SetAttribute("class", "test-class")
	
	if element.ID() != "test-id" {
		t.Errorf("Expected ID to be %q, got %q", "test-id", element.ID())
	}
	
	if element.ClassName() != "test-class" {
		t.Errorf("Expected ClassName to be %q, got %q", "test-class", element.ClassName())
	}
	
	if !element.HasAttribute("id") {
		t.Errorf("Expected HasAttribute('id') to be true")
	}
	
	if element.GetAttribute("id") != "test-id" {
		t.Errorf("Expected GetAttribute('id') to be %q, got %q", "test-id", element.GetAttribute("id"))
	}
	
	// Test children
	text := NewVText("Hello")
	element.AppendChild(text)
	
	if len(element.Children) != 1 {
		t.Errorf("Expected children count to be 1, got %d", len(element.Children))
	}
	
	if text.Parent() != element {
		t.Errorf("Expected text parent to be element")
	}
	
	// Test type assertion methods
	if !IsVElement(element) {
		t.Errorf("Expected IsVElement to return true for element")
	}
	
	if IsVElement(text) {
		t.Errorf("Expected IsVElement to return false for text")
	}
	
	if !IsVText(text) {
		t.Errorf("Expected IsVText to return true for text")
	}
	
	if IsVText(element) {
		t.Errorf("Expected IsVText to return false for element")
	}
	
	if el, ok := AsVElement(element); !ok || el != element {
		t.Errorf("Expected AsVElement to return element and true")
	}
	
	if _, ok := AsVElement(text); ok {
		t.Errorf("Expected AsVElement to return nil and false for text")
	}
	
	if tx, ok := AsVText(text); !ok || tx != text {
		t.Errorf("Expected AsVText to return text and true")
	}
	
	if _, ok := AsVText(element); ok {
		t.Errorf("Expected AsVText to return nil and false for element")
	}
}

func TestVDocument(t *testing.T) {
	html := NewVElement("html")
	body := NewVElement("body")
	html.AppendChild(body)
	
	doc := NewVDocument(html, body)
	
	if doc.DocumentElement != html {
		t.Errorf("Expected DocumentElement to be html")
	}
	
	if doc.Body != body {
		t.Errorf("Expected Body to be body")
	}
	
	// Test with URI
	doc.BaseURI = "https://example.com/"
	doc.DocumentURI = "https://example.com/page.html"
	
	if doc.BaseURI != "https://example.com/" {
		t.Errorf("Expected BaseURI to be %q, got %q", "https://example.com/", doc.BaseURI)
	}
	
	if doc.DocumentURI != "https://example.com/page.html" {
		t.Errorf("Expected DocumentURI to be %q, got %q", "https://example.com/page.html", doc.DocumentURI)
	}
}