package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	apiKey string
	apiURL string
	http   *http.Client
}

type apiResponse struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func NewClient(apiKey, apiURL string) *Client {
	return &Client{
		apiKey: apiKey,
		apiURL: apiURL,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Fetch returns a human-readable weather string for the given coordinates.
// Returns empty string (no error) when no API key is configured.
func (c *Client) Fetch(lat, lng float64) (string, error) {
	if c.apiKey == "" {
		return "", nil
	}

	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s&units=metric", c.apiURL, lat, lng, c.apiKey)
	resp, err := c.http.Get(url)
	if err != nil {
		return "", fmt.Errorf("weather api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("weather api returned %d", resp.StatusCode)
	}

	var w apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return "", err
	}
	if len(w.Weather) == 0 {
		return "", nil
	}

	desc := w.Weather[0].Description
	if len(desc) > 0 {
		desc = strings.ToUpper(desc[:1]) + desc[1:]
	}

	return fmt.Sprintf("%s, %.0f°C, Humidity: %d%%, Wind: %.1f m/s",
		desc, w.Main.Temp, w.Main.Humidity, w.Wind.Speed), nil
}
