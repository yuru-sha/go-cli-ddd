package repository

import (
	"github.com/yuru-sha/go-cli-ddd/internal/domain/model"
)

// NotificationRepository は通知機能のインターフェースを定義します
type NotificationRepository interface {
	// NotifyCommandResult はコマンド実行結果を通知します
	NotifyCommandResult(result *model.CommandResult) error

	// LogCommandResult はコマンド実行結果をログに出力します
	LogCommandResult(result *model.CommandResult)
}
