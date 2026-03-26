// Package cli provides the flag-based command boundary for the application.
package cli

import (
	"context"
	"log/slog"

	"github.com/yuru-sha/go-cli-ddd/internal/application/usecase"
)

type accountSyncUseCase interface {
	SyncAccounts(ctx context.Context) error
	SyncAccountsByIDs(ctx context.Context, accountIDs []int) error
}

type campaignSyncUseCase interface {
	SyncCampaigns(ctx context.Context) error
}

type masterSyncUseCase interface {
	SyncAll(ctx context.Context) error
}

type accountHandler struct {
	useCase accountSyncUseCase
}

// NewAccountHandler creates the account command handler.
func NewAccountHandler(useCase *usecase.AccountUseCase) AccountHandler {
	return &accountHandler{useCase: useCase}
}

func (h *accountHandler) Run(ctx context.Context, req AccountRequest) error {
	if len(req.AccountIDs) > 0 {
		slog.Info("指定されたアカウントのみ同期します", "account_ids", req.AccountIDs)
		return h.useCase.SyncAccountsByIDs(ctx, req.AccountIDs)
	}

	return h.useCase.SyncAccounts(ctx)
}

type campaignHandler struct {
	useCase campaignSyncUseCase
}

// NewCampaignHandler creates the campaign command handler.
func NewCampaignHandler(useCase *usecase.CampaignUseCase) CampaignHandler {
	return &campaignHandler{useCase: useCase}
}

func (h *campaignHandler) Run(ctx context.Context, req CampaignRequest) error {
	if len(req.AccountIDs) > 0 {
		slog.Info("指定されたアカウントのキャンペーンのみ同期します", "account_ids", req.AccountIDs)
	}
	if req.Status != "" {
		slog.Info("指定されたステータスのキャンペーンを同期します", "status", req.Status)
	}
	return h.useCase.SyncCampaigns(ctx)
}

type masterHandler struct {
	useCase masterSyncUseCase
}

// NewMasterHandler creates the master command handler.
func NewMasterHandler(useCase *usecase.MasterUseCase) MasterHandler {
	return &masterHandler{useCase: useCase}
}

func (h *masterHandler) Run(ctx context.Context, req MasterRequest) error {
	if len(req.AccountIDs) > 0 {
		slog.Info("指定されたアカウントのみ同期します", "account_ids", req.AccountIDs)
	}
	return h.useCase.SyncAll(ctx)
}
