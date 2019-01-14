package jsoncodec

import "github.com/orbs-network/orbs-client-sdk-go/codec"

type Event struct {
	ContractName string
	EventName    string
	Arguments    []*Arg
}

func MarshalEvents(events []*codec.Event) []*Event {
	res := []*Event{}
	for _, event := range events {
		res = append(res, &Event{
			ContractName: event.ContractName,
			EventName:    event.EventName,
			Arguments:    MarshalArgs(event.Arguments),
		})
	}
	return res
}
