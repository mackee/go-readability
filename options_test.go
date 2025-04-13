package readability_test

import (
	"testing"

	"github.com/mackee/go-readability"
	"github.com/mackee/go-readability/internal/dom"
)

func TestDefaultOptions(t *testing.T) {
	opts := readability.DefaultOptions()

	// Check default values
	if opts.CharThreshold != 500 {
		t.Errorf("Expected CharThreshold to be %d, got %d", 500, opts.CharThreshold)
	}

	if opts.NbTopCandidates != 5 {
		t.Errorf("Expected NbTopCandidates to be %d, got %d", 5, opts.NbTopCandidates)
	}

	if opts.GenerateAriaTree != false {
		t.Errorf("Expected GenerateAriaTree to be %v, got %v", false, opts.GenerateAriaTree)
	}

	if opts.ForcedPageType != "" {
		t.Errorf("Expected ForcedPageType to be empty, got %v", opts.ForcedPageType)
	}
}

func TestReadabilityArticleGetContentByPageType(t *testing.T) {
	// Create test elements
	div := dom.NewVElement("div")
	div.SetAttribute("id", "content")
	div.AppendChild(dom.NewVText("Article content"))

	header := dom.NewVElement("header")
	header.AppendChild(dom.NewVText("Header content"))

	footer := dom.NewVElement("footer")
	footer.AppendChild(dom.NewVText("Footer content"))

	// Create test significant nodes
	aside := dom.NewVElement("aside")
	aside.AppendChild(dom.NewVText("Sidebar content"))

	// Test case 1: Article page type
	articlePage := &readability.ReadabilityArticle{
		Title:                 "Test Article",
		Byline:                "John Doe",
		Root:                  div,
		NodeCount:             5,
		PageType:              readability.PageTypeArticle,
		Header:                header,
		Footer:                footer,
		OtherSignificantNodes: []*dom.VElement{aside},
	}

	content := articlePage.GetContentByPageType()
	articleContent, ok := content.(readability.ArticleContent)
	if !ok {
		t.Fatalf("Expected ArticleContent type, got %T", content)
	}

	if articleContent.Title != "Test Article" {
		t.Errorf("Expected title to be %q, got %q", "Test Article", articleContent.Title)
	}

	if articleContent.Byline != "John Doe" {
		t.Errorf("Expected byline to be %q, got %q", "John Doe", articleContent.Byline)
	}

	if articleContent.Root != div {
		t.Errorf("Expected root to be %v, got %v", div, articleContent.Root)
	}

	// Test case 2: Other page type
	otherPage := &readability.ReadabilityArticle{
		Title:                 "Test Page",
		Byline:                "",
		Root:                  nil,
		NodeCount:             5,
		PageType:              readability.PageTypeOther,
		Header:                header,
		Footer:                footer,
		OtherSignificantNodes: []*dom.VElement{aside},
		AriaTree: &readability.AriaTree{
			Root:      &readability.AriaNode{Type: readability.AriaNodeTypeGeneric},
			NodeCount: 1,
		},
	}

	content = otherPage.GetContentByPageType()
	otherContent, ok := content.(readability.OtherContent)
	if !ok {
		t.Fatalf("Expected OtherContent type, got %T", content)
	}

	if otherContent.Title != "Test Page" {
		t.Errorf("Expected title to be %q, got %q", "Test Page", otherContent.Title)
	}

	if otherContent.Header != header {
		t.Errorf("Expected header to be %v, got %v", header, otherContent.Header)
	}

	if otherContent.Footer != footer {
		t.Errorf("Expected footer to be %v, got %v", footer, otherContent.Footer)
	}

	if len(otherContent.OtherSignificantNodes) != 1 || otherContent.OtherSignificantNodes[0] != aside {
		t.Errorf("Expected OtherSignificantNodes to contain aside element")
	}

	if otherContent.AriaTree == nil || otherContent.AriaTree.Root.Type != readability.AriaNodeTypeGeneric {
		t.Errorf("Expected AriaTree with Generic root node")
	}
}
