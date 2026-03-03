package lorem

import (
	"math/rand"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	g := New(10)
	if g.Length != 10 {
		t.Errorf("Expected Length 10. Got %d", g.Length)
	}

	g2 := New(-5)
	if g2.Length != 1 {
		t.Errorf("Expected bounds to be clamped to 1. Got %d", g2.Length)
	}
}

func TestGenerateTokensLength(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	g := New(5).WithRand(r)

	for i := 0; i < 100; i++ {
		tokens := g.GenerateTokens()
		if len(tokens) != 5 {
			t.Errorf("Generated %d tokens, which is not equal to length 5", len(tokens))
		}
	}
}

func TestGenerateTokensFormat(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	g := New(3).WithRand(r)

	tokens := g.GenerateTokens()
	if len(tokens) != 3 {
		t.Fatalf("Expected 3 tokens, got %d", len(tokens))
	}

	if strings.HasPrefix(tokens[0], " ") {
		t.Errorf("Expected first token not to have a leading space, got %q", tokens[0])
	}

	for i := 1; i < len(tokens); i++ {
		if !strings.HasPrefix(tokens[i], " ") {
			t.Errorf("Expected token %d to have a leading space, got %q", i, tokens[i])
		}
	}
}

func TestGenerateString(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	g := New(3).WithRand(r)

	str := g.GenerateString()

	if str == "" {
		t.Errorf("Expected a non-empty string, got empty")
	}

	spaceCount := strings.Count(str, " ")
	if spaceCount != 2 {
		t.Errorf("Expected exactly 2 spaces in a 3-word string, got %d. String: %q", spaceCount, str)
	}
}
