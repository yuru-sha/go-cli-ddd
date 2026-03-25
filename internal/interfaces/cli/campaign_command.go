package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NewCampaignCommand はキャンペーンコマンドを作成します。
func NewCampaignCommand(handler CampaignHandler) *CampaignCommand {
	var accountIDs string
	flags := &CampaignRequest{}

	cmd := &cobra.Command{
		Use:   "campaign",
		Short: "キャンペーン情報を同期します",
		Long:  `アカウントごとに並列処理を行い、外部APIからキャンペーン情報を取得し、データベースに保存します。`,
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

			return runCommand(
				context.Background(),
				"キャンペーン同期コマンドを実行します",
				[]any{
					"account_ids", request.AccountIDs,
					"status", request.Status,
					"parallel", request.Parallel,
					"force", request.Force,
				},
				func(ctx context.Context) error {
					return handler.Run(ctx, request)
				},
			)
		},
	}

	cmd.Flags().StringVar(&accountIDs, "account-ids", "", "同期するアカウントID（カンマ区切り、例: '1,2,3'）、空の場合は全アカウント")
	cmd.Flags().StringVar(&flags.Status, "status", "", "同期するキャンペーンのステータス（active, paused, completedなど）")
	cmd.Flags().IntVar(&flags.Parallel, "parallel", 5, "並列処理数（1-10）")
	cmd.Flags().BoolVar(&flags.Force, "force", false, "強制同期フラグ（既存データを上書き）")

	return &CampaignCommand{Cmd: cmd}
}
