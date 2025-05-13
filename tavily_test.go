package tavily_test

import (
	"os"
	"testing"

	tavily "github.com/aiomni/tavily-go"
)

func TestTavilySearch(t *testing.T) {
	c := tavily.NewTavilyClient(os.Getenv("TAVILY_API_KEY"))

	resp, err := c.Search(t.Context(), &tavily.TavilySearchRequest{
		Query: "What is GitHub?",
	})

	if err != nil {
		t.Fatal("TavilyClient Search", err)
	}

	t.Logf("Search Resp: %+v", resp)
}

func TestTavilyExtract(t *testing.T) {
	c := tavily.NewTavilyClient(os.Getenv("TAVILY_API_KEY"))

	resp, err := c.Extract(t.Context(), &tavily.TavilyExtractRequest{
		URLs:          []string{"https://foreverz.cn"},
		IncludeImages: true,
		ExtractDepth:  "advanced",
	})
	
	if err != nil {
		t.Fatal("TavilyClient Extract", err)
	}

	t.Logf("Extract Resp: %+v", resp)
}