package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigDirName is the name of the config directory under UserConfigDir.
const ConfigDirName = "gum"

// ConfigFileName is the primary config filename (no extension; content is YAML).
const ConfigFileName = "config"

// FallbackConfigFile is the fallback config file path relative to home (~/.gumrc).
const FallbackConfigFile = ".gumrc"

var commandPrefix = map[string]string{
	"choose":  "GUM_CHOOSE_",
	"confirm": "GUM_CONFIRM_",
	"file":    "GUM_FILE_",
	"filter":  "GUM_FILTER_",
	"format":  "GUM_FORMAT_",
	"input":   "GUM_INPUT_",
	"join":    "GUM_JOIN_",
	"log":     "GUM_LOG_",
	"pager":   "GUM_PAGER_",
	"spin":    "GUM_SPIN_",
	"style":   "GUM_STYLE_",
	"table":   "GUM_TABLE_",
	"write":   "GUM_WRITE_",
}

// If GUM_DEBUG_CONFIG is set (to any value), Load writes the config file path
// and the number of env vars applied to stderr.
func Load() error {
	path := findConfigFile()
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}

	n := applyToEnv(raw)
	if os.Getenv("GUM_DEBUG_CONFIG") != "" {
		fmt.Fprintf(os.Stderr, "gum: loaded config from %s (%d defaults applied)\n", path, n)
	}
	return nil
}

// findConfigFile returns the path to the config file, or "" if none exists.
// It checks in order: UserConfigDir()/gum/config, XDG_CONFIG_HOME/gum/config
// (or ~/.config/gum/config), then ~/.gumrc. Each location is tried with
// .yaml, .yml, and no extension.
func findConfigFile() string {
	tryPaths := func(dir string) string {
		if dir == "" {
			return ""
		}
		gumDir := filepath.Join(dir, ConfigDirName)
		for _, name := range []string{ConfigFileName + ".yaml", ConfigFileName + ".yml", ConfigFileName} {
			p := filepath.Join(gumDir, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
		return ""
	}

	// Primary: platform config dir (e.g. ~/Library/Application Support on macOS)
	if dir, err := os.UserConfigDir(); err == nil {
		if p := tryPaths(dir); p != "" {
			return p
		}
	}

	// XDG-style: $XDG_CONFIG_HOME/gum or ~/.config/gum (so ~/.config works on macOS too)
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		if p := tryPaths(dir); p != "" {
			return p
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		if p := tryPaths(filepath.Join(home, ".config")); p != "" {
			return p
		}
	}

	// Fallback: ~/.gumrc[.yaml|.yml]
	if home, err := os.UserHomeDir(); err == nil {
		for _, name := range []string{FallbackConfigFile + ".yaml", FallbackConfigFile + ".yml", FallbackConfigFile} {
			p := filepath.Join(home, name)
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}

	return ""
}

// applyToEnv walks the parsed config and sets GUM_* env vars for any key not already set.
// It returns the number of env vars that were set.
func applyToEnv(raw map[string]interface{}) int {
	var n int
	for key, value := range raw {
		keyLower := strings.ToLower(key)
		prefix, ok := commandPrefix[keyLower]
		if !ok {
			continue
		}
		n += flattenToEnv(prefix, value)
	}
	return n
}

func flattenToEnv(prefix string, value interface{}) int {
	switch v := value.(type) {
	case map[string]interface{}:
		var n int
		for k, val := range v {
			segment := envSegment(k)
			n += flattenToEnv(prefix+segment+"_", val)
		}
		return n
	case []interface{}:
		return 0
	case string:
		envKey := strings.TrimSuffix(prefix, "_")
		if envKey != "" && os.Getenv(envKey) == "" {
			_ = os.Setenv(envKey, v)
			return 1
		}
		return 0
	case bool:
		envKey := strings.TrimSuffix(prefix, "_")
		if envKey != "" && os.Getenv(envKey) == "" {
			if v {
				_ = os.Setenv(envKey, "true")
			} else {
				_ = os.Setenv(envKey, "false")
			}
			return 1
		}
		return 0
	case int:
		envKey := strings.TrimSuffix(prefix, "_")
		if envKey != "" && os.Getenv(envKey) == "" {
			_ = os.Setenv(envKey, formatNumber(v))
			return 1
		}
		return 0
	case int64:
		envKey := strings.TrimSuffix(prefix, "_")
		if envKey != "" && os.Getenv(envKey) == "" {
			_ = os.Setenv(envKey, formatNumber(v))
			return 1
		}
		return 0
	case float64:
		envKey := strings.TrimSuffix(prefix, "_")
		if envKey != "" && os.Getenv(envKey) == "" {
			_ = os.Setenv(envKey, formatNumber(v))
			return 1
		}
		return 0
	default:
		return 0
	}
}

func envSegment(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToUpper(s)
}

func formatNumber(n interface{}) string {
	return fmt.Sprint(n)
}
