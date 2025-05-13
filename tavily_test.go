package tavily_test

import (
	"os"
	"testing"

	tavily "github.com/aiomni/tavily-go"
)

func TestTavilySearch(t *testing.T) {
	c := tavily.NewTavilyClient(os.Getenv("TAVILY_API_KEY"))

	resp, err := c.Search(t.Context(), &tavily.TavilySearchRequest{
		Query:                    "who is Leo Messi?",
		Topic:                    "general",
		SearchDepth:              "basic",
		ChunksPerSource:          3,
		MaxResults:               5,
		Days:                     7,
		IncludeAnswer:            "basic",
		IncludeRawContent:        false,
		IncludeImages:            false,
		IncludeImageDescriptions: false,
	})

	if err != nil {
		t.Fatalf("Tavily Search Fail: %e", err)
	}

	t.Logf("Tavily Search Resp: %+v", resp)
}

func TestTavilyExtract(t *testing.T) {
	c := tavily.NewTavilyClient(os.Getenv("TAVILY_API_KEY"))

	resp, err := c.Extract(t.Context(), &tavily.TavilyExtractRequest{
		URLs:          []string{"https://foreverz.cn/qwik"},
		IncludeImages: false,
		ExtractDepth:  "basic",
	})

	if err != nil {
		t.Fatalf("Tavily Extract Fail: %e", err)
	}

	t.Logf("Tavily Extract Resp: %+v", resp)
}

func TestTavilyCrawl(t *testing.T) {
	c := tavily.NewTavilyClient(os.Getenv("TAVILY_API_KEY"))

	resp, err := c.Crawl(t.Context(), &tavily.TavilyCrawlRequest{
		URL:           "https://foreverz.cn/qwik",
		MaxDepth:      1,
		MaxBreadth:    20,
		Limit:         50,
		AllowExternal: false,
		IncludeImages: false,
		ExtractDepth:  "basic",
	})

	if err != nil {
		t.Fatalf("Tavily Crawl Fail: %e", err)
	}

	t.Logf("Tavily Crawl Resp: %+v", resp)
}

func TestTavilyMap(t *testing.T) {
	c := tavily.NewTavilyClient(os.Getenv("TAVILY_API_KEY"))

	resp, err := c.Map(t.Context(), &tavily.TavilyMapRequest{
		URL:           "https://foreverz.cn",
		MaxDepth:      1,
		MaxBreadth:    20,
		Limit:         50,
		AllowExternal: false,
	})

	if err != nil {
		t.Fatalf("Tavily Map Fail: %e", err)
	}

	t.Logf("Tavily Map Resp: %+v", resp)
}
