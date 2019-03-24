// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import (
	"encoding/hex"
	"encoding/json"
	"github.com/orbs-network/gamma-cli/crypto/digest"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"strconv"
)

func MarshalTxProofResponse(r *codec.GetTransactionReceiptProofResponse) ([]byte, error) {
	return json.MarshalIndent(&struct {
		RequestStatus     codec.RequestStatus
		ExecutionResult   codec.ExecutionResult
		OutputArguments   []*Arg
		OutputEvents      []*Event
		TransactionStatus codec.TransactionStatus
		BlockHeight       string
		BlockTimestamp    string
		PackedProof       string
		PackedReceipt     string
		ProofSigners      []string
	}{
		RequestStatus:     r.RequestStatus,
		ExecutionResult:   r.ExecutionResult,
		OutputArguments:   MarshalArgs(r.OutputArguments),
		OutputEvents:      MarshalEvents(r.OutputEvents),
		TransactionStatus: r.TransactionStatus,
		BlockHeight:       strconv.FormatUint(r.BlockHeight, 10),
		BlockTimestamp:    r.BlockTimestamp.UTC().Format(codec.ISO_DATE_FORMAT),
		PackedProof:       "0x" + hex.EncodeToString(r.PackedProof),
		PackedReceipt:     "0x" + hex.EncodeToString(r.PackedReceipt),
		ProofSigners:      getProofSignersFromPackedProof(r.PackedProof),
	}, "", "  ")
}

func getProofSignersFromPackedProof(packedProof []byte) []string {
	nodeAddresses, err := digest.GetBlockSignersFromReceiptProof(packedProof)
	if err != nil {
		return nil
	}
	var res []string
	for _, nodeAddress := range nodeAddresses {
		signerString := "0x" + hex.EncodeToString(nodeAddress)
		res = append(res, signerString)
	}
	return res
}
