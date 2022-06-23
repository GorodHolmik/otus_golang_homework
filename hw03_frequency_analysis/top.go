package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var matchWordsRegular = regexp.MustCompile(`\pL+(-\pL)?\pL*`)

type wordOccurrence struct {
	word       string
	occurrence uint
}

func Top10(input string) []string {
	allWords := matchWordsRegular.FindAllString(input, -1)

	uniqueWordCount := map[string]uint{}
	for _, word := range allWords {
		uniqueWordCount[strings.ToLower(word)]++
	}

	wordOccurrences := make([]wordOccurrence, 0, len(uniqueWordCount))
	for word, occurrence := range uniqueWordCount {
		wordOccurrences = append(wordOccurrences,
			wordOccurrence{
				word:       word,
				occurrence: occurrence,
			},
		)
	}

	sort.Slice(wordOccurrences, func(i, j int) bool {
		if wordOccurrences[i].occurrence == wordOccurrences[j].occurrence {
			return wordOccurrences[i].word < wordOccurrences[j].word
		}
		return wordOccurrences[i].occurrence > wordOccurrences[j].occurrence
	})

	var resultLength int
	if 10 < len(wordOccurrences) {
		resultLength = 10
	} else {
		resultLength = len(wordOccurrences)
	}

	topWords := make([]string, resultLength)
	for i := range topWords {
		topWords[i] = wordOccurrences[i].word
	}

	return topWords
}
