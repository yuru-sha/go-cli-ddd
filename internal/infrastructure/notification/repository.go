package notification

import (
	"github.com/yuru-sha/go-cli-ddd/internal/domain/model"
	"github.com/yuru-sha/go-cli-ddd/internal/domain/repository"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// Repository は通知リポジトリの実装です
type Repository struct {
	slackNotifier *SlackNotifier
}

// NewRepository は新しい通知リポジトリを作成します
func NewRepository(config *config.Config) repository.NotificationRepository {
	return &Repository{
		slackNotifier: NewSlackNotifier(config.Notification.Slack),
	}
}

// NotifyCommandResult はコマンド実行結果を通知します
func (r *Repository) NotifyCommandResult(result *model.CommandResult) error {
	// Slackに通知
	if err := r.slackNotifier.NotifyCommandResult(result); err != nil {
		return err
	}

	// ログにも出力
	r.LogCommandResult(result)

	return nil
}

// LogCommandResult はコマンド実行結果をログに出力します
func (r *Repository) LogCommandResult(result *model.CommandResult) {
	r.slackNotifier.LogCommandResult(result)
}
