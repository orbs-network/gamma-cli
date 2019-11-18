// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_Arguments(t *testing.T) {
	cli := GammaCli().WithExperimentalServer().DownloadLatestGammaServer().StartGammaServerAndWait()
	defer cli.StopGammaServer()

	outputString, err := cli.Run("deploy", "./_arguments/arguments.go", "-name", "Arguments")
	t.Log(outputString)
	require.NoError(t, err, "deploy should succeed")
	require.True(t, strings.Contains(outputString, `"ExecutionResult": "SUCCESS"`))

	outputString, err = cli.Run("run-query", "arguments-get.json")
	t.Log(outputString)
	require.NoError(t, err, "get should succeed")
	require.True(t, strings.Contains(outputString, `"ExecutionResult": "SUCCESS"`))
	outputGetParsed := struct {
		OutputArguments []*struct {
			Type  string
			Value []string
		}
	}{}
	err = json.Unmarshal([]byte(outputString), &outputGetParsed)
	if err != nil {
		t.Log(err)
	}
	require.Len(t, outputGetParsed.OutputArguments, 8, "There should be 8 output arrays")

	outputString, err = cli.Run("run-query", "arguments-check.json")
	t.Log(outputString)
	require.NoError(t, err, "check should succeed")
	require.True(t, strings.Contains(outputString, `"ExecutionResult": "SUCCESS"`))
	outputCheckParsed := struct {
		OutputArguments []*struct {
			Type  string
			Value string
		}
	}{}
	err = json.Unmarshal([]byte(outputString), &outputCheckParsed)
	if err != nil {
		t.Log(err)
	}
	require.Len(t, outputCheckParsed.OutputArguments, 2, "There should be 2 output values")
	require.Equal(t, "bool", outputCheckParsed.OutputArguments[0].Type)
	require.Equal(t, "1", outputCheckParsed.OutputArguments[0].Value)
	require.Equal(t, "string", outputCheckParsed.OutputArguments[1].Type)
	require.Empty(t, outputCheckParsed.OutputArguments[1].Value)
}
