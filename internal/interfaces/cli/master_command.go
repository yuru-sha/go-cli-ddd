package cli

import (
	"context"
	"flag"
	"time"
)

// NewMasterCommand はマスターコマンドを作成します。
func NewMasterCommand(handler MasterHandler) *MasterCommand {
	return &MasterCommand{handler: handler}
}

// Execute parses flags and runs the master command.
func (c *MasterCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	request := MasterRequest{}
	var (
		accountIDs string
		timeoutSec int
	)

	fs.StringVar(&accountIDs, "account-ids", "", "同期するアカウントID（カンマ区切り、例: 1,2,3）")
	fs.IntVar(&request.Parallel, "parallel", 5, "並列処理数（1-10）")
	fs.IntVar(&timeoutSec, "timeout", 0, "タイムアウト時間（秒）")
	fs.BoolVar(&request.Force, "force", false, "強制同期フラグ")

	if err := parseCommandFlags(fs, args, c.Name()); err != nil {
		return err
	}

	parsedAccountIDs, err := parseCSVIntIDs(accountIDs)
	if err != nil {
		return err
	}
	if err := validateParallel(request.Parallel); err != nil {
		return err
	}

	request.AccountIDs = parsedAccountIDs
	request.Timeout = time.Duration(timeoutSec) * time.Second

	runCtx := ctx
	cancel := func() {}
	if request.Timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, request.Timeout)
	}
	defer cancel()

	return runCommand(
		runCtx,
		"マスター同期コマンドを実行します",
		[]any{
			"account_ids", request.AccountIDs,
			"parallel", request.Parallel,
			"timeout", request.Timeout,
			"force", request.Force,
		},
		func(execCtx context.Context) error {
			return c.handler.Run(execCtx, request)
		},
	)
}
