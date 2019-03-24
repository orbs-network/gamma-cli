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

func TestHelp(t *testing.T) {
	out, err := GammaCli().WithExperimentalServer().Run("help")
	t.Log(out)
	require.Error(t, err, "help should exit nonzero")
	require.NotEmpty(t, out, "help output should not be empty")
	require.True(t, strings.Contains(out, "start-local"))
	require.True(t, strings.Contains(out, "stop-local"))

	out2, err := GammaCli().WithExperimentalServer().Run()
	require.Error(t, err, "run without arguments should exit nonzero")
	require.Equal(t, out, out2, "help output should be equal")
}

func TestVersion(t *testing.T) {
	out, err := GammaCli().WithStableServer().Run("version")
	t.Log(out)
	require.NoError(t, err, "version should succeed")
	require.True(t, strings.Contains(out, "version"))
	require.False(t, strings.Contains(out, `version experimental (docker)`), "started Gamma server should not be experimental")
}
