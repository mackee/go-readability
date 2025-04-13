// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mackee/go-readability/internal/dom"
)

// escapeMarkdown escapes Markdown special characters in text.
// This ensures that special characters like asterisks and underscores are
// treated as literal characters rather than Markdown formatting.
//
// Parameters:
//   - text: The text to escape
//
// Returns:
//   - The escaped text with Markdown special characters escaped
func escapeMarkdown(text string) string {
	// Decode common HTML entities first
	decodedText := text
	decodedText = strings.ReplaceAll(decodedText, "&amp;", "&")
	decodedText = strings.ReplaceAll(decodedText, "&lt;", "<")
	decodedText = strings.ReplaceAll(decodedText, "&gt;", ">")
	decodedText = strings.ReplaceAll(decodedText, "&quot;", "\"")
	decodedText = strings.ReplaceAll(decodedText, "&#039;", "'")
	decodedText = strings.ReplaceAll(decodedText, "&nbsp;", " ")

	// Escape Markdown special characters
	re := regexp.MustCompile(`([*_\[\]\\` + "`" + `])`)
	return re.ReplaceAllString(decodedText, `\$1`)
}

// joinMarkdownParts joins an array of markdown strings, adding spaces where needed between inline elements/text.
// This handles the spacing between elements intelligently, avoiding double spaces
// and ensuring proper spacing around punctuation.
//
// Parameters:
//   - parts: An array of Markdown string parts to join
//
// Returns:
//   - A single Markdown string with proper spacing
func joinMarkdownParts(parts []string) string {
	var result strings.Builder

	for _, part := range parts {
		// Skip parts that are effectively empty after potential whitespace collapsing
		if part == "" || strings.TrimSpace(part) == "" {
			continue
		}

		if result.Len() == 0 {
			// For the first part, just add it
			result.WriteString(part)
		} else {
			// Check if previous result ends with whitespace
			endsWithWhitespace := regexp.MustCompile(`\s$`).MatchString(result.String())
			// Check if current part starts with whitespace
			startsWithWhitespace := regexp.MustCompile(`^\s`).MatchString(part)

			if !endsWithWhitespace && !startsWithWhitespace {
				// Don't add space if current part starts with punctuation
				firstChar := ""
				if len(part) > 0 {
					firstChar = string(part[0])
				}
				if !regexp.MustCompile(`[.,!?;:)]`).MatchString(firstChar) {
					result.WriteString(" ") // Add a single space
				}
			}
			result.WriteString(part)
		}
	}

	return result.String()
}

// getAllTextContent recursively gets all text content from a node.
// This extracts all text content from a node and its descendants,
// which is useful for code blocks and other elements where formatting
// should be preserved.
//
// Parameters:
//   - node: The node to extract text from
//
// Returns:
//   - A string containing all text content
func getAllTextContent(node dom.VNode) string {
	if textNode, ok := dom.AsVText(node); ok {
		return textNode.TextContent
	}

	if elementNode, ok := dom.AsVElement(node); ok {
		var result strings.Builder
		for _, child := range elementNode.Children {
			result.WriteString(getAllTextContent(child))
		}
		return result.String()
	}

	return ""
}

// convertNodeToMarkdown converts a VNode to Markdown string (recursive).
// This is the core function for HTML to Markdown conversion, handling
// different HTML elements appropriately to produce well-formatted Markdown.
//
// Parameters:
//   - node: The node to convert
//   - parentTagName: The tag name of the parent node
//   - depth: The current depth in the document tree
//   - isFirstChild: Whether this node is the first child of its parent
//
// Returns:
//   - A Markdown string representation of the node
func convertNodeToMarkdown(node dom.VNode, parentTagName string, depth int, isFirstChild bool) string {
	if textNode, ok := dom.AsVText(node); ok {
		if parentTagName == "pre" || parentTagName == "code" {
			return textNode.TextContent // Keep raw text
		}
		// Replace sequences of space/tab with a single space
		text := regexp.MustCompile(`[ \t]+`).ReplaceAllString(textNode.TextContent, " ")
		if text == "" {
			return ""
		}
		return escapeMarkdown(text)
	}

	elementNode, ok := dom.AsVElement(node)
	if !ok {
		return ""
	}

	tagName := strings.ToLower(elementNode.TagName)

	// Check if element is block
	isBlock := map[string]bool{
		"p": true, "h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
		"ul": true, "ol": true, "li": true, "pre": true, "blockquote": true, "hr": true,
		"table": true, "div": true,
	}[tagName]

	// Process children, store results in an array
	childrenResults := []string{}
	for i, child := range elementNode.Children {
		isCurrentChildFirst := i == 0
		childResult := convertNodeToMarkdown(child, tagName, func() int {
			if tagName == "ul" || tagName == "ol" || tagName == "blockquote" {
				return depth + 1
			}
			return depth
		}(), isCurrentChildFirst)
		childrenResults = append(childrenResults, childResult)
	}

	// Join children results using the helper function for smart spacing
	childrenMarkdown := joinMarkdownParts(childrenResults)

	// Trim children's markdown for block elements
	trimmedChildren := strings.TrimSpace(childrenMarkdown)

	switch tagName {
	// Headings
	case "h1":
		return fmt.Sprintf("# %s\n\n", trimmedChildren)
	case "h2":
		return fmt.Sprintf("## %s\n\n", trimmedChildren)
	case "h3":
		return fmt.Sprintf("### %s\n\n", trimmedChildren)
	case "h4":
		return fmt.Sprintf("#### %s\n\n", trimmedChildren)
	case "h5":
		return fmt.Sprintf("##### %s\n\n", trimmedChildren)
	case "h6":
		return fmt.Sprintf("###### %s\n\n", trimmedChildren)

	case "p":
		if trimmedChildren == "" {
			return ""
		}
		return fmt.Sprintf("%s\n\n", trimmedChildren)

	// Inline elements
	case "strong", "b":
		return fmt.Sprintf("**%s**", childrenMarkdown)
	case "em", "i":
		return fmt.Sprintf("*%s*", childrenMarkdown)
	case "code":
		if parentTagName != "pre" {
			// Inline code
			codeContent := childrenMarkdown

			// Find all backtick sequences to determine delimiter
			backtickRe := regexp.MustCompile("`+")
			backtickMatches := backtickRe.FindAllString(codeContent, -1)
			longestSequence := 0
			for _, match := range backtickMatches {
				if len(match) > longestSequence {
					longestSequence = len(match)
				}
			}
			delimiter := strings.Repeat("`", longestSequence+1)

			// Check if content consists only of backticks
			onlyBackticksRe := regexp.MustCompile("^`+$")
			if onlyBackticksRe.MatchString(codeContent) && len(codeContent) >= len(delimiter) {
				delimiter = strings.Repeat("`", len(codeContent)+1)
			}

			// Determine if padding is needed
			startsOrEndsWithBacktick := strings.HasPrefix(codeContent, "`") || strings.HasSuffix(codeContent, "`")
			consistsOnlyOfBackticks := onlyBackticksRe.MatchString(codeContent)
			isEmptyOrWhitespace := strings.TrimSpace(codeContent) == ""
			needsPadding := startsOrEndsWithBacktick || consistsOnlyOfBackticks || isEmptyOrWhitespace

			// Apply padding if needed
			finalContent := codeContent
			if needsPadding {
				finalContent = " " + codeContent + " "
			}

			return delimiter + finalContent + delimiter
		}
		// Code inside pre: Return raw content
		return childrenMarkdown

	case "pre":
		codeChild := func() *dom.VElement {
			for _, child := range elementNode.Children {
				if childElement, ok := dom.AsVElement(child); ok && strings.ToLower(childElement.TagName) == "code" {
					return childElement
				}
			}
			return nil
		}()

		// Get all text content recursively
		rawCodeContent := ""
		if codeChild != nil {
			rawCodeContent = getAllTextContent(codeChild)
		} else {
			rawCodeContent = getAllTextContent(elementNode)
		}

		// Extract language class
		lang := ""
		if codeChild != nil {
			classAttr := codeChild.Attributes["class"]
			langMatch := regexp.MustCompile(`language-([a-zA-Z0-9_-]+)`).FindStringSubmatch(classAttr)
			if len(langMatch) > 1 {
				lang = langMatch[1]
			}
		}

		// Clean code content
		cleanedCodeContent := regexp.MustCompile(`^\s*\n|\s+$`).ReplaceAllString(rawCodeContent, "")

		// Special handling for markdown code blocks
		if lang == "markdown" || lang == "md" {
			return fmt.Sprintf("````%s\n%s\n````", lang, cleanedCodeContent)
		}

		// Regular code block
		return fmt.Sprintf("```%s\n%s\n```", lang, cleanedCodeContent)

	case "blockquote":
		content := strings.TrimSpace(childrenMarkdown)
		if content == "" {
			return ""
		}
		lines := strings.Split(content, "\n")
		quotedLines := []string{}
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				quotedLines = append(quotedLines, ">")
			} else {
				quotedLines = append(quotedLines, "> "+line)
			}
		}
		return strings.Join(quotedLines, "\n") + "\n\n"

	case "ul", "ol":
		// Process only li children
		listItems := []string{}
		for _, child := range elementNode.Children {
			if childElement, ok := dom.AsVElement(child); ok && strings.ToLower(childElement.TagName) == "li" {
				childResult := convertNodeToMarkdown(childElement, tagName, depth+1, false)
				if strings.TrimSpace(childResult) != "" {
					listItems = append(listItems, childResult)
				}
			}
		}

		if len(listItems) == 0 {
			return ""
		}

		// Join list items
		listContent := strings.Join(listItems, "\n")

		// Indent the entire list block based on depth
		if depth > 1 {
			listIndent := strings.Repeat("  ", depth-1)
			indentedLines := []string{}
			for _, line := range strings.Split(listContent, "\n") {
				if strings.TrimSpace(line) != "" {
					indentedLines = append(indentedLines, listIndent+line)
				} else {
					indentedLines = append(indentedLines, line)
				}
			}
			listContent = strings.Join(indentedLines, "\n")
		}

		return listContent + "\n\n"

	case "li":
		// Determine marker based on parent
		marker := "-"
		if parentTagName == "ol" {
			marker = "1."
		}

		// Process children, separating main content from nested lists
		mainContentParts := []string{}
		nestedListParts := []string{}

		for _, child := range elementNode.Children {
			if childElement, ok := dom.AsVElement(child); ok {
				childTagName := strings.ToLower(childElement.TagName)
				if childTagName == "ul" || childTagName == "ol" {
					nestedListMd := convertNodeToMarkdown(childElement, tagName, depth+1, false)
					if nestedListMd != "" {
						nestedListParts = append(nestedListParts, regexp.MustCompile(`\n+$`).ReplaceAllString(nestedListMd, ""))
					}
				} else {
					mainContentParts = append(mainContentParts, convertNodeToMarkdown(childElement, tagName, depth, false))
				}
			} else {
				mainContentParts = append(mainContentParts, convertNodeToMarkdown(child, tagName, depth, false))
			}
		}

		// Join main content parts and trim
		mainContent := strings.TrimSpace(joinMarkdownParts(mainContentParts))

		// Format: Marker + Space + Content
		result := fmt.Sprintf("%s %s", marker, mainContent)

		// Append nested lists
		if len(nestedListParts) > 0 {
			if mainContent != "" {
				result += "\n"
			}
			result += strings.Join(nestedListParts, "\n")
		}

		return result

	case "a":
		href := elementNode.Attributes["href"]
		// Clean link content
		linkContent := strings.TrimSpace(strings.ReplaceAll(childrenMarkdown, "\n", " "))

		// Special handling for image links
		if len(elementNode.Children) == 1 {
			if childElement, ok := dom.AsVElement(elementNode.Children[0]); ok && strings.ToLower(childElement.TagName) == "img" {
				alt := childElement.Attributes["alt"]
				src := childElement.Attributes["src"]

				// Use alt if available, otherwise use src
				displayText := src
				if strings.TrimSpace(alt) != "" {
					displayText = alt
				}

				return fmt.Sprintf("[%s](%s)", displayText, href)
			}
		}

		// Regular link
		return fmt.Sprintf("[%s](%s)", linkContent, href)

	case "img":
		alt := escapeMarkdown(elementNode.Attributes["alt"])
		src := elementNode.Attributes["src"]
		title := ""
		if titleAttr, ok := elementNode.Attributes["title"]; ok && titleAttr != "" {
			title = fmt.Sprintf(` "%s"`, escapeMarkdown(titleAttr))
		}

		// If parent is an anchor, just return alt or src
		if parentTagName == "a" {
			if strings.TrimSpace(alt) != "" {
				return alt
			}
			return src
		}

		// Regular image
		return fmt.Sprintf("![%s](%s%s)", alt, src, title)

	case "hr":
		return "---\n\n"

	case "br":
		return "  \n"

	case "table":
		var headerRow []string
		var bodyRows [][]string
		maxColumns := 0

		// Find thead and tbody
		var thead, tbody *dom.VElement
		for _, child := range elementNode.Children {
			if childElement, ok := dom.AsVElement(child); ok {
				childTagName := strings.ToLower(childElement.TagName)
				switch childTagName {
				case "thead":
					thead = childElement
				case "tbody":
					tbody = childElement
				}
			}
		}

		// Process cell content
		processCell := func(cell *dom.VElement) string {
			return strings.TrimSpace(convertNodeToMarkdown(cell, strings.ToLower(cell.TagName), depth+1, false))
		}

		// Process header row
		if thead != nil {
			for _, child := range thead.Children {
				if trElement, ok := dom.AsVElement(child); ok && strings.ToLower(trElement.TagName) == "tr" {
					for _, thChild := range trElement.Children {
						if thElement, ok := dom.AsVElement(thChild); ok && strings.ToLower(thElement.TagName) == "th" {
							headerRow = append(headerRow, processCell(thElement))
						}
					}
					maxColumns = max(maxColumns, len(headerRow))
					break // Only process the first tr
				}
			}
		}

		// Process body rows
		rowsContainer := tbody
		if rowsContainer == nil {
			rowsContainer = elementNode
		}

		for _, child := range rowsContainer.Children {
			if trElement, ok := dom.AsVElement(child); ok && strings.ToLower(trElement.TagName) == "tr" {
				var row []string
				for _, tdChild := range trElement.Children {
					if tdElement, ok := dom.AsVElement(tdChild); ok {
						tdTagName := strings.ToLower(tdElement.TagName)
						if tdTagName == "td" || tdTagName == "th" {
							row = append(row, processCell(tdElement))
						}
					}
				}
				bodyRows = append(bodyRows, row)
				maxColumns = max(maxColumns, len(row))
			}
		}

		// Build Markdown table string
		var tableMd strings.Builder
		separator := strings.Join(func() []string {
			sep := make([]string, maxColumns)
			for i := range sep {
				sep[i] = "---"
			}
			return sep
		}(), " | ")

		if len(headerRow) > 0 {
			// Pad header row if needed
			for len(headerRow) < maxColumns {
				headerRow = append(headerRow, "")
			}
			tableMd.WriteString("| " + strings.Join(headerRow, " | ") + " |\n")
			tableMd.WriteString("| " + separator + " |\n")
		} else if len(bodyRows) > 0 && maxColumns > 0 {
			tableMd.WriteString("| " + separator + " |\n")
		}

		for _, row := range bodyRows {
			// Pad row if needed
			for len(row) < maxColumns {
				row = append(row, "")
			}
			tableMd.WriteString("| " + strings.Join(row, " | ") + " |\n")
		}

		if tableMd.Len() > 0 {
			return strings.TrimSpace(tableMd.String()) + "\n\n"
		}
		return ""

	// Ignored tags
	case "script", "style", "nav", "aside", "header", "footer", "form",
		"button", "iframe", "object", "embed", "applet", "link", "meta",
		"title", "svg":
		return ""

	// Default: Render children for unknown/other tags
	default:
		if isBlock {
			if trimmedChildren != "" {
				return trimmedChildren + "\n\n"
			}
			return ""
		}
		// Assume inline or unknown
		return childrenMarkdown
	}
}

// ToMarkdown converts a VElement to a Markdown string.
// This is the main entry point for HTML to Markdown conversion,
// which produces a well-formatted Markdown document from an HTML element.
//
// Parameters:
//   - element: The HTML element to convert to Markdown
//
// Returns:
//   - A Markdown string representation of the element
func ToMarkdown(element *dom.VElement) string {
	if element == nil {
		return ""
	}

	// Start conversion from the root element
	markdown := convertNodeToMarkdown(element, "", 0, true)

	// Final cleanup
	markdown = strings.TrimSpace(markdown)

	// Normalize block spacing: Replace 3 or more newlines with exactly two
	markdown = regexp.MustCompile(`\n{3,}`).ReplaceAllString(markdown, "\n\n")

	return markdown
}
