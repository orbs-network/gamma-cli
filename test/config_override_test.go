package test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGammaCli_StartWithConfigOverrides(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().DownloadLatestGammaServer().WithConfigOverrides(`{"virtual-chain-id":43}`).StartGammaServerAndWait()
	defer cli.StopGammaServer()

	out, err := cli.Run("start-local")
	t.Log(out)
	require.NoError(t, err, "start Gamma server should succeed")

	out, _ = cli.Run("deploy", "./_counter/contract.go", "-name", "CounterExample")
	require.Contains(t, out, "REJECTED_VIRTUAL_CHAIN_MISMATCH", "transaction was not rejected when sending a transaction for vcid 42 to a container running vcid 43")

}
