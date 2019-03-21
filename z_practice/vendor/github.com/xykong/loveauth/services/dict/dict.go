package dict

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/huichen/sego"
)

var segmenter sego.Segmenter

// Load dictionaries from dictPath
func Load(dictPath string) {
	segmenter.LoadDictionary(dictPath)
}

// ExistInvalidWord Check if text contains words defined in dictionary
func ExistInvalidWord(text string) (exist bool, dirtyWords []string) {
	segments := getSegments(text)

	for _, seg := range segments {
		token := seg.Token()
		if token.Frequency() > 1 {
			exist = true
			dirtyWords = append(dirtyWords, token.Text())
		}
	}
	return
}

// ReplaceInvalidWords Replace words defineds in dictionary
func ReplaceInvalidWords(text string) (string, bool, []string) {
	segments := getSegments(text)
	var exist bool
	var dirtyWords []string
	for _, seg := range segments {
		token := seg.Token()
		if token.Frequency() > 1 {
			oldText := token.Text()
			newText := strings.Repeat("*", utf8.RuneCountInString(oldText))
			text = regexp.
				MustCompile(fmt.Sprintf("(?i)%s", oldText)).
				ReplaceAllLiteralString(text, newText)
			exist = true
			dirtyWords = append(dirtyWords, oldText)
		}
	}
	return text, exist, dirtyWords
}

func getSegments(text string) []sego.Segment {
	return segmenter.Segment([]byte(text))
}
