package util

import "regexp"

// DefaultNTopCandidates は、候補を分析する際に考慮するトップ候補の数です。
const DefaultNTopCandidates = 5

// DefaultCharThreshold は、結果を返すために記事が持つべき最小文字数です。
const DefaultCharThreshold = 500

// DefaultTagsToScore はデフォルトでスコアリングする要素タグです。
var DefaultTagsToScore = []string{
	"section", "h2", "h3", "h4", "h5", "h6", "p", "td", "pre",
}

// Regexps は、readability内で使用されるすべての正規表現です。
var Regexps = struct {
	// UnlikelyCandidates は、コンテンツとして不適切な要素を識別するための正規表現です。
	UnlikelyCandidates *regexp.Regexp

	// OkMaybeItsACandidate は、UnlikelyCandidatesに一致しても候補として考慮される可能性のある要素を識別するための正規表現です。
	OkMaybeItsACandidate *regexp.Regexp

	// Positive は、コンテンツとして適切な要素を識別するための正規表現です。
	Positive *regexp.Regexp

	// Negative は、コンテンツとして不適切な要素を識別するための正規表現です。
	Negative *regexp.Regexp

	// Commas は、ラテン語、シンディ語、中国語、その他の様々なスクリプトで使用されるコンマを識別するための正規表現です。
	Commas *regexp.Regexp

	// Normalize は、空白を正規化するための正規表現です。
	Normalize *regexp.Regexp
}{
	UnlikelyCandidates: regexp.MustCompile(`-ad-|ai2html|banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|footer|gdpr|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote`),
	OkMaybeItsACandidate: regexp.MustCompile(`and|article|body|column|content|main|shadow`),
	Positive:             regexp.MustCompile(`article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story`),
	Negative:             regexp.MustCompile(`-ad-|hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|footer|gdpr|masthead|media|meta|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|widget`),
	Commas:               regexp.MustCompile(`,|،|﹐|︐|︑|⹁|⹔|⹒|，|、`),
	Normalize:            regexp.MustCompile(`\s{2,}`),
}

// DivToPElems は、hasChildBlockElementで使用される要素のセットです。
var DivToPElems = map[string]bool{
	"blockquote": true,
	"dl":         true,
	"div":        true,
	"img":        true,
	"ol":         true,
	"p":          true,
	"pre":        true,
	"table":      true,
	"ul":         true,
}

// PhrasingElems は、isPhrasingContentで使用される要素のスライスです。
var PhrasingElems = []string{
	"abbr", "audio", "b", "bdo", "br", "button", "cite", "code", "data",
	"datalist", "dfn", "em", "embed", "i", "img", "input", "kbd", "label",
	"mark", "math", "meter", "noscript", "object", "output", "progress", "q",
	"ruby", "samp", "script", "select", "small", "span", "strong", "sub",
	"sup", "textarea", "time", "var", "wbr",
}
