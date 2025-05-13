# tavily-go

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aiomni/tavily-go) [![Go Reference](https://pkg.go.dev/badge/github.com/aiomni/tavily-go.svg)](https://pkg.go.dev/github.com/aiomni/tavily-go)

A Go client library for the [Tavily API](https://tavily.com), providing programmatic access to Tavily's search, extraction, crawling, and mapping functionalities.

## Installation

```bash
go get github.com/aiomni/tavily-go
```

## Authentication

To use the Tavily API, you'll need an API key. Get yours by signing up at [Tavily](https://tavily.com).

## JSON Schema Generation

`tavily-go` supports JSON Schema generation for all request and response types using the [invopop/jsonschema](https://github.com/invopop/jsonschema) package. All structs include appropriate jsonschema tags for automatic schema generation.

## Usage

### Initializing the Client

```go
import tavily "github.com/aiomni/tavily-go"

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
    
    tavily "github.com/aiomni/tavily-go"
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

All Tavily API structs include jsonschema tags to provide detailed information about parameters and can be used with the [invopop/jsonschema](https://github.com/invopop/jsonschema) package to generate json schema, which can used for LLM.

### Search API

The Search API allows you to search the web with various filtering options.

```go
type TavilySearchRequest struct {
	Query                    string   `json:"query" jsonschema:"required,description=The search query to execute with Tavily,example=who is Leo Messi?"`
	Topic                    string   `json:"topic,omitempty" jsonschema:"enum=news,enum=general,default=general,description=The category of the search. 'news' is useful for retrieving real-time updates\\, particularly about politics\\, sports\\, and major current events covered by mainstream media sources. 'general' is for broader\\, more general-purpose searches that may include a wide range of sources."`
	SearchDepth              string   `json:"search_depth,omitempty" jsonschema:"enum=basic,enum=advanced,default=basic,description=The depth of the search. 'advanced' search is tailored to retrieve the most relevant sources and 'content' snippets for your query\\, while 'basic' search provides generic content snippets from each source. "`
	ChunksPerSource          int      `json:"chunks_per_source,omitempty" jsonschema:"default=3,maximum=3,minimum=1,description=The number of 'content' chunks to retrieve from each source. Each chunk's length is maximum 500 characters. Available only when 'search_depth' is 'advanced'."`
	MaxResults               int      `json:"max_results,omitempty" jsonschema:"default=5,maximum=20,minimum=0,description=The maximum number of search results to return."`
	TimeRange                string   `json:"time_range,omitempty" jsonschema:"enum=day,enum=week,enum=month,enum=year,description=The time range back from the current date to filter results. Useful when looking for sources that have published data."`
	Days                     int      `json:"days,omitempty"  jsonschema:"default=7,minimum=1,description=Number of days back from the current date to include. Available only if 'topic' is 'news'."`
	IncludeAnswer            string   `json:"include_answer,omitempty" jsonschema:"description=Include an LLM-generated answer to the provided query. 'basic' returns a quick answer. 'advanced' returns a more detailed answer."`
	IncludeRawContent        bool     `json:"include_raw_content,omitempty" jsonschema:"default=false,description=Include the cleaned and parsed HTML content of each search result."`
	IncludeImages            bool     `json:"include_images,omitempty" jsonschema:"default=false,description=Also perform an image search and include the results in the response."`
	IncludeImageDescriptions bool     `json:"include_image_descriptions,omitempty" jsonschema:"default=false,description=When 'include_images' is 'true'\\, also add a descriptive text for each image."`
	IncludeDomains           []string `json:"include_domains,omitempty" jsonschema:"description=A list of domains to specifically include in the search results."`
	ExcludeDomains           []string `json:"exclude_domains,omitempty" jsonschema:"description=A list of domains to specifically exclude from the search results."`
}
```

### Extract API

The Extract API allows you to extract content from specific URLs.

```go
type TavilyExtractRequest struct {
	URLs          []string `json:"urls" jsonschema:"required,description=The URLs to extract content from."`
	IncludeImages bool     `json:"include_images,omitempty" jsonschema:"default=false,description=Include a list of images extracted from the URLs in the response."`
	ExtractDepth  string   `json:"extract_depth,omitempty" jsonschema:"enum=basic,enum=advanced,default=basic,description=The depth of the extraction process. 'advanced' extraction retrieves more data\\, including tables and embedded content\\, with higher success but may increase latency."`
}
```

### Crawl API

The Crawl API allows you to recursively crawl a website starting from a base URL.

```go
type TavilyCrawlRequest struct {
	URL            string   `json:"url" jsonschema:"required,description=The root URL to begin the crawl."`
	MaxDepth       int      `json:"max_depth,omitempty" jsonschema:"default=1,minimum=1,description=Max depth of the crawl. Defines how far from the base URL the crawler can explore."`
	MaxBreadth     int      `json:"max_breadth,omitempty" jsonschema:"default=20,minimum=1,description=Max number of links to follow per level of the tree (i.e.\\, per page)."`
	Limit          int      `json:"limit,omitempty" jsonschema:"default=50,minimum=1,description=Total number of links the crawler will process before stopping."`
	Instructions   string   `json:"instructions,omitempty" jsonschema:"description=Natural language instructions for the crawler."`
	SelectPaths    []string `json:"select_paths,omitempty" jsonschema_description:"Regex patterns to select only URLs with specific path patterns (e.g., /docs/.*)."`
	SelectDomains  []string `json:"select_domains,omitempty" jsonschema_description:"Regex patterns to select crawling to specific domains or subdomains (e.g., ^docs\\.example\\.com$)."`
	ExcludePaths   []string `json:"exclude_paths,omitempty" jsonschema_description:"Regex patterns to exclude URLs with specific path patterns (e.g., /private/.*, /admin/.*)."`
	ExcludeDomains []string `json:"exclude_domains,omitempty" jsonschema_description:"Regex patterns to exclude specific domains or subdomains from crawling (e.g., ^private\\.example\\.com$)."`
	AllowExternal  bool     `json:"allow_external,omitempty" jsonschema:"default=false,description=Whether to allow following links that go to external domains."`
	IncludeImages  bool     `json:"include_images,omitempty" jsonschema:"default=false,description=Whether to include images in the crawl results."`
	Categories     []string `json:"categories,omitempty" jsonschema_description:"Filter URLs using predefined categories. Available options: Careers, Blog, Documentation, About, Pricing, Community, Developers, Contact, Media, API"`
	ExtractDepth   string   `json:"extract_depth,omitempty" jsonschema:"enum=basic,enum=advanced,default=basic,description=Advanced extraction retrieves more data\\, including tables and embedded content\\, with higher success but may increase latency."`
}
```

### Map API

The Map API allows you to generate a sitemap starting from a base URL.

```go
type TavilyMapRequest struct {
	URL            string   `json:"url" jsonschema:"required,description=The root URL to begin the mapping."`
	MaxDepth       int      `json:"max_depth,omitempty" jsonschema:"default=1,minimum=1,description=Max depth of the mapping. Defines how far from the base URL the crawler can explore."`
	MaxBreadth     int      `json:"max_breadth,omitempty" jsonschema:"default=20,minimum=1,description=Max number of links to follow per level of the tree (i.e.\\, per page)."`
	Limit          int      `json:"limit,omitempty" jsonschema:"default=50,minimum=1,description=Total number of links the crawler will process before stopping."`
	Instructions   string   `json:"instructions,omitempty" jsonschema_description:"Natural language instructions for the crawler."`
	SelectPaths    []string `json:"select_paths,omitempty" jsonschema_description:"Regex patterns to select only URLs with specific path patterns (e.g., /docs/.*)."`
	SelectDomains  []string `json:"select_domains,omitempty" jsonschema_description:"Regex patterns to select crawling to specific domains or subdomains (e.g., ^docs\\.example\\.com$)."`
	ExcludePaths   []string `json:"exclude_paths,omitempty" jsonschema_description:"Regex patterns to exclude URLs with specific path patterns (e.g., /private/.*, /admin/.*)."`
	ExcludeDomains []string `json:"exclude_domains,omitempty" jsonschema_description:"Regex patterns to exclude specific domains or subdomains from crawling (e.g., ^private\\.example\\.com$)."`
	AllowExternal  bool     `json:"allow_external,omitempty" jsonschema:"default=false,description=Whether to allow following links that go to external domains."`
	Categories     []string `json:"categories,omitempty" jsonschema_description:"Filter URLs using predefined categories. Available options: Careers, Blog, Documentation, About, Pricing, Community, Developers, Contact, Media, API"`
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
    
    tavily "github.com/aiomni/tavily-go"
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

## Generating JSON Schema

You can generate JSON Schema for any of the request types using the `jsonschema` package:

```go
import (
    "encoding/json"
    "fmt"
    
    "github.com/invopop/jsonschema"
    tavily "github.com/tavily/tavily-go"
)

func main() {
    // Create a reflector
    reflector := new(jsonschema.Reflector)
    
    // Generate schema for search request
    searchSchema := reflector.Reflect(&tavily.TavilySearchRequest{})
    
    // Marshal to JSON
    schemaJSON, err := json.MarshalIndent(searchSchema, "", "  ")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(schemaJSON))
}
```

This will generate a complete JSON Schema with all the properties, descriptions, defaults, and constraints defined in the struct tags:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/schema.json",
  "properties": {
    "query": {
      "type": "string",
      "description": "The search query to execute with Tavily",
      "examples": [
        "who is Leo Messi?"
      ]
    },
    "topic": {
      "type": "string",
      "enum": [
        "news",
        "general"
      ],
      "default": "general",
      "description": "The category of the search. 'news' is useful for retrieving real-time updates, particularly about politics, sports, and major current events covered by mainstream media sources. 'general' is for broader, more general-purpose searches that may include a wide range of sources."
    },
    "search_depth": {
      "type": "string",
      "enum": [
        "basic",
        "advanced"
      ],
      "default": "basic",
      "description": "The depth of the search. 'advanced' search is tailored to retrieve the most relevant sources and 'content' snippets for your query, while 'basic' search provides generic content snippets from each source. "
    }
    // ... other properties
  },
  "required": [
    "query"
  ],
  "additionalProperties": true
}
```

You can generate schemas for all request types:

```go
// Extract API schema
extractSchema := reflector.Reflect(&tavily.TavilyExtractRequest{})

// Crawl API schema
crawlSchema := reflector.Reflect(&tavily.TavilyCrawlRequest{})

// Map API schema
mapSchema := reflector.Reflect(&tavily.TavilyMapRequest{})
```

## License

This library is licensed under the MIT License.

## Support

For support, please open an issue on the GitHub repository.
