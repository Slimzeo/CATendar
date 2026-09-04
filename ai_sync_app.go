package main

import (
	"context"
	"errors"
	"time"

	"catendar-wails/core/ai_sync"
	"catendar-wails/core/email"
)

type AISyncSetup struct {
	Email email.Account              `json:"email"`
	Agent ai_sync.AgentConfiguration `json:"agent"`
}

type AISyncApp struct {
	ctx    context.Context
	email  *email.Service
	aiSync *ai_sync.Service
}

func NewAISyncApp(emailService *email.Service, aiSyncService *ai_sync.Service) *AISyncApp {
	return &AISyncApp{email: emailService, aiSync: aiSyncService}
}

func (a *AISyncApp) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *AISyncApp) GetAISyncSetup() (*AISyncSetup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	account, err := a.email.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	agentConfiguration, err := a.aiSync.GetAgentConfiguration(ctx)
	if err != nil {
		return nil, err
	}
	return &AISyncSetup{Email: account, Agent: agentConfiguration}, nil
}

func (a *AISyncApp) SaveEmailAccount(input email.AccountInput) (*email.Account, error) {
	if err := a.ensureSettingsWritable(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	account, err := a.email.SaveAccount(ctx, input)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AISyncApp) TestEmailAccount(input email.AccountInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return a.email.TestAccount(ctx, input)
}

func (a *AISyncApp) GetSourceEmail(sourceURL string) (*email.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return a.email.ReadSource(ctx, sourceURL)
}

func (a *AISyncApp) SaveAISettings(settings ai_sync.Settings) error {
	if err := a.ensureSettingsWritable(); err != nil {
		return err
	}
	return a.aiSync.SaveSettings(context.Background(), settings)
}

func (a *AISyncApp) StartAISync() (ai_sync.Status, error) {
	if a.ctx == nil {
		return ai_sync.Status{}, errors.New("CATendar is not ready")
	}
	return a.aiSync.Start(a.ctx)
}

func (a *AISyncApp) GetAISyncStatus() ai_sync.Status {
	return a.aiSync.GetStatus()
}

func (a *AISyncApp) CancelAISync() bool {
	return a.aiSync.Cancel()
}

func (a *AISyncApp) ensureSettingsWritable() error {
	if a.aiSync.GetStatus().State == ai_sync.StateRunning {
		return errors.New("cancel the current AI Sync before changing its settings")
	}
	return nil
}
