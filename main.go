package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jamesbehr/docshtest/extract"
	"github.com/jamesbehr/docshtest/shelltest"
	"mvdan.cc/sh/v3/interp"
)

type stringSlice []string

func (i stringSlice) String() string {
	return strings.Join(i, ",")
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var runHighlightedCodeFences stringSlice
var runCodeFences bool
var runCodeBlocks bool

func main() {
	flag.Var(&runHighlightedCodeFences, "run-highlighted-code-fences", "Run code fences (three backticks) with a specified syntax highlighting language. This flag can be provided multiple times to extract multiple languages")
	flag.BoolVar(&runCodeFences, "run-code-fences", false, "Run code fences (three backticks) with a no syntax highlighting language specified")
	flag.BoolVar(&runCodeBlocks, "run-code-blocks", false, "Run code blocks (code indented by four spaces).")
	flag.Parse()

	options := extract.Options{
		ExtractHighlightedCodeFences: runHighlightedCodeFences,
		ExtractCodeFences:            runCodeFences,
		ExtractCodeBlocks:            runCodeBlocks,
	}

	var blocks []extract.Block

	for _, name := range flag.Args() {
		b, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}

		blocks = append(blocks, extract.Bytes(name, b, options)...)
	}

	var tests []*shelltest.Test

	for _, block := range blocks {
		t, err := shelltest.ParseTests(block.Filename, block.Content, block.Line)
		if err != nil {
			log.Fatal(err)
		}

		tests = append(tests, t...)
	}

	output := bytes.NewBuffer([]byte{}) // combined stdout and stderr
	input := bytes.NewBuffer([]byte{})  // stdin will always be empty

	runner, err := interp.New(interp.StdIO(input, output, output))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	for _, test := range tests {
		if err := test.Run(ctx, runner); err != nil {
			log.Println(string(output.Bytes()))
			log.Fatal(err)
		}

		actualOutput := strings.TrimSpace(string(output.Bytes()))
		expectedOutput := strings.TrimSpace(test.ExpectedOutput)

		// The expected output doesn't have to contain all the output from the
		// command, but any output provided must match what comes out the
		// command in the same order that it is provided.
		if !strings.HasPrefix(actualOutput, expectedOutput) {
			log.Printf("Test failed: %s:%d", test.Filename, test.Line)
			log.Printf("expected: %q", expectedOutput)
			log.Printf("got: %q", actualOutput)
			os.Exit(1)
		}

		// Clear the buffer for the next test
		output.Reset()
	}
}
