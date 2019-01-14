package jsoncodec

import (
	"encoding/json"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
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
		OutputArguments: MarshalArgs(r.OutputArguments),
		OutputEvents:    MarshalEvents(r.OutputEvents),
		BlockHeight:     strconv.FormatUint(r.BlockHeight, 10),
		BlockTimestamp:  r.BlockTimestamp.UTC().Format(codec.ISO_DATE_FORMAT),
	}, "", "  ")
}
