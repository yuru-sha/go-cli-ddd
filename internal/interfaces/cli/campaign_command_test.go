package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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
	cmd := NewCampaignCommand(handler)

	err := cmd.Execute(context.Background(), []string{"--account-ids", "1, 2,3", "--status", "active", "--parallel", "3", "--force"})

	require.NoError(t, err)
	require.True(t, handler.called)
	require.Equal(t, []int{1, 2, 3}, handler.request.AccountIDs)
	require.Equal(t, "active", handler.request.Status)
	require.Equal(t, 3, handler.request.Parallel)
	require.True(t, handler.request.Force)
}

func TestCommandModulesRegisterToRoot(t *testing.T) {
	root := NewRootCommand()

	accountCmd := &AccountCommand{}
	campaignCmd := &CampaignCommand{}
	masterCmd := &MasterCommand{}

	for _, module := range []CommandModule{accountCmd, campaignCmd, masterCmd} {
		root.Register(module)
	}

	require.Len(t, root.modules, 3)
	require.Contains(t, root.modules, "account")
	require.Contains(t, root.modules, "campaign")
	require.Contains(t, root.modules, "master")
}
