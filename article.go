// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"github.com/mackee/go-readability/internal/dom"
)

// ReadabilityArticle represents the result of a readability extraction.
// It contains the extracted content, metadata, and structural information about the page.
type ReadabilityArticle struct {
	Title     string        // Extracted title
	Byline    string        // Extracted byline/author information
	Root      *dom.VElement // Main content root element (if score threshold is met)
	NodeCount int           // Total number of nodes
	PageType  PageType      // Classification of page type

	// Structural elements (set when PageType is ARTICLE but Root is nil)
	Header                *dom.VElement   // Page header element, if identified
	Footer                *dom.VElement   // Page footer element, if identified
	OtherSignificantNodes []*dom.VElement // Other semantically significant nodes

	// Fallback when article extraction fails
	AriaTree *AriaTree // ARIA tree representation
}

// ArticleContent represents the content of an article page.
// This is a simplified view of ReadabilityArticle focused on article-specific content.
type ArticleContent struct {
	Title  string        // Extracted title
	Byline string        // Extracted byline/author
	Root   *dom.VElement // Main content root element
}

// OtherContent represents the content of a non-article page.
// This is used for pages that don't fit the article pattern, such as index pages,
// landing pages, or other non-article content.
type OtherContent struct {
	Title                 string          // Extracted title
	Header                *dom.VElement   // Page header, if identified
	Footer                *dom.VElement   // Page footer, if identified
	OtherSignificantNodes []*dom.VElement // Other semantically significant nodes
	AriaTree              *AriaTree       // ARIA tree representation
}

// GetContentByPageType returns the appropriate content structure based on page type.
// It returns either ArticleContent or OtherContent depending on the page type.
// This allows consumers to handle different page types with type-specific structures.
//
// Returns:
//   - ArticleContent if the page is classified as an article
//   - OtherContent if the page is classified as any other type
func (r *ReadabilityArticle) GetContentByPageType() interface{} {
	if r.PageType == PageTypeArticle {
		return ArticleContent{
			Title:  r.Title,
			Byline: r.Byline,
			Root:   r.Root,
		}
	} else {
		return OtherContent{
			Title:                 r.Title,
			Header:                r.Header,
			Footer:                r.Footer,
			OtherSignificantNodes: r.OtherSignificantNodes,
			AriaTree:              r.AriaTree,
		}
	}
}
