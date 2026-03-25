package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// NewMasterCommand はマスターコマンドを作成します。
func NewMasterCommand(handler MasterHandler) *MasterCommand {
	var (
		accountIDs string
		timeoutSec int
	)
	flags := &MasterRequest{}

	cmd := &cobra.Command{
		Use:   "master",
		Short: "マスター情報を同期します",
		Long:  `アカウント情報とキャンペーン情報を順に同期します。`,
		RunE: func(_ *cobra.Command, _ []string) error {
			parsedAccountIDs, err := parseCSVIntIDs(accountIDs)
			if err != nil {
				return err
			}
			if flags.Parallel < 1 || flags.Parallel > 10 {
				return fmt.Errorf("parallel は 1 から 10 の範囲で指定してください")
			}

			request := *flags
			request.AccountIDs = parsedAccountIDs
			request.Timeout = time.Duration(timeoutSec) * time.Second

			ctx := context.Background()
			cancel := func() {}
			if request.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, request.Timeout)
			}
			defer cancel()

			return runCommand(
				ctx,
				"マスター同期コマンドを実行します",
				[]any{
					"account_ids", request.AccountIDs,
					"parallel", request.Parallel,
					"timeout", request.Timeout,
					"force", request.Force,
				},
				func(runCtx context.Context) error {
					return handler.Run(runCtx, request)
				},
			)
		},
	}

	cmd.Flags().StringVar(&accountIDs, "account-ids", "", "同期するアカウントID（カンマ区切り、例: '1,2,3'）、空の場合は全アカウント")
	cmd.Flags().IntVar(&flags.Parallel, "parallel", 5, "並列処理数（1-10）")
	cmd.Flags().IntVar(&timeoutSec, "timeout", 0, "タイムアウト時間（秒）、0の場合はタイムアウトなし")
	cmd.Flags().BoolVar(&flags.Force, "force", false, "強制同期フラグ（既存データを上書き）")

	return &MasterCommand{Cmd: cmd}
}
