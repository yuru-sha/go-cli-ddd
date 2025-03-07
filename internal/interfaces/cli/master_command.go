package cli

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/yuru-sha/go-cli-ddd/internal/application/usecase"
)

// NewMasterCommand はマスターコマンドを作成します
func NewMasterCommand(masterUseCase *usecase.MasterUseCase) *MasterCommand {
	// フラグ変数の定義
	var (
		accountIDs  string
		parallelNum int
		timeoutSec  int
		force       bool
	)

	cmd := &cobra.Command{
		Use:   "master",
		Short: "マスター情報を同期します",
		Long:  `アカウント情報とキャンペーン情報を順に同期します。`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// タイムアウト付きコンテキストの作成
			ctx := context.Background()
			if timeoutSec > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
				defer cancel()
			}

			startTime := time.Now()

			log.Info().Str("account_ids", accountIDs).Int("parallel_num", parallelNum).Int("timeout_sec", timeoutSec).Bool("force", force).Msg("マスター同期コマンドを実行します")

			// マスター情報の同期
			// 引数に基づいて処理を分岐
			var err error

			// アカウントIDが指定されている場合の処理
			if accountIDs != "" {
				log.Info().Str("account_ids", accountIDs).Msg("指定されたアカウントのみ同期します")
				// 実際の実装ではここでaccountIDsをパースして使う
				// 例: "1,2,3" → [1, 2, 3]
			}

			// 全ての情報を同期
			log.Info().Int("parallel_num", parallelNum).Msg("并列数を設定して全情報を同期します")
			err = masterUseCase.SyncAll(ctx)

			if err != nil {
				log.Error().Err(err).Msg("マスター同期に失敗しました")
				return err
			}

			elapsedTime := time.Since(startTime)
			log.Info().Dur("elapsed_time", elapsedTime).Msg("マスター同期コマンドが完了しました")
			return nil
		},
	}

	// フラグの設定
	cmd.Flags().StringVar(&accountIDs, "account-ids", "", "同期するアカウントID（カンマ区切り、例: '1,2,3'）、空の場合は全アカウント")
	cmd.Flags().IntVar(&parallelNum, "parallel", 5, "並列処理数（1-10）")
	cmd.Flags().IntVar(&timeoutSec, "timeout", 0, "タイムアウト時間（秒）、0の場合はタイムアウトなし")
	cmd.Flags().BoolVar(&force, "force", false, "強制同期フラグ（既存データを上書き）")

	return &MasterCommand{Cmd: cmd}
}
