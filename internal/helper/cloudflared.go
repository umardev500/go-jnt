package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const cloudflareURLAPI = "https://next-api-omega-ruby.vercel.app/api/urls?newer=true"

type CloudflareURL struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetCurrentCloudflareURL() (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(cloudflareURLAPI)
	if err != nil {
		return "", fmt.Errorf("fetch Cloudflare URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("URL API returned status %d", resp.StatusCode)
	}

	var urls []CloudflareURL
	if err := json.NewDecoder(resp.Body).Decode(&urls); err != nil {
		return "", fmt.Errorf("decode URL response: %w", err)
	}

	if len(urls) == 0 {
		return "", fmt.Errorf("no Cloudflare URL available")
	}

	// Use only one record: the first returned record.
	currentURL := strings.TrimRight(strings.TrimSpace(urls[0].URL), "/")
	if currentURL == "" {
		return "", fmt.Errorf("Cloudflare URL is empty")
	}

	if !strings.HasPrefix(currentURL, "https://") {
		return "", fmt.Errorf("invalid Cloudflare URL")
	}

	return currentURL, nil
}

func CallCurrentTunnel(path string) error {
	if path != "/granted" && path != "/revoke" {
		return fmt.Errorf("unsupported endpoint")
	}

	baseURL, err := GetCurrentCloudflareURL()
	if err != nil {
		return err
	}

	log.Info().Msgf("Calling %s", baseURL+path)

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+path,
		nil,
	)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("tunnel endpoint returned status %d", resp.StatusCode)
	}

	return nil
}
