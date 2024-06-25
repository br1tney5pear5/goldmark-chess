# goldmark-chess
=========================

This is an extension to [goldmark](http://github.com/yuin/goldmark) markdown that
that renders interactive previews for chess [PGN](https://en.wikipedia.org/wiki/Portable_Game_Notation)
embedded in a fenced code blocks.

Installation
--------------------
```
go get github.com/br1tney5pear5/goldmark-chess
```

Future work
--------------------

- I think we don't sanitise pgn atm so it'll be just pasted verbatim in
  the embed renderer's div. This is an issue bc I think goldmark has
  a safe mode where it'll refuse to render any raw html and this could
  be used to bypass it.

- Currently this extension renders pgn using [caissa](https://caissa.com/)
  embeded viewer. I'd be cool to have many different providers to choose
  from.

Hugo Integration
--------------------

```
diff --git a/markup/goldmark/convert.go b/markup/goldmark/convert.go
index 7c00433d5..037650545 100644
--- a/markup/goldmark/convert.go
+++ b/markup/goldmark/convert.go
@@ -39,6 +39,8 @@ import (
 
 	"github.com/gohugoio/hugo/markup/converter"
 	"github.com/gohugoio/hugo/markup/tableofcontents"
+
+	chess "github.com/br1tney5pear5/goldmark-chess"
 )
 
 const (
@@ -220,6 +222,10 @@ func newMarkdown(pcfg converter.ProviderConfig) goldmark.Markdown {
 		extensions = append(extensions, attributes.New())
 	}
 
+	if cfg.Extensions.Chess.Enable {
+		extensions = append(extensions, chess.Chess)
+        }
+
 	md := goldmark.New(
 		goldmark.WithExtensions(
 			extensions...,
diff --git a/markup/goldmark/goldmark_config/config.go b/markup/goldmark/goldmark_config/config.go
index 620475c48..e249f136c 100644
--- a/markup/goldmark/goldmark_config/config.go
+++ b/markup/goldmark/goldmark_config/config.go
@@ -70,6 +70,9 @@ var Default = Config{
 				Block:  [][]string{},
 			},
 		},
+		Chess: ChessConfig{
+			Enable: false,
+		},
 	},
 	Renderer: Renderer{
 		Unsafe: false,
@@ -128,6 +131,7 @@ type Extensions struct {
 	DefinitionList bool
 	Extras         Extras
 	Passthrough    Passthrough
+	Chess          ChessConfig
 
 	// GitHub flavored markdown
 	Table           bool
@@ -200,6 +204,11 @@ type Passthrough struct {
 	Delimiters DelimitersConfig
 }
 
+type ChessConfig struct {
+	Enable bool
+	// TODO: Add PGN display provider here
+}
+
 type DelimitersConfig struct {
 	// The delimiters to use for inline passthroughs. Each entry in the list
 	// is a size-2 list of strings, where the first string is the opening delimiter
```
Then set `markup.goldmark.extensions.chess.enable` to `true` in your config.
