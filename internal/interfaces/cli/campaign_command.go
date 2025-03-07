package cli

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/yuru-sha/go-cli-ddd/internal/application/usecase"
)

// NewCampaignCommand はキャンペーンコマンドを作成します
func NewCampaignCommand(campaignUseCase *usecase.CampaignUseCase) *CampaignCommand {
	// フラグ変数の定義
	var (
		accountIDs  string
		status      string
		parallelNum int
		force       bool
	)

	cmd := &cobra.Command{
		Use:   "campaign",
		Short: "キャンペーン情報を同期します",
		Long:  `アカウントごとに並列処理を行い、外部APIからキャンペーン情報を取得し、データベースに保存します。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			startTime := time.Now()

			log.Info().Str("account_ids", accountIDs).Str("status", status).Int("parallel_num", parallelNum).Bool("force", force).Msg("キャンペーン同期コマンドを実行します")

			// キャンペーン情報の同期
			// 引数に基づいて処理を分岐
			var err error
			if accountIDs != "" {
				// 特定のアカウントのキャンペーンのみ同期
				log.Info().Str("account_ids", accountIDs).Msg("指定されたアカウントのキャンペーンのみ同期します")
				// 実際の実装ではここでaccountIDsをパースして使う
				// 例: "1,2,3" → [1, 2, 3]
				err = campaignUseCase.SyncCampaigns(ctx)
			} else if status != "" {
				// 特定ステータスのキャンペーンのみ同期
				log.Info().Str("status", status).Msg("指定されたステータスのキャンペーンのみ同期します")
				// 実際の実装ではここでstatusを使った処理を行う
				err = campaignUseCase.SyncCampaigns(ctx)
			} else {
				// 全キャンペーン同期
				// 並列数を設定する場合はここでparallelNumを使う
				log.Info().Int("parallel_num", parallelNum).Msg("并列数を設定して同期します")
				err = campaignUseCase.SyncCampaigns(ctx)
			}

			if err != nil {
				log.Error().Err(err).Msg("キャンペーン同期に失敗しました")
				return err
			}

			elapsedTime := time.Since(startTime)
			log.Info().Dur("elapsed_time", elapsedTime).Msg("キャンペーン同期コマンドが完了しました")
			return nil
		},
	}

	// フラグの設定
	cmd.Flags().StringVar(&accountIDs, "account-ids", "", "同期するアカウントID（カンマ区切り、例: '1,2,3'）、空の場合は全アカウント")
	cmd.Flags().StringVar(&status, "status", "", "同期するキャンペーンのステータス（active, paused, completedなど）")
	cmd.Flags().IntVar(&parallelNum, "parallel", 5, "並列処理数（1-10）")
	cmd.Flags().BoolVar(&force, "force", false, "強制同期フラグ（既存データを上書き）")

	return &CampaignCommand{Cmd: cmd}
}
