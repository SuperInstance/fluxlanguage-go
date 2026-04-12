package fluxlanguage

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
	if len(e.Tokens) != 0 {
		t.Fatal("expected empty tokens")
	}
}

func TestTokenize(t *testing.T) {
	e := NewEngine()
	tokens := e.Tokenize("hello world")
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TokWord || tokens[0].Text != "hello" {
		t.Fatalf("token 0: expected Word 'hello', got %v %q", tokens[0].Type, tokens[0].Text)
	}
	if tokens[1].Type != TokSpace {
		t.Fatalf("token 1: expected Space, got %v", tokens[1].Type)
	}
	if tokens[2].Type != TokWord || tokens[2].Text != "world" {
		t.Fatalf("token 2: expected Word 'world', got %v %q", tokens[2].Type, tokens[2].Text)
	}
}

func TestTokenizeNumbers(t *testing.T) {
	e := NewEngine()
	tokens := e.Tokenize("abc 123 def")
	if len(tokens) != 5 {
		t.Fatalf("expected 5 tokens, got %d", len(tokens))
	}
	if tokens[2].Type != TokNumber {
		t.Fatalf("expected TokNumber, got %v", tokens[2].Type)
	}
}

func TestTokenizePunctuation(t *testing.T) {
	e := NewEngine()
	tokens := e.Tokenize("hello, world!")
	if len(tokens) != 5 { // hello , space world !
		t.Fatalf("expected 5 tokens, got %d: %v", len(tokens), tokens)
	}
	if tokens[1].Type != TokPunct {
		t.Fatalf("expected TokPunct for ',', got %v", tokens[1].Type)
	}
}

func TestTokenizeEmpty(t *testing.T) {
	e := NewEngine()
	tokens := e.Tokenize("")
	if len(tokens) != 0 {
		t.Fatalf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestTokenizeStrict(t *testing.T) {
	e := NewEngine()
	tokens := e.TokenizeStrict("hello world foo")
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	for _, tok := range tokens {
		if tok.Type == TokSpace {
			t.Fatal("found space token in strict mode")
		}
	}
}

func TestCountWords(t *testing.T) {
	e := NewEngine()
	e.Tokenize("hello world foo")
	if e.CountWords() != 3 {
		t.Fatalf("expected 3 words, got %d", e.CountWords())
	}
}

func TestCountWordsEmpty(t *testing.T) {
	e := NewEngine()
	e.Tokenize("")
	if e.CountWords() != 0 {
		t.Fatalf("expected 0 words, got %d", e.CountWords())
	}
}

func TestContains(t *testing.T) {
	e := NewEngine()
	e.Tokenize("hello world")
	if !e.Contains("hello") {
		t.Fatal("expected to contain 'hello'")
	}
	if !e.Contains("HELLO") {
		t.Fatal("expected case-insensitive match")
	}
	if e.Contains("foo") {
		t.Fatal("expected not to contain 'foo'")
	}
}

func TestAddPattern(t *testing.T) {
	e := NewEngine()
	e.AddPattern("test")
	if len(e.Patterns) != 1 || e.Patterns[0] != "test" {
		t.Fatal("pattern not added")
	}
}

func TestMatch(t *testing.T) {
	e := NewEngine()
	e.Tokenize("hello wonderful world")
	if !e.Match("wonder") {
		t.Fatal("expected match for 'wonder'")
	}
	if e.Match("xyz") {
		t.Fatal("expected no match for 'xyz'")
	}
}

func TestMatchExact(t *testing.T) {
	e := NewEngine()
	e.Tokenize("hello world")
	if !e.MatchExact("hello") {
		t.Fatal("expected exact match for 'hello'")
	}
	if e.MatchExact("hell") {
		t.Fatal("expected no exact match for 'hell'")
	}
}

func TestWordFrequency(t *testing.T) {
	e := NewEngine()
	e.Tokenize("hello hello world hello")
	if e.WordFrequency("hello") != 3 {
		t.Fatalf("expected frequency 3, got %d", e.WordFrequency("hello"))
	}
	if e.WordFrequency("world") != 1 {
		t.Fatalf("expected frequency 1, got %d", e.WordFrequency("world"))
	}
	if e.WordFrequency("foo") != 0 {
		t.Fatalf("expected frequency 0, got %d", e.WordFrequency("foo"))
	}
}

func TestMostCommon(t *testing.T) {
	e := NewEngine()
	e.Tokenize("a b a c a")
	if e.MostCommon() != "a" {
		t.Fatalf("expected 'a', got %q", e.MostCommon())
	}
}

func TestMostCommonEmpty(t *testing.T) {
	e := NewEngine()
	e.Tokenize("123 456")
	if e.MostCommon() != "" {
		t.Fatalf("expected empty string, got %q", e.MostCommon())
	}
}

func TestSimilarity(t *testing.T) {
	e := NewEngine()
	if e.Similarity("hello", "hello") != 100 {
		t.Fatal("identical strings should be 100")
	}
	if e.Similarity("abc", "xyz") < 0 || e.Similarity("abc", "xyz") > 100 {
		t.Fatal("similarity out of range")
	}
	if e.Similarity("kitten", "sitting") == 0 {
		t.Fatal("similar strings should not be 0")
	}
}

func TestIsNumber(t *testing.T) {
	if !IsNumber("123") {
		t.Fatal("123 should be a number")
	}
	if !IsNumber("3.14") {
		t.Fatal("3.14 should be a number")
	}
	if IsNumber("abc") {
		t.Fatal("abc should not be a number")
	}
	if IsNumber("") {
		t.Fatal("empty string should not be a number")
	}
}

func TestExtractNumber(t *testing.T) {
	e := NewEngine()
	e.Tokenize("foo 42 bar 7 baz")
	if e.ExtractNumber(0) != 42 {
		t.Fatalf("expected 42, got %v", e.ExtractNumber(0))
	}
	if e.ExtractNumber(1) != 7 {
		t.Fatalf("expected 7, got %v", e.ExtractNumber(1))
	}
}
