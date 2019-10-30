// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package test

import (
	"github.com/stretchr/testify/require"
	"os/exec"
	"strings"
	"testing"
)

func TestRestart(t *testing.T) {
	cli := GammaCli().WithStableServer()
	defer cli.StopGammaServer()

	out, err := cli.Run("start-local")
	t.Log(out)
	require.NoError(t, err, "start Gamma server should succeed")
	require.False(t, strings.Contains(out, `Orbs Gamma experimental personal blockchain`), "started Gamma server should not be experimental")
	require.False(t, strings.Contains(out, `Prism blockchain explorer experimental`), "started Prism server should not be experimental")

	_, err = cli.Run("stop-local")
	require.NoError(t, err, "stop Gamma server should succeed")

	_, err = cli.Run("stop-local")
	require.NoError(t, err, "second stop Gamma server should succeed")

	_, err = cli.Run("start-local")
	require.NoError(t, err, "start Gamma server should succeed")
}

func TestStopAfterCrashOfGammaServer(t *testing.T) {
	cli := GammaCli().WithExperimentalServer()
	defer cli.StopGammaServer()

	out, err := cli.Run("start-local")
	t.Log(out)
	require.NoError(t, err, "start Gamma server should succeed")
	require.True(t, strings.Contains(out, `Orbs Gamma personal blockchain experimental`), "started Gamma server should be experimental")
	require.True(t, strings.Contains(out, `Prism blockchain explorer experimental`), "started Prism server should be experimental")

	// stopping and removing gamma-server
	dockerOut, err := exec.Command("docker", "stop", "orbs-gamma-server").CombinedOutput()
	if err != nil {
		t.Fatalf("%s", dockerOut)
	}
	dockerOut, err = exec.Command("docker", "rm", "-f", "orbs-gamma-server").CombinedOutput()
	if err != nil {
		t.Fatalf("Could not remove docker container.\n\n%s", dockerOut)
	}

	// running the regular gamma stop
	out, err = cli.Run("stop-local")
	require.NoError(t, err, "stop Gamma server should succeed")
	require.True(t, strings.Contains(out, "Prism blockchain explorer stopped."), "Prism server should stop even if gamma crashed")
}
