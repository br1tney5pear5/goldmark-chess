package chess

import "github.com/yuin/goldmark/ast"

var KindChessBlock = ast.NewNodeKind("ChessBlock")

type ChessBlock struct {
	ast.BaseBlock
}

// IsRaw reports that this block should be rendered as-is.
// TODO: Is it raw or not?
func (*ChessBlock) IsRaw() bool { return true }

// Kind reports that this is a MermaidBlock.
func (*ChessBlock) Kind() ast.NodeKind { return KindChessBlock }

// Dump dumps the contents of this block to stdout.
func (b *ChessBlock) Dump(src []byte, level int) {
	ast.DumpHelper(b, src, level, nil, nil)
}
