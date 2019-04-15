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

func TestUpgradeStableServer(t *testing.T) {
	cli := GammaCli().WithStableServer()

	out, err := cli.Run("version")
	t.Log(out)
	require.NoError(t, err, "version should succeed")
	require.True(t, strings.Contains(out, `(docker)`), "version output should show docker")

	out, err = cli.Run("upgrade-server")
	t.Log(out)
	require.NoError(t, err, "upgrade server stable should succeed")
	require.True(t, strings.Contains(out, `does not require upgrade`), "upgrade same tag should not try to pull fresh copy")
}

func TestUpgradeExperimentalServer(t *testing.T) {
	cli := GammaCli().WithExperimentalServer()
	defer cli.StopGammaServer()

	out, err := cli.Run("version")
	t.Log(out)
	require.NoError(t, err, "version experimental should succeed")
	require.True(t, strings.Contains(out, `experimental (docker)`), "version output should show experimental docker")

	out, err = cli.Run("upgrade-server")
	t.Log(out)
	require.NoError(t, err, "upgrade server experimental should succeed")
	require.True(t, strings.Contains(out, `experimental: Pulling from orbsnetwork/gamma`), "experimental upgrade should always try to pull fresh copy")

	out, err = cli.Run("start-local")
	t.Log(out)
	require.NoError(t, err, "start Gamma server should succeed")
	require.True(t, strings.Contains(out, `Orbs Gamma personal blockchain experimental`), "started Gamma server should not be experimental")
}
