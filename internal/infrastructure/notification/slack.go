package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/model"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// SlackNotifier はSlackへの通知を担当します。
type SlackNotifier struct {
	config config.SlackConfig
}

// SlackMessage represents a Slack webhook payload.
type SlackMessage struct {
	Channel     string            `json:"channel,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment represents one Slack attachment block.
type SlackAttachment struct {
	Color  string `json:"color"`
	Text   string `json:"text"`
	Footer string `json:"footer,omitempty"`
	TS     int64  `json:"ts,omitempty"`
}

// NewSlackNotifier creates a Slack notifier.
func NewSlackNotifier(config config.SlackConfig) *SlackNotifier {
	return &SlackNotifier{config: config}
}

// NotifyCommandResult sends one command result to Slack.
func (n *SlackNotifier) NotifyCommandResult(result *model.CommandResult) error {
	if !n.config.Enabled {
		slog.Debug("Slack通知は無効化されています")
		return nil
	}

	if n.config.WebhookURL == "" {
		slog.Warn("Slack WebhookURLが設定されていません")
		return fmt.Errorf("slack webhook URLが設定されていません")
	}

	statusEmoji := n.config.SuccessEmoji
	statusColor := "good"
	if !result.IsSuccess() {
		statusEmoji = n.config.FailureEmoji
		statusColor = "danger"
	}

	headerText := fmt.Sprintf("%s %s: %s の処理が終了しました。", statusEmoji, n.config.Username, result.Process)

	argsText := fmt.Sprintf("```\nProcess: %s\n", result.Process)
	if len(result.AccountIDs) > 0 {
		argsText += fmt.Sprintf("AccountIds: %s\n", strings.Join(result.AccountIDs, ", "))
	}
	if result.DateFrom != "" {
		argsText += fmt.Sprintf("From: %s\n", result.DateFrom)
	}
	if result.DateTo != "" {
		argsText += fmt.Sprintf("To: %s\n", result.DateTo)
	}
	argsText += "```"

	resultText := fmt.Sprintf("```\nStatus: %s\nStart: %s\nEnd: %s\nTime: %s\nTotal: %d\nSuccess: %d\nError: %d\nTotal Records: %d\n```",
		result.Status,
		model.FormatJST(result.StartTime),
		model.FormatJST(result.EndTime),
		result.FormatDuration(),
		result.TotalCount,
		result.SuccessCount,
		result.ErrorCount,
		result.TotalRecords,
	)

	message := SlackMessage{
		Channel:   n.config.Channel,
		Username:  n.config.Username,
		IconEmoji: n.config.IconEmoji,
		Text:      headerText,
		Attachments: []SlackAttachment{
			{
				Color: statusColor,
				Text:  fmt.Sprintf("Args\n%s\n\nResult\n%s", argsText, resultText),
			},
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		slog.Error("Slackメッセージのシリアライズに失敗しました", "err", err)
		return err
	}

	if err := n.sendWebhook(context.Background(), jsonMessage); err != nil {
		slog.Error("Slack通知の送信に失敗しました", "err", err)
		return err
	}

	slog.Info("Slack通知を送信しました")
	return nil
}

// LogCommandResult logs one command result using slog.
func (n *SlackNotifier) LogCommandResult(result *model.CommandResult) {
	statusEmoji := n.config.SuccessEmoji
	logFn := slog.Info
	if !result.IsSuccess() {
		statusEmoji = n.config.FailureEmoji
		logFn = slog.Error
	}

	args := []any{
		"process", result.Process,
		"status", result.Status,
		"start_time", model.FormatJST(result.StartTime),
		"end_time", model.FormatJST(result.EndTime),
		"duration", result.FormatDuration(),
		"total", result.TotalCount,
		"success", result.SuccessCount,
		"error", result.ErrorCount,
		"total_records", result.TotalRecords,
	}
	if len(result.AccountIDs) > 0 {
		args = append(args, "account_ids", result.AccountIDs)
	}
	if result.DateFrom != "" {
		args = append(args, "date_from", result.DateFrom)
	}
	if result.DateTo != "" {
		args = append(args, "date_to", result.DateTo)
	}

	logFn(fmt.Sprintf("%s コマンド実行結果: %s", statusEmoji, result.Process), args...)
}

func (n *SlackNotifier) sendWebhook(ctx context.Context, jsonMessage []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.config.WebhookURL, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("slackリクエストの作成に失敗しました: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("slack通知の送信に失敗しました: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		slog.Error("Slack通知の送信に失敗しました", "status_code", resp.StatusCode)
		return fmt.Errorf("slack通知の送信に失敗しました: ステータスコード %d", resp.StatusCode)
	}

	return nil
}
