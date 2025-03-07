package cli

import (
	"github.com/spf13/cobra"
)

// RootCommand はルートコマンドを表します
type RootCommand struct {
	Cmd *cobra.Command
}

// AccountCommand はアカウントコマンドを表します
type AccountCommand struct {
	Cmd *cobra.Command
	// フラグ変数
	AccountID int
	SyncMode  string
	Force     bool
}

// CampaignCommand はキャンペーンコマンドを表します
type CampaignCommand struct {
	Cmd *cobra.Command
	// フラグ変数
	AccountID   int
	Status      string
	ParallelNum int
	Force       bool
}

// MasterCommand はマスターコマンドを表します
type MasterCommand struct {
	Cmd *cobra.Command
	// フラグ変数
	SyncMode    string
	ParallelNum int
	TimeoutSec  int
	Force       bool
}
