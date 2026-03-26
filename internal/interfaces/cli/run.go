package cli

import (
	"context"
	"log/slog"
	"time"
)

func runCommand(ctx context.Context, startMessage string, attrs []any, run func(context.Context) error) error {
	startedAt := time.Now()
	slog.Info(startMessage, attrs...)

	if err := run(ctx); err != nil {
		errorAttrs := append([]any{}, attrs...)
		errorAttrs = append(errorAttrs, "err", err)
		slog.Error("コマンド実行に失敗しました", errorAttrs...)
		return err
	}

	doneAttrs := append([]any{}, attrs...)
	doneAttrs = append(doneAttrs, "elapsed_time", time.Since(startedAt))
	slog.Info("コマンド実行が完了しました", doneAttrs...)
	return nil
}
