package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Update struct {
	Tool            string `json:"tool"`
	Status          string `json:"status"`
	PreviousVersion string `json:"previous_version"`
	CurrentVersion  string `json:"current_version"`
	Breaking        bool   `json:"breaking"`
}

type WebhookPayload struct {
	Event      string   `json:"event"`
	DetectedAt string   `json:"detected_at"`
	Updates    []Update `json:"updates"`
}

type ManifestTool struct {
	Version            string `json:"version"`
	PayloadURL         string `json:"payload_url"`
	GeneratedAt        string `json:"generated_at"`
	VerificationStatus string `json:"verification_status"`
	StaleAfter         string `json:"stale_after"`
}

type Manifest struct {
	Schema      string                  `json:"schema"`
	GeneratedAt string                  `json:"generated_at"`
	Tools       map[string]ManifestTool `json:"tools"`
}

type Config struct {
	ManifestURL string
	Interval    time.Duration
	JSONOutput  bool
	WebhookURL  string
}

type DiffFunc func(ctx context.Context, manifestURL string, jsonOutput bool) ([]Update, error)

func Run(cfg Config, diffFn DiffFunc) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nShutting down...")
		cancel()
	}()

	if !cfg.JSONOutput {
		fmt.Printf("clwatch watch — polling every %s\n", cfg.Interval)
		fmt.Printf("Manifest: %s\n\n", cfg.ManifestURL)
	}

	// Run immediately on start
	if err := runOnce(ctx, cfg, diffFn); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := runOnce(ctx, cfg, diffFn); err != nil {
				if ctx.Err() != nil {
					return nil
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}
}

func runOnce(ctx context.Context, cfg Config, diffFn DiffFunc) error {
	updates, err := diffFn(ctx, cfg.ManifestURL, cfg.JSONOutput)
	if err != nil {
		return err
	}

	if len(updates) > 0 && cfg.WebhookURL != "" {
		payload := WebhookPayload{
			Event:      "tools_updated",
			DetectedAt: time.Now().UTC().Format(time.RFC3339),
			Updates:    updates,
		}
		if err := postWebhook(ctx, cfg.WebhookURL, payload); err != nil {
			fmt.Fprintf(os.Stderr, "Webhook error: %v\n", err)
		}
	}

	return nil
}

func postWebhook(ctx context.Context, url string, payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned %d", resp.StatusCode)
	}
	return nil
}

func ParseInterval(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid interval %q: %w", s, err)
	}
	if d < 15*time.Minute {
		return 0, fmt.Errorf("interval must be at least 15m, got %s", d)
	}
	return d, nil
}
