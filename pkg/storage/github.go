package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/sitelogix/backend/pkg/config"
)

type GitHubStorage struct {
	token  string
	owner  string
	repo   string
	branch string
	client *http.Client
}

type uploadRequest struct {
	Message string `json:"message"`
	Content string `json:"content"`
	Branch  string `json:"branch"`
	SHA     string `json:"sha,omitempty"`
}

type uploadResponse struct {
	Content struct {
		SHA        string `json:"sha"`
		DownloadURL string `json:"download_url"`
		HTMLURL    string `json:"html_url"`
	} `json:"content"`
}

type fileInfo struct {
	SHA string `json:"sha"`
}

type UploadResult struct {
	URL string
	SHA string
	Key string
}

func NewGitHubStorage(cfg config.GitHubConfig) *GitHubStorage {
	if cfg.Token == "" || cfg.Owner == "" || cfg.Repo == "" {
		log.Println("[WARN] GitHub storage: GITHUB_TOKEN, GITHUB_OWNER, or GITHUB_REPO is not set — media uploads will fail")
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       60 * time.Second,
	}
	return &GitHubStorage{
		token:  cfg.Token,
		owner:  cfg.Owner,
		repo:   cfg.Repo,
		branch: cfg.Branch,
		client: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
	}
}

func (g *GitHubStorage) Upload(fileData []byte, originalName, folder string) (*UploadResult, error) {
	ext := filepath.Ext(originalName)
	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	encoded := base64.StdEncoding.EncodeToString(fileData)

	body := uploadRequest{
		Message: fmt.Sprintf("upload: %s", key),
		Content: encoded,
		Branch:  g.branch,
	}

	resp, err := g.apiCall(http.MethodPut, g.contentsURL(key), body)
	if err != nil {
		return nil, fmt.Errorf("github upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github upload failed (%d): %s", resp.StatusCode, string(raw))
	}

	var result uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode github response: %w", err)
	}

	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
		g.owner, g.repo, g.branch, key)

	return &UploadResult{
		URL: rawURL,
		SHA: result.Content.SHA,
		Key: key,
	}, nil
}

func (g *GitHubStorage) Delete(key, sha string) error {
	body := map[string]string{
		"message": fmt.Sprintf("delete: %s", key),
		"sha":     sha,
		"branch":  g.branch,
	}

	resp, err := g.apiCall(http.MethodDelete, g.contentsURL(key), body)
	if err != nil {
		return fmt.Errorf("github delete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github delete failed (%d): %s", resp.StatusCode, string(raw))
	}
	return nil
}

func (g *GitHubStorage) Download(key string) ([]byte, error) {
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
		g.owner, g.repo, g.branch, key)

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github fetch failed (%d)", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (g *GitHubStorage) GetSHA(key string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, g.contentsURL(key)+"?ref="+g.branch, nil)
	if err != nil {
		return "", err
	}
	g.setHeaders(req)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	var info fileInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}
	return info.SHA, nil
}

func (g *GitHubStorage) contentsURL(path string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", g.owner, g.repo, path)
}

func (g *GitHubStorage) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")
}

func (g *GitHubStorage) apiCall(method, url string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	g.setHeaders(req)

	return g.client.Do(req)
}
