package main

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/crypto/encoding"
	"github.com/orbs-network/orbs-client-sdk-go/gammacli/jsoncodec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"io/ioutil"
)

func commandGenerateTestKeys(requiredOptions []string) {
	keys := make(map[string]*jsoncodec.Key)
	for i := 0; i < 10; i++ {
		account, err := orbs.CreateAccount()
		if err != nil {
			die("Could not create Orbs account.")
		}
		user := fmt.Sprintf("user%d", i+1)
		keys[user] = &jsoncodec.Key{
			PrivateKey: encoding.EncodeHex(account.PrivateKey),
			PublicKey:  encoding.EncodeHex(account.PublicKey),
			Address:    account.Address,
		}
	}

	bytes, err := jsoncodec.MarshalKeys(keys)
	if err != nil {
		die("Could not encode keys to json.\n\n%s", err.Error())
	}

	filename := *flagKeyFile
	if filename == "" {
		filename = TEST_KEYS_FILENAME
	}
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		die("Could not write keys to file.\n\n%s", err.Error())
	}

	if !doesFileExist(filename) {
		die("File not found after write.")
	}

	log("10 new test keys written successfully to '%s'.\n", filename)
}

func getTestKeyFromFile(id string) *jsoncodec.RawKey {
	if !doesFileExist(*flagKeyFile) {
		commandGenerateTestKeys(nil)
	}

	bytes, err := ioutil.ReadFile(*flagKeyFile)
	if err != nil {
		die("Could not open keys file '%s'.\n\n%s", *flagKeyFile, err.Error())
	}

	keys, err := jsoncodec.UnmarshalKeys(bytes)
	if err != nil {
		die("Failed parsing keys json file '%s'. Try deleting the key file to have it automatically recreated.\n\n%s", *flagKeyFile, err.Error())
	}

	key, found := keys[id]
	if !found {
		die("Key with id '%s' not found in key file '%s'.", id, *flagKeyFile)
	}

	privateKey, err := encoding.DecodeHex(key.PrivateKey)
	if err != nil {
		die("Could not parse hex string '%s'. Try deleting the key file '%s' to have it automatically recreated.\n\n%s", privateKey, *flagKeyFile, err.Error())
	}

	publicKey, err := encoding.DecodeHex(key.PublicKey)
	if err != nil {
		die("Could not parse hex string '%s'. Try deleting the key file '%s' to have it automatically recreated.\n\n%s", publicKey, *flagKeyFile, err.Error())
	}

	address, err := encoding.DecodeHex(key.Address)
	if err != nil {
		die("Could not parse hex string '%s'. Try deleting the key file '%s' to have it automatically recreated.\n\n%s", address, *flagKeyFile, err.Error())
	}

	res := &jsoncodec.RawKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}

	return res
}
