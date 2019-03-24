// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import "encoding/json"

type ConfFile struct {
	Environments map[string]*ConfEnv
}

type ConfEnv struct {
	VirtualChain uint32
	Endpoints    []string
	Experimental bool
}

func UnmarshalConfFile(bytes []byte) (*ConfFile, error) {
	var confFile *ConfFile
	err := json.Unmarshal(bytes, &confFile)
	return confFile, err
}
