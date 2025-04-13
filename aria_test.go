package readability

import (
	"strings"
	"testing"

	"github.com/mackee/go-readability/internal/dom"
)

func TestGetAriaRole(t *testing.T) {
	tests := []struct {
		name     string
		element  *dom.VElement
		expected string
	}{
		{
			name: "explicit role",
			element: &dom.VElement{
				TagName: "div",
				Attributes: map[string]string{
					"role": "button",
				},
			},
			expected: "button",
		},
		{
			name: "implicit role for a with href",
			element: &dom.VElement{
				TagName: "a",
				Attributes: map[string]string{
					"href": "https://example.com",
				},
			},
			expected: "link",
		},
		{
			name: "implicit role for a without href",
			element: &dom.VElement{
				TagName: "a",
			},
			expected: "generic",
		},
		{
			name: "implicit role for heading",
			element: &dom.VElement{
				TagName: "h1",
			},
			expected: "heading",
		},
		{
			name: "implicit role for input checkbox",
			element: &dom.VElement{
				TagName: "input",
				Attributes: map[string]string{
					"type": "checkbox",
				},
			},
			expected: "checkbox",
		},
		{
			name: "implicit role for input text",
			element: &dom.VElement{
				TagName: "input",
			},
			expected: "textbox",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAriaRole(tt.element)
			if result != tt.expected {
				t.Errorf("GetAriaRole() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetAccessibleName(t *testing.T) {
	tests := []struct {
		name     string
		element  *dom.VElement
		expected string
	}{
		{
			name: "aria-label",
			element: &dom.VElement{
				TagName: "div",
				Attributes: map[string]string{
					"aria-label": "Test Label",
				},
			},
			expected: "Test Label",
		},
		{
			name: "alt for img",
			element: &dom.VElement{
				TagName: "img",
				Attributes: map[string]string{
					"alt": "Image Description",
				},
			},
			expected: "Image Description",
		},
		{
			name: "title",
			element: &dom.VElement{
				TagName: "div",
				Attributes: map[string]string{
					"title": "Title Text",
				},
			},
			expected: "Title Text",
		},
		{
			name: "text content for heading",
			element: &dom.VElement{
				TagName: "h1",
				Children: []dom.VNode{
					dom.NewVText("Heading Text"),
				},
			},
			expected: "Heading Text",
		},
		{
			name: "text content for paragraph",
			element: &dom.VElement{
				TagName: "p",
				Children: []dom.VNode{
					dom.NewVText("Paragraph Text"),
				},
			},
			expected: "Paragraph Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAccessibleName(tt.element)
			if result != tt.expected {
				t.Errorf("GetAccessibleName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBuildAriaNode(t *testing.T) {
	// Create a test element
	element := &dom.VElement{
		TagName: "h1",
		Attributes: map[string]string{
			"class": "title",
		},
		Children: []dom.VNode{
			dom.NewVText("Test Heading"),
		},
	}

	// Build AriaNode
	node := BuildAriaNode(element)

	// Verify basic properties
	if node.Type != AriaNodeTypeHeading {
		t.Errorf("Expected node type to be heading, got %v", node.Type)
	}

	if node.Name != "Test Heading" {
		t.Errorf("Expected node name to be 'Test Heading', got %v", node.Name)
	}

	if node.Level != 1 {
		t.Errorf("Expected heading level to be 1, got %v", node.Level)
	}

	if node.OriginalElement != element {
		t.Errorf("Expected original element to be preserved")
	}
}

func TestCountAriaNodes(t *testing.T) {
	// Create a simple tree
	root := &AriaNode{
		Type: AriaNodeTypeMain,
		Name: "Main Content",
		Children: []*AriaNode{
			{
				Type: AriaNodeTypeHeading,
				Name: "Title",
			},
			{
				Type: AriaNodeTypeText,
				Name: "Paragraph",
				Children: []*AriaNode{
					{
						Type: AriaNodeTypeLink,
						Name: "Link",
					},
				},
			},
		},
	}

	// Count nodes
	count := CountAriaNodes(root)

	// Expected: root + 2 children + 1 grandchild = 4
	if count != 4 {
		t.Errorf("Expected node count to be 4, got %d", count)
	}
}

func TestAriaTreeToString(t *testing.T) {
	// Create a simple tree
	tree := &AriaTree{
		Root: &AriaNode{
			Type: AriaNodeTypeMain,
			Name: "Main Content",
			Children: []*AriaNode{
				{
					Type:  AriaNodeTypeHeading,
					Name:  "Title",
					Level: 1,
				},
				{
					Type: AriaNodeTypeText,
					Name: "Paragraph text",
				},
			},
		},
		NodeCount: 3,
	}

	// Convert to string
	result := AriaTreeToString(tree)

	// Check that the output contains expected elements
	expectedSubstrings := []string{
		"- main",
		"\"Main Content\"",
		"- heading",
		"\"Title\"",
		"[level=1]",
		"- text",
		": Paragraph text",
	}

	for _, substr := range expectedSubstrings {
		if !containsSubstring(result, substr) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nOutput: %s", substr, result)
		}
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
}
