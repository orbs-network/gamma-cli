package jsoncodec

import (
	"encoding/json"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"strconv"
)

type SendTx struct {
	ContractName string
	MethodName   string
	Arguments    []*Arg
}

func UnmarshalSendTx(bytes []byte) (*SendTx, error) {
	var sendTx *SendTx
	err := json.Unmarshal(bytes, &sendTx)
	return sendTx, err
}

func MarshalSendTxResponse(r *codec.SendTransactionResponse, txId string) ([]byte, error) {
	return json.MarshalIndent(&struct {
		RequestStatus     codec.RequestStatus
		TxId              string
		ExecutionResult   codec.ExecutionResult
		OutputArguments   []*Arg
		OutputEvents      []*Event
		TransactionStatus codec.TransactionStatus
		BlockHeight       string
		BlockTimestamp    string
	}{
		RequestStatus:     r.RequestStatus,
		TxId:              txId,
		ExecutionResult:   r.ExecutionResult,
		OutputArguments:   MarshalArgs(r.OutputArguments),
		OutputEvents:      MarshalEvents(r.OutputEvents),
		TransactionStatus: r.TransactionStatus,
		BlockHeight:       strconv.FormatUint(r.BlockHeight, 10),
		BlockTimestamp:    r.BlockTimestamp.UTC().Format(codec.ISO_DATE_FORMAT),
	}, "", "  ")
}
