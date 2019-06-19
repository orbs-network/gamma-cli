// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/events"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
	"strconv"
)

var PUBLIC = sdk.Export(add, get, start)
var SYSTEM = sdk.Export(_init)
var EVENTS = sdk.Export(Log)

var COUNTER_KEY = []byte("count")

func Log(msg string) {}

func _init() {
	state.WriteUint64(COUNTER_KEY, 0)
}

func add(amount uint64) {
	count := state.ReadUint64(COUNTER_KEY)
	events.EmitEvent(Log, "previous count is " + strconv.FormatUint(count, 10))
	count += amount
	state.WriteUint64(COUNTER_KEY, count)
}

func get() uint64 {
	return state.ReadUint64(COUNTER_KEY)
}

func start() uint64 {
	return 0
}
