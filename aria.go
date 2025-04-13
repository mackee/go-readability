// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mackee/go-readability/internal/dom"
)

// AriaNodeType represents the type of an ARIA node.
type AriaNodeType string

// ARIA node types
const (
	// ARIA landmark roles
	AriaNodeTypeBanner        AriaNodeType = "banner"
	AriaNodeTypeComplementary AriaNodeType = "complementary"
	AriaNodeTypeContentInfo   AriaNodeType = "contentinfo"
	AriaNodeTypeForm          AriaNodeType = "form"
	AriaNodeTypeMain          AriaNodeType = "main"
	AriaNodeTypeNavigation    AriaNodeType = "navigation"
	AriaNodeTypeRegion        AriaNodeType = "region"
	AriaNodeTypeSearch        AriaNodeType = "search"
	
	// ARIA widget roles
	AriaNodeTypeArticle      AriaNodeType = "article"
	AriaNodeTypeButton       AriaNodeType = "button"
	AriaNodeTypeCell         AriaNodeType = "cell"
	AriaNodeTypeCheckbox     AriaNodeType = "checkbox"
	AriaNodeTypeColumnHeader AriaNodeType = "columnheader"
	AriaNodeTypeCombobox     AriaNodeType = "combobox"
	AriaNodeTypeDialog       AriaNodeType = "dialog"
	AriaNodeTypeFigure       AriaNodeType = "figure"
	AriaNodeTypeGrid         AriaNodeType = "grid"
	AriaNodeTypeGridCell     AriaNodeType = "gridcell"
	AriaNodeTypeHeading      AriaNodeType = "heading"
	AriaNodeTypeImg          AriaNodeType = "img"
	AriaNodeTypeLink         AriaNodeType = "link"
	AriaNodeTypeList         AriaNodeType = "list"
	AriaNodeTypeListItem     AriaNodeType = "listitem"
	AriaNodeTypeMenuItem     AriaNodeType = "menuitem"
	AriaNodeTypeOption       AriaNodeType = "option"
	AriaNodeTypeProgressBar  AriaNodeType = "progressbar"
	AriaNodeTypeRadio        AriaNodeType = "radio"
	AriaNodeTypeRadioGroup   AriaNodeType = "radiogroup"
	AriaNodeTypeRow          AriaNodeType = "row"
	AriaNodeTypeRowGroup     AriaNodeType = "rowgroup"
	AriaNodeTypeRowHeader    AriaNodeType = "rowheader"
	AriaNodeTypeSearchBox    AriaNodeType = "searchbox"
	AriaNodeTypeSeparator    AriaNodeType = "separator"
	AriaNodeTypeSlider       AriaNodeType = "slider"
	AriaNodeTypeSpinButton   AriaNodeType = "spinbutton"
	AriaNodeTypeSwitch       AriaNodeType = "switch"
	AriaNodeTypeTab          AriaNodeType = "tab"
	AriaNodeTypeTable        AriaNodeType = "table"
	AriaNodeTypeTabList      AriaNodeType = "tablist"
	AriaNodeTypeTabPanel     AriaNodeType = "tabpanel"
	AriaNodeTypeTextBox      AriaNodeType = "textbox"
	AriaNodeTypeText         AriaNodeType = "text"
	AriaNodeTypeGeneric      AriaNodeType = "generic" // Any other role
)

// AriaNode represents a node in an accessibility tree.
// It contains information about the accessibility properties of an element,
// such as its role, name, state, and children, which is useful for understanding
// the semantic structure of a document from an accessibility perspective.
type AriaNode struct {
	Type            AriaNodeType  // Type of the ARIA node
	Name            string        // Accessible name
	Role            string        // Explicit ARIA role
	Level           int           // Heading level, etc.
	Checked         *bool         // Checkbox state (pointer to allow nil for "not applicable")
	Selected        *bool         // Selection state
	Expanded        *bool         // Expansion state
	Disabled        *bool         // Disabled state
	Required        *bool         // Required state
	ValueMin        *float64      // Minimum value
	ValueMax        *float64      // Maximum value
	ValueText       string        // Text representation of value
	Children        []*AriaNode   // Child nodes
	OriginalElement *dom.VElement // Reference to the original DOM element
}

// AriaTree represents an accessibility tree.
// This is a hierarchical representation of a document's accessibility structure,
// which can be used as a fallback when traditional content extraction fails.
type AriaTree struct {
	Root      *AriaNode // Root node of the ARIA tree
	NodeCount int       // Total number of nodes in the tree
}

// GetAriaRole returns the ARIA role of an element.
// It returns the explicit role attribute or an implicit role based on the tag name.
// ARIA roles provide semantic meaning to elements for accessibility purposes.
//
// Parameters:
//   - element: The element to get the role for
//
// Returns:
//   - The ARIA role as a string
func GetAriaRole(element *dom.VElement) string {
	// Prioritize explicit role attribute
	if explicitRole := dom.GetAttribute(element, "role"); explicitRole != "" {
		return strings.ToLower(explicitRole)
	}

	// Implicit role based on tag name
	tagName := strings.ToLower(element.TagName)

	// Mapping of common HTML elements to implicit roles
	implicitRoles := map[string]string{
		"a":        "generic", // Default to generic, will be updated to "link" if href exists
		"article":  "article",
		"aside":    "complementary",
		"button":   "button",
		"footer":   "contentinfo",
		"form":     "form",
		"h1":       "heading",
		"h2":       "heading",
		"h3":       "heading",
		"h4":       "heading",
		"h5":       "heading",
		"h6":       "heading",
		"header":   "banner",
		"img":      "img",
		"li":       "listitem",
		"main":     "main",
		"nav":      "navigation",
		"ol":       "list",
		"option":   "option",
		"progress": "progressbar",
		"section":  "region",
		"select":   "combobox",
		"table":    "table",
		"textarea": "textbox",
		"ul":       "list",
	}

	// Special case for <a> with href
	if tagName == "a" && dom.GetAttribute(element, "href") != "" {
		return "link"
	}

	// Special case for <input> based on type
	if tagName == "input" {
		inputType := strings.ToLower(dom.GetAttribute(element, "type"))
		if inputType == "" {
			inputType = "text" // Default input type
		}

		switch inputType {
		case "checkbox":
			return "checkbox"
		case "radio":
			return "radio"
		case "button":
			return "button"
		case "search":
			return "searchbox"
		default:
			return "textbox"
		}
	}

	if role, ok := implicitRoles[tagName]; ok {
		return role
	}

	return "generic"
}

// GetAccessibleName returns the accessible name of an element.
// It follows the accessible name calculation algorithm, prioritizing aria-label,
// aria-labelledby, alt, title, and text content. The accessible name is what would
// be announced by screen readers and other assistive technologies.
//
// Parameters:
//   - element: The element to get the accessible name for
//
// Returns:
//   - The accessible name as a string
func GetAccessibleName(element *dom.VElement) string {
	// Prioritize aria-label attribute
	if ariaLabel := dom.GetAttribute(element, "aria-label"); ariaLabel != "" {
		return ariaLabel
	}

	// Alt attribute for images
	if element.TagName == "img" {
		if alt := dom.GetAttribute(element, "alt"); alt != "" {
			return alt
		}
	}

	// Title attribute
	if title := dom.GetAttribute(element, "title"); title != "" {
		return title
	}

	// Use text content for headings, links, buttons, etc.
	isNameFromContent := map[string]bool{
		"a":      true,
		"button": true,
		"h1":     true,
		"h2":     true,
		"h3":     true,
		"h4":     true,
		"h5":     true,
		"h6":     true,
		"label":  true,
	}

	if isNameFromContent[strings.ToLower(element.TagName)] {
		text := dom.GetInnerText(element, true)
		if text != "" {
			// Truncate if too long
			if len(text) > 50 {
				return text[:47] + "..."
			}
			return text
		}
	}

	// For paragraphs and divs with short text
	if element.TagName == "p" || element.TagName == "div" {
		text := dom.GetInnerText(element, true)
		if text != "" && len(text) < 100 {
			return text
		}
	}

	return ""
}

// GetAriaNodeType determines the AriaNodeType of an element based on its role.
// This maps ARIA roles to their corresponding AriaNodeType enum values.
//
// Parameters:
//   - element: The element to determine the node type for
//
// Returns:
//   - The AriaNodeType corresponding to the element's role
func GetAriaNodeType(element *dom.VElement) AriaNodeType {
	role := GetAriaRole(element)

	// Map role to AriaNodeType
	roleToType := map[string]AriaNodeType{
		"banner":        AriaNodeTypeBanner,
		"complementary": AriaNodeTypeComplementary,
		"contentinfo":   AriaNodeTypeContentInfo,
		"form":          AriaNodeTypeForm,
		"main":          AriaNodeTypeMain,
		"navigation":    AriaNodeTypeNavigation,
		"region":        AriaNodeTypeRegion,
		"search":        AriaNodeTypeSearch,
		"article":       AriaNodeTypeArticle,
		"button":        AriaNodeTypeButton,
		"cell":          AriaNodeTypeCell,
		"checkbox":      AriaNodeTypeCheckbox,
		"columnheader":  AriaNodeTypeColumnHeader,
		"combobox":      AriaNodeTypeCombobox,
		"dialog":        AriaNodeTypeDialog,
		"figure":        AriaNodeTypeFigure,
		"grid":          AriaNodeTypeGrid,
		"gridcell":      AriaNodeTypeGridCell,
		"heading":       AriaNodeTypeHeading,
		"img":           AriaNodeTypeImg,
		"link":          AriaNodeTypeLink,
		"list":          AriaNodeTypeList,
		"listitem":      AriaNodeTypeListItem,
		"menuitem":      AriaNodeTypeMenuItem,
		"option":        AriaNodeTypeOption,
		"progressbar":   AriaNodeTypeProgressBar,
		"radio":         AriaNodeTypeRadio,
		"radiogroup":    AriaNodeTypeRadioGroup,
		"row":           AriaNodeTypeRow,
		"rowgroup":      AriaNodeTypeRowGroup,
		"rowheader":     AriaNodeTypeRowHeader,
		"searchbox":     AriaNodeTypeSearchBox,
		"separator":     AriaNodeTypeSeparator,
		"slider":        AriaNodeTypeSlider,
		"spinbutton":    AriaNodeTypeSpinButton,
		"switch":        AriaNodeTypeSwitch,
		"tab":           AriaNodeTypeTab,
		"table":         AriaNodeTypeTable,
		"tablist":       AriaNodeTypeTabList,
		"tabpanel":      AriaNodeTypeTabPanel,
		"textbox":       AriaNodeTypeTextBox,
	}

	// If it's a generic role but has text children, treat it as text
	if role == "generic" {
		for _, child := range element.Children {
			if _, ok := dom.AsVText(child); ok {
				return AriaNodeTypeText
			}
		}
	}

	if nodeType, ok := roleToType[role]; ok {
		return nodeType
	}

	return AriaNodeTypeGeneric
}

// BuildAriaNode builds an AriaNode from a DOM element.
// This recursively constructs an accessibility tree node from a DOM element,
// including its properties and children.
//
// Parameters:
//   - element: The DOM element to build an AriaNode from
//
// Returns:
//   - An AriaNode representing the element and its children
func BuildAriaNode(element *dom.VElement) *AriaNode {
	nodeType := GetAriaNodeType(element)
	name := GetAccessibleName(element)
	role := GetAriaRole(element)

	// Create basic AriaNode
	node := &AriaNode{
		Type:            nodeType,
		Role:            role,
		OriginalElement: element,
	}

	// Add name if available
	if name != "" {
		node.Name = name
	}

	// Add heading level
	if nodeType == AriaNodeTypeHeading {
		if headingMatch := strings.ToLower(element.TagName); len(headingMatch) == 2 && headingMatch[0] == 'h' {
			if level, err := strconv.Atoi(string(headingMatch[1])); err == nil && level >= 1 && level <= 6 {
				node.Level = level
			}
		}
	}

	// Checkbox or radio state
	if nodeType == AriaNodeTypeCheckbox || nodeType == AriaNodeTypeRadio {
		checked := false
		if _, exists := element.Attributes["checked"]; exists {
			checked = true
		} else if dom.GetAttribute(element, "aria-checked") == "true" {
			checked = true
		}
		node.Checked = &checked
	}

	// Selected state for options and tabs
	if nodeType == AriaNodeTypeOption || nodeType == AriaNodeTypeTab {
		selected := false
		if _, exists := element.Attributes["selected"]; exists {
			selected = true
		} else if dom.GetAttribute(element, "aria-selected") == "true" {
			selected = true
		}
		node.Selected = &selected
	}

	// Expanded state
	if ariaExpanded := dom.GetAttribute(element, "aria-expanded"); ariaExpanded != "" {
		expanded := ariaExpanded == "true"
		node.Expanded = &expanded
	}

	// Disabled state
	if _, exists := element.Attributes["disabled"]; exists || dom.GetAttribute(element, "aria-disabled") == "true" {
		disabled := true
		node.Disabled = &disabled
	}

	// Required state
	if _, exists := element.Attributes["required"]; exists || dom.GetAttribute(element, "aria-required") == "true" {
		required := true
		node.Required = &required
	}

	// Value range (for sliders, etc.)
	if valueMin := dom.GetAttribute(element, "aria-valuemin"); valueMin != "" {
		if min, err := strconv.ParseFloat(valueMin, 64); err == nil {
			node.ValueMin = &min
		}
	} else if min := dom.GetAttribute(element, "min"); min != "" {
		if minVal, err := strconv.ParseFloat(min, 64); err == nil {
			node.ValueMin = &minVal
		}
	}

	if valueMax := dom.GetAttribute(element, "aria-valuemax"); valueMax != "" {
		if max, err := strconv.ParseFloat(valueMax, 64); err == nil {
			node.ValueMax = &max
		}
	} else if max := dom.GetAttribute(element, "max"); max != "" {
		if maxVal, err := strconv.ParseFloat(max, 64); err == nil {
			node.ValueMax = &maxVal
		}
	}

	if valueText := dom.GetAttribute(element, "aria-valuetext"); valueText != "" {
		node.ValueText = valueText
	} else if value := dom.GetAttribute(element, "value"); value != "" {
		node.ValueText = value
	}

	// Build child nodes recursively
	var childNodes []*AriaNode

	for _, child := range element.Children {
		childElement, ok := dom.AsVElement(child)
		if !ok {
			continue
		}

		// Skip invisible elements
		if !dom.IsProbablyVisible(childElement) {
			continue
		}

		childNode := BuildAriaNode(childElement)
		
		// Only add meaningful child nodes
		if childNode.Name != "" || childNode.Type != AriaNodeTypeGeneric || (childNode.Children != nil && len(childNode.Children) > 0) {
			childNodes = append(childNodes, childNode)
		}
	}

	// Add children if any
	if len(childNodes) > 0 {
		node.Children = childNodes
	}

	return node
}

// isInsignificantNode determines if a node is insignificant.
// Insignificant nodes are those that don't contribute meaningful information
// to the accessibility tree and can be pruned during tree compression.
//
// Parameters:
//   - node: The node to check
//
// Returns:
//   - true if the node is insignificant, false otherwise
func isInsignificantNode(node *AriaNode) bool {
	return node.Name == "" && node.Type == AriaNodeTypeGeneric && (node.Children == nil || len(node.Children) == 0)
}

// CountAriaNodes counts the total number of nodes in an AriaNode tree.
// This includes the node itself and all its descendants.
//
// Parameters:
//   - node: The root node to count from
//
// Returns:
//   - The total number of nodes in the tree
func CountAriaNodes(node *AriaNode) int {
	if node == nil {
		return 0
	}

	count := 1 // Count the node itself
	if node.Children != nil {
		for _, child := range node.Children {
			count += CountAriaNodes(child)
		}
	}
	return count
}

// CompressAriaTree compresses an AriaTree by removing insignificant nodes,
// merging similar nodes, and simplifying the structure. This produces a more
// concise and meaningful representation of the document's accessibility structure.
//
// Parameters:
//   - node: The root node of the tree to compress
//
// Returns:
//   - The compressed tree's root node
func CompressAriaTree(node *AriaNode) *AriaNode {
	if node == nil {
		return nil
	}

	// If no children, return as is (with possible text content check)
	if node.Children == nil || len(node.Children) == 0 {
		// Remove empty text nodes
		if node.Type == AriaNodeTypeText && (node.Name == "" || strings.TrimSpace(node.Name) == "") {
			return &AriaNode{
				Type:            AriaNodeTypeGeneric,
				Role:            "generic",
				OriginalElement: node.OriginalElement,
			}
		}
		return node
	}

	// First, recursively compress all children
	var processedChildren []*AriaNode
	for _, child := range node.Children {
		compressed := CompressAriaTree(child)
		if compressed != nil && !isInsignificantNode(compressed) {
			// Filter out empty text nodes
			if !(compressed.Type == AriaNodeTypeText && (compressed.Name == "" || strings.TrimSpace(compressed.Name) == "")) {
				processedChildren = append(processedChildren, compressed)
			}
		}
	}

	// Special case: text node with one significant child
	if node.Type == AriaNodeTypeText && len(processedChildren) == 1 {
		significantChild := processedChildren[0]
		significantTypes := map[AriaNodeType]bool{
			AriaNodeTypeMain:          true,
			AriaNodeTypeArticle:       true,
			AriaNodeTypeRegion:        true,
			AriaNodeTypeNavigation:    true,
			AriaNodeTypeBanner:        true,
			AriaNodeTypeContentInfo:   true,
		}

		if significantTypes[significantChild.Type] {
			// Merge parent name to child if needed
			if node.Name != "" && significantChild.Name == "" {
				significantChild.Name = node.Name
			}
			return significantChild
		}
	}

	// If text node with only generic children, merge them
	if node.Type == AriaNodeTypeText && len(processedChildren) > 0 {
		allGeneric := true
		for _, child := range processedChildren {
			if child.Type != AriaNodeTypeGeneric {
				allGeneric = false
				break
			}
		}

		if allGeneric {
			var newChildren []*AriaNode
			for _, child := range processedChildren {
				if child.Children != nil && len(child.Children) > 0 {
					newChildren = append(newChildren, child.Children...)
				}
			}
			if len(newChildren) > 0 {
				result := *node // Create a copy
				result.Children = newChildren
				return &result
			}
		}
	}

	// General case: if only one child, consider merging
	if len(processedChildren) == 1 {
		child := processedChildren[0]

		// If parent is generic with no name, or parent and child are same type
		if (node.Type == AriaNodeTypeGeneric && node.Name == "") || node.Type == child.Type {
			// Merge names if needed
			if node.Name != "" {
				if child.Name == "" {
					child.Name = node.Name
				} else {
					child.Name = node.Name + " " + child.Name
				}
			}
			return child
		}
	}

	// Check if this is a significant structural node
	isSignificantNode := map[AriaNodeType]bool{
		AriaNodeTypeMain:          true,
		AriaNodeTypeArticle:       true,
		AriaNodeTypeRegion:        true,
		AriaNodeTypeNavigation:    true,
		AriaNodeTypeBanner:        true,
		AriaNodeTypeContentInfo:   true,
		AriaNodeTypeForm:          true,
		AriaNodeTypeSearch:        true,
	}[node.Type]

	// Handle generic children under significant nodes
	if len(processedChildren) > 0 {
		hasGenericChildren := false
		for _, child := range processedChildren {
			if child.Type == AriaNodeTypeGeneric {
				hasGenericChildren = true
				break
			}
		}

		if hasGenericChildren && (isSignificantNode || func() bool {
			for _, child := range processedChildren {
				if child.Type != AriaNodeTypeGeneric {
					return false
				}
			}
			return true
		}()) {
			var newChildren []*AriaNode
			for _, child := range processedChildren {
				if child.Type == AriaNodeTypeGeneric {
					if child.Children != nil && len(child.Children) > 0 {
						newChildren = append(newChildren, child.Children...)
					}
				} else {
					newChildren = append(newChildren, child)
				}
			}

			if len(newChildren) > 0 {
				result := *node // Create a copy
				result.Children = newChildren
				return &result
			}
		}
	}

	// Group similar nodes
	var mergedChildren []*AriaNode
	var currentGroup *AriaNode
	groupByType := make(map[AriaNodeType][]*AriaNode)

	// Group specific types of nodes
	for _, child := range processedChildren {
		if child.Type == AriaNodeTypeArticle || child.Type == AriaNodeTypeRegion || 
		   child.Type == AriaNodeTypeListItem || child.Type == AriaNodeTypeImg {
			groupByType[child.Type] = append(groupByType[child.Type], child)
			continue
		}

		// Start a new group if needed
		if currentGroup == nil || currentGroup.Type != child.Type {
			// Deep copy the child to create a new group
			currentGroup = &AriaNode{}
			*currentGroup = *child
			mergedChildren = append(mergedChildren, currentGroup)
			continue
		}

		// Merge with current group if same type
		if child.Name != "" {
			if currentGroup.Name != "" {
				currentGroup.Name += " " + child.Name
			} else {
				currentGroup.Name = child.Name
			}
		}

		// Merge children
		if child.Children != nil && len(child.Children) > 0 {
			if currentGroup.Children == nil {
				currentGroup.Children = make([]*AriaNode, 0, len(child.Children))
			}
			currentGroup.Children = append(currentGroup.Children, child.Children...)
		}
	}

	// Add grouped nodes
	for nodeType, nodes := range groupByType {
		if len(nodes) > 1 {
			// Create a parent node for grouped nodes
			parentNode := &AriaNode{
				Type:            nodeType,
				Role:            string(nodeType),
				OriginalElement: node.OriginalElement,
				Children:        nodes,
			}
			mergedChildren = append(mergedChildren, parentNode)
		} else if len(nodes) == 1 {
			mergedChildren = append(mergedChildren, nodes[0])
		}
	}

	// Flatten nested structures
	for i := 0; i < len(mergedChildren); i++ {
		child := mergedChildren[i]

		// Flatten single-child nodes
		if child.Children != nil && len(child.Children) == 1 {
			grandchild := child.Children[0]

			// Merge if same type or special case
			if child.Type == grandchild.Type || 
			   (child.Type == AriaNodeTypeText && (grandchild.Type == AriaNodeTypeMain || 
												  grandchild.Type == AriaNodeTypeArticle || 
												  grandchild.Type == AriaNodeTypeRegion)) {
				// Merge names
				if grandchild.Name != "" {
					if child.Name != "" {
						child.Name += " " + grandchild.Name
					} else {
						child.Name = grandchild.Name
					}
				}

				// Move grandchild's children to child
				if grandchild.Children != nil && len(grandchild.Children) > 0 {
					child.Children = grandchild.Children
					i-- // Process this node again
					continue
				} else {
					child.Children = nil
				}
			}
		}

		// Handle multiple children with same type as parent
		if child.Children != nil && len(child.Children) > 1 {
			var sameTypeChildren []*AriaNode
			var otherChildren []*AriaNode

			for _, c := range child.Children {
				if c.Type == child.Type {
					sameTypeChildren = append(sameTypeChildren, c)
				} else {
					otherChildren = append(otherChildren, c)
				}
			}

			if len(sameTypeChildren) > 0 {
				var newChildren []*AriaNode

				// Merge names from same-type children
				for _, sameChild := range sameTypeChildren {
					if sameChild.Name != "" {
						if child.Name != "" {
							child.Name += " " + sameChild.Name
						} else {
							child.Name = sameChild.Name
						}
					}

					// Add same-type children's children
					if sameChild.Children != nil && len(sameChild.Children) > 0 {
						newChildren = append(newChildren, sameChild.Children...)
					}
				}

				// Add other children
				newChildren = append(newChildren, otherChildren...)

				// Update children
				child.Children = newChildren
				i-- // Process this node again
				continue
			}
		}
	}

	// Create result with compressed children
	result := *node // Create a copy
	if len(mergedChildren) > 0 {
		result.Children = mergedChildren
	} else {
		result.Children = nil
	}

	return &result
}

// BuildAriaTree builds an AriaTree from a DOM document.
// This constructs a complete accessibility tree from a document, then compresses
// it to produce a more concise and meaningful representation.
//
// Parameters:
//   - doc: The DOM document to build an AriaTree from
//
// Returns:
//   - An AriaTree representing the document's accessibility structure
func BuildAriaTree(doc *dom.VDocument) *AriaTree {
	// Build tree from document body
	rootNode := BuildAriaNode(doc.Body)

	// Compress the tree
	compressedRoot := CompressAriaTree(rootNode)

	// Handle special case for root level nesting
	if compressedRoot.Type == AriaNodeTypeText && compressedRoot.Children != nil && len(compressedRoot.Children) > 0 {
		// Look for significant child nodes
		var significantChild *AriaNode
		for _, child := range compressedRoot.Children {
			if child.Type == AriaNodeTypeMain || child.Type == AriaNodeTypeArticle || 
			   child.Type == AriaNodeTypeRegion || child.Type == AriaNodeTypeNavigation || 
			   child.Type == AriaNodeTypeBanner || child.Type == AriaNodeTypeContentInfo {
				significantChild = child
				break
			}
		}

		// If found, make it the root
		if significantChild != nil {
			// Merge names if needed
			if compressedRoot.Name != "" && significantChild.Name == "" {
				significantChild.Name = compressedRoot.Name
			}
			compressedRoot = significantChild
		} else if len(compressedRoot.Children) == 1 {
			// If only one child, merge it with root
			child := compressedRoot.Children[0]

			// Merge names
			if child.Name != "" {
				if compressedRoot.Name != "" {
					compressedRoot.Name += " " + child.Name
				} else {
					compressedRoot.Name = child.Name
				}
			}

			// Move child's children to root
			if child.Children != nil && len(child.Children) > 0 {
				compressedRoot.Children = child.Children
			} else {
				compressedRoot.Children = nil
			}
		}
	}

	// Count nodes
	nodeCount := CountAriaNodes(compressedRoot)

	return &AriaTree{
		Root:      compressedRoot,
		NodeCount: nodeCount,
	}
}

// AriaTreeToString converts an AriaTree to a YAML-like string format
// compatible with Playwright's snapshot format. This is useful for debugging
// and visualizing the accessibility structure of a document.
//
// Parameters:
//   - tree: The AriaTree to convert to a string
//
// Returns:
//   - A YAML-like string representation of the tree
func AriaTreeToString(tree *AriaTree) string {
	if tree == nil || tree.Root == nil {
		return ""
	}
	
	var sb strings.Builder
	nodeToString(tree.Root, 0, &sb)
	return sb.String()
}

// nodeToString recursively converts an AriaNode to a string with proper indentation.
// This is a helper function for AriaTreeToString that handles the formatting of
// individual nodes in the tree.
//
// Parameters:
//   - node: The node to convert to a string
//   - indent: The current indentation level
//   - sb: A string builder to append the result to
func nodeToString(node *AriaNode, indent int, sb *strings.Builder) {
	if node == nil {
		return
	}

	// Skip empty nodes
	if (node.Name == "" && node.Children == nil && node.Type == AriaNodeTypeGeneric) ||
	   node.Name == "" ||
	   (node.Type == AriaNodeTypeList && (node.Children == nil || len(node.Children) == 0 || func() bool {
			for _, child := range node.Children {
				if child.Name != "" || (child.Children != nil && len(child.Children) > 0) {
					return false
				}
			}
			return true
		}())) {
		return
	}

	// Indent
	indentStr := strings.Repeat("  ", indent)
	sb.WriteString(fmt.Sprintf("%s- %s", indentStr, node.Type))

	// Get href for links, src and alt for images
	var href, src, alt string
	if node.OriginalElement != nil {
		if node.Type == AriaNodeTypeLink {
			href = node.OriginalElement.Attributes["href"]
		} else if node.Type == AriaNodeTypeImg {
			src = node.OriginalElement.Attributes["src"]
			alt = node.OriginalElement.Attributes["alt"]
		}
	}

	// Add name based on node type
	if node.Name != "" {
		// Check if name is a regex pattern
		if strings.HasPrefix(node.Name, "/") && strings.HasSuffix(node.Name, "/") {
			sb.WriteString(fmt.Sprintf(" %s", node.Name))
		} else if node.Type == AriaNodeTypeText || node.Type == AriaNodeTypeListItem || node.Type == AriaNodeTypeTextBox {
			// Colon format for text nodes
			sb.WriteString(fmt.Sprintf(": %s", node.Name))
		} else {
			// Quote format for other nodes
			sb.WriteString(fmt.Sprintf(" \"%s\"", node.Name))
		}
	}

	// Add href, src, alt attributes
	if href != "" {
		sb.WriteString(fmt.Sprintf(" [href=\"%s\"]", href))
	}
	if src != "" {
		sb.WriteString(fmt.Sprintf(" [src=\"%s\"]", src))
	}
	if alt != "" {
		sb.WriteString(fmt.Sprintf(" [alt=\"%s\"]", alt))
	}

	// Add other attributes
	var attributes []string

	// Heading level
	if node.Level > 0 {
		attributes = append(attributes, fmt.Sprintf("level=%d", node.Level))
	}

	// Checkbox state
	if node.Checked != nil {
		attributes = append(attributes, fmt.Sprintf("checked=%t", *node.Checked))
	}

	// Selection state
	if node.Selected != nil {
		attributes = append(attributes, fmt.Sprintf("selected=%t", *node.Selected))
	}

	// Expansion state
	if node.Expanded != nil {
		attributes = append(attributes, fmt.Sprintf("expanded=%t", *node.Expanded))
	}

	// Disabled state
	if node.Disabled != nil && *node.Disabled {
		attributes = append(attributes, "disabled")
	}

	// Required state
	if node.Required != nil && *node.Required {
		attributes = append(attributes, "required")
	}

	// Value range
	if node.ValueMin != nil || node.ValueMax != nil {
		min := ""
		max := ""
		if node.ValueMin != nil {
			min = fmt.Sprintf("%g", *node.ValueMin)
		}
		if node.ValueMax != nil {
			max = fmt.Sprintf("%g", *node.ValueMax)
		}
		attributes = append(attributes, fmt.Sprintf("range=%s-%s", min, max))
	}

	// Value text
	if node.ValueText != "" {
		attributes = append(attributes, fmt.Sprintf("value=\"%s\"", node.ValueText))
	}

	// Add attributes if any
	if len(attributes) > 0 {
		sb.WriteString(fmt.Sprintf(" [%s]", strings.Join(attributes, ", ")))
	}

	sb.WriteString("\n")

	// Process children
	if node.Children != nil && len(node.Children) > 0 {
		// Special case: if only one child with same type
		if len(node.Children) == 1 && node.Children[0].Type == node.Type && node.Type != AriaNodeTypeGeneric {
			child := node.Children[0]
			
			// Combine names
			combinedName := ""
			if node.Name != "" && child.Name != "" {
				combinedName = node.Name + " " + child.Name
			} else {
				combinedName = node.Name + child.Name
			}
			
			// Create temporary node with combined properties
			tempNode := *node // Copy the node
			tempNode.Name = combinedName
			
			// Set children from the child's children
			if child.Children != nil && len(child.Children) > 0 {
				tempNode.Children = child.Children
			} else {
				tempNode.Children = nil
			}
			
			// Process this node instead
			nodeToString(&tempNode, indent, sb)
			return
		}
		
		// Normal child processing
		for _, child := range node.Children {
			nodeToString(child, indent+1, sb)
		}
	}
}
