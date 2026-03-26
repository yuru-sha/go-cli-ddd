package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubAccountHandler struct {
	request AccountRequest
	called  bool
}

func (h *stubAccountHandler) Run(_ context.Context, req AccountRequest) error {
	h.called = true
	h.request = req
	return nil
}

func TestAccountCommandMapsFlagsToRequest(t *testing.T) {
	handler := &stubAccountHandler{}
	cmd := NewAccountCommand(handler)

	err := cmd.Execute(context.Background(), []string{"--id", "1,2", "--mode", "diff", "--force"})

	require.NoError(t, err)
	require.True(t, handler.called)
	require.Equal(t, []int{1, 2}, handler.request.AccountIDs)
	require.Equal(t, "diff", handler.request.SyncMode)
	require.True(t, handler.request.Force)
}
