// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"io"

	"github.com/mackee/go-readability/internal/dom"
	"github.com/mackee/go-readability/internal/parser"
)

// ParseHTML parses an HTML string and returns a virtual DOM document.
// It uses golang.org/x/net/html for parsing and converts the result to our internal DOM structure.
// The baseURI parameter is used to resolve relative URLs in the document.
//
// Parameters:
//   - htmlContent: The HTML string to parse
//   - baseURI: The base URI for resolving relative URLs (can be empty)
//
// Returns:
//   - A pointer to a VDocument representing the parsed HTML
//   - An error if parsing fails
func ParseHTML(htmlContent string, baseURI string) (*dom.VDocument, error) {
	return parser.ParseHTML(htmlContent, baseURI)
}

// SerializeToHTML converts a virtual DOM element to an HTML string.
// This is useful for converting a VNode back to an HTML string after processing.
//
// Parameters:
//   - node: The VNode to serialize
//
// Returns:
//   - An HTML string representation of the node
func SerializeToHTML(node dom.VNode) string {
	return parser.SerializeToHTML(node)
}

// SerializeDocumentToHTML converts a virtual DOM document to an HTML string.
// This serializes an entire document, including the doctype and HTML structure.
//
// Parameters:
//   - doc: The VDocument to serialize
//
// Returns:
//   - An HTML string representation of the document
func SerializeDocumentToHTML(doc *dom.VDocument) string {
	return parser.SerializeDocumentToHTML(doc)
}

// SerializeToWriter writes the HTML representation of a node to a writer.
// This is useful for streaming HTML output to a file or response writer.
//
// Parameters:
//   - node: The VNode to serialize
//   - w: The io.Writer to write to
//
// Returns:
//   - An error if writing fails
func SerializeToWriter(node dom.VNode, w io.Writer) error {
	return parser.SerializeToWriter(node, w)
}

// SerializeDocumentToWriter writes the HTML representation of a document to a writer.
// This serializes an entire document to a writer, which is useful for streaming
// HTML output to a file or response writer.
//
// Parameters:
//   - doc: The VDocument to serialize
//   - w: The io.Writer to write to
//
// Returns:
//   - An error if writing fails
func SerializeDocumentToWriter(doc *dom.VDocument, w io.Writer) error {
	return parser.SerializeDocumentToWriter(doc, w)
}
