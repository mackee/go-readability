package readability

import (
	"strconv"
	"testing"

	"github.com/mackee/go-readability/internal/dom"
)

func TestFindMainCandidates(t *testing.T) {
	// Test cases for FindMainCandidates
	testCases := []struct {
		name           string
		setupDoc       func() *dom.VDocument
		nbTopCandidates int
		expectedCount  int
		checkFirst     func(t *testing.T, element *dom.VElement)
	}{
		{
			name: "single article tag",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				article := dom.NewVElement("article")
				article.AppendChild(dom.NewVText("This is an article with enough text to be considered."))
				body.AppendChild(article)
				
				return dom.NewVDocument(html, body)
			},
			nbTopCandidates: 5,
			expectedCount:   1,
			checkFirst: func(t *testing.T, element *dom.VElement) {
				if element.TagName != "article" {
					t.Errorf("Expected first candidate to be 'article', got '%s'", element.TagName)
				}
			},
		},
		{
			name: "single main tag",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				main := dom.NewVElement("main")
				main.AppendChild(dom.NewVText("This is a main section with enough text to be considered."))
				body.AppendChild(main)
				
				return dom.NewVDocument(html, body)
			},
			nbTopCandidates: 5,
			expectedCount:   1,
			checkFirst: func(t *testing.T, element *dom.VElement) {
				if element.TagName != "main" {
					t.Errorf("Expected first candidate to be 'main', got '%s'", element.TagName)
				}
			},
		},
		{
			name: "multiple candidates",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Create a div with content class (high score)
				contentDiv := dom.NewVElement("div")
				contentDiv.SetAttribute("class", "content")
				p1 := dom.NewVElement("p")
				p1.AppendChild(dom.NewVText("This is a paragraph with enough text to be considered. It has commas, and more text."))
				contentDiv.AppendChild(p1)
				body.AppendChild(contentDiv)
				
				// Create a div with sidebar class (low score)
				sidebarDiv := dom.NewVElement("div")
				sidebarDiv.SetAttribute("class", "sidebar")
				p2 := dom.NewVElement("p")
				p2.AppendChild(dom.NewVText("This is another paragraph with enough text to be considered."))
				sidebarDiv.AppendChild(p2)
				body.AppendChild(sidebarDiv)
				
				return dom.NewVDocument(html, body)
			},
			nbTopCandidates: 2,
			expectedCount:   2,
			checkFirst: func(t *testing.T, element *dom.VElement) {
				if element.ClassName() != "content" {
					t.Errorf("Expected first candidate to have class 'content', got '%s'", element.ClassName())
				}
			},
		},
		{
			name: "no candidates returns body",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Add some elements but with very little text
				div := dom.NewVElement("div")
				div.AppendChild(dom.NewVText("Short text."))
				body.AppendChild(div)
				
				return dom.NewVDocument(html, body)
			},
			nbTopCandidates: 5,
			expectedCount:   1,
			checkFirst: func(t *testing.T, element *dom.VElement) {
				if element.TagName != "body" {
					t.Errorf("Expected first candidate to be 'body', got '%s'", element.TagName)
				}
			},
		},
		{
			name: "limit by nbTopCandidates",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Create multiple divs with content
				for i := 0; i < 5; i++ {
					div := dom.NewVElement("div")
					p := dom.NewVElement("p")
					p.AppendChild(dom.NewVText("This is a paragraph with enough text to be considered. It has commas, and more text."))
					div.AppendChild(p)
					body.AppendChild(div)
				}
				
				return dom.NewVDocument(html, body)
			},
			nbTopCandidates: 3,
			expectedCount:   3,
			checkFirst: nil, // No specific check for first element
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := tc.setupDoc()
			candidates := FindMainCandidates(doc, tc.nbTopCandidates)
			
			if len(candidates) != tc.expectedCount {
				t.Errorf("Expected %d candidates, got %d", tc.expectedCount, len(candidates))
			}
			
			if len(candidates) > 0 && tc.checkFirst != nil {
				tc.checkFirst(t, candidates[0])
			}
		})
	}
}

func TestIsProbablyContent(t *testing.T) {
	// Test cases for IsProbablyContent
	testCases := []struct {
		name           string
		setupElement   func() *dom.VElement
		expected       bool
	}{
		{
			name: "visible element with good content",
			setupElement: func() *dom.VElement {
				div := dom.NewVElement("div")
				div.SetAttribute("class", "content")
				
				// Add enough text to pass the length check
				longText := "This is a long text that should be considered as content. " +
					"It has multiple sentences and is definitely longer than 140 characters. " +
					"This should be enough to pass the text length check in the IsProbablyContent function. " +
					"We need to make sure it's long enough."
				div.AppendChild(dom.NewVText(longText))
				
				return div
			},
			expected: true,
		},
		{
			name: "invisible element",
			setupElement: func() *dom.VElement {
				div := dom.NewVElement("div")
				div.SetAttribute("style", "display: none;")
				
				// Add enough text to pass the length check
				longText := "This is a long text that should be considered as content. " +
					"It has multiple sentences and is definitely longer than 140 characters."
				div.AppendChild(dom.NewVText(longText))
				
				return div
			},
			expected: false,
		},
		{
			name: "element with unlikely class",
			setupElement: func() *dom.VElement {
				div := dom.NewVElement("div")
				div.SetAttribute("class", "sidebar")
				
				// Add enough text to pass the length check
				longText := "This is a long text that should be considered as content. " +
					"It has multiple sentences and is definitely longer than 140 characters."
				div.AppendChild(dom.NewVText(longText))
				
				return div
			},
			expected: false,
		},
		{
			name: "element with short text",
			setupElement: func() *dom.VElement {
				div := dom.NewVElement("div")
				div.SetAttribute("class", "content")
				
				// Add short text that won't pass the length check
				shortText := "This is a short text."
				div.AppendChild(dom.NewVText(shortText))
				
				return div
			},
			expected: false,
		},
		{
			name: "element with high link density",
			setupElement: func() *dom.VElement {
				div := dom.NewVElement("div")
				
				// Add a paragraph with some text
				p := dom.NewVElement("p")
				p.AppendChild(dom.NewVText("This is some text. "))
				div.AppendChild(p)
				
				// Add many links to increase link density
				for i := 0; i < 10; i++ {
					a := dom.NewVElement("a")
					a.AppendChild(dom.NewVText("Link text that is quite long to increase the link density."))
					div.AppendChild(a)
				}
				
				return div
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			element := tc.setupElement()
			result := IsProbablyContent(element)
			
			if result != tc.expected {
				t.Errorf("Expected IsProbablyContent to return %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestInitializeNode(t *testing.T) {
	testCases := []struct {
		name          string
		tagName       string
		className     string
		id            string
		expectedScore float64
	}{
		{
			name:          "div element with no class or id",
			tagName:       "div",
			expectedScore: 5, // div gets +5 base score
		},
		{
			name:          "pre element with no class or id",
			tagName:       "pre",
			expectedScore: 3, // pre gets +3 base score
		},
		{
			name:          "h1 element with no class or id",
			tagName:       "h1",
			expectedScore: -5, // h1 gets -5 base score
		},
		{
			name:          "div with positive class",
			tagName:       "div",
			className:     "article content",
			expectedScore: 30, // div(+5) + positive class(+25)
		},
		{
			name:          "div with negative class",
			tagName:       "div",
			className:     "comment sidebar",
			expectedScore: -20, // div(+5) + negative class(-25)
		},
		{
			name:          "div with positive id",
			tagName:       "div",
			id:            "main-content",
			expectedScore: 30, // div(+5) + positive id(+25)
		},
		{
			name:          "div with negative id",
			tagName:       "div",
			id:            "sidebar",
			expectedScore: -20, // div(+5) + negative id(-25)
		},
		{
			name:          "div with both positive class and negative id",
			tagName:       "div",
			className:     "article",
			id:            "sidebar",
			expectedScore: 5, // div(+5) + positive class(+25) + negative id(-25)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test element
			element := dom.NewVElement(tc.tagName)
			if tc.className != "" {
				element.SetAttribute("class", tc.className)
			}
			if tc.id != "" {
				element.SetAttribute("id", tc.id)
			}

			// Initialize the node
			InitializeNode(element)

			// Check if readability data was set
			if element.GetReadabilityData() == nil {
				t.Fatal("ReadabilityData was not set")
			}

			// Check if the score matches the expected value
			if element.GetReadabilityData().ContentScore != tc.expectedScore {
				t.Errorf("Expected score %f, got %f", tc.expectedScore, element.GetReadabilityData().ContentScore)
			}
		})
	}
}

func TestGetClassWeight(t *testing.T) {
	testCases := []struct {
		name          string
		className     string
		id            string
		expectedWeight float64
	}{
		{
			name:           "no class or id",
			expectedWeight: 0,
		},
		{
			name:           "positive class only",
			className:      "article content",
			expectedWeight: 25,
		},
		{
			name:           "negative class only",
			className:      "comment sidebar",
			expectedWeight: -25,
		},
		{
			name:           "positive id only",
			id:             "main-content",
			expectedWeight: 25,
		},
		{
			name:           "negative id only",
			id:             "sidebar",
			expectedWeight: -25,
		},
		{
			name:           "positive class and positive id",
			className:      "article",
			id:             "content",
			expectedWeight: 50, // positive class(+25) + positive id(+25)
		},
		{
			name:           "negative class and negative id",
			className:      "comment",
			id:             "sidebar",
			expectedWeight: -50, // negative class(-25) + negative id(-25)
		},
		{
			name:           "positive class and negative id",
			className:      "article",
			id:             "sidebar",
			expectedWeight: 0, // positive class(+25) + negative id(-25)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test element
			element := dom.NewVElement("div")
			if tc.className != "" {
				element.SetAttribute("class", tc.className)
			}
			if tc.id != "" {
				element.SetAttribute("id", tc.id)
			}

			// Get the class weight
			weight := GetClassWeight(element)

			// Check if the weight matches the expected value
			if weight != tc.expectedWeight {
				t.Errorf("Expected weight %f, got %f", tc.expectedWeight, weight)
			}
		})
	}
}

func TestFindStructuralElements(t *testing.T) {
	testCases := []struct {
		name                    string
		setupDoc                func() *dom.VDocument
		expectHeader            bool
		expectFooter            bool
		expectSignificantNodes  int
		checkHeader             func(t *testing.T, header *dom.VElement)
		checkFooter             func(t *testing.T, footer *dom.VElement)
		checkSignificantNodes   func(t *testing.T, nodes []*dom.VElement)
	}{
		{
			name: "single header and footer tags",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				header := dom.NewVElement("header")
				header.AppendChild(dom.NewVText("This is a header"))
				body.AppendChild(header)
				
				content := dom.NewVElement("div")
				content.SetAttribute("class", "content")
				content.AppendChild(dom.NewVText("This is the main content"))
				body.AppendChild(content)
				
				footer := dom.NewVElement("footer")
				footer.AppendChild(dom.NewVText("This is a footer"))
				body.AppendChild(footer)
				
				return dom.NewVDocument(html, body)
			},
			expectHeader: true,
			expectFooter: true,
			expectSignificantNodes: 1, // content div
			checkHeader: func(t *testing.T, header *dom.VElement) {
				if header.TagName != "header" {
					t.Errorf("Expected header tag to be 'header', got '%s'", header.TagName)
				}
			},
			checkFooter: func(t *testing.T, footer *dom.VElement) {
				if footer.TagName != "footer" {
					t.Errorf("Expected footer tag to be 'footer', got '%s'", footer.TagName)
				}
			},
			checkSignificantNodes: func(t *testing.T, nodes []*dom.VElement) {
				if len(nodes) == 0 {
					t.Fatalf("Expected at least one significant node")
				}
				if nodes[0].ClassName() != "content" {
					t.Errorf("Expected significant node to have class 'content', got '%s'", nodes[0].ClassName())
				}
			},
		},
		{
			name: "header and footer by class/id",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				header := dom.NewVElement("div")
				header.SetAttribute("id", "header")
				header.AppendChild(dom.NewVText("This is a header div"))
				body.AppendChild(header)
				
				content := dom.NewVElement("main")
				content.AppendChild(dom.NewVText("This is the main content"))
				body.AppendChild(content)
				
				footer := dom.NewVElement("div")
				footer.SetAttribute("class", "footer")
				footer.AppendChild(dom.NewVText("This is a footer div"))
				body.AppendChild(footer)
				
				return dom.NewVDocument(html, body)
			},
			expectHeader: true,
			expectFooter: true,
			expectSignificantNodes: 1, // main element
			checkHeader: func(t *testing.T, header *dom.VElement) {
				if header.ID() != "header" {
					t.Errorf("Expected header id to be 'header', got '%s'", header.ID())
				}
			},
			checkFooter: func(t *testing.T, footer *dom.VElement) {
				if footer.ClassName() != "footer" {
					t.Errorf("Expected footer class to be 'footer', got '%s'", footer.ClassName())
				}
			},
			checkSignificantNodes: func(t *testing.T, nodes []*dom.VElement) {
				if len(nodes) == 0 {
					t.Fatalf("Expected at least one significant node")
				}
				if nodes[0].TagName != "main" {
					t.Errorf("Expected significant node to be 'main', got '%s'", nodes[0].TagName)
				}
			},
		},
		{
			name: "header and footer by role",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				header := dom.NewVElement("div")
				header.SetAttribute("role", "banner")
				header.AppendChild(dom.NewVText("This is a header with role"))
				body.AppendChild(header)
				
				article := dom.NewVElement("article")
				article.AppendChild(dom.NewVText("This is an article"))
				body.AppendChild(article)
				
				footer := dom.NewVElement("div")
				footer.SetAttribute("role", "contentinfo")
				footer.AppendChild(dom.NewVText("This is a footer with role"))
				body.AppendChild(footer)
				
				return dom.NewVDocument(html, body)
			},
			expectHeader: true,
			expectFooter: true,
			expectSignificantNodes: 1, // article element
			checkHeader: func(t *testing.T, header *dom.VElement) {
				role := GetAttribute(header, "role")
				if role != "banner" {
					t.Errorf("Expected header role to be 'banner', got '%s'", role)
				}
			},
			checkFooter: func(t *testing.T, footer *dom.VElement) {
				role := GetAttribute(footer, "role")
				if role != "contentinfo" {
					t.Errorf("Expected footer role to be 'contentinfo', got '%s'", role)
				}
			},
			checkSignificantNodes: func(t *testing.T, nodes []*dom.VElement) {
				if len(nodes) == 0 {
					t.Fatalf("Expected at least one significant node")
				}
				if nodes[0].TagName != "article" {
					t.Errorf("Expected significant node to be 'article', got '%s'", nodes[0].TagName)
				}
			},
		},
		{
			name: "multiple significant nodes",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				header := dom.NewVElement("header")
				header.AppendChild(dom.NewVText("This is a header"))
				body.AppendChild(header)
				
				main := dom.NewVElement("main")
				main.AppendChild(dom.NewVText("This is the main content"))
				body.AppendChild(main)
				
				article := dom.NewVElement("article")
				article.AppendChild(dom.NewVText("This is an article"))
				body.AppendChild(article)
				
				section := dom.NewVElement("section")
				section.SetAttribute("class", "content")
				section.AppendChild(dom.NewVText("This is a content section"))
				body.AppendChild(section)
				
				aside := dom.NewVElement("aside")
				aside.AppendChild(dom.NewVText("This is an aside"))
				body.AppendChild(aside)
				
				footer := dom.NewVElement("footer")
				footer.AppendChild(dom.NewVText("This is a footer"))
				body.AppendChild(footer)
				
				return dom.NewVDocument(html, body)
			},
			expectHeader: true,
			expectFooter: true,
			expectSignificantNodes: 4, // main, article, section, aside
			checkHeader: func(t *testing.T, header *dom.VElement) {
				if header.TagName != "header" {
					t.Errorf("Expected header tag to be 'header', got '%s'", header.TagName)
				}
			},
			checkFooter: func(t *testing.T, footer *dom.VElement) {
				if footer.TagName != "footer" {
					t.Errorf("Expected footer tag to be 'footer', got '%s'", footer.TagName)
				}
			},
			checkSignificantNodes: func(t *testing.T, nodes []*dom.VElement) {
				if len(nodes) < 4 {
					t.Fatalf("Expected at least 4 significant nodes, got %d", len(nodes))
				}
				
				// Check if all expected elements are present
				foundMain := false
				foundArticle := false
				foundSection := false
				foundAside := false
				
				for _, node := range nodes {
					switch node.TagName {
					case "main":
						foundMain = true
					case "article":
						foundArticle = true
					case "section":
						foundSection = true
					case "aside":
						foundAside = true
					}
				}
				
				if !foundMain {
					t.Errorf("Expected to find 'main' element in significant nodes")
				}
				if !foundArticle {
					t.Errorf("Expected to find 'article' element in significant nodes")
				}
				if !foundSection {
					t.Errorf("Expected to find 'section' element in significant nodes")
				}
				if !foundAside {
					t.Errorf("Expected to find 'aside' element in significant nodes")
				}
			},
		},
		{
			name: "no header or footer",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				main := dom.NewVElement("main")
				main.AppendChild(dom.NewVText("This is the main content"))
				body.AppendChild(main)
				
				article := dom.NewVElement("article")
				article.AppendChild(dom.NewVText("This is an article"))
				body.AppendChild(article)
				
				return dom.NewVDocument(html, body)
			},
			expectHeader: false,
			expectFooter: false,
			expectSignificantNodes: 2, // main, article
			checkHeader: nil,
			checkFooter: nil,
			checkSignificantNodes: func(t *testing.T, nodes []*dom.VElement) {
				if len(nodes) < 2 {
					t.Fatalf("Expected at least 2 significant nodes, got %d", len(nodes))
				}
				
				// Check if all expected elements are present
				foundMain := false
				foundArticle := false
				
				for _, node := range nodes {
					switch node.TagName {
					case "main":
						foundMain = true
					case "article":
						foundArticle = true
					}
				}
				
				if !foundMain {
					t.Errorf("Expected to find 'main' element in significant nodes")
				}
				if !foundArticle {
					t.Errorf("Expected to find 'article' element in significant nodes")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := tc.setupDoc()
			header, footer, significantNodes := FindStructuralElements(doc)
			
			// Check header
			if tc.expectHeader {
				if header == nil {
					t.Errorf("Expected to find a header, but got nil")
				} else if tc.checkHeader != nil {
					tc.checkHeader(t, header)
				}
			} else {
				if header != nil {
					t.Errorf("Expected no header, but got one: %s", header.TagName)
				}
			}
			
			// Check footer
			if tc.expectFooter {
				if footer == nil {
					t.Errorf("Expected to find a footer, but got nil")
				} else if tc.checkFooter != nil {
					tc.checkFooter(t, footer)
				}
			} else {
				if footer != nil {
					t.Errorf("Expected no footer, but got one: %s", footer.TagName)
				}
			}
			
			// Check significant nodes
			if len(significantNodes) != tc.expectSignificantNodes {
				t.Errorf("Expected %d significant nodes, got %d", tc.expectSignificantNodes, len(significantNodes))
			}
			
			if tc.checkSignificantNodes != nil && len(significantNodes) > 0 {
				tc.checkSignificantNodes(t, significantNodes)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	testCases := []struct {
		name        string
		html        string
		options     ReadabilityOptions
		checkResult func(t *testing.T, result ReadabilityArticle, err error)
	}{
		{
			name: "simple article",
			html: `<!DOCTYPE html>
<html>
<head>
  <title>Test Article</title>
  <meta name="author" content="Test Author">
</head>
<body>
  <article>
    <h1>Article Heading</h1>
    <p>This is a test article with enough content to be considered an article.
    It has multiple sentences and paragraphs to ensure it passes the content threshold.
    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor 
    incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud 
    exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
    <p>Second paragraph with more content to ensure it's long enough.
    Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
    Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
  </article>
</body>
</html>`,
			options: ReadabilityOptions{
				CharThreshold:   500,
				NbTopCandidates: 5,
				GenerateAriaTree: false,
			},
			checkResult: func(t *testing.T, result ReadabilityArticle, err error) {
				// Check for errors
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				
				// Check title and byline
				if result.Title != "Test Article" {
					t.Errorf("Expected title 'Test Article', got '%s'", result.Title)
				}
				
				if result.Byline != "Test Author" {
					t.Errorf("Expected byline 'Test Author', got '%s'", result.Byline)
				}
				
				// Check page type
				if result.PageType != PageTypeArticle {
					t.Errorf("Expected page type 'article', got '%s'", result.PageType)
				}
				
				// Check content extraction
				if result.Root == nil {
					t.Errorf("Expected content to be extracted, but Root is nil")
				} else {
					if result.Root.TagName != "article" {
						t.Errorf("Expected root element to be 'article', got '%s'", result.Root.TagName)
					}
				}
				
				// Check node count
				if result.NodeCount <= 0 {
					t.Errorf("Expected positive node count, got %d", result.NodeCount)
				}
			},
		},
		// Note: The parser is quite forgiving and can handle invalid HTML,
		// so we're not testing for errors with invalid HTML anymore
		{
			name: "non-article page",
			html: `<!DOCTYPE html>
<html>
<head>
  <title>Index Page</title>
</head>
<body>
  <div class="navigation">
    <ul>
      <li><a href="#">Link 1</a></li>
      <li><a href="#">Link 2</a></li>
      <li><a href="#">Link 3</a></li>
    </ul>
  </div>
  <div class="items">
    <div class="item">
      <h2><a href="#">Item 1</a></h2>
      <p>Short description</p>
    </div>
    <div class="item">
      <h2><a href="#">Item 2</a></h2>
      <p>Short description</p>
    </div>
    <div class="item">
      <h2><a href="#">Item 3</a></h2>
      <p>Short description</p>
    </div>
  </div>
</body>
</html>`,
			options: ReadabilityOptions{
				ForcedPageType: PageTypeOther,
			},
			checkResult: func(t *testing.T, result ReadabilityArticle, err error) {
				// Check for errors
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				
				// Check title
				if result.Title != "Index Page" {
					t.Errorf("Expected title 'Index Page', got '%s'", result.Title)
				}
				
				// Check page type
				if result.PageType != PageTypeOther {
					t.Errorf("Expected page type 'other', got '%s'", result.PageType)
				}
				
				// Check content extraction
				if result.Root != nil {
					t.Errorf("Expected content to not be extracted for non-article page, but Root is not nil")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Extract(tc.html, tc.options)
			
			if tc.checkResult != nil {
				tc.checkResult(t, result, err)
			}
		})
	}
}

func TestCreateExtractor(t *testing.T) {
	// Create an extractor with custom options
	options := ReadabilityOptions{
		CharThreshold:   300,
		NbTopCandidates: 3,
		ForcedPageType:  PageTypeArticle,
	}
	
	extractor := CreateExtractor(options)
	
	// Test HTML with enough content to pass the threshold
	html := `<!DOCTYPE html>
<html>
<head>
  <title>Test Article</title>
</head>
<body>
  <article>
    <h1>Article Heading</h1>
    <p>This is a test article with enough content to be considered an article.
    It has multiple sentences to ensure it passes the content threshold.
    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor 
    incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud 
    exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure 
    dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.</p>
  </article>
</body>
</html>`
	
	// Extract content using the custom extractor
	result, err := extractor(html)
	
	// Check results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", result.Title)
	}
	
	if result.PageType != PageTypeArticle {
		t.Errorf("Expected page type 'article', got '%s'", result.PageType)
	}
	
	if result.Root == nil {
		t.Errorf("Expected content to be extracted, but Root is nil")
	} else {
		if result.Root.TagName != "article" {
			t.Errorf("Expected root element to be 'article', got '%s'", result.Root.TagName)
		}
	}
}

func TestExtractContent(t *testing.T) {
	testCases := []struct {
		name           string
		setupDoc       func() *dom.VDocument
		options        ReadabilityOptions
		checkResult    func(t *testing.T, result ReadabilityArticle)
	}{
		{
			name: "article page with good content",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Add title
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Test Article Title"))
				head := dom.NewVElement("head")
				head.AppendChild(title)
				html.AppendChild(head)
				
				// Add author meta
				meta := dom.NewVElement("meta")
				meta.SetAttribute("name", "author")
				meta.SetAttribute("content", "Test Author")
				head.AppendChild(meta)
				
				// Add main content
				article := dom.NewVElement("article")
				article.SetAttribute("class", "content")
				
				// Add enough text to pass the threshold
				longText := "This is a long article text that should be considered as content. " +
					"It has multiple sentences and is definitely longer than the default threshold. " +
					"This should be enough to pass the text length check in the ExtractContent function. " +
					"We need to make sure it's long enough to be considered as an article. " +
					"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
					"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
					"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
					"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. " +
					"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
				
				p := dom.NewVElement("p")
				p.AppendChild(dom.NewVText(longText))
				article.AppendChild(p)
				body.AppendChild(article)
				
				return dom.NewVDocument(html, body)
			},
			options: ReadabilityOptions{
				CharThreshold: 500,
				NbTopCandidates: 5,
				GenerateAriaTree: false,
			},
			checkResult: func(t *testing.T, result ReadabilityArticle) {
				// Check title and byline
				if result.Title != "Test Article Title" {
					t.Errorf("Expected title 'Test Article Title', got '%s'", result.Title)
				}
				
				if result.Byline != "Test Author" {
					t.Errorf("Expected byline 'Test Author', got '%s'", result.Byline)
				}
				
				// Check page type
				if result.PageType != PageTypeArticle {
					t.Errorf("Expected page type 'article', got '%s'", result.PageType)
				}
				
				// Check content extraction
				if result.Root == nil {
					t.Errorf("Expected content to be extracted, but Root is nil")
				} else {
					if result.Root.TagName != "article" {
						t.Errorf("Expected root element to be 'article', got '%s'", result.Root.TagName)
					}
					
					if result.Root.ClassName() != "content" {
						t.Errorf("Expected root element to have class 'content', got '%s'", result.Root.ClassName())
					}
				}
				
				// Check node count
				if result.NodeCount <= 0 {
					t.Errorf("Expected positive node count, got %d", result.NodeCount)
				}
				
				// Check that structural elements are not set
				if result.Header != nil || result.Footer != nil || len(result.OtherSignificantNodes) > 0 {
					t.Errorf("Expected structural elements to be nil for article with content")
				}
				
				// Check that AriaTree is not generated
				if result.AriaTree != nil {
					t.Errorf("Expected AriaTree to be nil when GenerateAriaTree is false")
				}
			},
		},
		{
			name: "article page with content below threshold",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Add title
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Short Article"))
				head := dom.NewVElement("head")
				head.AppendChild(title)
				html.AppendChild(head)
				
				// Add main content with short text
				article := dom.NewVElement("article")
				article.SetAttribute("class", "content")
				
				// Add text below the threshold
				shortText := "This is a short article text."
				
				p := dom.NewVElement("p")
				p.AppendChild(dom.NewVText(shortText))
				article.AppendChild(p)
				body.AppendChild(article)
				
				// Add header and footer
				header := dom.NewVElement("header")
				header.AppendChild(dom.NewVText("This is a header"))
				body.AppendChild(header)
				
				footer := dom.NewVElement("footer")
				footer.AppendChild(dom.NewVText("This is a footer"))
				body.AppendChild(footer)
				
				return dom.NewVDocument(html, body)
			},
			options: ReadabilityOptions{
				CharThreshold: 500,
				NbTopCandidates: 5,
				GenerateAriaTree: false,
				ForcedPageType: PageTypeArticle, // Force article type
			},
			checkResult: func(t *testing.T, result ReadabilityArticle) {
				// Check title
				if result.Title != "Short Article" {
					t.Errorf("Expected title 'Short Article', got '%s'", result.Title)
				}
				
				// Check page type
				if result.PageType != PageTypeArticle {
					t.Errorf("Expected page type 'article', got '%s'", result.PageType)
				}
				
				// Check content extraction failed due to threshold
				if result.Root != nil {
					t.Errorf("Expected content to not be extracted due to threshold, but Root is not nil")
				}
				
				// Check structural elements are set
				if result.Header == nil {
					t.Errorf("Expected header to be set for article without content")
				} else if result.Header.TagName != "header" {
					t.Errorf("Expected header tag to be 'header', got '%s'", result.Header.TagName)
				}
				
				if result.Footer == nil {
					t.Errorf("Expected footer to be set for article without content")
				} else if result.Footer.TagName != "footer" {
					t.Errorf("Expected footer tag to be 'footer', got '%s'", result.Footer.TagName)
				}
			},
		},
		{
			name: "non-article page",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Add title
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Index Page"))
				head := dom.NewVElement("head")
				head.AppendChild(title)
				html.AppendChild(head)
				
				// Add multiple links to make it look like an index page
				for i := 0; i < 20; i++ {
					div := dom.NewVElement("div")
					div.SetAttribute("class", "item")
					
					a := dom.NewVElement("a")
					a.SetAttribute("href", "#")
					a.AppendChild(dom.NewVText("Link " + strconv.Itoa(i)))
					
					div.AppendChild(a)
					body.AppendChild(div)
				}
				
				return dom.NewVDocument(html, body)
			},
			options: ReadabilityOptions{
				CharThreshold: 500,
				NbTopCandidates: 5,
				GenerateAriaTree: false,
				ForcedPageType: PageTypeOther, // Force other type
			},
			checkResult: func(t *testing.T, result ReadabilityArticle) {
				// Check title
				if result.Title != "Index Page" {
					t.Errorf("Expected title 'Index Page', got '%s'", result.Title)
				}
				
				// Check page type
				if result.PageType != PageTypeOther {
					t.Errorf("Expected page type 'other', got '%s'", result.PageType)
				}
				
				// Check content extraction
				if result.Root != nil {
					t.Errorf("Expected content to not be extracted for non-article page, but Root is not nil")
				}
				
				// Check node count
				if result.NodeCount != 0 {
					t.Errorf("Expected node count to be 0 for non-article page, got %d", result.NodeCount)
				}
			},
		},
		{
			name: "with AriaTree generation",
			setupDoc: func() *dom.VDocument {
				html := dom.NewVElement("html")
				body := dom.NewVElement("body")
				html.AppendChild(body)
				
				// Add title
				title := dom.NewVElement("title")
				title.AppendChild(dom.NewVText("Test Article"))
				head := dom.NewVElement("head")
				head.AppendChild(title)
				html.AppendChild(head)
				
				// Add main content
				article := dom.NewVElement("article")
				article.SetAttribute("class", "content")
				
				// Add enough text to pass the threshold
				longText := "This is a long article text that should be considered as content. " +
					"It has multiple sentences and is definitely longer than the default threshold."
				
				p := dom.NewVElement("p")
				p.AppendChild(dom.NewVText(longText))
				article.AppendChild(p)
				body.AppendChild(article)
				
				return dom.NewVDocument(html, body)
			},
			options: ReadabilityOptions{
				CharThreshold: 100, // Lower threshold
				NbTopCandidates: 5,
				GenerateAriaTree: true, // Enable AriaTree generation
			},
			checkResult: func(t *testing.T, result ReadabilityArticle) {
				// Check title
				if result.Title != "Test Article" {
					t.Errorf("Expected title 'Test Article', got '%s'", result.Title)
				}
				
				// Check page type
				if result.PageType != PageTypeArticle {
					t.Errorf("Expected page type 'article', got '%s'", result.PageType)
				}
				
				// Check content extraction
				if result.Root == nil {
					t.Errorf("Expected content to be extracted, but Root is nil")
				}
				
				// AriaTree would be nil in our implementation since we haven't implemented it yet
				// This test will need to be updated when AriaTree is implemented
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc := tc.setupDoc()
			result := ExtractContent(doc, tc.options)
			
			if tc.checkResult != nil {
				tc.checkResult(t, result)
			}
		})
	}
}

func TestAddSignificantElementsByClassOrId(t *testing.T) {
	testCases := []struct {
		name           string
		setupDoc       func() (*dom.VElement, []*dom.VElement)
		expectedCount  int
		checkElements  func(t *testing.T, elements []*dom.VElement)
	}{
		{
			name: "elements with significant class names",
			setupDoc: func() (*dom.VElement, []*dom.VElement) {
				body := dom.NewVElement("body")
				
				// Create elements with significant class names
				content := dom.NewVElement("div")
				content.SetAttribute("class", "content")
				body.AppendChild(content)
				
				article := dom.NewVElement("div")
				article.SetAttribute("class", "article")
				body.AppendChild(article)
				
				main := dom.NewVElement("div")
				main.SetAttribute("class", "main-container")
				body.AppendChild(main)
				
				// Create element with non-significant class name
				other := dom.NewVElement("div")
				other.SetAttribute("class", "other")
				body.AppendChild(other)
				
				return body, []*dom.VElement{}
			},
			expectedCount: 3,
			checkElements: func(t *testing.T, elements []*dom.VElement) {
				// Check if all expected elements are present
				foundContent := false
				foundArticle := false
				foundMain := false
				
				for _, el := range elements {
					className := el.ClassName()
					if className == "content" {
						foundContent = true
					} else if className == "article" {
						foundArticle = true
					} else if className == "main-container" {
						foundMain = true
					}
				}
				
				if !foundContent {
					t.Errorf("Expected to find element with class 'content'")
				}
				if !foundArticle {
					t.Errorf("Expected to find element with class 'article'")
				}
				if !foundMain {
					t.Errorf("Expected to find element with class 'main-container'")
				}
			},
		},
		{
			name: "elements with significant IDs",
			setupDoc: func() (*dom.VElement, []*dom.VElement) {
				body := dom.NewVElement("body")
				
				// Create elements with significant IDs
				content := dom.NewVElement("div")
				content.SetAttribute("id", "content")
				body.AppendChild(content)
				
				article := dom.NewVElement("div")
				article.SetAttribute("id", "article")
				body.AppendChild(article)
				
				blog := dom.NewVElement("div")
				blog.SetAttribute("id", "blog-post")
				body.AppendChild(blog)
				
				// Create element with non-significant ID
				other := dom.NewVElement("div")
				other.SetAttribute("id", "other")
				body.AppendChild(other)
				
				return body, []*dom.VElement{}
			},
			expectedCount: 3,
			checkElements: func(t *testing.T, elements []*dom.VElement) {
				// Check if all expected elements are present
				foundContent := false
				foundArticle := false
				foundBlog := false
				
				for _, el := range elements {
					id := el.ID()
					if id == "content" {
						foundContent = true
					} else if id == "article" {
						foundArticle = true
					} else if id == "blog-post" {
						foundBlog = true
					}
				}
				
				if !foundContent {
					t.Errorf("Expected to find element with id 'content'")
				}
				if !foundArticle {
					t.Errorf("Expected to find element with id 'article'")
				}
				if !foundBlog {
					t.Errorf("Expected to find element with id 'blog-post'")
				}
			},
		},
		{
			name: "mixed elements with existing potentialNodes",
			setupDoc: func() (*dom.VElement, []*dom.VElement) {
				body := dom.NewVElement("body")
				
				// Create elements with significant class names and IDs
				content := dom.NewVElement("div")
				content.SetAttribute("class", "content")
				body.AppendChild(content)
				
				article := dom.NewVElement("div")
				article.SetAttribute("id", "article")
				body.AppendChild(article)
				
				// Create element with non-significant class and ID
				other := dom.NewVElement("div")
				other.SetAttribute("class", "other")
				other.SetAttribute("id", "other")
				body.AppendChild(other)
				
				// Create existing potentialNodes
				main := dom.NewVElement("main")
				body.AppendChild(main)
				
				existingNodes := []*dom.VElement{main}
				
				return body, existingNodes
			},
			expectedCount: 3, // main (existing) + content + article
			checkElements: func(t *testing.T, elements []*dom.VElement) {
				// Check if all expected elements are present
				foundMain := false
				foundContent := false
				foundArticle := false
				
				for _, el := range elements {
					if el.TagName == "main" {
						foundMain = true
					} else if el.ClassName() == "content" {
						foundContent = true
					} else if el.ID() == "article" {
						foundArticle = true
					}
				}
				
				if !foundMain {
					t.Errorf("Expected to find 'main' element")
				}
				if !foundContent {
					t.Errorf("Expected to find element with class 'content'")
				}
				if !foundArticle {
					t.Errorf("Expected to find element with id 'article'")
				}
			},
		},
		{
			name: "no significant elements",
			setupDoc: func() (*dom.VElement, []*dom.VElement) {
				body := dom.NewVElement("body")
				
				// Create elements with non-significant class names and IDs
				div1 := dom.NewVElement("div")
				div1.SetAttribute("class", "other")
				body.AppendChild(div1)
				
				div2 := dom.NewVElement("div")
				div2.SetAttribute("id", "other")
				body.AppendChild(div2)
				
				return body, []*dom.VElement{}
			},
			expectedCount: 0,
			checkElements: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, potentialNodes := tc.setupDoc()
			
			// Call the function
			AddSignificantElementsByClassOrId(body, &potentialNodes)
			
			// Check the results
			if len(potentialNodes) != tc.expectedCount {
				t.Errorf("Expected %d elements, got %d", tc.expectedCount, len(potentialNodes))
			}
			
			if tc.checkElements != nil && len(potentialNodes) > 0 {
				tc.checkElements(t, potentialNodes)
			}
		})
	}
}
