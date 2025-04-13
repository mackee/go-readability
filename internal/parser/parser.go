// Package parser provides HTML parsing functionality for the readability library.
package parser

import (
	"bytes"
	"io"
	"strings"

	"github.com/mackee/go-readability/internal/dom"
	"golang.org/x/net/html"
)

// ParseHTML parses an HTML string and returns a virtual DOM document.
// It uses golang.org/x/net/html for parsing and converts the result to our internal DOM structure.
func ParseHTML(htmlContent string, baseURI string) (*dom.VDocument, error) {
	// Parse HTML using golang.org/x/net/html
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// Find the html and body elements in the parsed document
	var htmlNode, bodyNode *html.Node
	
	// Helper function to find html and body nodes
	var findNodes func(*html.Node)
	findNodes = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if strings.ToLower(n.Data) == "html" {
				htmlNode = n
			} else if strings.ToLower(n.Data) == "body" {
				bodyNode = n
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNodes(c)
		}
	}
	
	findNodes(doc)
	
	// Create virtual DOM elements
	htmlElement := dom.NewVElement("html")
	var bodyElement *dom.VElement
	
	// Process the document structure
	if htmlNode != nil {
		// Process only the children of the html node to avoid duplication
		for child := htmlNode.FirstChild; child != nil; child = child.NextSibling {
			processNode(child, htmlElement)
		}
		
		// Find the body element in our processed structure
		for _, child := range htmlElement.Children {
			if element, ok := dom.AsVElement(child); ok && element.TagName == "body" {
				bodyElement = element
				break
			}
		}
	} else {
		// If no html element is found, process all children of the document
		for c := doc.FirstChild; c != nil; c = c.NextSibling {
			processNode(c, htmlElement)
		}
	}
	
	// If no body element is found, create one
	if bodyElement == nil {
		bodyElement = dom.NewVElement("body")
		
		// If bodyNode was found, process its children
		if bodyNode != nil {
			for child := bodyNode.FirstChild; child != nil; child = child.NextSibling {
				processNode(child, bodyElement)
			}
		}
		
		// Add body to html
		htmlElement.AppendChild(bodyElement)
	}
	
	// Create the document
	vdoc := dom.NewVDocument(htmlElement, bodyElement)
	vdoc.BaseURI = baseURI
	vdoc.DocumentURI = baseURI
	
	return vdoc, nil
}

// processNode recursively processes an HTML node and its children,
// converting them to our virtual DOM structure.
func processNode(node *html.Node, parent *dom.VElement) {
	switch node.Type {
	case html.ElementNode:
		// Create a new element
		element := dom.NewVElement(strings.ToLower(node.Data))
		
		// Process attributes
		for _, attr := range node.Attr {
			element.SetAttribute(attr.Key, attr.Val)
		}
		
		// Add to parent
		parent.AppendChild(element)
		
		// Process children
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			processNode(child, element)
		}
		
	case html.TextNode:
		// Create a text node and add to parent
		text := dom.NewVText(node.Data)
		parent.AppendChild(text)
		
	case html.DocumentNode:
		// Process children of document node
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			processNode(child, parent)
		}
		
	// Other node types (comments, etc.) are ignored
	}
}

// SerializeToHTML converts a virtual DOM element to an HTML string.
func SerializeToHTML(node dom.VNode) string {
	if node == nil {
		return ""
	}

	if textNode, ok := dom.AsVText(node); ok {
		return html.EscapeString(textNode.TextContent)
	}

	element, ok := dom.AsVElement(node)
	if !ok {
		return ""
	}

	var buf bytes.Buffer

	// Open tag
	buf.WriteString("<")
	buf.WriteString(element.TagName)

	// Attributes
	for key, value := range element.Attributes {
		buf.WriteString(" ")
		buf.WriteString(key)
		buf.WriteString("=\"")
		buf.WriteString(html.EscapeString(value))
		buf.WriteString("\"")
	}

	// Self-closing tags
	selfClosingTags := map[string]bool{
		"area":   true,
		"base":   true,
		"br":     true,
		"col":    true,
		"embed":  true,
		"hr":     true,
		"img":    true,
		"input":  true,
		"link":   true,
		"meta":   true,
		"param":  true,
		"source": true,
		"track":  true,
		"wbr":    true,
	}

	if selfClosingTags[element.TagName] && len(element.Children) == 0 {
		buf.WriteString("/>")
		return buf.String()
	}

	buf.WriteString(">")

	// Children
	for _, child := range element.Children {
		buf.WriteString(SerializeToHTML(child))
	}

	// Close tag
	buf.WriteString("</")
	buf.WriteString(element.TagName)
	buf.WriteString(">")

	return buf.String()
}

// SerializeDocumentToHTML converts a virtual DOM document to an HTML string.
func SerializeDocumentToHTML(doc *dom.VDocument) string {
	if doc == nil || doc.DocumentElement == nil {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>\n")
	buf.WriteString(SerializeToHTML(doc.DocumentElement))
	return buf.String()
}

// SerializeToWriter writes the HTML representation of a node to a writer.
func SerializeToWriter(node dom.VNode, w io.Writer) error {
	_, err := io.WriteString(w, SerializeToHTML(node))
	return err
}

// SerializeDocumentToWriter writes the HTML representation of a document to a writer.
func SerializeDocumentToWriter(doc *dom.VDocument, w io.Writer) error {
	_, err := io.WriteString(w, SerializeDocumentToHTML(doc))
	return err
}
