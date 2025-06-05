package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/BrunodsLilly/Summarizer/cmd/web/templates"
	"github.com/BrunodsLilly/Summarizer/pkg/core"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/summarize", summarizeHandler)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	component := templates.Index()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
		return
	}
}

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	if strings.TrimSpace(url) == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	// Simple summarization logic (placeholder)
	summary := generateSummary(url)

	component := templates.SummaryResult(summary)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template rendering error: %v", err)
		return
	}
}

func generateSummary(url string) string {
	res, err := core.SummarizeURL(url)
	if err != nil {
		log.Printf("Error summarizing URL: %v", err)
		// return "Error generating summary"
		return fmt.Sprintf("Error generating summary for URL: %s\n%s", url, err.Error())
	}
	return markdownToHTML(res)
}

func markdownToHTML(markdown string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		log.Printf("Error converting markdown to HTML: %v", err)
		return markdown // fallback to original text
	}
	
	return buf.String()
}

