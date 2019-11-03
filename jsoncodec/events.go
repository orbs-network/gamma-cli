// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/pkg/errors"
)

type Event struct {
	ContractName string
	EventName    string
	Arguments    []*Arg
}

func MarshalEvents(events []*codec.Event) ([]*Event, error) {
	res := []*Event{}
	for i, event := range events {
		eventArgs, err := MarshalArgs(event.Arguments)
		if err != nil {
			return nil, errors.Errorf("Event %d arguments marshaling failed with %s \n", i+1, err.Error())
		}
		res = append(res, &Event{
			ContractName: event.ContractName,
			EventName:    event.EventName,
			Arguments:    eventArgs,
		})
	}
	return res, nil
}
