package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

//go:embed index.html
var indexHtml string // index.html served on "/"

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		panic(err)
	}
}

// app is the main component of our application.
type app struct {
	weaver.Implements[weaver.Main]
	searcher weaver.Ref[Searcher]
	lis      weaver.Listener `weaver:"emojis"`
}

// run implements the application main.
func run(ctx context.Context, a *app) error {
	a.Logger(ctx).Info("emojis listener available.", "addr", a.lis)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if _, err := fmt.Fprint(w, indexHtml); err != nil {
			a.Logger(ctx).Error("error writing index.html", "err", err)
		}
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		a.handleSearch(a.searcher.Get().Search, w, r)
	})
	http.HandleFunc("/search_chatgpt", func(w http.ResponseWriter, r *http.Request) {
		a.handleSearch(a.searcher.Get().SearchChatGPT, w, r)
	})
	return http.Serve(a.lis, nil)
}

// handleSearch handles HTTP requests to the /search?q=<query> and
// /search_chatgpt?q=<query> endpoints.
func (a *app) handleSearch(search func(context.Context, string) ([]string, error), w http.ResponseWriter, r *http.Request) {
	// Search for the list of matching emojis.
	query := r.URL.Query().Get("q")
	emojis, err := search(r.Context(), query)
	if err != nil {
		a.Logger(r.Context()).Error("error getting search results", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSON serialize the results.
	bytes, err := json.Marshal(emojis)
	if err != nil {
		a.Logger(r.Context()).Error("error marshaling search results", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprintln(w, string(bytes)); err != nil {
		a.Logger(r.Context()).Error("error writing search results", "err", err)
	}
}
