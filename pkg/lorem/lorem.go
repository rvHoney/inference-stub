package lorem

import (
	"math/rand"
	"strings"
)

var defaultTokens = []string{
	"lorem", "ipsum", "dolor", "sit", "amet,", "consectetur", "adipiscing", "elit.",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua.", "ut", "enim", "ad", "minim", "veniam,", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea", "commodo",
	"consequat.", "duis", "aute", "irure", "dolor", "in", "reprehenderit", "in",
	"voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla",
	"pariatur.", "excepteur", "sint", "occaecat", "cupidatat", "non", "proident,",
	"sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id", "est", "laborum.",
}

// Generator provides randomized Lorem Ipsum tokens.
type Generator struct {
	Length  int
	randSrc *rand.Rand
}

// New creates a new Generator with the given length.
func New(length int) *Generator {
	if length < 1 {
		length = 1
	}

	return &Generator{
		Length: length,
	}
}

// WithRand injects a custom *rand.Rand for deterministic testing.
func (g *Generator) WithRand(r *rand.Rand) *Generator {
	g.randSrc = r
	return g
}

// randIntn returns a random integer in [0, n).
func (g *Generator) randIntn(n int) int {
	if g.randSrc != nil {
		return g.randSrc.Intn(n)
	}
	return rand.Intn(n)
}

// GenerateTokens creates a slice of randomized Lorem Ipsum tokens.
func (g *Generator) GenerateTokens() []string {
	tokens := make([]string, g.Length)
	for i := 0; i < g.Length; i++ {
		word := defaultTokens[g.randIntn(len(defaultTokens))]

		// Add leading space to all tokens
		if i > 0 {
			tokens[i] = " " + word
		} else {
			// Title case the first word
			if len(word) > 0 {
				tokens[i] = strings.ToUpper(word[:1]) + word[1:]
			}
		}
	}
	return tokens
}

// GenerateString creates a single concatenated string of the randomized tokens.
func (g *Generator) GenerateString() string {
	tokens := g.GenerateTokens()
	return strings.Join(tokens, "")
}
