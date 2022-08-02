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

func (t Test) Empty() bool {
	return t.Command == ""
}

var ErrInvalidTest = errors.New("shelltest: command output appeared before command")

const prompt = "$ "

func ParseTests(filename, source string, lineOffset int) ([]*Test, error) {
	lines := strings.Split(source, "\n")

	var current Test
	var tests []*Test

	for i, line := range lines {
		if strings.HasPrefix(line, prompt) {
			if !current.Empty() {
				tests = append(tests, &current)
			}

			current.Command = strings.TrimPrefix(line, prompt)
			current.Filename = filename
			current.Line = lineOffset + i
		} else {
			if current.Empty() {
				return nil, ErrInvalidTest
			}

			if current.ExpectedOutput == "" {
				current.ExpectedOutput = line
			} else {
				current.ExpectedOutput += "\n" + line
			}
		}
	}

	if !current.Empty() {
		tests = append(tests, &current)
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
