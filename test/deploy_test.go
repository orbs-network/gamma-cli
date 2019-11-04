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

func TestDeployCounter(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().DownloadLatestGammaServer().StartGammaServerAndWait()
	defer cli.StopGammaServer()

	out, err := cli.Run("deploy", "./_counter/contract.go", "-name", "CounterExample")
	t.Log(out)
	require.NoError(t, err, "deploy should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))

	out, err = cli.Run("run-query", "counter-get.json")
	t.Log(out)
	require.NoError(t, err, "get should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))
	require.True(t, strings.Contains(out, `"Value": "0"`))

	out, err = cli.Run("send-tx", "counter-add.json")
	t.Log(out)
	require.NoError(t, err, "add should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))
	require.True(t, strings.Contains(out, `"Value": "previous count is 0"`))

	out, err = cli.Run("run-query", "counter-get.json")
	t.Log(out)
	require.NoError(t, err, "get should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))
	require.True(t, strings.Contains(out, `"Value": "25"`))
}

func TestDeployCorruptContract(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().DownloadLatestGammaServer().StartGammaServerAndWait()
	defer cli.StopGammaServer()

	out, err := cli.Run("deploy", "./_corrupt/corrupt.go", "-name", "CounterExample")
	t.Log(out)
	require.NoError(t, err, "deploy should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "ERROR_SMART_CONTRACT"`))
	require.True(t, strings.Contains(out, `compilation of deployable contract 'CounterExample' failed`))
}

func TestDeployOfAlreadyDeployed(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().DownloadLatestGammaServer().StartGammaServerAndWait()
	defer cli.StopGammaServer()

	out, err := cli.Run("deploy", "./_counter/contract.go", "-name", "CounterExample")
	t.Log(out)
	require.NoError(t, err, "deploy should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "SUCCESS"`))

	out, err = cli.Run("deploy", "./_counter/contract.go", "-name", "CounterExample")
	t.Log(out)
	require.NoError(t, err, "deploy should succeed")
	require.True(t, strings.Contains(out, `"ExecutionResult": "ERROR_SMART_CONTRACT"`))
	require.True(t, strings.Contains(out, `a contract with same name (case insensitive) already exists`))
}

func TestRunMethodWithoutDeploy(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().StartGammaServerAndWait()
	defer cli.StopGammaServer()

	out, err := cli.Run("send-tx", "counter-add.json")
	t.Log(out)
	require.NoError(t, err, "add should succeed")
	require.True(t, strings.Contains(out, `"RequestStatus": "BAD_REQUEST"`))
	require.True(t, strings.Contains(out, `"ExecutionResult": "ERROR_CONTRACT_NOT_DEPLOYED"`))
}
