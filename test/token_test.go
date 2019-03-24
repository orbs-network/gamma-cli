// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

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
	require.True(t, strings.Contains(out, `"ProofSigners": [`))
	require.True(t, strings.Contains(out, `"0xa328846cd5b4979d68a8c58a9bdfeee657b34de7"`))

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
