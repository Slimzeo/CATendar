package agent

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Provider string

const (
	ProviderCodex  Provider = "codex"
	ProviderClaude Provider = "claude"
)

type RunRequest struct {
	BinaryPath  string
	WorkDir     string
	Prompt      string
	CLICommand  string
	Environment []string
}

type Invocation struct {
	Path  string
	Args  []string
	Stdin string
}

type Availability struct {
	Provider  Provider `json:"provider"`
	Installed bool     `json:"installed"`
	Path      string   `json:"path"`
	Version   string   `json:"version,omitempty"`
	Error     string   `json:"error,omitempty"`
}

type Adapter interface {
	Provider() Provider
	BinaryName() string
	Build(RunRequest) Invocation
}

func ForProvider(provider Provider) (Adapter, error) {
	switch provider {
	case ProviderCodex:
		return codexAdapter{}, nil
	case ProviderClaude:
		return claudeAdapter{}, nil
	default:
		return nil, fmt.Errorf("unsupported AI provider %q", provider)
	}
}

func Resolve(adapter Adapter, configuredPath string) Availability {
	availability := Availability{Provider: adapter.Provider()}
	path, err := findExecutable(adapter.BinaryName(), configuredPath)
	if err != nil {
		availability.Error = err.Error()
		return availability
	}
	availability.Installed = true
	availability.Path = path
	return availability
}

func findExecutable(name, configuredPath string) (string, error) {
	if configuredPath != "" {
		path, err := filepath.Abs(configuredPath)
		if err != nil {
			return "", fmt.Errorf("resolve configured %s path: %w", name, err)
		}
		if isExecutable(path) {
			return path, nil
		}
		return "", fmt.Errorf("configured %s executable was not found", name)
	}

	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join("/opt/homebrew/bin", name),
		filepath.Join("/usr/local/bin", name),
		filepath.Join(home, ".local", "bin", name),
		filepath.Join(home, ".npm-global", "bin", name),
	}
	if runtime.GOOS == "windows" {
		for index := range candidates {
			candidates[index] += ".exe"
		}
	}
	for _, path := range candidates {
		if isExecutable(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("%s CLI is not installed or not discoverable", name)
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Ext(path), ".exe") || filepath.Ext(path) == ""
	}
	return info.Mode()&0o111 != 0
}

func ValidateProvider(value string) (Provider, error) {
	provider := Provider(strings.ToLower(strings.TrimSpace(value)))
	if provider != ProviderCodex && provider != ProviderClaude {
		return "", errors.New("AI provider must be codex or claude")
	}
	return provider, nil
}
