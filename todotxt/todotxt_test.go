package todotxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseList(t *testing.T) {
	fixtures := []struct {
		lines    string
		expected []string
	}{
		{"Single Line", []string{"Single Line"}},
		{"First Line\nSecond Line", []string{"First Line", "Second Line"}},
	}

	for _, fixture := range fixtures {
		actual := parseList(fixture.lines)
		if !assert.Equal(t, fixture.expected, actual) {
			t.Errorf("error: parsed lists don't match: expected %v, got %v", fixture.expected, actual)
		}
	}
}
