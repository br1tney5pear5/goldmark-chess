package chess

import (
	"bytes"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Transformer struct {
}

var _chess = []byte("chess")

// Transform transforms the provided Markdown AST.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	var (
		chessBlocks []*ast.FencedCodeBlock
	)

	// Collect all blocks to be replaced without modifying the tree.
	_ = ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}

		cb, ok := node.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		lang := cb.Language(reader.Source())

		if !bytes.Equal(lang, _chess) {
			return ast.WalkContinue, nil
		}

		chessBlocks = append(chessBlocks, cb)
		return ast.WalkContinue, nil
	})

	if len(chessBlocks) == 0 {
		return
	}

	for _, cb := range chessBlocks {
		b := new(ChessBlock)
		b.SetLines(cb.Lines())

		parent := cb.Parent()
		if parent != nil {
			parent.ReplaceChild(parent, cb, b)
		}
	}
}
