// Package cli provides the flag-based command boundary for the application.
package cli

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CommandModule represents one executable CLI command.
type CommandModule interface {
	Name() string
	Execute(ctx context.Context, args []string) error
}

// RootCommand executes the top-level CLI and dispatches subcommands.
type RootCommand struct {
	configPath string
	env        string
	modules    map[string]CommandModule
}

// AccountCommand wraps the account subcommand.
type AccountCommand struct {
	handler AccountHandler
}

// Name returns the subcommand name.
func (c *AccountCommand) Name() string {
	return "account"
}

// CampaignCommand wraps the campaign subcommand.
type CampaignCommand struct {
	handler CampaignHandler
}

// Name returns the subcommand name.
func (c *CampaignCommand) Name() string {
	return "campaign"
}

// MasterCommand wraps the master subcommand.
type MasterCommand struct {
	handler MasterHandler
}

// Name returns the subcommand name.
func (c *MasterCommand) Name() string {
	return "master"
}

// AccountRequest is the CLI input for the account command.
type AccountRequest struct {
	AccountIDs []int
	SyncMode   string
	Force      bool
}

// CampaignRequest is the CLI input for the campaign command.
type CampaignRequest struct {
	AccountIDs []int
	Status     string
	Parallel   int
	Force      bool
}

// MasterRequest is the CLI input for the master command.
type MasterRequest struct {
	AccountIDs []int
	Parallel   int
	Timeout    time.Duration
	Force      bool
}

// AccountHandler executes the account command.
type AccountHandler interface {
	Run(ctx context.Context, req AccountRequest) error
}

// CampaignHandler executes the campaign command.
type CampaignHandler interface {
	Run(ctx context.Context, req CampaignRequest) error
}

// MasterHandler executes the master command.
type MasterHandler interface {
	Run(ctx context.Context, req MasterRequest) error
}

func parseCSVIntIDs(value string) ([]int, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parts := strings.Split(value, ",")
	ids := make([]int, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		id, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, fmt.Errorf("ID %q の解析に失敗しました: %w", trimmed, err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func parseCommandFlags(fs *flag.FlagSet, args []string, commandName string) error {
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("%s コマンドに未対応の引数があります: %v", commandName, fs.Args())
	}

	return nil
}

func validateParallel(parallel int) error {
	if parallel < 1 || parallel > 10 {
		return fmt.Errorf("parallel は 1 から 10 の範囲で指定してください")
	}

	return nil
}
