// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

// PageType represents the type of a page (article, other, etc.)
// This is used to classify pages based on their content structure and characteristics.
type PageType string

const (
	// PageTypeArticle represents a standard article page
	PageTypeArticle PageType = "article"
	// PageTypeOther represents any page that is not a standard article (e.g., index, list, error)
	PageTypeOther PageType = "other"
	// Future types like INDEX, LIST, ERROR can be added here
)

// ReadabilityOptions contains configuration options for the readability extraction process.
// These options control various aspects of the content extraction algorithm, such as
// thresholds, candidate selection, and output format.
type ReadabilityOptions struct {
	// CharThreshold is the minimum number of characters an article must have
	CharThreshold int
	// NbTopCandidates is the number of top candidates to consider
	NbTopCandidates int
	// GenerateAriaTree indicates whether to generate ARIA tree representation
	GenerateAriaTree bool
	// ForcedPageType allows forcing a specific page type classification
	ForcedPageType PageType
	// Parser is a custom HTML parser function (not used in the Go implementation as we use golang.org/x/net/html)
	// This is kept as a placeholder to match the TypeScript API
	// Parser func(string) (*dom.VDocument, error)
}

// DefaultOptions returns a ReadabilityOptions struct with default values.
// This provides a convenient way to get a pre-configured options object
// with reasonable defaults for most extraction scenarios.
//
// Returns:
//   - A ReadabilityOptions struct initialized with default values
func DefaultOptions() ReadabilityOptions {
	return ReadabilityOptions{
		CharThreshold:    500,   // Default minimum character threshold
		NbTopCandidates:  5,     // Default number of top candidates
		GenerateAriaTree: false, // By default, don't generate ARIA tree
	}
}
