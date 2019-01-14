package main

import (
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/events"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/state"
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
	events.EmitEvent(Log, fmt.Sprintf("previous count is %d", count))
	count += amount
	state.WriteUint64(COUNTER_KEY, count)
}

func get() uint64 {
	return state.ReadUint64(COUNTER_KEY)
}

func start() uint64 {
	return 0
}
