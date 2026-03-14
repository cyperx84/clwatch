package workspace

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ValidationError represents a config validation error with context.
type ValidationError struct {
	Field   string
	Message string
	Hint    string
}

func (e ValidationError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s: %s (hint: %s)", e.Field, e.Message, e.Hint)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}
	var lines []string
	for _, e := range errs {
		lines = append(lines, "  - "+e.Error())
	}
	return "config validation errors:\n" + strings.Join(lines, "\n")
}

// Validate checks the config for common errors and returns helpful messages.
func (c Config) Validate(configPath string) ValidationErrors {
	var errs ValidationErrors

	// Schema validation
	if c.Schema == "" {
		errs = append(errs, ValidationError{
			Field:   "schema",
			Message: "required field missing",
			Hint:    "set to \"clwatch.config.v1\"",
		})
	} else if c.Schema != "clwatch.config.v1" {
		errs = append(errs, ValidationError{
			Field:   "schema",
			Message: fmt.Sprintf("unsupported schema version %q", c.Schema),
			Hint:    "only \"clwatch.config.v1\" is supported",
		})
	}

	// Tools validation
	if len(c.Tools) == 0 {
		errs = append(errs, ValidationError{
			Field:   "tools",
			Message: "no tools configured",
			Hint:    "add at least one tool (e.g. \"claude-code\", \"codex-cli\")",
		})
	} else {
		knownTools := map[string]bool{
			"claude-code": true,
			"codex-cli":   true,
			"gemini-cli":  true,
			"opencode":    true,
			"openclaw":    true,
		}
		for _, tool := range c.Tools {
			if tool == "" {
				errs = append(errs, ValidationError{
					Field:   "tools",
					Message: "empty tool name in list",
					Hint:    "remove empty strings from the tools array",
				})
				continue
			}
			if !knownTools[tool] {
				errs = append(errs, ValidationError{
					Field:   "tools",
					Message: fmt.Sprintf("unknown tool %q", tool),
					Hint:    "known tools: claude-code, codex-cli, gemini-cli, opencode, openclaw",
				})
			}
		}
	}

	// ManifestURL validation
	if c.ManifestURL != "" {
		u, err := url.Parse(c.ManifestURL)
		if err != nil {
			errs = append(errs, ValidationError{
				Field:   "manifestUrl",
				Message: fmt.Sprintf("invalid URL: %v", err),
				Hint:    "must be a valid HTTP/HTTPS URL",
			})
		} else if u.Scheme != "http" && u.Scheme != "https" {
			errs = append(errs, ValidationError{
				Field:   "manifestUrl",
				Message: fmt.Sprintf("unsupported scheme %q", u.Scheme),
				Hint:    "use http:// or https://",
			})
		}
	}

	// ReferenceDir validation
	if c.ReferenceDir != "" {
		if !strings.HasSuffix(c.ReferenceDir, "/") {
			errs = append(errs, ValidationError{
				Field:   "referenceDir",
				Message: "must end with /",
				Hint:    fmt.Sprintf("change to %q", c.ReferenceDir+"/"),
			})
		}
		// Check if directory exists (relative to config file)
		if configPath != "" {
			configDir := filepath.Dir(configPath)
			refPath := filepath.Join(configDir, c.ReferenceDir)
			if _, err := os.Stat(refPath); os.IsNotExist(err) {
				errs = append(errs, ValidationError{
					Field:   "referenceDir",
					Message: fmt.Sprintf("directory %q does not exist", c.ReferenceDir),
					Hint:    "run `clwatch init` to create it, or create the directory manually",
				})
			}
		}
	}

	// Tier2Threshold validation
	if c.Tier2Threshold != "" {
		valid := map[string]bool{"small": true, "medium": true, "large": true}
		if !valid[c.Tier2Threshold] {
			errs = append(errs, ValidationError{
				Field:   "tier2Threshold",
				Message: fmt.Sprintf("invalid threshold %q", c.Tier2Threshold),
				Hint:    "must be one of: small, medium, large",
			})
		}
	}

	return errs
}

// LoadAndValidate loads config and returns helpful errors if invalid.
func LoadAndValidate(configPath string) (Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading %s: %w", configPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		// Try to give a helpful parse error
		return Config{}, fmt.Errorf("parsing %s: %w (check for missing commas, unquoted strings, or trailing commas)", configPath, err)
	}

	// Validate
	if errs := cfg.Validate(configPath); len(errs) > 0 {
		return cfg, errs
	}

	return cfg, nil
}

// ValidateConfigFile is a convenience function for CLI use.
func ValidateConfigFile(configPath string) (bool, []string) {
	cfg, err := LoadAndValidate(configPath)
	if err != nil {
		if errs, ok := err.(ValidationErrors); ok {
			var msgs []string
			for _, e := range errs {
				msgs = append(msgs, e.Error())
			}
			return false, msgs
		}
		return false, []string{err.Error()}
	}

	// Valid config
	var msgs []string
	msgs = append(msgs, fmt.Sprintf("✓ Schema: %s", cfg.Schema))
	msgs = append(msgs, fmt.Sprintf("✓ Tools: %s", strings.Join(cfg.Tools, ", ")))
	msgs = append(msgs, fmt.Sprintf("✓ Manifest URL: %s", cfg.ManifestURL))
	msgs = append(msgs, fmt.Sprintf("✓ Reference dir: %s", cfg.ReferenceDir))
	if cfg.NotifyOnBreaking {
		msgs = append(msgs, "✓ Notify on breaking changes: enabled")
	}

	return true, msgs
}
