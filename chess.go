package chess

import (
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"math/rand"
	"strings"
)

type ChessHTMLRenderer struct {
	html.Config
}

func NewChessHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &ChessHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func renderCaissaPGN(w util.BufWriter, source []byte, lines *text.Segments) {
	id := RandStringRunes(6)
	vif := fmt.Sprintf("cwvif_%s", id)
	vip := fmt.Sprintf("cwvip_%s", id)
	vfm := fmt.Sprintf("cwvfm_%s", id)
	vpg := fmt.Sprintf("cwvpg_%s", id)

	header := []string{
		"<div style='width:100%;height:auto;aspect-ratio:1.5;border:1px solid black'>",
		"<iframe name=", vif, " style='width:100%;height:100%;border:0'></iframe>",
		"<form id=", vfm, " method=post action='https://caissa.com/bin/gameview' target=", vif, ">",
		"<input id=", vip, " name=edat type=hidden>",
		"</form>",
		"</div><div id=", vpg, " style='display:none'>",
	}
	for _, s := range header {
		_, _ = w.WriteString(s)
	}

	//// TODO: of course go to escape here so we can't like e.g inject html
	for i := 0; i < lines.Len(); i++ {
		s := lines.At(i)
		_, _ = w.WriteString(strings.Trim(string(source[s.Start:s.Stop]), "\n "))
		_, _ = w.WriteString("<br>")
	}

	footer := []string{
		"</div>", "<script>",
		"document.getElementById('", vip, "').value=document.getElementById('", vpg, "').innerHTML;",
		"document.getElementById('", vfm, "').submit();",
		"</script>", "\n",
	}
	for _, s := range footer {
		_, _ = w.WriteString(s)
	}
}

func (r *ChessHTMLRenderer) renderChess(
	w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ChessBlock)
	if !entering {
		return ast.WalkContinue, nil
	}

	renderCaissaPGN(w, source, n.Lines())
	return ast.WalkContinue, nil
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *ChessHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindChessBlock, r.renderChess)
}

type chess struct {
}

var Chess = &chess{}

// Extend implements goldmark.Extender.
func (e *chess) Extend(m goldmark.Markdown) {
	// Use high (low) prio so we get to transform chess blocks
	// before code highlighter gets to them.
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{}, 50),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewChessHTMLRenderer(), 200),
	))
}
