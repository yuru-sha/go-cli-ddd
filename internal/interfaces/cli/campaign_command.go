package cli

import (
	"context"
	"flag"
)

// NewCampaignCommand はキャンペーンコマンドを作成します。
func NewCampaignCommand(handler CampaignHandler) *CampaignCommand {
	return &CampaignCommand{handler: handler}
}

// Execute parses flags and runs the campaign command.
func (c *CampaignCommand) Execute(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet(c.Name(), flag.ContinueOnError)

	request := CampaignRequest{}
	var accountIDs string

	fs.StringVar(&accountIDs, "account-ids", "", "同期するアカウントID（カンマ区切り、例: 1,2,3）")
	fs.StringVar(&request.Status, "status", "", "同期するキャンペーンのステータス")
	fs.IntVar(&request.Parallel, "parallel", 5, "並列処理数（1-10）")
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

	return runCommand(
		ctx,
		"キャンペーン同期コマンドを実行します",
		[]any{
			"account_ids", request.AccountIDs,
			"status", request.Status,
			"parallel", request.Parallel,
			"force", request.Force,
		},
		func(runCtx context.Context) error {
			return c.handler.Run(runCtx, request)
		},
	)
}
