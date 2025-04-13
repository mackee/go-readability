package util

import (
	"testing"
)

func TestUnlikelyCandidates(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"ad-banner", true},
		{"sidebar", true},
		{"comment-section", true},
		{"footer", true},
		{"header", true},
		{"main-content", false},
		{"article", false},
		{"body", false},
	}

	for _, test := range tests {
		result := Regexps.UnlikelyCandidates.MatchString(test.input)
		if result != test.expected {
			t.Errorf("UnlikelyCandidates.MatchString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestOkMaybeItsACandidate(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"article", true},
		{"body", true},
		{"main-content", true},
		{"content", true},
		{"shadow-root", true},
		{"footer", false},
		{"sidebar", false},
		{"comment", false},
	}

	for _, test := range tests {
		result := Regexps.OkMaybeItsACandidate.MatchString(test.input)
		if result != test.expected {
			t.Errorf("OkMaybeItsACandidate.MatchString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestPositive(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"article", true},
		{"body", true},
		{"content", true},
		{"main", true},
		{"blog-post", true},
		{"story", true},
		{"footer", false},
		{"sidebar", false},
		{"comment", false},
	}

	for _, test := range tests {
		result := Regexps.Positive.MatchString(test.input)
		if result != test.expected {
			t.Errorf("Positive.MatchString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestNegative(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"footer", true},
		{"sidebar", true},
		{"comment", true},
		{"hidden", true},
		{"hid", true},
		{"article", false},
		{"content", false},
		{"main", false},
	}

	for _, test := range tests {
		result := Regexps.Negative.MatchString(test.input)
		if result != test.expected {
			t.Errorf("Negative.MatchString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestCommas(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{",", true},        // U+002C: COMMA
		{"،", true},        // U+060C: ARABIC COMMA
		{"﹐", true},        // U+FE50: SMALL COMMA
		{"，", true},        // U+FF0C: FULLWIDTH COMMA
		{"、", true},        // U+3001: IDEOGRAPHIC COMMA
		{"abc,def", true},  // Contains comma
		{"abc def", false}, // No comma
	}

	for _, test := range tests {
		result := Regexps.Commas.MatchString(test.input)
		if result != test.expected {
			t.Errorf("Commas.MatchString(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a  b", "a b"},
		{"a   b", "a b"},
		{"a    b", "a b"},
		{"a\t\tb", "a b"},
		{"a\n\nb", "a b"},
		{"a\r\rb", "a b"},
		{"a b", "a b"}, // No change
	}

	for _, test := range tests {
		result := Regexps.Normalize.ReplaceAllString(test.input, " ")
		if result != test.expected {
			t.Errorf("Normalize.ReplaceAllString(%q, \" \") = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestDefaultTagsToScore(t *testing.T) {
	expected := []string{"section", "h2", "h3", "h4", "h5", "h6", "p", "td", "pre"}
	if len(DefaultTagsToScore) != len(expected) {
		t.Errorf("DefaultTagsToScore has %d elements, expected %d", len(DefaultTagsToScore), len(expected))
	}

	for i, tag := range expected {
		if i >= len(DefaultTagsToScore) || DefaultTagsToScore[i] != tag {
			t.Errorf("DefaultTagsToScore[%d] = %q, expected %q", i, DefaultTagsToScore[i], tag)
		}
	}
}

func TestDivToPElems(t *testing.T) {
	expected := []string{"blockquote", "dl", "div", "img", "ol", "p", "pre", "table", "ul"}
	if len(DivToPElems) != len(expected) {
		t.Errorf("DivToPElems has %d elements, expected %d", len(DivToPElems), len(expected))
	}

	for _, tag := range expected {
		if !DivToPElems[tag] {
			t.Errorf("DivToPElems[%q] = false, expected true", tag)
		}
	}
}

func TestPhrasingElems(t *testing.T) {
	// サンプルとして一部の要素をチェック
	expectedSample := []string{"abbr", "audio", "b", "br", "code", "em", "i", "img", "span", "strong"}
	for _, tag := range expectedSample {
		found := false
		for _, elem := range PhrasingElems {
			if elem == tag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PhrasingElems does not contain %q", tag)
		}
	}

	// 要素数の確認
	expectedCount := 39 // 元のTypeScriptコードの要素数
	if len(PhrasingElems) != expectedCount {
		t.Errorf("PhrasingElems has %d elements, expected %d", len(PhrasingElems), expectedCount)
	}
}
