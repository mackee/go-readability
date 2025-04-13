// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"regexp"
	"strings"

	"github.com/mackee/go-readability/internal/dom"
)

// List of semantic tags to remove (lowercase)
var tagsToRemove = []string{
	"aside",    // Supplementary information not directly related to the main content, like sidebars
	"nav",      // Navigation menus
	"header",   // Page headers
	"footer",   // Page footers
	"script",   // JavaScript
	"style",    // CSS
	"noscript", // Alternative content for when JavaScript is disabled
	"iframe",   // Embedded frames (e.g., ads, social media widgets)
	"form",     // Form elements (e.g., login forms)
	"button",   // Button elements
	"object",   // Embedded objects
	"embed",    // Embedded content
	"applet",   // Old embedded Java applets
	"map",      // Image maps
	"dialog",   // Dialog boxes
	// "audio",  // Audio players
	// "video",  // Video players
	// Excluded because they might be necessary for the main content
	// "figure",  // Figures (with captions)
	// "canvas",  // Canvas elements
	// "details", // Collapsible details information
}

// Patterns for class names or ID names likely indicating ads
var adPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ad-`),
	regexp.MustCompile(`(?i)^ad$`),
	regexp.MustCompile(`(?i)^ads$`),
	regexp.MustCompile(`(?i)advert`),
	regexp.MustCompile(`(?i)banner`),
	regexp.MustCompile(`(?i)sponsor`),
	regexp.MustCompile(`(?i)promo`),
	regexp.MustCompile(`(?i)google-ad`),
	regexp.MustCompile(`(?i)adsense`),
	regexp.MustCompile(`(?i)doubleclick`),
	regexp.MustCompile(`(?i)amazon`),
	regexp.MustCompile(`(?i)affiliate`),
	regexp.MustCompile(`(?i)commercial`),
	regexp.MustCompile(`(?i)paid`),
	regexp.MustCompile(`(?i)shopping`),
	regexp.MustCompile(`(?i)recommendation`),
}

// PreprocessDocument removes noise elements from the document.
// This includes removing semantic tags, unnecessary tags, and ad elements.
// Preprocessing is an important step to clean up the document before content extraction.
//
// Parameters:
//   - doc: The parsed HTML document to preprocess
//
// Returns:
//   - The same document after preprocessing (for method chaining)
func PreprocessDocument(doc *dom.VDocument) *dom.VDocument {
	// 1. Remove semantic tags and unnecessary tags
	removeUnwantedTags(doc)

	// 2. Remove ad elements
	removeAds(doc)

	return doc
}

// removeUnwantedTags removes unwanted tags from the document.
// This removes elements that are unlikely to contain main content, such as
// navigation, scripts, styles, and other non-content elements.
//
// Parameters:
//   - doc: The document to process
func removeUnwantedTags(doc *dom.VDocument) {
	for _, tagName := range tagsToRemove {
		elements := dom.GetElementsByTagName(doc.DocumentElement, tagName)

		// Remove elements from their parent
		for _, element := range elements {
			if parent := element.Parent(); parent != nil {
				for i, child := range parent.Children {
					if child == element {
						parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
						break
					}
				}
			}
		}
	}
}

// removeAds removes ad elements from the document.
// This identifies and removes elements that are likely to be advertisements
// based on class names, IDs, and other attributes.
//
// Parameters:
//   - doc: The document to process
func removeAds(doc *dom.VDocument) {
	// Get all elements under body
	allElements := dom.GetElementsByTagName(doc.Body, "*")

	// Remove elements that seem to be ads
	for _, element := range allElements {
		if isLikelyAd(element) && element.Parent() != nil {
			parent := element.Parent()
			for i, child := range parent.Children {
				if child == element {
					parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
					break
				}
			}
		}
	}
}

// isLikelyAd determines if an element is likely an ad.
// It checks various properties of an element to determine if it's likely
// to be an advertisement, including class names, IDs, and attributes.
//
// Parameters:
//   - element: The element to check
//
// Returns:
//   - true if the element is likely an advertisement, false otherwise
func isLikelyAd(element *dom.VElement) bool {
	// Check class name and ID
	className := element.ClassName()
	id := element.ID()
	combinedString := className + " " + id

	// Check if it matches ad patterns
	for _, pattern := range adPatterns {
		if pattern.MatchString(combinedString) {
			return true
		}
	}

	// Check ad-related attributes
	if element.GetAttribute("role") == "advertisement" ||
		element.HasAttribute("data-ad") ||
		element.HasAttribute("data-ad-client") ||
		element.HasAttribute("data-ad-slot") {
		return true
	}

	return false
}

// isVisible determines if an element is visible.
// It checks CSS properties and attributes to determine if an element
// would be visible to a user in a browser.
//
// Parameters:
//   - element: The element to check
//
// Returns:
//   - true if the element is likely visible, false otherwise
func isVisible(element *dom.VElement) bool {
	// Check style attribute
	style := element.GetAttribute("style")
	if strings.Contains(style, "display: none") ||
		strings.Contains(style, "visibility: hidden") ||
		strings.Contains(style, "opacity: 0") {
		return false
	}

	// Check hidden attribute
	if element.HasAttribute("hidden") {
		return false
	}

	// Check aria-hidden attribute
	if element.GetAttribute("aria-hidden") == "true" {
		return false
	}

	return true
}
