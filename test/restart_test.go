package test

import (
	"github.com/stretchr/testify/require"
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

	_, err = cli.Run("stop-local")
	require.NoError(t, err, "stop Gamma server should succeed")

	_, err = cli.Run("stop-local")
	require.NoError(t, err, "second stop Gamma server should succeed")

	_, err = cli.Run("start-local")
	require.NoError(t, err, "start Gamma server should succeed")
}

func TestStartedButNotReadyMessage(t *testing.T) {
	cli := GammaCli().WithExperimentalServer()
	defer cli.StopGammaServer()

	_, err := cli.Run("start-local") // without -wait
	require.NoError(t, err, "start Gamma server should succeed")

	out, err := cli.Run("send-tx", "transfer.json")
	t.Log(out)

	require.True(t, strings.Contains(out, `may need a second to initialize`))
}
