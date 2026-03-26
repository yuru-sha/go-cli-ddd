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

// NewRepository creates a notification repository implementation.
func NewRepository(config *config.Config) repository.NotificationRepository {
	return &Repository{
		slackNotifier: NewSlackNotifier(config.Notification.Slack),
	}
}

// NotifyCommandResult sends the command result to configured notifiers.
func (r *Repository) NotifyCommandResult(result *model.CommandResult) error {
	// Slackに通知
	if err := r.slackNotifier.NotifyCommandResult(result); err != nil {
		return err
	}

	// ログにも出力
	r.LogCommandResult(result)

	return nil
}

// LogCommandResult writes the command result to the logger.
func (r *Repository) LogCommandResult(result *model.CommandResult) {
	r.slackNotifier.LogCommandResult(result)
}
