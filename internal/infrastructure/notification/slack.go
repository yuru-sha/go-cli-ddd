package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/yuru-sha/go-cli-ddd/internal/domain/model"
	"github.com/yuru-sha/go-cli-ddd/internal/infrastructure/config"
)

// SlackNotifier はSlackへの通知を担当します
type SlackNotifier struct {
	config config.SlackConfig
}

// SlackMessage はSlack APIに送信するメッセージ構造体です
type SlackMessage struct {
	Channel     string            `json:"channel,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment はSlackメッセージの添付ファイル構造体です
type SlackAttachment struct {
	Color  string `json:"color"`
	Text   string `json:"text"`
	Footer string `json:"footer,omitempty"`
	Ts     int64  `json:"ts,omitempty"`
}

// NewSlackNotifier は新しいSlackNotifierを作成します
func NewSlackNotifier(config config.SlackConfig) *SlackNotifier {
	return &SlackNotifier{
		config: config,
	}
}

// NotifyCommandResult はコマンド実行結果をSlackに通知します
func (n *SlackNotifier) NotifyCommandResult(result *model.CommandResult) error {
	if !n.config.Enabled {
		log.Debug().Msg("Slack通知は無効化されています")
		return nil
	}

	if n.config.WebhookURL == "" {
		log.Warn().Msg("Slack WebhookURLが設定されていません")
		return fmt.Errorf("Slack WebhookURLが設定されていません")
	}

	// 成功/失敗に応じたアイコンとカラーを設定
	statusEmoji := n.config.SuccessEmoji
	statusColor := "good" // green
	if !result.IsSuccess() {
		statusEmoji = n.config.FailureEmoji
		statusColor = "danger" // red
	}

	// メッセージのヘッダーテキスト
	headerText := fmt.Sprintf("%s %s: %s の処理が終了しました。", statusEmoji, n.config.Username, result.Process)

	// 引数部分のテキスト
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

	// 結果部分のテキスト
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

	// Slackメッセージの構築
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

	// メッセージをJSONに変換
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Slackメッセージのシリアライズに失敗しました")
		return err
	}

	// Slackに通知を送信
	err = n.sendWebhook(context.Background(), jsonMessage)
	if err != nil {
		log.Error().Err(err).Msg("Slack通知の送信に失敗しました")
		return err
	}

	log.Info().Msg("Slack通知を送信しました")
	return nil
}

// LogCommandResult はコマンド実行結果をログに出力します
func (n *SlackNotifier) LogCommandResult(result *model.CommandResult) {
	// 成功/失敗に応じたアイコンを設定
	statusEmoji := n.config.SuccessEmoji
	if !result.IsSuccess() {
		statusEmoji = n.config.FailureEmoji
	}

	// ログメッセージの構築
	logEvent := log.Info()
	if !result.IsSuccess() {
		logEvent = log.Error()
	}

	// 基本情報をログに追加
	logEvent.
		Str("process", result.Process).
		Str("status", result.Status).
		Str("start_time", model.FormatJST(result.StartTime)).
		Str("end_time", model.FormatJST(result.EndTime)).
		Str("duration", result.FormatDuration()).
		Int("total", result.TotalCount).
		Int("success", result.SuccessCount).
		Int("error", result.ErrorCount).
		Int("total_records", result.TotalRecords)

	// オプション情報をログに追加
	if len(result.AccountIDs) > 0 {
		logEvent.Strs("account_ids", result.AccountIDs)
	}
	if result.DateFrom != "" {
		logEvent.Str("date_from", result.DateFrom)
	}
	if result.DateTo != "" {
		logEvent.Str("date_to", result.DateTo)
	}

	// ログメッセージを出力
	logEvent.Msgf("%s コマンド実行結果: %s", statusEmoji, result.Process)
}

// sendWebhook はSlackのWebhookにJSONメッセージを送信します
func (n *SlackNotifier) sendWebhook(ctx context.Context, jsonMessage []byte) error {
	// リクエストを作成
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.config.WebhookURL, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("Slackリクエストの作成に失敗しました: %w", err)
	}

	// Content-Typeを設定
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Slack通知の送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスのステータスコードを確認
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Msg("Slack通知の送信に失敗しました")
		return fmt.Errorf("Slack通知の送信に失敗しました: ステータスコード %d", resp.StatusCode)
	}

	return nil
}
