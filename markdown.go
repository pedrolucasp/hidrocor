// This package holds the context that given a route, it'll pipe to wiki.Lookup
// and output a markdown file converted to html. We'll use it later to feed
// into our main HTML template.

package hidrocor

import (
	"context"
	"os"
	"net/http"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type contextKey struct {
	name string
}

var mdCtxKey = &contextKey{"markdown"}

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			md := goldmark.New(
				goldmark.WithExtensions(
					extension.GFM,
					extension.Linkify,
					extension.Table,
					extension.TaskList,
					extension.DefinitionList,
					extension.Strikethrough,
					extension.Footnote,
				),
				goldmark.WithParserOptions(
					parser.WithAutoHeadingID(),
				),
				goldmark.WithRendererOptions(
					html.WithHardWraps(),
				),
			)

			r = r.WithContext(Context(r.Context(), md))

			next.ServeHTTP(w, r)
		})
	}
}

func Context(ctx context.Context, md goldmark.Markdown) context.Context {
	return context.WithValue(ctx, mdCtxKey, md)
}

func ForContext(ctx context.Context) *goldmark.Markdown {
	raw, ok := ctx.Value(mdCtxKey).(*goldmark.Markdown)

	if !ok {
		panic(errors.New("Invalid markdown context"))
	}

	return raw
}
