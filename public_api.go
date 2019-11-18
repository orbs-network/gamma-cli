// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/orbs-network/gamma-cli/jsoncodec"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path"
	"strings"
	"syscall"
)

const DEPLOY_SYSTEM_CONTRACT_NAME = "_Deployments"
const DEPLOY_SYSTEM_METHOD_NAME = "deployService"
const PROCESSOR_TYPE_NATIVE = uint32(1)
const PROCESSOR_TYPE_JAVASCRIPT = uint32(2)

func _getSource(name string) (code [][]byte, err error) {
	if info, err := os.Stat(name); err != nil {
		return nil, err
	} else if info.IsDir() {
		return orbs.ReadSourcesFromDir(name)
	} else {
		singleFile, err := ioutil.ReadFile(name)
		if err != nil {
			return nil, err
		}

		return [][]byte{singleFile}, err
	}
}

func commandDeploy(requiredOptions []string) {
	codeFile := requiredOptions[0]

	if *flagContractName == "" {
		*flagContractName = getFilenameWithoutExtension(codeFile)
	}

	code, err := _getSource(codeFile)
	if err != nil {
		die("Could not find path\n\n%s", err.Error())
	}

	signer := getTestKeyFromFile(*flagSigner)

	client := createOrbsClient()

	payload, txId, err := client.CreateDeployTransaction(signer.PublicKey, signer.PrivateKey, string(*flagContractName), orbs.PROCESSOR_TYPE_NATIVE, code...)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.SendTransaction(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalSendTxResponse(response, txId)
		if err != nil {
			die("Could not encode send-tx response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
		exit()
	}

	if clientErr != nil {
		die("Request transaction failed on server.\n\n%s", clientErr.Error())
	}
}

func commandSendTx(requiredOptions []string) {
	inputFile := requiredOptions[0]

	signer := getTestKeyFromFile(*flagSigner)

	bytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		die("Could not open input file.\n\n%s", err.Error())
	}

	sendTx, err := jsoncodec.UnmarshalSendTx(bytes)
	if err != nil {
		die("Failed parsing input json file '%s'.\n\n%s", inputFile, err.Error())
	}

	// override contract name
	if *flagContractName != "" {
		sendTx.ContractName = *flagContractName
	}

	overrideArgsWithFlags(sendTx.Arguments)
	inputArgs, err := jsoncodec.UnmarshalArgs(sendTx.Arguments, getTestKeyFromFile)
	if err != nil {
		die(err.Error())
	}

	client := createOrbsClient()

	payload, txId, err := client.CreateTransaction(signer.PublicKey, signer.PrivateKey, sendTx.ContractName, sendTx.MethodName, inputArgs...)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.SendTransaction(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalSendTxResponse(response, txId)
		if err != nil {
			die("Could not encode send-tx response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
		exit()
	}

	if clientErr != nil {
		die("Request send-tx failed on server.\n\n%s", clientErr.Error())
	}
}

func commandRunQuery(requiredOptions []string) {
	inputFile := requiredOptions[0]

	signer := getTestKeyFromFile(*flagSigner)

	bytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		die("Could not open input file.\n\n%s", err.Error())
	}

	runQuery, err := jsoncodec.UnmarshalRead(bytes)
	if err != nil {
		die("Failed parsing input json file '%s'.\n\n%s", inputFile, err.Error())
	}

	// override contract name
	if *flagContractName != "" {
		runQuery.ContractName = *flagContractName
	}

	overrideArgsWithFlags(runQuery.Arguments)
	inputArgs, err := jsoncodec.UnmarshalArgs(runQuery.Arguments, getTestKeyFromFile)
	if err != nil {
		die(err.Error())
	}

	client := createOrbsClient()

	payload, err := client.CreateQuery(signer.PublicKey, runQuery.ContractName, runQuery.MethodName, inputArgs...)
	if err != nil {
		die("Could not encode payload of the message about to be sent to server.\n\n%s", err.Error())
	}

	response, clientErr := client.SendQuery(payload)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalReadResponse(response)
		if err != nil {
			die("Could not encode run-query response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
		exit()
	}

	if clientErr != nil {
		die("Request run-query failed on server.\n\n%s", clientErr.Error())
	}
}

func commandTxStatus(requiredOptions []string) {
	txId := requiredOptions[0]

	client := createOrbsClient()

	response, clientErr := client.GetTransactionStatus(txId)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalTxStatusResponse(response)
		if err != nil {
			die("Could not encode status response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
		exit()
	}

	if clientErr != nil {
		die("Request status failed on server.\n\n%s", clientErr.Error())
	}
}

func commandTxProof(requiredOptions []string) {
	txId := requiredOptions[0]

	client := createOrbsClient()

	response, clientErr := client.GetTransactionReceiptProof(txId)
	handleNoConnectionGracefully(clientErr, client)
	if response != nil {
		output, err := jsoncodec.MarshalTxProofResponse(response)
		if err != nil {
			die("Could not encode tx proof response to json.\n\n%s", err.Error())
		}

		log("%s\n", string(output))
		exit()
	}

	if clientErr != nil {
		die("Request status failed on server.\n\n%s", clientErr.Error())
	}
}

func createOrbsClient() *orbs.OrbsClient {
	env := getEnvironmentFromConfigFile(*flagEnv)
	if len(env.Endpoints) == 0 {
		die("Environment Endpoints key does not contain any endpoints.")
	}

	endpoint := env.Endpoints[0]
	if endpoint == "localhost" {
		if !isDockerContainerRunning(gammaHandlerOptions().containerName) && !isPortListening(gammaHandlerOptions().port) {
			die("Local Gamma server is not running, use 'gamma-cli start-local' to start it.")
		}
		endpoint = fmt.Sprintf("http://localhost:%d", gammaHandlerOptions().port)
	}

	return orbs.NewClient(endpoint, env.VirtualChain, codec.NETWORK_TYPE_TEST_NET)
}

// Will get to it when we implement JS
func getProcessorTypeFromFilename(filename string) uint32 {
	if strings.HasSuffix(filename, ".go") {
		return PROCESSOR_TYPE_NATIVE
	}
	if strings.HasSuffix(filename, ".js") {
		return PROCESSOR_TYPE_JAVASCRIPT
	}
	die("Unsupported code file type.\n\nSupported code file extensions are: .go .js")
	return 0
}

// TODO: this needs to be simplified
func handleNoConnectionGracefully(err error, client *orbs.OrbsClient) {
	msg := fmt.Sprintf("Cannot connect to server at endpoint %s\n\nPlease check that:\n - The server is started and running (if just started, may need a second to initialize).\n - The server is accessible over the network.\n - The endpoint is properly configured if a config file is used.", client.Endpoint)
	switch err := errors.Cause(err).(type) {
	case *url.Error:
		die(msg)
	case *net.OpError:
		if err.Op == "dial" || err.Op == "read" {
			die(msg)
		}
	case net.Error:
		if err.Timeout() {
			die(msg)
		}
	case syscall.Errno:
		if err == syscall.ECONNREFUSED {
			die(msg)
		}
	default:
		if err == orbs.NoConnectionError {
			die(msg)
		}
		return
	}
}

func getFilenameWithoutExtension(filename string) string {
	return strings.Split(path.Base(filename), ".")[0]
}

func overrideArgWithPossibleArray(arg *jsoncodec.Arg, value string) {
	if strings.HasSuffix(arg.Type, "Array") {
		var valueAsArray []interface{}
		err := json.Unmarshal([]byte(value), &valueAsArray)
		if err != nil {
			arg.Value = []interface{}{value}
		} else {
			arg.Value = valueAsArray
		}
	} else {
		arg.Value = value
	}
}

func overrideArgsWithFlags(args []*jsoncodec.Arg) {
	if *flagArg1 != "" {
		overrideArgWithPossibleArray(args[0], *flagArg1)
	}
	if *flagArg2 != "" {
		overrideArgWithPossibleArray(args[1], *flagArg2)
	}
	if *flagArg3 != "" {
		overrideArgWithPossibleArray(args[2], *flagArg3)
	}
	if *flagArg4 != "" {
		overrideArgWithPossibleArray(args[3], *flagArg4)
	}
	if *flagArg5 != "" {
		overrideArgWithPossibleArray(args[4], *flagArg5)
	}
	if *flagArg6 != "" {
		overrideArgWithPossibleArray(args[5], *flagArg6)
	}
	if *flagArg7 != "" {
		overrideArgWithPossibleArray(args[6], *flagArg7)
	}
	if *flagArg8 != "" {
		overrideArgWithPossibleArray(args[7], *flagArg8)
	}
	if *flagArg9 != "" {
		overrideArgWithPossibleArray(args[8], *flagArg9)
	}
}
