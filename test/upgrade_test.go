// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os/exec"
	"regexp"
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
	require.True(t, strings.Contains(out, `Current Orbs Gamma`), "upgrade worked on gamma-server")
	require.True(t, strings.Contains(out, `Current Prism`), "upgrade worked on prism")
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
	require.True(t, strings.Contains(out, `experimental: Pulling from orbsnetwork/gamma`), "experimental upgrade should always try to pull fresh copy (gamma)")
	require.True(t, strings.Contains(out, `experimental: Pulling from orbsnetwork/prism`), "experimental upgrade should always try to pull fresh copy (prism)")

	out, err = cli.Run("start-local")
	t.Log(out)
	require.NoError(t, err, "start Gamma server should succeed")
	require.True(t, strings.Contains(out, `Orbs Gamma personal blockchain experimental`), "started Gamma server should be experimental")
	require.True(t, strings.Contains(out, `Prism blockchain explorer experimental`), "started Prism server should be experimental")
}

func TestStableServerDoesNotRestartWhenVersionUpgradeNotRequired(t *testing.T) {
	cli := GammaCli().WithStableServer()
	defer cli.StopGammaServer()

	out, err := cli.Run("version")
	t.Log(out)
	require.NoError(t, err, "version should succeed")
	require.True(t, strings.Contains(out, `(docker)`), "version output should show docker")

	out, err = cli.Run("start-local")
	require.NoError(t, err, "start Gamma server should succeed")

	out, err = cli.Run("upgrade-server")
	require.NoError(t, err, "upgrade server stable should succeed")
	require.True(t, strings.Contains(out, `does not require upgrade`), "upgrade same tag should not try to pull fresh copy")
	require.True(t, strings.Contains(out, `Current Orbs Gamma`), "upgrade worked on gamma-server")
	require.True(t, strings.Contains(out, `Current Prism`), "upgrade worked on prism")
	require.False(t, strings.Contains(out, "Orbs Gamma personal blockchain stopped."), "Gamma should not restart if no upgrade happened")
	require.False(t, strings.Contains(out, "Prism blockchain explorer stopped."), "Prism should not restart if no upgrade happened")
}

func TestExperimentalServerDoesNotRestartWhenVersionUpgradeNotRequired(t *testing.T) {
	cli := GammaCli().WithExperimentalServer()
	defer cli.StopGammaServer()

	out, err := cli.Run("version")
	t.Log(out)
	require.NoError(t, err, "version should succeed")
	require.True(t, strings.Contains(out, `(docker)`), "version output should show docker")

	out, err = cli.Run("start-local")
	require.NoError(t, err, "start Gamma server should succeed")

	out, err = cli.Run("upgrade-server")
	require.NoError(t, err, "upgrade server stable should succeed")
	require.True(t, strings.Contains(out, `Image is up to date`), "re-upgrade experimental should not try to pull fresh copy")
	require.False(t, strings.Contains(out, "Orbs Gamma personal blockchain stopped."), "Gamma should not restart if no upgrade happened")
	require.False(t, strings.Contains(out, "Prism blockchain explorer stopped."), "Prism should not restart if no upgrade happened")
}

func TestPrismNotRestartedOnUpgradeRestartWhenNotRunningBefore(t *testing.T) {
	// can only test this with stable, because of how docker/gamma automation works
	cli := GammaCli().WithStableServer().WithNoPrism()
	defer cli.StopGammaServer()

	out, err := cli.Run("version")
	t.Log(out)
	require.NoError(t, err, "version should succeed")
	require.True(t, strings.Contains(out, `(docker)`), "version output should show docker")

	// get an old version
	dockerOut, err := exec.Command("docker", "pull", "orbsnetwork/gamma:v1.1.1").CombinedOutput()
	if err != nil {
		t.Fatalf("%s", dockerOut)
	}

	// get latest tag and remove it (so upgrade will happen)
	pattern := fmt.Sprintf(`%s\s+(v\S+)`, regexp.QuoteMeta("Gamma server version"))
	re := regexp.MustCompile(pattern)
	res := re.FindStringSubmatch(out)
	dockerOut, err = exec.Command("docker", "rmi", "orbsnetwork/gamma:"+res[1]).CombinedOutput()
	if err != nil {
		t.Fatalf("%s", dockerOut)
	}

	cli = cli.StartGammaServerAndWait()

	out, err = cli.Run("upgrade-server")
	t.Log(out)
	require.NoError(t, err, "upgrade server stable should succeed")
	require.True(t, strings.Contains(out, `Downloading latest`), "we are forcing an upgrade, it should download something")
	require.True(t, strings.Contains(out, "Orbs Gamma personal blockchain stopped."), "Gamma was upgraded and needs to restart")
	require.False(t, strings.Contains(out, "Prism blockchain explorer stopped."), "Prism should not restart as it was not running before")
}