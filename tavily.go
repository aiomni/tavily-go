package tavily

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultTavilyBaseURL = "https://api.tavily.com"

type TavilyClient struct {
	BaseURL    string
	HttpClient *http.Client
	APIKey     string
}

type TavilyClientOption func(*TavilyClient)

func NewTavilyClient(apiKey string, opts ...TavilyClientOption) *TavilyClient {
	c := &TavilyClient{
		APIKey:     apiKey,
		BaseURL:    DefaultTavilyBaseURL,
		HttpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func TavilyClientWithBaseURL(BaseURL string) TavilyClientOption {
	return func(c *TavilyClient) {
		c.BaseURL = BaseURL
	}
}

func TavilyClientWithHttpClient(httpClient *http.Client) TavilyClientOption {
	return func(c *TavilyClient) {
		c.HttpClient = httpClient
	}
}

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

type TavilySearchResult struct {
	Title      string  `json:"title,omitempty"`
	URL        string  `json:"url,omitempty"`
	Content    string  `json:"content,omitempty"`
	Score      float64 `json:"score,omitempty"`
	RawContent any     `json:"raw_content,omitempty"`
}

type TavilySearchResponse struct {
	Query             string               `json:"query"`
	FollowUpQuestions []string             `json:"follow_up_questions,omitempty"`
	Answer            string               `json:"answer,omitempty"`
	Images            TavilyImages         `json:"images,omitempty"`
	Results           []TavilySearchResult `json:"results,omitempty"`
	ResponseTime      float64              `json:"response_time"`
}

type TavilySearchImage struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}


type TavilyImages []TavilySearchImage

func (m *TavilyImages) UnmarshalJSON(data []byte) error {
	var urls []string
	if err := json.Unmarshal(data, &urls); err == nil {
		var images []TavilySearchImage
		for _, url := range urls {
			images = append(images, TavilySearchImage{URL: url})
		}
		*m = images
		return nil
	}

	var detailed []TavilySearchImage
	if err := json.Unmarshal(data, &detailed); err == nil {
		*m = detailed
		return nil
	}

	return fmt.Errorf("images: invalid format")
}

func (c *TavilyClient) do(ctx context.Context, path string, requestBody []byte) (responseBody []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("tavily client search, build request: %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, call /search api: %w", err)
	}
	defer response.Body.Close()

	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily client search, response status code: %d, response body: %s", response.StatusCode, string(responseBody))
	}
	return
}

func (c *TavilyClient) Search(ctx context.Context, searchRequest *TavilySearchRequest) (*TavilySearchResponse, error) {
	requestJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, marshal search request: %w", err)
	}

	responseBody, err := c.do(ctx, "/search", requestJSON)
	if err != nil {
		return nil, err
	}

	result := TavilySearchResponse{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, parse response: %w", err)
	}

	return &result, nil
}


type TavilyExtractRequest struct {
	URLs          []string `json:"urls" jsonschema:"required,description=The URLs to extract content from."`
	IncludeImages bool     `json:"include_images,omitempty" jsonschema:"default=false,description=Include a list of images extracted from the URLs in the response."`
	ExtractDepth  string   `json:"extract_depth,omitempty" jsonschema:"enum=basic,enum=advanced,default=basic,description=The depth of the extraction process. 'advanced' extraction retrieves more data\\, including tables and embedded content\\, with higher success but may increase latency."`
}

type TavilyExtractResult struct {
	URL        string   `json:"url" jsonschema:"description=The URL from which the content was extracted."`
	RawContent string   `json:"raw_content,omitempty" jsonschema:"description=The full content extracted from the page."`
	Images     []string `json:"images" jsonschema:"description=A list of image URLs extracted from the page."`
}

type TavilyExtractFailedResult struct {
	URL   string `json:"url" jsonschema:"description=The URL that failed to be processed."`
	Error string `json:"error,omitempty" jsonschema:"description=An error message describing why the URL couldn't be processed."`
}

type TavilyExtractResponse struct {
	Results       []TavilyExtractResult       `json:"results" jsonschema:"description=A list of extracted content from the provided URLs."`
	FailedResults []TavilyExtractFailedResult `json:"failed_results,omitempty" jsonschema:"description=A list of URLs that could not be processed."`
	ResponseTime  float64                     `json:"response_time"`
}

func (c *TavilyClient) Extract(ctx context.Context, request *TavilyExtractRequest) (*TavilyExtractResponse, error) {
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("tavily client extract, marshal request: %w", err)
	}

	responseBody, err := c.do(ctx, "/extract", requestJSON)
	if err != nil {
		return nil, err
	}

	result := TavilyExtractResponse{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("tavily client extract, parse response: %w", err)
	}

	return &result, nil
}
