package metarank

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/junichi-fukushima/tech-flow/backend/domain/metaRank"
	"github.com/junichi-fukushima/tech-flow/backend/infrastructure/external/initializer"
	"io"
	"net/http"
)

// Client represents a Metarank API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Metarank client.
func NewClient() *Client {
	//host := os.Getenv("METARANK_HOST")
	//if host == "" {
	//	initializer.Logger.Info("Metarank is disabled due to missing METARANK_HOST")
	//}
	// TODO: 一時的にdebug用に直打ちしている
	host := "http://3.112.51.68:8080"
	return &Client{
		BaseURL:    host,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) isDisabled() bool {
	return c.BaseURL == ""
}

// SendFeedback sends a feedback to the Metarank API.
func (c *Client) SendFeedback(event any) error {
	if c.isDisabled() {
		initializer.Logger.Info("Skip sending feedback to Metarank because it is disabled")
		return nil
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	url := fmt.Sprintf("%s/feedback", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		initializer.Logger.Info(err.Error())
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("API error: %s", string(body)))
	}

	return nil
}

// GetTrending get trending articles from Metarank API.
func (c *Client) GetTrending(event any) (*metaRank.Trending, error) {
	if c.isDisabled() {
		initializer.Logger.Info("Skip getting trending articles from Metarank because it is disabled")
		return &metaRank.Trending{
			Items: []metaRank.ItemDetail{},
		}, nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	url := fmt.Sprintf("%s/recommend/trending", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		initializer.Logger.Info(err.Error())
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(fmt.Sprintf("API error: %s", string(body)))
	}

	var trending metaRank.Trending
	if err := json.NewDecoder(resp.Body).Decode(&trending); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &trending, nil
}

// GetRecommendation get recommendation from Metarank API.
func (c *Client) GetRecommendation(event any) (*metaRank.Recommend, error) {
	if c.isDisabled() {
		initializer.Logger.Info("Skip getting recommendation from Metarank because it is disabled")
		return &metaRank.Recommend{
			Items: []metaRank.RecommendItemDetail{},
		}, nil
	}

	body, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	url := fmt.Sprintf("%s/rank/xgboost", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		initializer.Logger.Info(err.Error())
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(fmt.Sprintf("API error: %s", string(body)))
	}

	var recommendation metaRank.Recommend
	if err := json.NewDecoder(resp.Body).Decode(&recommendation); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &recommendation, nil
}
