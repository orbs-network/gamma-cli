// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import "encoding/json"

type Key struct {
	PrivateKey string // hex string starting with 0x
	PublicKey  string // hex string starting with 0x
	Address    string // hex string starting with 0x
}

type RawKey struct {
	PrivateKey []byte
	PublicKey  []byte
	Address    []byte
}

func UnmarshalKeys(bytes []byte) (map[string]*Key, error) {
	keys := make(map[string]*Key)
	err := json.Unmarshal(bytes, &keys)
	return keys, err
}

func MarshalKeys(keys map[string]*Key) ([]byte, error) {
	return json.MarshalIndent(keys, "", "  ")
}
