// Package cli provides the Cobra-based command boundary for the application.
package cli

import (
	"context"

	"github.com/spf13/cobra"
)

// NewAccountCommand はアカウントコマンドを作成します。
func NewAccountCommand(handler AccountHandler) *AccountCommand {
	flags := &AccountRequest{}

	cmd := &cobra.Command{
		Use:   "account",
		Short: "アカウント情報を同期します",
		Long:  `外部APIからアカウント情報を取得し、データベースに保存します。`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runCommand(
				context.Background(),
				"アカウント同期コマンドを実行します",
				[]any{
					"account_ids", flags.AccountIDs,
					"sync_mode", flags.SyncMode,
					"force", flags.Force,
				},
				func(ctx context.Context) error {
					return handler.Run(ctx, *flags)
				},
			)
		},
	}

	cmd.Flags().IntSliceVar(&flags.AccountIDs, "id", []int{}, "同期するアカウントID（指定しない場合は全アカウント）")
	cmd.Flags().StringVar(&flags.SyncMode, "mode", "full", "同期モード（full: 全同期, diff: 差分同期）")
	cmd.Flags().BoolVar(&flags.Force, "force", false, "強制同期フラグ（既存データを上書き）")

	return &AccountCommand{Cmd: cmd}
}
