package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/spf13/cobra"
)

type stubCampaignHandler struct {
	request CampaignRequest
	called  bool
}

func (h *stubCampaignHandler) Run(_ context.Context, req CampaignRequest) error {
	h.called = true
	h.request = req
	return nil
}

func TestCampaignCommandParsesAccountIDs(t *testing.T) {
	handler := &stubCampaignHandler{}
	cmd := NewCampaignCommand(handler).Cmd
	cmd.SetArgs([]string{"--account-ids", "1, 2,3", "--status", "active", "--parallel", "3", "--force"})

	err := cmd.Execute()

	require.NoError(t, err)
	require.True(t, handler.called)
	require.Equal(t, []int{1, 2, 3}, handler.request.AccountIDs)
	require.Equal(t, "active", handler.request.Status)
	require.Equal(t, 3, handler.request.Parallel)
	require.True(t, handler.request.Force)
}

func TestCommandModulesRegisterToRoot(t *testing.T) {
	root := &cobra.Command{Use: "root"}

	accountCmd := &AccountCommand{Cmd: &cobra.Command{Use: "account"}}
	campaignCmd := &CampaignCommand{Cmd: &cobra.Command{Use: "campaign"}}
	masterCmd := &MasterCommand{Cmd: &cobra.Command{Use: "master"}}

	for _, module := range []CommandModule{accountCmd, campaignCmd, masterCmd} {
		module.Register(root)
	}

	require.Len(t, root.Commands(), 3)
	require.Equal(t, "account", root.Commands()[0].Use)
	require.Equal(t, "campaign", root.Commands()[1].Use)
	require.Equal(t, "master", root.Commands()[2].Use)
}
