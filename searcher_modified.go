
package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/metrics"
	"github.com/rivo/uniseg"
	"golang.org/x/exp/slices"
)

// Two counters that count cache hits and misses in the Search method.
var (
	cacheHits = metrics.NewCounter(
		"search_cache_hits",
		"Number of Search cache hits",
	)
	cacheMisses = metrics.NewCounter(
		"search_cache_misses",
		"Number of Search cache misses",
	)
)

// Searcher is an emoji search engine component.
type Searcher interface {
	// Search returns the set of emojis that match the provided query.
	Search(ctx context.Context, query string) ([]string, error)

	// SearchChatGPT returns the set of emojis that ChatGPT thinks match the
	// provided query.
	SearchChatGPT(ctx context.Context, query string) ([]string, error)
}

// searcher is the implementation of the Searcher component.
type searcher struct {
	weaver.Implements[Searcher]
	cache   weaver.Ref[Cache]
	chatgpt weaver.Ref[ChatGPT]
}

// Search method implementation that retrieves relevant emojis from emoji.go.
func (s *searcher) Search(ctx context.Context, query string) ([]string, error) {
	// Example logic to search relevant emojis from the emoji stack in emoji.go.
	emojiStack := s.getEmojiStack()

	var matchedEmojis []string
	for _, emoji := range emojiStack {
		if strings.Contains(strings.ToLower(query), strings.ToLower(emoji.Keyword)) {
			matchedEmojis = append(matchedEmojis, emoji.Symbol)
		}
	}

	if len(matchedEmojis) == 0 {
		cacheMisses.Add(ctx, 1)
	} else {
		cacheHits.Add(ctx, 1)
	}

	return matchedEmojis, nil
}

// getEmojiStack retrieves the list of all emojis and their associated keywords from emoji.go
func (s *searcher) getEmojiStack() []Emoji {
	// This function should interact with emoji.go to retrieve all emojis
	// For now, let's assume emoji.go has a simple struct with Emoji and Keyword
	return []Emoji{
		{Symbol: "😊", Keyword: "happy"},
		{Symbol: "🌍", Keyword: "world"},
		// Add more emojis and keywords as needed
	}
}

// SearchChatGPT returns emojis that ChatGPT suggests (implementation left as-is)
func (s *searcher) SearchChatGPT(ctx context.Context, query string) ([]string, error) {
	// Assuming there is existing logic here to handle this functionality
	return nil, fmt.Errorf("Not implemented")
}

// Emoji represents a simple emoji struct with symbol and keyword
type Emoji struct {
	Symbol  string
	Keyword string
}
