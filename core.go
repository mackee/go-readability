// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"strings"

	"github.com/mackee/go-readability/internal/dom"
	"github.com/mackee/go-readability/internal/util"
)

// Extract extracts the article content from HTML.
// This is the main entry point for the readability extraction process.
// It parses the HTML, preprocesses the document, and extracts the main content
// based on the provided options.
//
// Parameters:
//   - html: The HTML string to extract content from
//   - options: Configuration options for the extraction process
//
// Returns:
//   - A ReadabilityArticle containing the extracted content and metadata
//   - An error if the HTML parsing fails
func Extract(html string, options ReadabilityOptions) (ReadabilityArticle, error) {
	// Parse HTML to create virtual DOM
	doc, err := ParseHTML(html, "")
	if err != nil {
		return ReadabilityArticle{}, err
	}

	// Execute preprocessing
	PreprocessDocument(doc)

	// Set default values if not provided
	if options.CharThreshold <= 0 {
		options.CharThreshold = util.DefaultCharThreshold
	}

	if options.NbTopCandidates <= 0 {
		options.NbTopCandidates = util.DefaultNTopCandidates
	}

	// Set default page type if not specified
	if options.ForcedPageType == "" {
		options.ForcedPageType = PageTypeArticle
	}

	// Extract content
	return ExtractContent(doc, options), nil
}

// ExtractContent extracts the main content from a document.
// This is the core function for content extraction that implements the main
// readability algorithm to identify and extract the primary content.
//
// Parameters:
//   - doc: The parsed HTML document as a VDocument
//   - options: Configuration options for the extraction process
//
// Returns:
//   - A ReadabilityArticle containing the extracted content and metadata
func ExtractContent(doc *dom.VDocument, options ReadabilityOptions) ReadabilityArticle {
	// Set default values if not provided
	charThreshold := options.CharThreshold
	if charThreshold <= 0 {
		charThreshold = util.DefaultCharThreshold
	}

	nbTopCandidates := options.NbTopCandidates
	if nbTopCandidates <= 0 {
		nbTopCandidates = util.DefaultNTopCandidates
	}

	generateAriaTree := options.GenerateAriaTree

	// Find content candidates
	candidates := FindMainCandidates(doc, nbTopCandidates)
	var topCandidate *dom.VElement
	var articleContent *dom.VElement

	// Select the best candidate if any exist
	if len(candidates) > 0 {
		topCandidate = candidates[0] // Highest scoring candidate

		// Check if the candidate contains meaningful content
		textLength := len(GetInnerText(topCandidate, false))
		linkDensity := GetLinkDensity(topCandidate)

		// If the candidate has enough text and low link density, it's probably content
		if textLength >= charThreshold && linkDensity <= 0.5 {
			articleContent = topCandidate
		}
	}

	// Determine page type (forced or auto-detected)
	pageType := options.ForcedPageType
	if pageType == "" {
		// If we found content, it's probably an article
		if articleContent != nil {
			pageType = PageTypeArticle
		} else {
			pageType = ClassifyPageType(doc, candidates, charThreshold, "")
		}
	}

	// Get metadata
	title := GetArticleTitle(doc)
	byline := GetArticleByline(doc)

	// Detect structural elements if needed (for ARTICLE type but no content found)
	var header *dom.VElement
	var footer *dom.VElement
	var otherSignificantNodes []*dom.VElement

	if pageType == PageTypeArticle && articleContent == nil {
		header, footer, otherSignificantNodes = FindStructuralElements(doc)
	}

	// Generate AriaTree if requested or if no content was found
	var ariaTree *AriaTree
	if generateAriaTree || (articleContent == nil && pageType == PageTypeArticle) {
		// AriaTree generation would be implemented here
		// For now, we'll leave it as nil
		ariaTree = nil
	}

	// Create and return the article
	return ReadabilityArticle{
		Title:                 title,
		Byline:                byline,
		Root:                  articleContent,
		NodeCount:             CountNodes(articleContent),
		PageType:              pageType,
		Header:                header,
		Footer:                footer,
		OtherSignificantNodes: otherSignificantNodes,
		AriaTree:              ariaTree,
	}
}

// FindStructuralElements detects header, footer, and other significant structural elements in a document.
// This is particularly useful for pages that are classified as articles but where the main content
// extraction fails to meet the threshold. It uses semantic tags, ARIA roles, and common class/ID patterns
// to identify important page structures.
//
// Parameters:
//   - doc: The parsed HTML document
//
// Returns:
//   - header: The identified page header element, if found
//   - footer: The identified page footer element, if found
//   - otherSignificantNodes: Other semantically significant elements found in the document
func FindStructuralElements(doc *dom.VDocument) (
	header *dom.VElement,
	footer *dom.VElement,
	otherSignificantNodes []*dom.VElement,
) {
	body := doc.Body

	// 1. Look for header candidates
	headerTags := GetElementsByTagName(doc.DocumentElement, "header")
	if len(headerTags) == 1 {
		header = headerTags[0]
	} else {
		// Look for role="banner" or common ID/class names
		allElements := GetElementsByTagName(body, "*")
		for _, el := range allElements {
			role := strings.ToLower(GetAttribute(el, "role"))
			id := strings.ToLower(el.ID())
			className := strings.ToLower(el.ClassName())

			if role == "banner" ||
				id == "header" ||
				id == "masthead" ||
				strings.Contains(className, "header") ||
				strings.Contains(className, "masthead") {
				// Prefer elements closer to the body (higher in the DOM)
				if header == nil || (el.Parent() == body && header.Parent() != body) {
					header = el
				}
			}
		}
	}

	// 2. Look for footer candidates
	footerTags := GetElementsByTagName(doc.DocumentElement, "footer")
	if len(footerTags) == 1 {
		footer = footerTags[0]
	} else {
		// Look for role="contentinfo" or common ID/class names
		allElements := GetElementsByTagName(body, "*")
		// Search from the bottom of the DOM as footers are typically at the bottom
		for i := len(allElements) - 1; i >= 0; i-- {
			el := allElements[i]
			role := strings.ToLower(GetAttribute(el, "role"))
			id := strings.ToLower(el.ID())
			className := strings.ToLower(el.ClassName())

			if role == "contentinfo" ||
				id == "footer" ||
				id == "colophon" ||
				strings.Contains(className, "footer") ||
				strings.Contains(className, "site-info") {
				// Only set if not already found
				if footer == nil {
					// Exclude footers inside headers
					isInsideHeader := false
					current := el
					for current != nil && current != body {
						if current == header {
							isInsideHeader = true
							break
						}
						current = current.Parent()
					}
					if !isInsideHeader {
						footer = el
					}
				}
			}
		}
	}

	// 3. Find other significant nodes (<main>, <article>, <section>, <aside>, <nav>, etc.)
	mainTags := GetElementsByTagName(body, "main")
	articleTags := GetElementsByTagName(body, "article")
	sectionTags := GetElementsByTagName(body, "section")
	asideTags := GetElementsByTagName(body, "aside")
	navTags := GetElementsByTagName(body, "nav")

	potentialNodes := []*dom.VElement{}
	potentialNodes = append(potentialNodes, mainTags...)
	potentialNodes = append(potentialNodes, articleTags...)
	potentialNodes = append(potentialNodes, sectionTags...)
	potentialNodes = append(potentialNodes, asideTags...)
	potentialNodes = append(potentialNodes, navTags...)

	// Add elements with significant class names or IDs
	AddSignificantElementsByClassOrId(body, &potentialNodes)

	// Filter out nodes inside header or footer
	for _, node := range potentialNodes {
		// Check if node is inside header or footer
		isInsideHeaderOrFooter := false
		current := node
		for current != nil && current != body {
			if current == header || current == footer {
				isInsideHeaderOrFooter = true
				break
			}
			current = current.Parent()
		}

		// Check if node is already in the list
		alreadyIncluded := false
		for _, n := range otherSignificantNodes {
			if n == node {
				alreadyIncluded = true
				break
			}
		}

		if !isInsideHeaderOrFooter && !alreadyIncluded {
			// Check if node is visible and has significant content
			if IsProbablyVisible(node) && (IsSignificantNode(node) || IsSemanticTag(node)) {
				otherSignificantNodes = append(otherSignificantNodes, node)
			}
		}
	}

	return header, footer, otherSignificantNodes
}

// AddSignificantElementsByClassOrId detects elements with meaningful class names or IDs
// and adds them to the potentialNodes slice. This helps identify content containers
// that might not use semantic HTML tags but follow common naming conventions.
//
// Parameters:
//   - body: The body element to search within
//   - potentialNodes: A pointer to a slice where identified elements will be added
func AddSignificantElementsByClassOrId(body *dom.VElement, potentialNodes *[]*dom.VElement) {
	allElements := GetElementsByTagName(body, "*")

	// Patterns for significant class names or IDs
	significantPatterns := []string{
		"content",
		"main",
		"article",
		"post",
		"entry",
		"body",
		"text",
		"story",
		"container",
		"wrapper",
		"page",
		"blog",
		"section",
	}

	for _, el := range allElements {
		className := strings.ToLower(el.ClassName())
		id := strings.ToLower(el.ID())
		combinedString := className + " " + id

		// Check if element has a significant class name or ID
		for _, pattern := range significantPatterns {
			if strings.Contains(combinedString, pattern) {
				// Check if element is already in the list
				alreadyIncluded := false
				for _, n := range *potentialNodes {
					if n == el {
						alreadyIncluded = true
						break
					}
				}

				if !alreadyIncluded {
					*potentialNodes = append(*potentialNodes, el)
				}
				break
			}
		}
	}
}

// min returns the smaller of x or y.
// A simple utility function for integer comparison.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// minFloat returns the smaller of x or y.
// A simple utility function for float64 comparison.
func minFloat(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

// FindMainCandidates detects nodes that are likely to be the main content candidates, sorted by score.
// It implements the core scoring algorithm of readability, analyzing elements based on
// content length, tag types, class names, and other heuristics to identify the most
// likely content containers.
//
// Parameters:
//   - doc: The parsed HTML document
//   - nbTopCandidates: The number of top candidates to return
//
// Returns:
//   - A slice of the top N candidate elements, sorted by score in descending order
func FindMainCandidates(doc *dom.VDocument, nbTopCandidates int) []*dom.VElement {
	// Use default value if nbTopCandidates is not provided
	if nbTopCandidates <= 0 {
		nbTopCandidates = util.DefaultNTopCandidates
	}

	// 1. First, look for semantic tags (simple method)
	semanticTags := []string{"article", "main"}
	for _, tag := range semanticTags {
		elements := GetElementsByTagName(doc.DocumentElement, tag)
		if len(elements) == 1 {
			// If a single semantic tag is found, return it as the only candidate
			return []*dom.VElement{elements[0]}
		}
	}

	// 2. Scoring-based detection
	body := doc.Body
	candidates := []*dom.VElement{}
	elementsToScore := []*dom.VElement{}

	// Collect elements to score
	for _, tag := range util.DefaultTagsToScore {
		elements := GetElementsByTagName(body, tag)
		elementsToScore = append(elementsToScore, elements...)
	}

	// Score each element
	for _, elementToScore := range elementsToScore {
		// Ignore elements with less than 25 characters
		innerText := GetInnerText(elementToScore, false)
		if len(innerText) < 25 {
			continue
		}

		// Get ancestor elements (up to 3 levels)
		ancestors := GetNodeAncestors(elementToScore, 3)
		if len(ancestors) == 0 {
			continue
		}

		// Calculate base score
		contentScore := 1.0                                                            // Base points
		contentScore += float64(len(util.Regexps.Commas.FindAllString(innerText, -1))) // Number of commas
		contentScore += float64(min(len(innerText)/100, 3))                            // Text length (max 3 points)

		// Add score to ancestor elements
		for level, ancestor := range ancestors {
			if ancestor.GetReadabilityData() == nil {
				InitializeNode(ancestor)
				candidates = append(candidates, ancestor)
			}

			// Decrease score for deeper levels
			scoreDivider := 1
			if level == 1 {
				scoreDivider = 2
			} else if level >= 2 {
				scoreDivider = level * 3
			}

			if ancestor.GetReadabilityData() != nil {
				ancestor.GetReadabilityData().ContentScore += contentScore / float64(scoreDivider)
			}
		}
	}

	// Score and select candidates
	type scoredCandidate struct {
		element *dom.VElement
		score   float64
	}
	scoredCandidates := []scoredCandidate{}

	for _, candidate := range candidates {
		// Adjust score based on link density
		if candidate.GetReadabilityData() != nil {
			linkDensity := GetLinkDensity(candidate)
			candidate.GetReadabilityData().ContentScore *= (1.0 - linkDensity)

			// Also consider text density
			// Elements with high text density are more likely to contain more text content
			textDensity := GetTextDensity(candidate)
			if textDensity > 0 {
				// Slightly increase the score for higher text density (up to 10%)
				candidate.GetReadabilityData().ContentScore *= (1.0 + minFloat(textDensity/10.0, 0.1))
			}

			// Check parent node score - the parent might be a better candidate
			currentCandidate := candidate
			parentOfCandidate := currentCandidate.Parent()
			for parentOfCandidate != nil && strings.ToLower(parentOfCandidate.TagName) != "body" {
				if parentOfCandidate.GetReadabilityData() != nil && currentCandidate.GetReadabilityData() != nil &&
					parentOfCandidate.GetReadabilityData().ContentScore > currentCandidate.GetReadabilityData().ContentScore {
					currentCandidate = parentOfCandidate
				}
				parentOfCandidate = parentOfCandidate.Parent()
			}

			// Avoid adding duplicates if parent check resulted in the same element
			// Also ensure readability property exists before accessing contentScore
			if currentCandidate.GetReadabilityData() != nil {
				isDuplicate := false
				for _, sc := range scoredCandidates {
					if sc.element == currentCandidate {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					scoredCandidates = append(scoredCandidates, scoredCandidate{
						element: currentCandidate,
						score:   currentCandidate.GetReadabilityData().ContentScore,
					})
				}
			}
		}
	}

	// Sort candidates by score in descending order
	// Using a simple bubble sort for clarity (can be optimized if needed)
	for i := 0; i < len(scoredCandidates)-1; i++ {
		for j := 0; j < len(scoredCandidates)-i-1; j++ {
			if scoredCandidates[j].score < scoredCandidates[j+1].score {
				scoredCandidates[j], scoredCandidates[j+1] = scoredCandidates[j+1], scoredCandidates[j]
			}
		}
	}

	// Return top N candidates
	topCandidates := []*dom.VElement{}
	for i := 0; i < min(len(scoredCandidates), nbTopCandidates); i++ {
		topCandidates = append(topCandidates, scoredCandidates[i].element)
	}

	// Return body if no candidate is found and body exists
	if len(topCandidates) == 0 && doc.Body != nil {
		return []*dom.VElement{doc.Body}
	}

	return topCandidates
}

// IsProbablyContent determines content probability (simplified version similar to isProbablyReaderable).
// It checks various properties of an element to determine if it's likely to contain
// meaningful content, including visibility, class/ID patterns, text length, and link density.
//
// Parameters:
//   - element: The element to evaluate
//
// Returns:
//   - true if the element is likely to contain meaningful content, false otherwise
func IsProbablyContent(element *dom.VElement) bool {
	// Visibility check
	if !IsProbablyVisible(element) {
		return false
	}

	// Check class name and ID
	className := element.ClassName()
	id := element.ID()
	matchString := className + " " + id

	if util.Regexps.UnlikelyCandidates.MatchString(matchString) &&
		!util.Regexps.OkMaybeItsACandidate.MatchString(matchString) {
		return false
	}

	// Check text length
	textLength := len(GetInnerText(element, false))
	if textLength < 140 {
		return false
	}

	// Check link density
	linkDensity := GetLinkDensity(element)
	if linkDensity > 0.5 {
		return false
	}

	// Check text density
	// If text density is extremely low, it's unlikely to be the main content
	textDensity := GetTextDensity(element)
	return textDensity >= 0.1
}

// InitializeNode initializes a node with a readability score.
// It sets an initial score based on the tag name and adjusts it based on class name and ID.
// This is a key part of the content scoring algorithm, establishing baseline scores
// for different HTML elements.
//
// Parameters:
//   - node: The element to initialize with a readability score
func InitializeNode(node *dom.VElement) {
	// Create a new ReadabilityData with initial score of 0
	node.SetReadabilityData(&dom.ReadabilityData{
		ContentScore: 0,
	})

	// Initial score based on tag name (case-insensitive)
	switch strings.ToLower(node.TagName) {
	case "div":
		node.GetReadabilityData().ContentScore += 5
	case "pre", "td", "blockquote":
		node.GetReadabilityData().ContentScore += 3
	case "address", "ol", "ul", "dl", "dd", "dt", "li", "form":
		node.GetReadabilityData().ContentScore -= 3
	case "h1", "h2", "h3", "h4", "h5", "h6", "th":
		node.GetReadabilityData().ContentScore -= 5
	}

	// Score adjustment based on class name and ID
	node.GetReadabilityData().ContentScore += GetClassWeight(node)
}

// CreateExtractor creates a custom extractor function with specific options.
// This is useful when you want to reuse the same extraction configuration multiple times.
// The returned function can be called with HTML strings to extract content using the
// predefined options.
//
// Parameters:
//   - options: The readability options to use for all extractions
//
// Returns:
//   - A function that takes an HTML string and returns a ReadabilityArticle and error
func CreateExtractor(options ReadabilityOptions) func(string) (ReadabilityArticle, error) {
	return func(html string) (ReadabilityArticle, error) {
		return Extract(html, options)
	}
}

// GetClassWeight calculates a score adjustment based on the class name and ID of an element.
// It returns a positive score for elements likely to contain content and a negative score for elements
// likely to be noise. This helps the algorithm prioritize content-rich elements and
// deprioritize elements that typically contain non-content material.
//
// Parameters:
//   - node: The element to calculate a class weight for
//
// Returns:
//   - A float64 score adjustment (positive for likely content, negative for likely noise)
func GetClassWeight(node *dom.VElement) float64 {
	var weight float64 = 0

	// Check class name
	className := node.ClassName()
	if className != "" {
		if util.Regexps.Negative.MatchString(className) {
			weight -= 25
		}
		if util.Regexps.Positive.MatchString(className) {
			weight += 25
		}
	}

	// Check ID
	id := node.ID()
	if id != "" {
		if util.Regexps.Negative.MatchString(id) {
			weight -= 25
		}
		if util.Regexps.Positive.MatchString(id) {
			weight += 25
		}
	}

	return weight
}
