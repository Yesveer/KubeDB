package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type PrometheusService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewPrometheusService - Initialize service
func NewPrometheusService(prometheusURL string) *PrometheusService {
	return &PrometheusService{
		BaseURL: prometheusURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Query - Execute Prometheus query
func (ps *PrometheusService) Query(ctx context.Context, query string) (interface{}, error) {
	
	// Build URL
	queryURL := fmt.Sprintf("%s/api/v1/query", ps.BaseURL)
	
	// Add query parameter
	params := url.Values{}
	params.Add("query", query)
	fullURL := fmt.Sprintf("%s?%s", queryURL, params.Encode())
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Execute request
	resp, err := ps.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Parse JSON
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	return result, nil
}

// QueryRange - Execute Prometheus range query
func (ps *PrometheusService) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (interface{}, error) {
	
	// Build URL
	queryURL := fmt.Sprintf("%s/api/v1/query_range", ps.BaseURL)
	
	// Add parameters
	params := url.Values{}
	params.Add("query", query)
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	params.Add("end", fmt.Sprintf("%d", end.Unix()))
	params.Add("step", fmt.Sprintf("%.0fs", step.Seconds()))
	fullURL := fmt.Sprintf("%s?%s", queryURL, params.Encode())
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Execute request
	resp, err := ps.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Parse JSON
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	return result, nil
}