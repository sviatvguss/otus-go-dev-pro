package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type freqDict struct {
	Word  string
	Count int
}

const LEN = 10

func Top10(input string) []string {
	dict := make(map[string]int)
	for _, w := range strings.Fields(input) {
		dict[w]++
	}

	words := make([]freqDict, 0, len(dict))
	for wrd, cnt := range dict {
		words = append(words, freqDict{Word: wrd, Count: cnt})
	}
	sort.Slice(words, func(i, j int) bool {
		if words[i].Count == words[j].Count {
			return words[i].Word < words[j].Word
		}
		return words[i].Count > words[j].Count
	})

	end := LEN
	if len(words) < end {
		end = len(words)
	}
	result := make([]string, end)
	for i, w := range words {
		if i == end {
			break
		}
		result[i] = w.Word
	}

	return result
}
