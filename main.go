package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func requestError(msg string, w http.ResponseWriter) {
	log.Printf("500 on " + msg)
	w.WriteHeader(500)
	w.Write([]byte("500 Internal Server Error \n"))
}

func main() {
	wikiPath := "/home/eletrotupi/sources/cirandas.net-docs"

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

	router := chi.NewRouter()
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		var (
			buf    bytes.Buffer
			source []byte
		)

		route := chi.URLParam(r, "*")

		// TODO: Move this into a context and pass this along with the main
		// template
		if route == "" {
			route = "README.md"
		}

		routePath := path.Join(wikiPath, route)
		fileInfo, err := os.Stat(routePath)

		if err != nil {
			requestError("On Stating File", w)

			return
		}

		if fileInfo.IsDir() {
			files, err := os.ReadDir(routePath)
			if err != nil {
				requestError("On Reading Dir", w)

				return
			}

			for _, file := range files {
				log.Printf("fileName: ", file.Name())

				switch file.Name() {
				case
					"index.md",
					"INDEX.md",
					"readme.md",
					"README.md":
					source, err = os.ReadFile(path.Join(routePath, file.Name()))

					if err != nil {
						requestError("On Opening file inside folder", w)

						return
					}
				}
			}
		} else {
			log.Printf("Not dir")
			source, err = os.ReadFile(routePath)

			if err != nil {
				requestError("On Opening file", w)

				return
			}
		}

		if err := md.Convert(source, &buf); err != nil {
			log.Printf("500 on Converting Markdown")
			w.WriteHeader(500)
			w.Write([]byte("500 Internal Server Error \n"))

			return
		}

		// XXX: How to deal with images? Also move the charset to the HTML
		// template
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write(buf.Bytes())
	})

	http.ListenAndServe(":3000", router)
}
