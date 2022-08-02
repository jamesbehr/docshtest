package extract

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Options struct {
	ExtractHighlightedCodeFences []string
	ExtractCodeFences            bool
	ExtractCodeBlocks            bool
}

type Block struct {
	Filename string
	Content  string
	Line     int
	Offset   int
}

func codeBlockContents(filename string, source []byte, node ast.Node) Block {
	var sb strings.Builder

	segments := node.Lines()

	reader := text.NewReader(source)
	reader.Advance(segments.At(0).Start)
	offset := reader.LineOffset()
	line, _ := reader.Position()

	for i := 0; i < segments.Len(); i++ {
		segment := segments.At(i)

		sb.Write(segment.Value(source))
	}

	return Block{
		Filename: filename,
		Content:  sb.String(),
		Line:     line,
		Offset:   offset,
	}
}

func Bytes(filename string, b []byte, options Options) []Block {
	reader := text.NewReader(b)

	md := goldmark.New()
	node := md.Parser().Parse(reader)

	var blocks []Block

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindFencedCodeBlock {
			block := n.(*ast.FencedCodeBlock)

			if block.Info == nil {
				if options.ExtractCodeFences {
					blocks = append(blocks, codeBlockContents(filename, b, n))
				}
			} else {
				language := string(block.Info.Text(b))

				for _, extractedLanguage := range options.ExtractHighlightedCodeFences {
					if language == extractedLanguage {
						blocks = append(blocks, codeBlockContents(filename, b, n))
					}
				}
			}
		} else if n.Kind() == ast.KindCodeBlock {
			if options.ExtractCodeBlocks {
				blocks = append(blocks, codeBlockContents(filename, b, n))
			}
		}

		return ast.WalkContinue, nil
	})

	return blocks
}
