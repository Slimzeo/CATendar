package ai_sync

import (
	"context"
	"fmt"
	"strings"
	"time"

	"catendar-wails/core/ai_sync/agent"
)

type Settings struct {
	Provider   agent.Provider `json:"provider"`
	CodexPath  string         `json:"codexPath"`
	ClaudePath string         `json:"claudePath"`
}

type AgentConfiguration struct {
	Settings Settings             `json:"settings"`
	Agents   []agent.Availability `json:"agents"`
}

func (s *Service) GetAgentConfiguration(ctx context.Context) (AgentConfiguration, error) {
	settings, err := s.getSettings(ctx)
	if err != nil {
		return AgentConfiguration{}, err
	}

	result := AgentConfiguration{Settings: settings}
	for _, provider := range []agent.Provider{agent.ProviderCodex, agent.ProviderClaude} {
		adapter, _ := agent.ForProvider(provider)
		availability := agent.Resolve(adapter, configuredPath(settings, provider))
		if availability.Installed {
			versionContext, cancel := context.WithTimeout(ctx, 3*time.Second)
			availability = agent.Version(versionContext, availability)
			cancel()
		}
		result.Agents = append(result.Agents, availability)
	}
	return result, nil
}

func (s *Service) SaveSettings(ctx context.Context, settings Settings) error {
	provider, err := agent.ValidateProvider(string(settings.Provider))
	if err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `
		UPDATE ai_settings
		SET provider = ?, codex_path = ?, claude_path = ?
		WHERE id = 1
	`, provider, strings.TrimSpace(settings.CodexPath), strings.TrimSpace(settings.ClaudePath)); err != nil {
		return fmt.Errorf("save AI settings: %w", err)
	}
	return nil
}

func (s *Service) getSettings(ctx context.Context) (Settings, error) {
	settings := Settings{Provider: agent.ProviderCodex}
	if err := s.db.QueryRowContext(ctx, `
		SELECT provider, codex_path, claude_path FROM ai_settings WHERE id = 1
	`).Scan(&settings.Provider, &settings.CodexPath, &settings.ClaudePath); err != nil {
		return Settings{}, fmt.Errorf("read AI settings: %w", err)
	}
	if _, err := agent.ValidateProvider(string(settings.Provider)); err != nil {
		settings.Provider = agent.ProviderCodex
	}
	return settings, nil
}

func configuredPath(settings Settings, provider agent.Provider) string {
	if provider == agent.ProviderClaude {
		return settings.ClaudePath
	}
	return settings.CodexPath
}
