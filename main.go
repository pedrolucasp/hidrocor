package main

import (
	"bytes"
	"embed"
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
	"github.com/yuin/goldmark-meta"
)

type Document struct {
	Title   string
	Content template.HTML
	Meta map[string]interface{}
}

var (
	wikiPath string
)

//go:embed templates/*
var templateData embed.FS

func requestError(msg string, w http.ResponseWriter) {
	log.Printf("500 on " + msg)
	w.WriteHeader(500)
	w.Write([]byte("500 Internal Server Error \n"))
}

func main() {
	flag.StringVar(&wikiPath, "wiki", wikiPath, "Wiki path")
	flag.Parse()

	// TODO: Write an actual man page
	if wikiPath == "" {
		log.Fatalf("You'll need to point to the location of the wiki (--wiki)")

		return
	}

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	parsedTemplate, _ := template.ParseFS(templateData, "templates/layout.html")

	// TODO: Make this configurable
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
				meta.Meta,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
		)

		doc, err := Lookup(wikiPath, route)

		if err != nil {
			requestError(fmt.Sprintf("Couldn't lookup file: %v", err), w)

			return
		}

		document := &Document{}
		context := parser.NewContext()

		if err := md.Convert(doc, &buf, parser.WithContext(context)); err != nil {
			requestError(fmt.Sprintf("Couldn't convert markdown due: %v", err), w)

			return
		}

		document.Content = template.HTML(buf.String())
		document.Meta = meta.Get(context)

		title, ok := document.Meta["title"].(string)
		if ok {
			document.Title = title
		}

		err = parsedTemplate.Execute(w, document)
		if err != nil {
			requestError(fmt.Sprintf("Couldn't execute the template due %v", err), w)

			return
		}
	})

	// TODO: Make this configurable
	log.Printf("Alive, at :3000")
	http.ListenAndServe(":3000", router)
}
