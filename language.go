package fluxlanguage

import (
	"strings"
	"unicode"
)

type TokenType int

const (
	TokWord TokenType = iota
	TokNumber
	TokPunct
	TokSpace
	TokUnknown
)

type Token struct {
	Type TokenType
	Text string
}

type Engine struct {
	Tokens   []Token
	Patterns []string
}

func NewEngine() *Engine {
	return &Engine{}
}

func classify(r rune) TokenType {
	if unicode.IsSpace(r) {
		return TokSpace
	}
	if unicode.IsDigit(r) || (r == '.' && false) {
		return TokNumber
	}
	if unicode.IsPunct(r) {
		return TokPunct
	}
	if unicode.IsLetter(r) {
		return TokWord
	}
	return TokUnknown
}

func tokenType(s string) TokenType {
	hasLetter := false
	hasDigit := false
	hasPunct := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
		} else if unicode.IsDigit(r) {
			hasDigit = true
		} else if unicode.IsPunct(r) {
			hasPunct = true
		}
	}
	if hasPunct {
		return TokPunct
	}
	if hasDigit && !hasLetter {
		return TokNumber
	}
	if hasLetter {
		return TokWord
	}
	if hasDigit {
		return TokNumber
	}
	return TokUnknown
}

func (e *Engine) Tokenize(text string) []Token {
	var tokens []Token
	runes := []rune(text)
	i := 0
	for i < len(runes) {
		ct := classify(runes[i])
		j := i + 1
		for j < len(runes) && classify(runes[j]) == ct {
			j++
		}
		s := string(runes[i:j])
		tokens = append(tokens, Token{Type: ct, Text: s})
		i = j
	}
	e.Tokens = tokens
	return tokens
}

func (e *Engine) TokenizeStrict(text string) []Token {
	all := e.Tokenize(text)
	var tokens []Token
	for _, t := range all {
		if t.Type != TokSpace {
			tokens = append(tokens, t)
		}
	}
	e.Tokens = tokens
	return tokens
}

func (e *Engine) CountWords() int {
	count := 0
	for _, t := range e.Tokens {
		if t.Type == TokWord {
			count++
		}
	}
	return count
}

func (e *Engine) Contains(word string) bool {
	lower := strings.ToLower(word)
	for _, t := range e.Tokens {
		if strings.ToLower(t.Text) == lower {
			return true
		}
	}
	return false
}

func (e *Engine) AddPattern(pattern string) {
	e.Patterns = append(e.Patterns, pattern)
}

func (e *Engine) Match(pattern string) bool {
	for _, t := range e.Tokens {
		if strings.Contains(strings.ToLower(t.Text), strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (e *Engine) MatchExact(word string) bool {
	lower := strings.ToLower(word)
	for _, t := range e.Tokens {
		if strings.ToLower(t.Text) == lower {
			return true
		}
	}
	return false
}

func (e *Engine) WordFrequency(word string) int {
	count := 0
	lower := strings.ToLower(word)
	for _, t := range e.Tokens {
		if strings.ToLower(t.Text) == lower {
			count++
		}
	}
	return count
}

func (e *Engine) MostCommon() string {
	freq := make(map[string]int)
	for _, t := range e.Tokens {
		if t.Type == TokWord {
			freq[strings.ToLower(t.Text)]++
		}
	}
	if len(freq) == 0 {
		return ""
	}
	best := ""
	max := 0
	for w, c := range freq {
		if c > max {
			max = c
			best = w
		}
	}
	return best
}

func (e *Engine) Similarity(a, b string) int {
	a, b = strings.ToLower(a), strings.ToLower(b)
	if a == b {
		return 100
	}
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	// Levenshtein distance
	la, lb := len(a), len(b)
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+cost))
		}
	}
	maxLen := la
	if lb > maxLen {
		maxLen = lb
	}
	sim := (maxLen - d[la][lb]) * 100 / maxLen
	return sim
}

func IsNumber(text string) bool {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return false
	}
	hasDigit := false
	for _, r := range text {
		if unicode.IsDigit(r) {
			hasDigit = true
		} else if r != '.' && r != '-' && r != '+' {
			return false
		}
	}
	return hasDigit
}

func (e *Engine) ExtractNumber(index int) float64 {
	count := 0
	for _, t := range e.Tokens {
		if t.Type == TokNumber {
			if count == index {
				var f float64
				for _, r := range t.Text {
					if r >= '0' && r <= '9' {
						f = f*10 + float64(r-'0')
					}
				}
				return f
			}
			count++
		}
	}
	return 0
}
