# tavily-go

A Go client library for the [Tavily API](https://tavily.com), providing programmatic access to Tavily's search, extraction, crawling, and mapping functionalities.

## Installation

```bash
go get github.com/aiomni/tavily-go
```

## Authentication

To use the Tavily API, you'll need an API key. Get yours by signing up at [Tavily](https://tavily.com).

## Usage

### Initializing the Client

```go
import "github.com/aiomni/tavily-go/tavily"

// Create a client with default settings
client := tavily.NewTavilyClient("your-api-key")

// Or customize the client with options
client := tavily.NewTavilyClient(
    "your-api-key",
    tavily.TavilyClientWithBaseURL("https://custom-api.tavily.com"),
    tavily.TavilyClientWithHttpClient(customHttpClient),
)
```

### Search API

Perform web searches with various filtering options:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/aiomni/tavily-go/tavily"
)

func main() {
    client := tavily.NewTavilyClient("your-api-key")
    
    searchRequest := &tavily.TavilySearchRequest{
        Query:           "Who is Leo Messi?",
        SearchDepth:     "advanced",
        MaxResults:      5,
        IncludeAnswer:   "advanced",
        IncludeImages:   true,
    }
    
    ctx := context.Background()
    response, err := client.Search(ctx, searchRequest)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Answer: %s\n\n", response.Answer)
    
    fmt.Println("Search Results:")
    for i, result := range response.Results {
        fmt.Printf("%d. %s\n   URL: %s\n\n", i+1, result.Title, result.URL)
    }
    
    if len(response.Images) > 0 {
        fmt.Println("Images:")
        for i, image := range response.Images {
            fmt.Printf("%d. %s\n", i+1, image.URL)
            if image.Description != "" {
                fmt.Printf("   Description: %s\n", image.Description)
            }
        }
    }
}
```

### Extract API

Extract content from specific URLs:

```go
extractRequest := &tavily.TavilyExtractRequest{
    URLs:          []string{"https://tavily.com/blog/post1", "https://tavily.com/blog/post2"},
    IncludeImages: true,
    ExtractDepth:  "advanced",
}

response, err := client.Extract(ctx, extractRequest)
if err != nil {
    log.Fatal(err)
}

for _, result := range response.Results {
    fmt.Printf("URL: %s\n", result.URL)
    fmt.Printf("Content: %s\n\n", result.RawContent[:100]) // First 100 chars
    
    if len(result.Images) > 0 {
        fmt.Println("Images found:", len(result.Images))
    }
}

if len(response.FailedResults) > 0 {
    fmt.Println("Failed URLs:")
    for _, fail := range response.FailedResults {
        fmt.Printf("- %s: %s\n", fail.URL, fail.Error)
    }
}
```

### Crawl API (Beta)

Recursively crawl a website starting from a base URL:

```go
crawlRequest := &tavily.TavilyCrawlRequest{
    URL:           "docs.tavily.com",
    MaxDepth:      2,
    MaxBreadth:    10,
    Limit:         30,
    Instructions:  "Focus on API documentation",
    SelectPaths:   []string{"/api/.*", "/docs/.*"},
    AllowExternal: false,
    ExtractDepth:  "basic",
}

response, err := client.Crawl(ctx, crawlRequest)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Crawled base URL: %s\n", response.BaseURL)
fmt.Printf("Found %d pages\n\n", len(response.Results))

for i, result := range response.Results {
    fmt.Printf("%d. %s\n", i+1, result.Title)
    fmt.Printf("   URL: %s\n", result.URL)
    fmt.Printf("   Found %d links on this page\n\n", len(result.Links))
}
```

### Map API (Beta)

Generate a sitemap starting from a base URL:

```go
mapRequest := &tavily.TavilyMapRequest{
    URL:           "tavily.com",
    MaxDepth:      3,
    MaxBreadth:    15,
    Categories:    []string{"Documentation", "Blog"},
    AllowExternal: false,
}

response, err := client.Map(ctx, mapRequest)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Mapped base URL: %s\n", response.BaseURL)
fmt.Printf("Found %d URLs\n\n", len(response.Results))

for i, url := range response.Results {
    fmt.Printf("%d. %s\n", i+1, url)
}
```

## API Reference

### Search API

The Search API allows you to search the web with various filtering options.

```go
type TavilySearchRequest struct {
    Query                    string   `json:"query"`
    Topic                    string   `json:"topic,omitempty"`                     // "news" or "general"
    SearchDepth              string   `json:"search_depth,omitempty"`              // "basic" or "advanced"
    ChunksPerSource          int      `json:"chunks_per_source,omitempty"`
    MaxResults               int      `json:"max_results,omitempty"`
    TimeRange                string   `json:"time_range,omitempty"`                // "day", "week", "month", or "year"
    Days                     int      `json:"days,omitempty"`
    IncludeAnswer            string   `json:"include_answer,omitempty"`            // "basic" or "advanced"
    IncludeRawContent        bool     `json:"include_raw_content,omitempty"`
    IncludeImages            bool     `json:"include_images,omitempty"`
    IncludeImageDescriptions bool     `json:"include_image_descriptions,omitempty"`
    IncludeDomains           []string `json:"include_domains,omitempty"`
    ExcludeDomains           []string `json:"exclude_domains,omitempty"`
}
```

### Extract API

The Extract API allows you to extract content from specific URLs.

```go
type TavilyExtractRequest struct {
    URLs          []string `json:"urls"`
    IncludeImages bool     `json:"include_images,omitempty"`
    ExtractDepth  string   `json:"extract_depth,omitempty"`      // "basic" or "advanced"
}
```

### Crawl API (Beta)

The Crawl API allows you to recursively crawl a website starting from a base URL.

```go
type TavilyCrawlRequest struct {
    URL            string   `json:"url"`
    MaxDepth       int      `json:"max_depth,omitempty"`
    MaxBreadth     int      `json:"max_breadth,omitempty"`
    Limit          int      `json:"limit,omitempty"`
    Instructions   string   `json:"instructions,omitempty"`
    SelectPaths    []string `json:"select_paths,omitempty"`
    SelectDomains  []string `json:"select_domains,omitempty"`
    ExcludePaths   []string `json:"exclude_paths,omitempty"`
    ExcludeDomains []string `json:"exclude_domains,omitempty"`
    AllowExternal  bool     `json:"allow_external,omitempty"`
    IncludeImages  bool     `json:"include_images,omitempty"`
    Categories     []string `json:"categories,omitempty"`
    ExtractDepth   string   `json:"extract_depth,omitempty"`      // "basic" or "advanced"
}
```

### Map API (Beta)

The Map API allows you to generate a sitemap starting from a base URL.

```go
type TavilyMapRequest struct {
    URL            string   `json:"url"`
    MaxDepth       int      `json:"max_depth,omitempty"`
    MaxBreadth     int      `json:"max_breadth,omitempty"`
    Limit          int      `json:"limit,omitempty"`
    Instructions   string   `json:"instructions,omitempty"`
    SelectPaths    []string `json:"select_paths,omitempty"`
    SelectDomains  []string `json:"select_domains,omitempty"`
    ExcludePaths   []string `json:"exclude_paths,omitempty"`
    ExcludeDomains []string `json:"exclude_domains,omitempty"`
    AllowExternal  bool     `json:"allow_external,omitempty"`
    Categories     []string `json:"categories,omitempty"`
}
```

## Error Handling

All API methods return both a response and an error. The error will be non-nil if:

1. The request could not be created (invalid parameters)
2. The HTTP request failed (network issues)
3. The API returned a non-200 status code
4. The response could not be unmarshaled (invalid JSON)

Example:

```go
response, err := client.Search(ctx, searchRequest)
if err != nil {
    if strings.Contains(err.Error(), "response status code: 401") {
        // Handle authentication error
        log.Fatal("Invalid API key")
    } else if strings.Contains(err.Error(), "response status code: 429") {
        // Handle rate limiting
        log.Fatal("Rate limit exceeded")
    } else {
        // Handle other errors
        log.Fatal(err)
    }
}
```

## Advanced Configuration

### Custom HTTP Client

You can provide your own HTTP client to handle specific requirements like custom timeouts, retries, or logging:

```go
import (
    "net/http"
    "time"
)

httpClient := &http.Client{
    Timeout: 30 * time.Second,
    // Add other customizations as needed
}

client := tavily.NewTavilyClient(
    "your-api-key", 
    tavily.TavilyClientWithHttpClient(httpClient),
)
```

### Custom Base URL

You can override the default API endpoint:

```go
client := tavily.NewTavilyClient(
    "your-api-key",
    tavily.TavilyClientWithBaseURL("https://your-custom-endpoint.example.com"),
)
```

## License

This library is licensed under the MIT License.

## Support

For support, please contact [Tavily support](https://tavily.com/support) or open an issue on the GitHub repository.