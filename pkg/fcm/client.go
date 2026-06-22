package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const fcmEndpoint = "https://fcm.googleapis.com/fcm/send"

type Client struct {
	serverKey string
	http      *http.Client
}

type message struct {
	To           string            `json:"to"`
	Notification notification      `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
	Priority     string            `json:"priority"`
}

type notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

type fcmResponse struct {
	Success int `json:"success"`
	Failure int `json:"failure"`
}

func NewClient(serverKey string) *Client {
	return &Client{
		serverKey: serverKey,
		http:      &http.Client{Timeout: 15 * time.Second},
	}
}

// Send pushes a notification to a single device token.
// Silently no-ops when no server key is configured.
func (c *Client) Send(deviceToken, title, body string, data map[string]string) error {
	if c.serverKey == "" || deviceToken == "" {
		return nil
	}

	payload := message{
		To:       deviceToken,
		Priority: "high",
		Notification: notification{
			Title: title,
			Body:  body,
			Sound: "default",
		},
		Data: data,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fcmEndpoint, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "key="+c.serverKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("fcm request: %w", err)
	}
	defer resp.Body.Close()

	var result fcmResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Failure > 0 {
		return fmt.Errorf("fcm: %d message(s) failed to deliver", result.Failure)
	}
	return nil
}

// SendMulticast sends to multiple device tokens.
func (c *Client) SendMulticast(tokens []string, title, body string, data map[string]string) []error {
	errs := make([]error, 0)
	for _, t := range tokens {
		if err := c.Send(t, title, body, data); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
