package test

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSimpleTransfer(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().DownloadLatestGammaServer().StartGammaServer()
	defer cli.StopGammaServer()

	out, err := cli.Run("run-query", "get-balance.json")
	t.Log(out)
	require.NoError(t, err, "get balance should not fail (although not deployed)")
	require.True(t, strings.Contains(out, `"ExecutionResult": "ERROR_CONTRACT_NOT_DEPLOYED"`))

	out, err = cli.Run("send-tx", "transfer.json")
	t.Log(out)
	require.NoError(t, err, "transfer should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))

	txId := extractTxIdFromSendTxOutput(out)
	t.Log(txId)

	out, err = cli.Run("tx-status", txId)
	t.Log(out)
	require.NoError(t, err, "get tx status should succeed")
	require.True(t, strings.Contains(out, `"RequestStatus": "COMPLETED"`))

	out, err = cli.Run("tx-proof", txId)
	t.Log(out)
	require.NoError(t, err, "get tx proof should succeed")
	require.True(t, strings.Contains(out, `"RequestStatus": "COMPLETED"`))
	require.True(t, strings.Contains(out, `"PackedProof": "0x`))
	require.True(t, strings.Contains(out, `"PackedReceipt": "0x`))

	out, err = cli.Run("send-tx", "transfer.json", "-arg1", "2")
	t.Log(out)
	require.NoError(t, err, "transfer should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))

	out, err = cli.Run("run-query", "get-balance.json")
	t.Log(out)
	require.NoError(t, err, "get balance should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))
	require.True(t, strings.Contains(out, `"Value": "19"`))

	out, err = cli.Run("send-tx", "transfer-direct.json")
	t.Log(out)
	require.NoError(t, err, "transfer should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))
}
