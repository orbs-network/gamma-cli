package jsoncodec

import (
	"encoding/hex"
	"encoding/json"
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
	}, "", "  ")
}
