package shelltest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTests(t *testing.T) {
	tests := []struct {
		Error      error
		Expected   []*Test
		Filename   string
		Source     string
		LineOffset int
	}{
		{
			Filename: "invalid.md",
			Source: `find foo/bar
foo/bar
foo/bar/baz
foo/bar/quux`,
			Error: ErrInvalidTest,
		},
		{
			Filename:   "single_command.md",
			LineOffset: 10,
			Source: `$ find foo/bar
foo/bar
foo/bar/baz
foo/bar/quux`,
			Expected: []*Test{
				{
					Filename:       "single_command.md",
					Line:           10,
					Command:        "find foo/bar",
					ExpectedOutput: "foo/bar\nfoo/bar/baz\nfoo/bar/quux",
				},
			},
		},
		{
			Filename:   "multiple_commands.md",
			LineOffset: 7,
			Source: `$ foo
bar
baz

$ another
test

$ one more`,
			Expected: []*Test{
				{
					Filename:       "multiple_commands.md",
					Line:           7,
					Command:        "foo",
					ExpectedOutput: "bar\nbaz\n",
				},
				{
					Filename:       "multiple_commands.md",
					Line:           11,
					Command:        "another",
					ExpectedOutput: "test\n",
				},
				{
					Filename:       "multiple_commands.md",
					Line:           14,
					Command:        "one more",
					ExpectedOutput: "",
				},
			},
		},
	}

	for _, test := range tests {
		actual, err := ParseTests(test.Filename, test.Source, test.LineOffset)
		assert.Equal(t, test.Error, err)
		assert.Equal(t, test.Expected, actual)
	}
}
