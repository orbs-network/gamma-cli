# Gamma - Personal Orbs Blockchain for Developers 

Gamma is a personal Orbs blockchain that allows developers to easily test, run and deploy smart contracts.

`Gamma server` - runs an in-memory virtual chain on top of an Orbs blockchain with several nodes on your local machine.

`gamma-cli` - command line tool for developers to interact with a Gamma server instance running on their machine. 

&nbsp;

## Detailed documentation

The detailed documentation website for Gamma is available here:

https://orbs.gitbook.io

&nbsp;

## Quick start

### Prerequisites (Mac)

* Make sure [brew](https://brew.sh/) is available on your machine.

* Make sure [Docker](https://docs.docker.com/docker-for-mac/install/) is installed on your machine.

* If you're planning to develop your own smart contracts in Go, install the [Orbs Smart Contract SDK](https://github.com/orbs-network/orbs-contract-sdk#installation).

### Installation

1. To install the command line tool, run the following command in terminal:

    ```
    brew install orbs-network/devtools/gamma-cli
    ```
    
    > To verify the installation, run in terminal `gamma-cli version`
    
2. Gamma server will automatically be installed the first time you start it with `gamma-cli start-local`

### Starting and stopping Gamma server 

* Start Gamma server by running in terminal:

    ```
    gamma-cli start-local
    ```

* When finished working with the server, stop it by running in terminal:

    ```
    gamma-cli stop-local
    ```
    
    > Note: The local blockchain instance is running in-memory. The next time you start the instance, all contracts and state will disappear from memory and you will need to deploy them again.

&nbsp;

## Commands

```
Usage:

gamma-cli COMMAND [OPTIONS]

Commands:

  start-local      start a local Orbs personal blockchain instance listening on port
                   options: -port <PORT>
                   example: gamma-cli start-local -port 8080

  stop-local       stop a locally running Orbs personal blockchain instance

  gen-test-keys    generate a new batch of 10 test keys and store in orbs-test-keys.json (default filename)
                   options: -keys [OUTPUT_FILE]
                   example: gamma-cli gen-test-keys -keys orbs-test-keys.json

  deploy           deploy a smart contract with the code specified in the source file <CODE_FILE>
                   options: <CODE_FILE> -name [CONTRACT_NAME] -signer [ID_FROM_KEYS_JSON]
                   example: gamma-cli deploy MyToken.go -signer user1
                            gamma-cli deploy contract.go -name MyToken

  send-tx          sign and send the transaction specified in the JSON file <INPUT_FILE>
                   options: <INPUT_FILE> -arg# [OVERRIDE_ARG_#] -signer [ID_FROM_KEYS_JSON]
                   example: gamma-cli send-tx transfer.json -signer user1
                            gamma-cli send-tx transfer.json -arg2 0x5B63Ca66637316A0D7f84Ebf60E50963c10059aD

  run-query        read state or run a read-only contract method as specified in the JSON file <INPUT_FILE>
                   options: <INPUT_FILE> -arg# [OVERRIDE_ARG_#] -signer [ID_FROM_KEYS_JSON]
                   example: gamma-cli run-query get-balance.json -signer user1
                            gamma-cli run-query get-balance.json -arg1 0x5B63Ca66637316A0D7f84Ebf60E50963c10059aD

  tx-status        get the current status of a sent transaction with txid <TX_ID> (from send-tx response)
                   options: <TX_ID>
                   example: gamma-cli tx-status 0xB68fa95B7f397815Ddf41150d79b27a888448a22e08DeAf8600E7a495c406303659f8C3782614660

  tx-proof         get cryptographic proof for transaction receipt with txid <TX_ID> (from send-tx response)
                   options: <TX_ID>
                   example: gamma-cli tx-proof 0xB68fa95B7f397815Ddf41150d79b27a888448a22e08DeAf8600E7a495c406303659f8C3782614660

  upgrade-server   upgrade to the latest stable version of Gamma server
                   example: gamma-cli upgrade-server
                            gamma-cli upgrade-server -env experimental

  version          print gamma-cli and Gamma server versions

  help             print this help screen


Options:

  -config string
      path to config file (default "orbs-gamma-config.json")
  -env string
      environment from config file containing server connection details (default "local")
  -keys string
      name of the json file containing test keys (default "orbs-test-keys.json")
  -name string
      name of the smart contract being deployed
  -port int
      listening port for Gamma server (default "8080")
  -signer string
      id of the signing key from the test key json (default "user1")
  -wait
      wait until Gamma server is ready and listening

Multiple environments (eg. local and testnet) can be defined in orbs-gamma-config.json configuration file.
See https://orbs.gitbook.io for more info.
```

&nbsp;

## Upgrading to latest stable versions

* Upgrade to the latest version of `gamma-cli` by running in terminal:

    ```
    brew upgrade gamma-cli
    ```

* Upgrade to the latest version of Gamma server by running in terminal:

    ```
    gamma-cli upgrade-server
    ```

&nbsp;

## Working with latest dev versions (experimental)

* Upgrade to the latest dev version of `gamma-cli` by running in terminal:

    ```
    brew upgrade gamma-cli --devel
    ```

* Upgrade to the latest dev version of Gamma server by running in terminal:

    ```
    gamma-cli upgrade-server -env experimental
    ```
    
* Start the experimental Gamma server by running in terminal:

    ```
    gamma-cli start-local -env experimental
    ```

## License

MIT
