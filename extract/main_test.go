package extract

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	f, err := os.ReadFile("testdata/markdown.md")
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("extract highlighted code fences", func(t *testing.T) {
		blocks := Bytes(f, Options{
			ExtractHighlightedCodeFences: []string{"foo"},
		})

		assert.Equal(t, []Block{{"1\n2\n3\n", 9, 0}, {"13\n", 65, 0}}, blocks)
	})

	t.Run("extract multiple highlighted code fences", func(t *testing.T) {
		blocks := Bytes(f, Options{
			ExtractHighlightedCodeFences: []string{"foo", "bar"},
		})

		assert.Equal(t, []Block{
			{"1\n2\n3\n", 9, 0},
			{"4\n5\n6\n", 15, 0},
			{"13\n", 65, 0},
			{"14\n", 69, 0},
		}, blocks)
	})

	t.Run("extract unhighlighted code fences", func(t *testing.T) {
		blocks := Bytes(f, Options{
			ExtractCodeFences: true,
		})

		assert.Equal(t, []Block{{"10\n11\n12\n", 49, 0}, {"16\n", 86, 0}}, blocks)
	})

	t.Run("extract code blocks", func(t *testing.T) {
		blocks := Bytes(f, Options{
			ExtractCodeBlocks: true,
		})

		assert.Equal(t, []Block{{"7\n8\n9\n", 31, 4}, {"15\n", 83, 4}}, blocks)
	})
}
