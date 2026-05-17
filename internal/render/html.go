package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

//go:embed assets/*
var assetsFS embed.FS

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{.Title}}</title>
{{- if .HasMermaid}}
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Caveat:wght@400..700&display=swap" rel="stylesheet">
{{- end}}
<style>{{.CSS}}</style>
<style>
:root { color-scheme: dark light; }
body.markdown-body {
  box-sizing: border-box;
  min-width: 200px;
  max-width: 980px;
  margin: 0 auto;
  padding: 45px;
}
@media (max-width: 767px) { body.markdown-body { padding: 15px; } }
{{- if .HasMermaid}}
.mermaid { display: flex; justify-content: center; margin: 1.25em 0; }
.mermaid svg {
  max-width: 100%;
  height: auto;
  font-family: "Caveat", "Virgil", cursive;
  font-size: 18px;
}
.mermaid .nodeLabel,
.mermaid .edgeLabel,
.mermaid .cluster-label,
.mermaid text { font-family: inherit; }
{{- end}}
</style>
</head>
<body class="markdown-body">
{{.Body}}
{{- if .HasMermaid}}
<script type="module">
import mermaid from "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";
const dark = window.matchMedia("(prefers-color-scheme: dark)").matches;
mermaid.initialize({
  startOnLoad: false,
  look: "handDrawn",
  theme: dark ? "dark" : "neutral",
  securityLevel: "loose",
  fontFamily: "Caveat, Virgil, cursive",
});
document.querySelectorAll("pre > code.language-mermaid").forEach(code => {
  const div = document.createElement("div");
  div.className = "mermaid";
  div.textContent = code.textContent;
  code.parentElement.replaceWith(div);
});
mermaid.run();
</script>
{{- end}}
</body>
</html>`

type htmlTemplateData struct {
	Title      string
	CSS        template.CSS
	Body       template.HTML
	HasMermaid bool
}

func ToHTML(filename, markdown string) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.DefinitionList,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github-dark"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(false),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var body bytes.Buffer
	if err := md.Convert([]byte(markdown), &body); err != nil {
		return nil, fmt.Errorf("goldmark convert: %w", err)
	}

	css, err := assetsFS.ReadFile("assets/github-markdown.css")
	if err != nil {
		return nil, fmt.Errorf("read css: %w", err)
	}

	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, htmlTemplateData{
		Title:      filepath.Base(filename),
		CSS:        template.CSS(css),
		Body:       template.HTML(body.String()),
		HasMermaid: hasMermaidBlock(markdown),
	})
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// hasMermaidBlock reports whether the markdown contains at least one fenced
// code block declared as `mermaid`. Used to gate mermaid.js + Caveat font
// loading so plain markdown previews stay fully local.
func hasMermaidBlock(markdown string) bool {
	source := []byte(markdown)
	root := goldmark.New().Parser().Parse(text.NewReader(source))
	found := false
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		cb, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}
		if string(cb.Language(source)) == "mermaid" {
			found = true
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}

func OpenInBrowser(filename, content string) error {
	var htmlBytes []byte
	var err error
	if IsMarkdown(filename) {
		htmlBytes, err = ToHTML(filename, content)
	} else {
		htmlBytes, err = SourceToHTML(filename, content)
	}
	if err != nil {
		return err
	}

	f, err := os.CreateTemp("", "miru-*.html")
	if err != nil {
		return err
	}
	if _, err := f.Write(htmlBytes); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return openCmd(f.Name())
}

func openCmd(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}
