// api/client.go
package movieExternalApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, // Default timeout for requests
		},
	}
}

func (c *APIClient) Get(path string, queryParams map[string]string, result interface{}) error {
	fullURL, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters if provided
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		fullURL.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil) // nil for body on GET
	if err != nil {
		return fmt.Errorf("error creating GET request: %w", err)
	}

	// Set common headers. Authorization header is no longer set.
	req.Header.Set("X-RapidAPI-Host", "movie-database-api1.p.rapidapi.com")
	req.Header.Set("X-RapidAPI-Key", "8918ef8442msh0541d1e1ae87ed5p1117b7jsn73b4ffb5dba4")
	req.Header.Set("Accept", "application/json") // Request JSON response

	return c.doRequest(req, result)
}

func (c *APIClient) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body) // Read body for error details
		return fmt.Errorf("API returned non-success status: %d %s, Body: %s", resp.StatusCode, resp.Status, string(bodyBytes))
	}

	if resp.StatusCode == http.StatusNoContent { // Handle 204 No Content
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if result != nil && len(body) > 0 { // Only unmarshal if a result pointer is provided and body is not empty
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("error unmarshaling JSON response: %w (body: %s)", err, string(body))
		}
	}

	return nil
}
