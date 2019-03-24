// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const GAMMA_CLI_VERSION = "0.6.6"
const CONFIG_FILENAME = "orbs-gamma-config.json"
const TEST_KEYS_FILENAME = "orbs-test-keys.json"
const LOCAL_ENV_ID = "local"
const EXPERIMENTAL_ENV_ID = "experimental"

type command struct {
	desc            string
	args            string
	example         string
	example2        string
	handler         func([]string)
	sort            int
	requiredOptions []string
}

var commands = map[string]*command{
	"start-local": {
		desc:            "start a local Orbs personal blockchain instance listening on port",
		args:            "-port <PORT>",
		example:         "gamma-cli start-local -port 8080",
		handler:         commandStartLocal,
		sort:            0,
		requiredOptions: nil,
	},
	"stop-local": {
		desc:            "stop a locally running Orbs personal blockchain instance",
		handler:         commandStopLocal,
		sort:            1,
		requiredOptions: nil,
	},
	"gen-test-keys": {
		desc:            "generate a new batch of 10 test keys and store in " + TEST_KEYS_FILENAME + " (default filename)",
		args:            "-keys [OUTPUT_FILE]",
		example:         "gamma-cli gen-test-keys -keys " + TEST_KEYS_FILENAME,
		handler:         commandGenerateTestKeys,
		sort:            2,
		requiredOptions: nil,
	},
	"deploy": {
		desc:            "deploy a smart contract with the code specified in the source file <CODE_FILE>",
		args:            "<CODE_FILE> -name [CONTRACT_NAME] -signer [ID_FROM_KEYS_JSON]",
		example:         "gamma-cli deploy MyToken.go -signer user1",
		example2:        "gamma-cli deploy contract.go -name MyToken",
		handler:         commandDeploy,
		sort:            3,
		requiredOptions: []string{"<CODE_FILE> - path of file with source code"},
	},
	"send-tx": {
		desc:            "sign and send the transaction specified in the JSON file <INPUT_FILE>",
		args:            "<INPUT_FILE> -arg# [OVERRIDE_ARG_#] -signer [ID_FROM_KEYS_JSON]",
		example:         "gamma-cli send-tx transfer.json -signer user1",
		example2:        "gamma-cli send-tx transfer.json -arg2 0x5B63Ca66637316A0D7f84Ebf60E50963c10059aD",
		handler:         commandSendTx,
		sort:            4,
		requiredOptions: []string{"<INPUT_FILE> - path of JSON file with transaction details"},
	},
	"run-query": {
		desc:            "read state or run a read-only contract method as specified in the JSON file <INPUT_FILE>",
		args:            "<INPUT_FILE> -arg# [OVERRIDE_ARG_#] -signer [ID_FROM_KEYS_JSON]",
		example:         "gamma-cli run-query get-balance.json -signer user1",
		example2:        "gamma-cli run-query get-balance.json -arg1 0x5B63Ca66637316A0D7f84Ebf60E50963c10059aD",
		handler:         commandRunQuery,
		sort:            5,
		requiredOptions: []string{"<INPUT_FILE> - path of JSON file with query details"},
	},
	"tx-status": {
		desc:            "get the current status of a sent transaction with txid <TX_ID> (from send-tx response)",
		args:            "<TX_ID>",
		example:         "gamma-cli tx-status 0xB68fa95B7f397815Ddf41150d79b27a888448a22e08DeAf8600E7a495c406303659f8C3782614660",
		handler:         commandTxStatus,
		sort:            6,
		requiredOptions: []string{"<TX_ID> - txid of previously sent transaction, from send-tx response"},
	},
	"tx-proof": {
		desc:            "get cryptographic proof for transaction receipt with txid <TX_ID> (from send-tx response)",
		args:            "<TX_ID>",
		example:         "gamma-cli tx-proof 0xB68fa95B7f397815Ddf41150d79b27a888448a22e08DeAf8600E7a495c406303659f8C3782614660",
		handler:         commandTxProof,
		sort:            7,
		requiredOptions: []string{"<TX_ID> - txid of previously sent transaction, from send-tx response"},
	},
	"upgrade-server": {
		desc:            "upgrade to the latest stable version of Gamma server",
		example:         "gamma-cli upgrade-server",
		example2:        "gamma-cli upgrade-server -env experimental",
		handler:         commandUpgradeServer,
		sort:            8,
		requiredOptions: nil,
	},
	"version": {
		desc:            "print gamma-cli and Gamma server versions",
		handler:         commandVersion,
		sort:            9,
		requiredOptions: nil,
	},
	"help": {
		desc:            "print this help screen",
		sort:            10,
		requiredOptions: nil,
	},
}

var (
	flagPort         = flag.Int("port", 8080, "listening port for Gamma server")
	flagSigner       = flag.String("signer", "user1", "id of the signing key from the test key json")
	flagContractName = flag.String("name", "", "name of the smart contract being deployed")
	flagKeyFile      = flag.String("keys", TEST_KEYS_FILENAME, "name of the json file containing test keys")
	flagConfigFile   = flag.String("config", CONFIG_FILENAME, "path to config file")
	flagEnv          = flag.String("env", LOCAL_ENV_ID, "environment from config file containing server connection details")
	flagWait         = flag.Bool("wait", false, "wait until Gamma server is ready and listening")

	// args (hidden from help)
	flagArg1 = flag.String("arg1", "", "")
	flagArg2 = flag.String("arg2", "", "")
	flagArg3 = flag.String("arg3", "", "")
	flagArg4 = flag.String("arg4", "", "")
	flagArg5 = flag.String("arg5", "", "")
	flagArg6 = flag.String("arg6", "", "")
)

func main() {
	flag.Usage = func() { commandShowHelp(nil) }
	commands["help"].handler = commandShowHelp

	if len(os.Args) <= 1 {
		commandShowHelp(nil)
	}
	cmdName := os.Args[1]
	cmd, found := commands[cmdName]
	if !found {
		die("Command '%s' not found, run 'gamma-cli help' to see available commands.", cmdName)
	}

	requiredOptions := []string{}
	if len(cmd.requiredOptions) > 0 {
		if len(os.Args) < 2+len(cmd.requiredOptions) {
			die("Command '%s' is missing required arguments %v.", cmdName, cmd.requiredOptions)
		}
		requiredOptions = os.Args[2 : 2+len(cmd.requiredOptions)]
		for i, requiredOption := range requiredOptions {
			if strings.HasPrefix(requiredOption, "-") {
				die("Command '%s' argument %d should be %s.", cmdName, i+1, cmd.requiredOptions[i])
			}
		}
	}

	os.Args = os.Args[2+len(cmd.requiredOptions)-1:]
	flag.Parse()

	cmd.handler(requiredOptions)
}

func log(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprintf(os.Stdout, "\n")
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR:\n  ")
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n\n")
	os.Exit(2)
}

func exit() {
	os.Exit(0)
}

func doesFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
