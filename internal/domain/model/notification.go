package model

import (
	"fmt"
	"strings"
	"time"
)

// CommandResult はコマンド実行結果を表します
type CommandResult struct {
	Process      string    // 実行したコマンド
	AccountIDs   []string  // 処理対象のアカウントID
	DateFrom     string    // 処理対象期間（開始）
	DateTo       string    // 処理対象期間（終了）
	Status       string    // 処理結果のステータス（success/failed）
	StartTime    time.Time // 処理開始時刻
	EndTime      time.Time // 処理終了時刻
	TotalCount   int       // 処理したアカウントの件数
	SuccessCount int       // 処理したアカウントの成功した件数
	ErrorCount   int       // 処理したアカウントの失敗した件数
	TotalRecords int       // 登録/更新したレコード数
}

// NewCommandResult はCommandResultの新しいインスタンスを作成します
func NewCommandResult(process string) *CommandResult {
	return &CommandResult{
		Process:   process,
		StartTime: time.Now(),
		Status:    "success", // デフォルトは成功
	}
}

// SetAccountIDs は処理対象のアカウントIDを設定します
func (r *CommandResult) SetAccountIDs(accountIDs []string) {
	r.AccountIDs = accountIDs
}

// SetDateRange は処理対象期間を設定します
func (r *CommandResult) SetDateRange(from, to string) {
	r.DateFrom = from
	r.DateTo = to
}

// Complete は処理を完了し、終了時刻を記録します
func (r *CommandResult) Complete() {
	r.EndTime = time.Now()
}

// SetFailed は処理を失敗としてマークします
func (r *CommandResult) SetFailed() {
	r.Status = "failed"
}

// AddCounts は処理結果のカウントを追加します
func (r *CommandResult) AddCounts(success, errorCount, records int) {
	r.SuccessCount += success
	r.ErrorCount += errorCount
	r.TotalCount += (success + errorCount)
	r.TotalRecords += records
}

// FormatDuration は処理時間を人間が読みやすい形式でフォーマットします
func (r *CommandResult) FormatDuration() string {
	duration := r.EndTime.Sub(r.StartTime)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	parts := []string{}

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d時間", hours))
	}

	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d分", minutes))
	}

	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d秒", seconds))
	}

	return strings.Join(parts, " ")
}

// IsSuccess は処理が成功したかどうかを返します
func (r *CommandResult) IsSuccess() bool {
	return r.Status == "success"
}

// FormatJST は時刻をJST形式でフォーマットします
func FormatJST(t time.Time) string {
	jst := time.FixedZone("JST", 9*60*60)
	return t.In(jst).Format("2006-01-02 15:04:05 MST")
}
