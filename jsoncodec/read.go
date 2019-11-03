// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import (
	"encoding/json"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/pkg/errors"
	"strconv"
)

type Read struct {
	ContractName string
	MethodName   string
	Arguments    []*Arg
}

func UnmarshalRead(bytes []byte) (*Read, error) {
	var read *Read
	err := json.Unmarshal(bytes, &read)
	return read, err
}

func MarshalReadResponse(r *codec.RunQueryResponse) ([]byte, error) {
	outputArgs, err := MarshalArgs(r.OutputArguments)
	if err != nil {
		return nil, errors.Errorf("Read response marshaling output arguments failed with %s \n", err.Error())
	}
	outputEvents, err := MarshalEvents(r.OutputEvents)
	if err != nil {
		return nil, errors.Errorf("Read response marshaling output events failed with %s \n", err.Error())
	}
	return json.MarshalIndent(&struct {
		RequestStatus   codec.RequestStatus
		ExecutionResult codec.ExecutionResult
		OutputArguments []*Arg
		OutputEvents    []*Event
		BlockHeight     string
		BlockTimestamp  string
	}{
		RequestStatus:   r.RequestStatus,
		ExecutionResult: r.ExecutionResult,
		OutputArguments: outputArgs,
		OutputEvents:    outputEvents,
		BlockHeight:     strconv.FormatUint(r.BlockHeight, 10),
		BlockTimestamp:  r.BlockTimestamp.UTC().Format(codec.ISO_DATE_FORMAT),
	}, "", "  ")
}
