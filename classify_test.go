package readability

import (
	"strings"
	"testing"

	"github.com/mackee/go-readability/internal/dom"
	"github.com/mackee/go-readability/internal/parser"
)

func TestGetExpectedPageTypeByUrl(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected PageType
	}{
		{
			name:     "記事URLパターン - /articles/を含む",
			url:      "https://example.com/articles/123",
			expected: PageTypeArticle,
		},
		{
			name:     "記事URLパターン - 3階層以上のパス",
			url:      "https://example.com/blog/category/post-title",
			expected: PageTypeArticle,
		},
		{
			name:     "記事URLパターン - 数字のみのID",
			url:      "https://example.com/post/12345",
			expected: PageTypeArticle,
		},
		{
			name:     "記事URLパターン - 英数字混合のID",
			url:      "https://example.com/post/abc123def",
			expected: PageTypeArticle,
		},
		{
			name:     "非記事URLパターン - トップページ",
			url:      "https://example.com/",
			expected: PageTypeOther,
		},
		{
			name:     "非記事URLパターン - 1階層のパス",
			url:      "https://example.com/about",
			expected: PageTypeOther,
		},
		{
			name:     "非記事URLパターン - 2階層のパス（英字のみ）",
			url:      "https://example.com/blog/about",
			expected: PageTypeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetExpectedPageTypeByUrl(tt.url)
			if result != tt.expected {
				t.Errorf("GetExpectedPageTypeByUrl(%s) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestAnalyzeUrlPattern(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		contains string // 結果に含まれるべき文字列
	}{
		{
			name:     "数字のみのパターン",
			url:      "https://example.com/post/12345",
			contains: "数字のみ",
		},
		{
			name:     "英数字混合のパターン",
			url:      "https://example.com/post/abc123def",
			contains: "英数字混合",
		},
		{
			name:     "英字のみのパターン",
			url:      "https://example.com/post/abcdef",
			contains: "英字のみ",
		},
		{
			name:     "末尾なしのパターン",
			url:      "https://example.com/",
			contains: "末尾なし",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnalyzeUrlPattern(tt.url)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("AnalyzeUrlPattern(%s) = %v, want to contain %v", tt.url, result, tt.contains)
			}
		})
	}
}

func TestIsSemanticTag(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "main要素",
			html:     "<main>メインコンテンツ</main>",
			expected: true,
		},
		{
			name:     "article要素",
			html:     "<article>記事コンテンツ</article>",
			expected: true,
		},
		{
			name:     "contentクラスを持つdiv",
			html:     "<div class=\"content\">コンテンツ</div>",
			expected: true,
		},
		{
			name:     "contentIDを持つdiv",
			html:     "<div id=\"content\">コンテンツ</div>",
			expected: true,
		},
		{
			name:     "子要素にmainを持つdiv",
			html:     "<div><main>メインコンテンツ</main></div>",
			expected: true,
		},
		{
			name:     "子要素にarticleを持つdiv",
			html:     "<div><article>記事コンテンツ</article></div>",
			expected: true,
		},
		{
			name:     "セマンティックでない要素",
			html:     "<div>通常のdiv</div>",
			expected: false,
		},
		{
			name:     "セマンティックでないクラス",
			html:     "<div class=\"wrapper\">ラッパー</div>",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.ParseHTML(tt.html, "https://example.com")
			if err != nil {
				t.Fatalf("HTML解析エラー: %v", err)
			}

			var element *dom.VElement
			if doc.Body != nil && len(doc.Body.Children) > 0 {
				if elem, ok := doc.Body.Children[0].(*dom.VElement); ok {
					element = elem
				}
			}

			if element == nil {
				t.Fatalf("テスト要素が見つかりません")
			}

			result := IsSemanticTag(element)
			if result != tt.expected {
				t.Errorf("IsSemanticTag() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsSignificantNode(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "header要素",
			html:     "<header>ヘッダー</header>",
			expected: true,
		},
		{
			name:     "footer要素",
			html:     "<footer>フッター</footer>",
			expected: true,
		},
		{
			name:     "main要素",
			html:     "<main>メインコンテンツ</main>",
			expected: true,
		},
		{
			name:     "article要素",
			html:     "<article>記事</article>",
			expected: true,
		},
		{
			name:     "aside要素",
			html:     "<aside>サイドバー</aside>",
			expected: true,
		},
		{
			name:     "nav要素",
			html:     "<nav>ナビゲーション</nav>",
			expected: true,
		},
		{
			name:     "role=banner属性を持つdiv",
			html:     "<div role=\"banner\">バナー</div>",
			expected: true,
		},
		{
			name:     "role=contentinfo属性を持つdiv",
			html:     "<div role=\"contentinfo\">コンテンツ情報</div>",
			expected: true,
		},
		{
			name:     "role=main属性を持つdiv",
			html:     "<div role=\"main\">メイン</div>",
			expected: true,
		},
		{
			name:     "headerクラスを持つdiv",
			html:     "<div class=\"header\">ヘッダー</div>",
			expected: true,
		},
		{
			name:     "footerクラスを持つdiv",
			html:     "<div class=\"footer\">フッター</div>",
			expected: true,
		},
		{
			name:     "mainContentクラスを持つdiv",
			html:     "<div class=\"mainContent\">メインコンテンツ</div>",
			expected: true,
		},
		{
			name:     "有意でない要素",
			html:     "<div>通常のdiv</div>",
			expected: false,
		},
		{
			name:     "有意でないクラス",
			html:     "<div class=\"wrapper\">ラッパー</div>",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.ParseHTML(tt.html, "https://example.com")
			if err != nil {
				t.Fatalf("HTML解析エラー: %v", err)
			}

			var element *dom.VElement
			if doc.Body != nil && len(doc.Body.Children) > 0 {
				if elem, ok := doc.Body.Children[0].(*dom.VElement); ok {
					element = elem
				}
			}

			if element == nil {
				t.Fatalf("テスト要素が見つかりません")
			}

			result := IsSignificantNode(element)
			if result != tt.expected {
				t.Errorf("IsSignificantNode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClassifyPageType(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		url           string
		charThreshold int
		expected      PageType
	}{
		{
			name: "記事ページ - セマンティックタグと十分なテキスト",
			html: `
				<html>
					<body>
						<article>
							<h1>記事タイトル</h1>
							<p>` + strings.Repeat("これは記事の本文です。", 100) + `</p>
						</article>
					</body>
				</html>
			`,
			url:           "https://example.com/articles/123",
			charThreshold: 500,
			expected:      PageTypeArticle,
		},
		{
			name: "記事ページ - 十分なテキストと低いリンク密度",
			html: `
				<html>
					<body>
						<div class="content">
							<h1>記事タイトル</h1>
							<p>` + strings.Repeat("これは記事の本文です。", 100) + `</p>
							<a href="#">リンク1</a>
							<a href="#">リンク2</a>
						</div>
					</body>
				</html>
			`,
			url:           "https://example.com/post/12345",
			charThreshold: 500,
			expected:      PageTypeArticle,
		},
		{
			name: "非記事ページ - リスト要素が多い",
			html: `
				<html>
					<body>
						<div>
							<h1>記事一覧</h1>
							` + strings.Repeat("<article class=\"card\"><h2>記事タイトル</h2><p>概要</p></article>", 15) + `
						</div>
					</body>
				</html>
			`,
			url:           "https://example.com/blog",
			charThreshold: 500,
			expected:      PageTypeOther,
		},
		{
			name: "非記事ページ - リンクが多く本文が少ない",
			html: `
				<html>
					<body>
						<div>
							<h1>リンク集</h1>
							<p>いくつかのリンク</p>
							` + strings.Repeat("<a href=\"#\">リンク</a>", 40) + `
						</div>
					</body>
				</html>
			`,
			url:           "https://example.com/links",
			charThreshold: 500,
			expected:      PageTypeOther,
		},
		{
			name: "非記事ページ - テキストが少ない",
			html: `
				<html>
					<body>
						<div>
							<h1>短いページ</h1>
							<p>これは短いテキストです。</p>
						</div>
					</body>
				</html>
			`,
			url:           "https://example.com/short",
			charThreshold: 500,
			expected:      PageTypeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.ParseHTML(tt.html, "https://example.com")
			if err != nil {
				t.Fatalf("HTML解析エラー: %v", err)
			}

			// 候補を見つける
			candidates := FindMainCandidates(doc, 5)

			// ページタイプを分類
			result := ClassifyPageType(doc, candidates, tt.charThreshold, tt.url)
			if result != tt.expected {
				t.Errorf("ClassifyPageType() = %v, want %v", result, tt.expected)
			}
		})
	}
}
