// Package cli provides the Cobra-based command boundary for the application.
package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// CommandModule represents a CLI module that can register itself on the root command.
type CommandModule interface {
	Register(root *cobra.Command)
}

// RootCommand wraps the Cobra root command.
type RootCommand struct {
	Cmd *cobra.Command
}

// AccountCommand wraps the account subcommand.
type AccountCommand struct {
	Cmd *cobra.Command
}

// Register attaches the account command to the root command.
func (c *AccountCommand) Register(root *cobra.Command) {
	root.AddCommand(c.Cmd)
}

// CampaignCommand wraps the campaign subcommand.
type CampaignCommand struct {
	Cmd *cobra.Command
}

// Register attaches the campaign command to the root command.
func (c *CampaignCommand) Register(root *cobra.Command) {
	root.AddCommand(c.Cmd)
}

// MasterCommand wraps the master subcommand.
type MasterCommand struct {
	Cmd *cobra.Command
}

// Register attaches the master command to the root command.
func (c *MasterCommand) Register(root *cobra.Command) {
	root.AddCommand(c.Cmd)
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
