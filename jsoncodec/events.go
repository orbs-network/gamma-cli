// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

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
