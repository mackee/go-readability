// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"github.com/mackee/go-readability/internal/dom"
)

// GetElementsByTagName returns all elements with the specified tag name in the element tree.
// If tagName is "*", it returns all elements.
//
// Parameters:
//   - element: The root element to search from
//   - tagName: The tag name to search for, or "*" for all elements
//
// Returns:
//   - A slice of elements matching the tag name
func GetElementsByTagName(element *dom.VElement, tagName string) []*dom.VElement {
	return dom.GetElementsByTagName(element, tagName)
}

// GetElementsByTagNames returns all elements with any of the specified tag names in the element tree.
// This is useful for finding elements of multiple types in a single pass.
//
// Parameters:
//   - element: The root element to search from
//   - tagNames: A slice of tag names to search for
//
// Returns:
//   - A slice of elements matching any of the tag names
func GetElementsByTagNames(element *dom.VElement, tagNames []string) []*dom.VElement {
	return dom.GetElementsByTagNames(element, tagNames)
}

// IsProbablyVisible checks if an element is likely to be visible based on its attributes.
// This helps filter out hidden elements that shouldn't be included in the extracted content.
//
// Parameters:
//   - node: The element to check
//
// Returns:
//   - true if the element is likely visible, false otherwise
func IsProbablyVisible(node *dom.VElement) bool {
	return dom.IsProbablyVisible(node)
}

// GetNodeAncestors returns the ancestor elements of a node up to a specified depth.
// If maxDepth is less than or equal to 0, all ancestors are returned.
// This is useful for traversing up the DOM tree to find parent elements.
//
// Parameters:
//   - node: The element to get ancestors for
//   - maxDepth: The maximum number of ancestors to return, or <= 0 for all
//
// Returns:
//   - A slice of ancestor elements, ordered from closest to furthest
func GetNodeAncestors(node *dom.VElement, maxDepth int) []*dom.VElement {
	return dom.GetNodeAncestors(node, maxDepth)
}

// CreateElement creates a new element with the given tag name.
// This is useful for creating new elements to insert into the DOM.
//
// Parameters:
//   - tagName: The tag name for the new element
//
// Returns:
//   - A new VElement with the specified tag name
func CreateElement(tagName string) *dom.VElement {
	return dom.CreateElement(tagName)
}

// CreateTextNode creates a new text node with the given content.
// This is useful for creating text nodes to insert into the DOM.
//
// Parameters:
//   - content: The text content for the new node
//
// Returns:
//   - A new VText node with the specified content
func CreateTextNode(content string) *dom.VText {
	return dom.CreateTextNode(content)
}

// GetAttribute gets the value of an attribute on an element.
// Returns an empty string if the attribute doesn't exist.
//
// Parameters:
//   - element: The element to get the attribute from
//   - name: The name of the attribute to get
//
// Returns:
//   - The attribute value, or an empty string if not found
func GetAttribute(element *dom.VElement, name string) string {
	return dom.GetAttribute(element, name)
}

// HasAncestorTag checks if a node has an ancestor with the specified tag name.
// If maxDepth is less than or equal to 0, all ancestors are checked.
// This is useful for determining if an element is contained within a specific type of element.
//
// Parameters:
//   - node: The node to check ancestors for
//   - tagName: The tag name to look for in ancestors
//   - maxDepth: The maximum depth to check, or <= 0 for unlimited
//
// Returns:
//   - true if an ancestor with the specified tag name is found, false otherwise
func HasAncestorTag(node dom.VNode, tagName string, maxDepth int) bool {
	return dom.HasAncestorTag(node, tagName, maxDepth)
}

// GetInnerText returns the inner text of an element or text node.
// If normalizeSpaces is true, consecutive whitespace is normalized to a single space.
// This extracts all text content from an element and its descendants.
//
// Parameters:
//   - node: The node to get text from
//   - normalizeSpaces: Whether to normalize whitespace
//
// Returns:
//   - The combined text content of the node and its descendants
func GetInnerText(node dom.VNode, normalizeSpaces bool) string {
	return dom.GetInnerText(node, normalizeSpaces)
}

// GetLinkDensity calculates the ratio of link text to all text in an element.
// Returns a value between 0 and 1, where higher values indicate more links.
// This is useful for identifying navigation areas and other link-heavy sections
// that are unlikely to be main content.
//
// Parameters:
//   - element: The element to calculate link density for
//
// Returns:
//   - A float64 between 0 and 1 representing the link density
func GetLinkDensity(element *dom.VElement) float64 {
	return dom.GetLinkDensity(element)
}

// GetTextDensity calculates the ratio of text to child elements in an element.
// Returns a value where higher values indicate more text-dense content.
// This helps identify content-rich elements that are likely to be the main content.
//
// Parameters:
//   - element: The element to calculate text density for
//
// Returns:
//   - A float64 representing the text density
func GetTextDensity(element *dom.VElement) float64 {
	return dom.GetTextDensity(element)
}
