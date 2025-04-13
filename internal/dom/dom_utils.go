// Package dom provides virtual DOM structures and operations for HTML parsing and manipulation.
package dom

import (
	"regexp"
	"strings"
)

// normalizeRegexp is used to normalize whitespace in text
var normalizeRegexp = regexp.MustCompile(`\s{2,}`)

// GetElementsByTagName returns all elements with the specified tag name(s) in the element tree.
// If tagName is "*", it returns all elements.
func GetElementsByTagName(element *VElement, tagName string) []*VElement {
	return getElementsByTagNameInternal(element, []string{strings.ToLower(tagName)})
}

// GetElementsByTagNames returns all elements with any of the specified tag names in the element tree.
func GetElementsByTagNames(element *VElement, tagNames []string) []*VElement {
	// Convert all tag names to lowercase for case-insensitive matching
	lowerTagNames := make([]string, len(tagNames))
	for i, tag := range tagNames {
		lowerTagNames[i] = strings.ToLower(tag)
	}
	return getElementsByTagNameInternal(element, lowerTagNames)
}

// getElementsByTagNameInternal is the internal implementation for GetElementsByTagName and GetElementsByTagNames.
func getElementsByTagNameInternal(element *VElement, tagNames []string) []*VElement {
	var result []*VElement

	// Check if this element matches (using lowercase)
	for _, tag := range tagNames {
		if tag == "*" || tag == element.TagName {
			result = append(result, element)
			break
		}
	}

	// Recursively check child elements
	for _, child := range element.Children {
		if childElement, ok := AsVElement(child); ok {
			result = append(result, getElementsByTagNameInternal(childElement, tagNames)...)
		}
	}

	return result
}

// IsProbablyVisible checks if an element is likely to be visible based on its attributes.
func IsProbablyVisible(node *VElement) bool {
	style := node.GetAttribute("style")
	hidden := node.HasAttribute("hidden")
	ariaHidden := node.GetAttribute("aria-hidden") == "true"

	return !(strings.Contains(style, "display: none") ||
		strings.Contains(style, "visibility: hidden") ||
		hidden ||
		ariaHidden)
}

// GetNodeAncestors returns the ancestor elements of a node up to a specified depth.
// If maxDepth is less than or equal to 0, all ancestors are returned.
func GetNodeAncestors(node *VElement, maxDepth int) []*VElement {
	ancestors := make([]*VElement, 0)
	currentNode := node.Parent()
	depth := 0

	for currentNode != nil && (maxDepth <= 0 || depth < maxDepth) {
		ancestors = append(ancestors, currentNode)
		currentNode = currentNode.Parent()
		depth++
	}

	return ancestors
}

// CreateElement creates a new element with the given tag name.
func CreateElement(tagName string) *VElement {
	return NewVElement(strings.ToLower(tagName))
}

// CreateTextNode creates a new text node with the given content.
func CreateTextNode(content string) *VText {
	return NewVText(content)
}

// GetAttribute gets the value of an attribute on an element.
// Returns an empty string if the attribute doesn't exist.
func GetAttribute(element *VElement, name string) string {
	return element.GetAttribute(name)
}

// HasAncestorTag checks if a node has an ancestor with the specified tag name.
// If maxDepth is less than or equal to 0, all ancestors are checked.
func HasAncestorTag(node VNode, tagName string, maxDepth int) bool {
	tagName = strings.ToLower(tagName)
	depth := 0
	
	var currentNode *VElement
	if element, ok := AsVElement(node); ok {
		currentNode = element.Parent()
	} else if text, ok := AsVText(node); ok {
		currentNode = text.Parent()
	} else {
		return false
	}

	for currentNode != nil {
		if maxDepth > 0 && depth >= maxDepth {
			return false
		}

		if currentNode.TagName == tagName {
			return true
		}

		currentNode = currentNode.Parent()
		depth++
	}

	return false
}

// GetInnerText returns the inner text of an element or text node.
// If normalizeSpaces is true, consecutive whitespace is normalized to a single space.
func GetInnerText(node VNode, normalizeSpaces bool) string {
	var text string

	switch n := node.(type) {
	case *VText:
		text = n.TextContent
	case *VElement:
		for i, child := range n.Children {
			// Add space between text nodes if not the first child
			if i > 0 && text != "" {
				text += " "
			}
			
			if childText, ok := AsVText(child); ok {
				text += childText.TextContent
			} else if childElement, ok := AsVElement(child); ok {
				childText := GetInnerText(childElement, false)
				if childText != "" {
					text += childText
				}
			}
		}
	}

	text = strings.TrimSpace(text)

	if normalizeSpaces {
		text = normalizeRegexp.ReplaceAllString(text, " ")
	}

	return text
}

// GetLinkDensity calculates the ratio of link text to all text in an element.
// Returns a value between 0 and 1, where higher values indicate more links.
func GetLinkDensity(element *VElement) float64 {
	textLength := len(GetInnerText(element, true))
	if textLength == 0 {
		return 0
	}

	var linkLength int
	links := GetElementsByTagName(element, "a")

	for _, link := range links {
		href := link.GetAttribute("href")
		coefficient := 1.0
		if strings.HasPrefix(href, "#") {
			coefficient = 0.3
		}
		linkLength += int(float64(len(GetInnerText(link, true))) * coefficient)
	}

	return float64(linkLength) / float64(textLength)
}

// GetTextDensity calculates the ratio of text to child elements in an element.
// Returns a value where higher values indicate more text-dense content.
func GetTextDensity(element *VElement) float64 {
	text := GetInnerText(element, true)
	textLength := len(text)
	if textLength == 0 {
		return 0
	}

	var childElementCount int
	for _, child := range element.Children {
		if _, ok := AsVElement(child); ok {
			childElementCount++
		}
	}

	if childElementCount == 0 {
		childElementCount = 1
	}

	return float64(textLength) / float64(childElementCount)
}
