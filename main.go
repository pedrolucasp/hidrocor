package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Document struct {
	Title   string
	Content template.HTML
}

var (
	wikiPath string
)

func requestError(msg string, w http.ResponseWriter) {
	log.Printf("500 on " + msg)
	w.WriteHeader(500)
	w.Write([]byte("500 Internal Server Error \n"))
}

func main() {
	flag.StringVar(&wikiPath, "wiki", wikiPath, "Wiki path")
	flag.Parse()

	if wikiPath == "" {
		log.Fatalf("You'll need to point to the location of the wiki (--wiki)")

		return
	}

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		var (
			buf bytes.Buffer
		)

		route := chi.URLParam(r, "*")
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

		doc, err := Lookup(wikiPath, route)

		if err != nil {
			requestError(fmt.Sprintf("Couldn't lookup file: %v", err), w)

			return
		}

		document := &Document{Title: route}

		if err := md.Convert(doc, &buf); err != nil {
			requestError(fmt.Sprintf("Couldn't convert markdown due: %v", err), w)

			return
		}

		document.Content = template.HTML(buf.String())
		parsedTemplate, _ := template.ParseFiles("templates/layout.html")
		err = parsedTemplate.Execute(w, document)
		if err != nil {
			requestError(fmt.Sprintf("Couldn't execute the template due %v", err), w)

			return
		}
	})

	http.ListenAndServe(":3000", router)
}
