// Package readability provides functionality to extract readable content from HTML documents.
// It implements an algorithm similar to Mozilla's Readability.js to identify and extract
// the main content from web pages, removing clutter, navigation, ads, and other non-content elements.
package readability

import (
	"regexp"
	"strings"

	"github.com/mackee/go-readability/internal/dom"
	"github.com/mackee/go-readability/internal/util"
)

// ClassifyPageType classifies a document as an article or other type of page.
// It uses various heuristics including URL pattern, semantic tags, text length,
// link density, and more to determine the page type. This classification helps
// the extraction process decide how to handle different types of content.
//
// Parameters:
//   - doc: The parsed HTML document
//   - candidates: The list of content candidates found by the scoring algorithm
//   - charThreshold: The minimum character threshold for article content
//   - url: The URL of the page (optional, used for URL pattern analysis)
//
// Returns:
//   - PageType: Either PageTypeArticle or PageTypeOther
func ClassifyPageType(
	doc *dom.VDocument,
	candidates []*dom.VElement,
	charThreshold int,
	url string,
) PageType {
	// If charThreshold is not provided, use the default
	if charThreshold <= 0 {
		charThreshold = util.DefaultCharThreshold
	}

	// URLパターンによる判定（URLが提供された場合）
	if url != "" {
		// URLパターンが強い指標になる場合は、それを優先
		if strings.Contains(url, "/articles/") {
			// 候補がある場合のみ ARTICLE として扱う
			if len(candidates) > 0 {
				return PageTypeArticle
			}
			return PageTypeOther
		}

		// 追加: 末尾に英単語ではなさそうなハッシュ・連番・UUIDのような文字列を含む場合
		urlParts := strings.Split(url, "/")
		lastPart := urlParts[len(urlParts)-1]

		// 末尾の部分が存在し、.htmlなどの拡張子を含む場合はその前の部分を取得
		lastPartWithoutExt := strings.Split(lastPart, ".")[0]

		// 数字のみ、または数字と英字の混合で、かつ5文字以上の場合は記事IDと判断
		digitOnlyPattern := regexp.MustCompile(`^\d+$`)
		alphaNumericPattern := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
		hasDigitPattern := regexp.MustCompile(`\d`)

		if digitOnlyPattern.MatchString(lastPartWithoutExt) || // 数字のみ
			(alphaNumericPattern.MatchString(lastPartWithoutExt) && // 英数字のみ
				hasDigitPattern.MatchString(lastPartWithoutExt) && // 少なくとも1つの数字を含む
				len(lastPartWithoutExt) >= 5) { // 5文字以上
			// 候補がある場合のみ ARTICLE として扱う
			if len(candidates) > 0 {
				return PageTypeArticle
			}
			return PageTypeOther
		}

		// トップレベルドメインやユーザーページは OTHER の可能性が高い
		topLevelPattern := regexp.MustCompile(`^https?://[^/]+/?$`)
		userPagePattern := regexp.MustCompile(`^https?://[^/]+/[^/]+/?$`)

		if topLevelPattern.MatchString(url) || userPagePattern.MatchString(url) {
			// ただし、内容が明らかに記事の場合は例外
			if len(candidates) > 0 {
				textLength := GetInnerText(candidates[0], false)
				// 非常に長いテキストがあり、リンク密度が低い場合のみ ARTICLE
				if len(textLength) > charThreshold*2 && GetLinkDensity(candidates[0]) < 0.3 {
					return PageTypeArticle
				}
			}
			return PageTypeOther
		}
	}

	// 候補がない場合は OTHER
	if len(candidates) == 0 {
		return PageTypeOther
	}

	topCandidate := candidates[0]

	// 1. ページ構造の分析
	// 見出し数をカウント
	h1Elements := GetElementsByTagName(doc.Body, "h1")
	h2Elements := GetElementsByTagName(doc.Body, "h2")
	h3Elements := GetElementsByTagName(doc.Body, "h3")
	headingCount := len(h1Elements) + len(h2Elements) + len(h3Elements)

	// 画像数をカウント
	imgElements := GetElementsByTagName(doc.Body, "img")
	imageCount := len(imgElements)

	// リンク数をカウント
	aElements := GetElementsByTagName(doc.Body, "a")
	linkCount := len(aElements)

	// 記事リスト要素をカウント
	articleElements := GetElementsByTagName(doc.Body, "article")
	listItemElements := GetElementsByTagName(doc.Body, "li")

	// カード要素をカウント
	cardElements := []*dom.VElement{}
	for _, child := range doc.Body.Children {
		if childElem, ok := child.(*dom.VElement); ok {
			className := strings.ToLower(childElem.ClassName())
			if strings.Contains(className, "card") ||
				strings.Contains(className, "item") ||
				strings.Contains(className, "entry") {
				cardElements = append(cardElements, childElem)
			}
		}
	}

	listElementCount := len(articleElements) + len(listItemElements) + len(cardElements)

	// 2. トップページの特徴を検出
	// - 多数の記事/カードリスト要素
	// - 多数のリンク
	// - 多数の画像
	// - 見出しが少ない、または多すぎる
	hasIndexPageCharacteristics :=
		listElementCount > 10 || // 多数のリスト要素
			(linkCount > 50 && imageCount > 20) || // 多数のリンクと画像
			headingCount > 10 ||
			headingCount == 0 // 見出しが多すぎるか、まったくない

	if hasIndexPageCharacteristics {
		// トップページの特徴が強い場合は OTHER
		return PageTypeOther
	}

	// 3. セマンティックタグの確認 + テキスト長チェック
	isSemanticTag := IsSemanticTag(topCandidate)

	if isSemanticTag {
		textLength := GetInnerText(topCandidate, false)
		linkDensity := GetLinkDensity(topCandidate)

		// セマンティックタグでも、テキスト長が短すぎる場合は OTHER
		if len(textLength) >= charThreshold/2 && linkDensity <= 0.5 {
			// 記事リスト要素が多い場合は OTHER
			if listElementCount > 10 {
				return PageTypeOther
			}
			return PageTypeArticle
		}

		// テキスト長が非常に短い場合は OTHER
		if len(textLength) < 100 {
			return PageTypeOther
		}
	}

	// 4. テキスト長とリンク密度の確認
	textLength := GetInnerText(topCandidate, false)
	linkDensity := GetLinkDensity(topCandidate)

	// 記事の特徴: 十分なテキスト長、低いリンク密度、適切な見出し数
	if len(textLength) >= charThreshold &&
		linkDensity <= 0.5 &&
		headingCount >= 1 &&
		headingCount <= 10 {
		return PageTypeArticle
	}

	// 5. 候補のスコア差を確認（平衡性）
	if len(candidates) >= 2 {
		topScore := 0.0
		if topCandidate.GetReadabilityData() != nil {
			topScore = topCandidate.GetReadabilityData().ContentScore
		}

		secondScore := 0.0
		if candidates[1].GetReadabilityData() != nil {
			secondScore = candidates[1].GetReadabilityData().ContentScore
		}

		scoreRatio := 1.0
		if topScore > 0 {
			scoreRatio = secondScore / topScore
		}

		if scoreRatio > 0.8 {
			// 候補が平衡している場合、リンク密度と全体のリンク数を確認
			bodyTextLength := len(GetInnerText(doc.Body, false))
			var bodyLinkDensity float64 = 0
			if bodyTextLength > 0 {
				bodyLinkDensity = float64(linkCount) / float64(bodyTextLength)
			}

			// リンク密度が高い場合は OTHER（リスト/インデックスページの可能性）
			if bodyLinkDensity > 0.25 || linkDensity > 0.3 {
				return PageTypeOther
			}
		}
	}

	// 6. 全体のリンク数と本文の比率を確認
	bodyTextLength := len(GetInnerText(doc.Body, false))

	// リンクが多く、本文が少ない場合は OTHER
	if linkCount > 30 && bodyTextLength < int(float64(charThreshold)*1.5) {
		return PageTypeOther
	}

	// 7. 最終判定
	// ある程度のテキスト量があり、リンク密度が低い場合は ARTICLE
	if len(textLength) >= 140 && linkDensity <= 0.5 {
		// 記事リスト要素が多い場合は OTHER
		if listElementCount > 10 {
			return PageTypeOther
		}
		return PageTypeArticle
	}

	// それ以外の場合は OTHER
	return PageTypeOther
}

// IsSignificantNode determines if a node is semantically significant.
// This includes elements like header, footer, main, article, etc.
// Significant nodes are important structural elements that help understand
// the page's organization even when the main content extraction fails.
//
// Parameters:
//   - node: The element to check
//
// Returns:
//   - true if the node is semantically significant, false otherwise
func IsSignificantNode(node *dom.VElement) bool {
	// Check tag name
	tagName := strings.ToLower(node.TagName)
	if tagName == "header" || tagName == "footer" || tagName == "main" ||
		tagName == "article" || tagName == "aside" || tagName == "nav" {
		return true
	}

	// Check role attribute
	role := strings.ToLower(GetAttribute(node, "role"))
	if role == "banner" || role == "contentinfo" || role == "main" ||
		role == "navigation" || role == "complementary" {
		return true
	}

	// Check class and ID
	className := strings.ToLower(node.ClassName())
	id := strings.ToLower(node.ID())

	// Common class/id patterns for significant elements
	significantPatterns := []string{
		"header", "footer", "main", "content", "article", "navigation",
		"nav", "sidebar", "menu", "banner", "mainContent", "mainContainer",
	}

	for _, pattern := range significantPatterns {
		if strings.Contains(className, pattern) || strings.Contains(id, pattern) {
			return true
		}
	}

	return false
}

// IsSemanticTag checks if an element is a semantic tag or contains semantic tags.
// Semantic tags include main, article, and elements with content-related classes/IDs.
// These tags provide structural meaning to the content and are strong indicators
// of meaningful content areas.
//
// Parameters:
//   - element: The element to check
//
// Returns:
//   - true if the element is or contains semantic tags, false otherwise
func IsSemanticTag(element *dom.VElement) bool {
	// Check if the element itself is a semantic tag
	tagName := strings.ToLower(element.TagName)
	if tagName == "main" || tagName == "article" {
		return true
	}

	// Check class and ID for content indicators
	className := strings.ToLower(element.ClassName())
	id := strings.ToLower(element.ID())
	if strings.Contains(className, "content") || strings.Contains(id, "content") {
		return true
	}

	// Check if any child elements are semantic tags
	for _, child := range element.Children {
		if childElem, ok := child.(*dom.VElement); ok {
			childTagName := strings.ToLower(childElem.TagName)
			if childTagName == "main" || childTagName == "article" {
				return true
			}
		}
	}

	return false
}

// GetExpectedPageTypeByUrl determines the expected page type based on URL patterns.
// This is a helper function that can be used before full page analysis to get
// a preliminary classification based solely on URL patterns.
//
// Parameters:
//   - url: The URL of the page to analyze
//
// Returns:
//   - PageType: Either PageTypeArticle or PageTypeOther based on URL patterns
func GetExpectedPageTypeByUrl(url string) PageType {
	// URLパターンに基づく判定
	// 記事ページのパターン: /articles/ を含む、または特定のパターンに一致
	if strings.Contains(url, "/articles/") {
		return PageTypeArticle
	}

	// 3階層以上の深さを持つパス（少なくとも3つのスラッシュで区切られたパス）
	threeDepthPattern := regexp.MustCompile(`^https?://[^/]+/[^/]+/[^/]+/[^/]*$`)
	if threeDepthPattern.MatchString(url) {
		return PageTypeArticle
	}

	// 追加: 末尾に英単語ではなさそうなハッシュ・連番・UUIDのような文字列を含む場合
	urlParts := strings.Split(url, "/")
	lastPart := urlParts[len(urlParts)-1]

	// 末尾の部分が存在し、.htmlなどの拡張子を含む場合はその前の部分を取得
	lastPartWithoutExt := strings.Split(lastPart, ".")[0]

	// 数字のみ、または数字と英字の混合で、かつ5文字以上の場合は記事IDと判断
	digitOnlyPattern := regexp.MustCompile(`^\d+$`)
	alphaNumericPattern := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	hasDigitPattern := regexp.MustCompile(`\d`)

	if digitOnlyPattern.MatchString(lastPartWithoutExt) || // 数字のみ
		(alphaNumericPattern.MatchString(lastPartWithoutExt) && // 英数字のみ
			hasDigitPattern.MatchString(lastPartWithoutExt) && // 少なくとも1つの数字を含む
			len(lastPartWithoutExt) >= 5) { // 5文字以上
		return PageTypeArticle
	}

	// トップページやユーザーページなど
	return PageTypeOther
}

// AnalyzeUrlPattern analyzes the pattern of the URL's last part.
// This is a helper function for debugging and understanding URL patterns.
// It categorizes the last part of a URL into patterns like "numeric only",
// "alphanumeric", etc.
//
// Parameters:
//   - url: The URL to analyze
//
// Returns:
//   - A string describing the pattern of the URL's last part
func AnalyzeUrlPattern(url string) string {
	urlParts := strings.Split(url, "/")
	lastPart := urlParts[len(urlParts)-1]

	// 末尾の部分が存在し、.htmlなどの拡張子を含む場合はその前の部分を取得
	lastPartWithoutExt := strings.Split(lastPart, ".")[0]

	if lastPartWithoutExt == "" {
		return "末尾なし"
	}

	digitOnlyPattern := regexp.MustCompile(`^\d+$`)
	if digitOnlyPattern.MatchString(lastPartWithoutExt) {
		return "数字のみ (" + lastPartWithoutExt + ")"
	}

	alphaNumericPattern := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	hasDigitPattern := regexp.MustCompile(`\d`)
	if alphaNumericPattern.MatchString(lastPartWithoutExt) && hasDigitPattern.MatchString(lastPartWithoutExt) {
		return "英数字混合 (" + lastPartWithoutExt + ")"
	}

	alphaOnlyPattern := regexp.MustCompile(`^[a-zA-Z-_]+$`)
	if alphaOnlyPattern.MatchString(lastPartWithoutExt) {
		return "英字のみ (" + lastPartWithoutExt + ")"
	}

	return "その他 (" + lastPartWithoutExt + ")"
}
