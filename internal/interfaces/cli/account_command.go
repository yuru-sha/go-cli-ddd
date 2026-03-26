package cli

import (
	"context"
	"flag"
)

// NewAccountCommand はアカウントコマンドを作成します。
func NewAccountCommand(handler AccountHandler) *AccountCommand {
	return &AccountCommand{handler: handler}
}

// Execute parses flags and runs the account command.
func (c *AccountCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	var (
		accountIDs string
		syncMode   string
		force      bool
	)

	fs.StringVar(&accountIDs, "id", "", "同期するアカウントID（カンマ区切り、未指定なら全件）")
	fs.StringVar(&syncMode, "mode", "full", "同期モード（full, diff）")
	fs.BoolVar(&force, "force", false, "強制同期フラグ")

	if err := parseCommandFlags(fs, args, c.Name()); err != nil {
		return err
	}

	parsedIDs, err := parseCSVIntIDs(accountIDs)
	if err != nil {
		return err
	}

	request := AccountRequest{
		AccountIDs: parsedIDs,
		SyncMode:   syncMode,
		Force:      force,
	}

	return runCommand(
		ctx,
		"アカウント同期コマンドを実行します",
		[]any{
			"account_ids", request.AccountIDs,
			"sync_mode", request.SyncMode,
			"force", request.Force,
		},
		func(runCtx context.Context) error {
			return c.handler.Run(runCtx, request)
		},
	)
}
