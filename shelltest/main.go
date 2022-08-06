package shelltest

import (
	"context"
	"errors"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type Test struct {
	Filename       string
	Line           int
	Command        string
	ExpectedOutput string
}

var ErrInvalidTest = errors.New("shelltest: command output appeared before command")

const prompt = "$ "

func ParseTests(filename, source string, lineOffset int) ([]*Test, error) {
	lines := strings.Split(source, "\n")

	var current *Test
	var tests []*Test

	for i, line := range lines {
		if strings.HasPrefix(line, prompt) {
			if current != nil {
				tests = append(tests, current)
				current = nil
			}

			current = &Test{
				Command:        strings.TrimPrefix(line, prompt),
				ExpectedOutput: "",
				Filename:       filename,
				Line:           lineOffset + i,
			}
		} else {
			if current == nil {
				return nil, ErrInvalidTest
			}

			if current.ExpectedOutput == "" {
				current.ExpectedOutput = line
			} else {
				current.ExpectedOutput += "\n" + line
			}
		}
	}

	if current != nil {
		tests = append(tests, current)
	}

	return tests, nil
}

func (t Test) Run(ctx context.Context, runner *interp.Runner) error {
	parser := syntax.NewParser()
	r := strings.NewReader(t.Command)

	program, err := parser.Parse(r, "")
	if err != nil {
		return err
	}

	return runner.Run(ctx, program)
}
