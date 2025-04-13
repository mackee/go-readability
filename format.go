// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"strings"

	"github.com/mackee/go-readability/internal/dom"
)

// selfClosingTags is a set of HTML tags that are self-closing.
var selfClosingTags = map[string]bool{
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

// blockElements is a set of HTML tags that are block-level elements.
var blockElements = map[string]bool{
	"address":    true,
	"article":    true,
	"aside":      true,
	"blockquote": true,
	"details":    true,
	"dialog":     true,
	"dd":         true,
	"div":        true,
	"dl":         true,
	"dt":         true,
	"fieldset":   true,
	"figcaption": true,
	"figure":     true,
	"footer":     true,
	"form":       true,
	"h1":         true,
	"h2":         true,
	"h3":         true,
	"h4":         true,
	"h5":         true,
	"h6":         true,
	"header":     true,
	"hgroup":     true,
	"hr":         true,
	"li":         true,
	"main":       true,
	"nav":        true,
	"ol":         true,
	"p":          true,
	"pre":        true,
	"section":    true,
	"table":      true,
	"ul":         true,
}

// ToHTML generates HTML string from VElement, omitting span tags and class attributes.
// This produces a cleaner HTML representation of the extracted content by removing
// unnecessary styling and presentation elements.
//
// Parameters:
//   - element: The element to convert to HTML
//
// Returns:
//   - A string containing the HTML representation of the element
func ToHTML(element *dom.VElement) string {
	if element == nil {
		return ""
	}

	tagName := strings.ToLower(element.TagName)

	// Omit span tags, process children directly
	if tagName == "span" {
		var result strings.Builder
		for _, child := range element.Children {
			if text, ok := dom.AsVText(child); ok {
				result.WriteString(escapeHTML(text.TextContent))
			} else if elem, ok := dom.AsVElement(child); ok {
				result.WriteString(ToHTML(elem))
			}
		}
		return result.String()
	}

	// Generate attribute string, excluding 'class'
	var attrs strings.Builder
	for key, value := range element.Attributes {
		if key != "class" { // Exclude class attribute
			if attrs.Len() > 0 {
				attrs.WriteString(" ")
			}
			attrs.WriteString(key)
			attrs.WriteString("=\"")
			attrs.WriteString(escapeHTML(value))
			attrs.WriteString("\"")
		}
	}

	// For self-closing tags
	if selfClosingTags[tagName] && len(element.Children) == 0 {
		if attrs.Len() > 0 {
			return "<" + tagName + " " + attrs.String() + "/>"
		}
		return "<" + tagName + "/>"
	}

	// Start tag
	var result strings.Builder
	if attrs.Len() > 0 {
		result.WriteString("<" + tagName + " " + attrs.String() + ">")
	} else {
		result.WriteString("<" + tagName + ">")
	}

	// Process child elements
	for _, child := range element.Children {
		if text, ok := dom.AsVText(child); ok {
			result.WriteString(escapeHTML(text.TextContent))
		} else if elem, ok := dom.AsVElement(child); ok {
			result.WriteString(ToHTML(elem))
		}
	}

	// End tag
	result.WriteString("</" + tagName + ">")

	return result.String()
}

// escapeHTML escapes HTML special characters.
// This prevents XSS and other security issues when outputting HTML content.
//
// Parameters:
//   - str: The string to escape
//
// Returns:
//   - The escaped string with HTML special characters replaced with entities
func escapeHTML(str string) string {
	result := strings.ReplaceAll(str, "&", "&amp;")         // Must be first
	result = strings.ReplaceAll(result, "\u00a0", "&nbsp;") // Handle non-breaking space
	result = strings.ReplaceAll(result, "<", "&lt;")
	result = strings.ReplaceAll(result, ">", "&gt;")
	result = strings.ReplaceAll(result, "\"", "&quot;")
	result = strings.ReplaceAll(result, "'", "&#039;")
	return result
}

// Stringify converts VElement to a readable string format.
// Removes tags while applying line breaks considering block and inline elements.
// Aligns all text to the shallowest indent.
// Merges consecutive line breaks into one.
//
// Parameters:
//   - element: The element to convert to a string
//
// Returns:
//   - A plain text representation of the element's content
func Stringify(element *dom.VElement) string {
	if element == nil {
		return ""
	}

	tagName := strings.ToLower(element.TagName)
	isBlock := blockElements[tagName]

	// Handle special tags
	if tagName == "br" {
		return "\n"
	}

	if tagName == "hr" {
		return "\n----------\n"
	}

	var result strings.Builder

	// Insert line break before block elements
	if isBlock {
		result.WriteString("\n")
	}

	// Process child elements
	for _, child := range element.Children {
		if text, ok := dom.AsVText(child); ok {
			// Append text node directly
			trimmedText := strings.TrimSpace(text.TextContent)
			if trimmedText != "" {
				result.WriteString(trimmedText)
				result.WriteString(" ")
			}
		} else if elem, ok := dom.AsVElement(child); ok {
			// Recursively process element nodes
			childResult := Stringify(elem)

			// Add the child result to our result
			result.WriteString(childResult)

			// Add a space after the child content if it doesn't end with a space or newline
			if len(childResult) > 0 &&
				!strings.HasSuffix(childResult, " ") &&
				!strings.HasSuffix(childResult, "\n") {
				result.WriteString(" ")
			}
		}
	}

	// Remove trailing space
	resultStr := result.String()
	if len(resultStr) > 0 && resultStr[len(resultStr)-1] == ' ' {
		resultStr = resultStr[:len(resultStr)-1]
	}

	// Insert line break after block elements
	if isBlock {
		resultStr += "\n"
	}

	// Merge consecutive line breaks into one
	resultStr = strings.ReplaceAll(resultStr, "\n\n", "\n")
	for strings.Contains(resultStr, "\n\n") {
		resultStr = strings.ReplaceAll(resultStr, "\n\n", "\n")
	}

	return resultStr
}

// FormatDocument formats the entire document.
// Merges consecutive line breaks into one, removes extra line breaks at the beginning and end.
// This produces a cleaner, more readable text output.
//
// Parameters:
//   - text: The text to format
//
// Returns:
//   - The formatted text
func FormatDocument(text string) string {
	// Merge consecutive line breaks
	result := text
	for strings.Contains(result, "\n\n") {
		result = strings.ReplaceAll(result, "\n\n", "\n")
	}

	// Remove leading line breaks
	result = strings.TrimLeft(result, "\n")

	// Remove trailing line breaks
	result = strings.TrimRight(result, "\n")

	// Remove leading/trailing whitespace
	return strings.TrimSpace(result)
}

// ExtractTextContent extracts text content from VElement.
// This returns only the text nodes' content, without any HTML formatting.
//
// Parameters:
//   - element: The element to extract text from
//
// Returns:
//   - A string containing all text content from the element and its descendants
func ExtractTextContent(element *dom.VElement) string {
	if element == nil {
		return ""
	}

	var result strings.Builder
	for _, child := range element.Children {
		if text, ok := dom.AsVText(child); ok {
			result.WriteString(text.TextContent)
		} else if elem, ok := dom.AsVElement(child); ok {
			result.WriteString(ExtractTextContent(elem))
		}
	}
	return result.String()
}

// CountNodes counts the number of nodes within a VElement.
// This includes the element itself and all its descendants (both elements and text nodes).
//
// Parameters:
//   - element: The element to count nodes for
//
// Returns:
//   - The total number of nodes
func CountNodes(element *dom.VElement) int {
	if element == nil {
		return 0
	}

	// Count itself as 1
	count := 1

	// Recursively count child elements
	for _, child := range element.Children {
		if elem, ok := dom.AsVElement(child); ok {
			count += CountNodes(elem)
		} else {
			// Count text nodes as 1
			count++
		}
	}

	return count
}
