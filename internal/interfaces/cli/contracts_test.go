package cli

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommandFlagsRejectsUnexpectedArgs(t *testing.T) {
	fs := flag.NewFlagSet("account", flag.ContinueOnError)
	var force bool
	fs.BoolVar(&force, "force", false, "force sync")

	err := parseCommandFlags(fs, []string{"--force", "extra"}, "account")

	require.EqualError(t, err, "account コマンドに未対応の引数があります: [extra]")
}

func TestValidateParallel(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		require.NoError(t, validateParallel(5))
	})

	t.Run("invalid", func(t *testing.T) {
		err := validateParallel(0)
		require.EqualError(t, err, "parallel は 1 から 10 の範囲で指定してください")
	})
}
