package cli

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/yuru-sha/go-cli-ddd/internal/application/usecase"
)

// NewAccountCommand はアカウントコマンドを作成します
func NewAccountCommand(accountUseCase *usecase.AccountUseCase) *AccountCommand {
	// フラグ変数の定義
	var (
		accountIDs []int
		syncMode   string
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "account",
		Short: "アカウント情報を同期します",
		Long:  `外部APIからアカウント情報を取得し、データベースに保存します。`,
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			startTime := time.Now()

			log.Info().Ints("account_ids", accountIDs).Str("sync_mode", syncMode).Bool("force", force).Msg("アカウント同期コマンドを実行します")

			// アカウント情報の同期
			// 引数に基づいて処理を分岐
			var err error
			if len(accountIDs) > 0 {
				// 特定のアカウントのみ同期
				log.Info().Ints("account_ids", accountIDs).Msg("指定されたアカウントのみ同期します")
				err = accountUseCase.SyncAccountsByIDs(ctx, accountIDs)
			} else {
				// 全アカウント同期
				err = accountUseCase.SyncAccounts(ctx)
			}

			if err != nil {
				log.Error().Err(err).Msg("アカウント同期に失敗しました")
				return err
			}

			elapsedTime := time.Since(startTime)
			log.Info().Dur("elapsed_time", elapsedTime).Msg("アカウント同期コマンドが完了しました")
			return nil
		},
	}

	// フラグの設定
	cmd.Flags().IntSliceVar(&accountIDs, "id", []int{}, "同期するアカウントID（指定しない場合は全アカウント）")
	cmd.Flags().StringVar(&syncMode, "mode", "full", "同期モード（full: 全同期, diff: 差分同期）")
	cmd.Flags().BoolVar(&force, "force", false, "強制同期フラグ（既存データを上書き）")

	return &AccountCommand{Cmd: cmd}
}
