// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestArgumentsUnMarshalingTypes(t *testing.T) {
	expectedBigInt := big.NewInt(0)
	expectedBigInt.SetBytes([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})

	tests := []struct {
		name      string
		shouldErr bool
		arg       *Arg
		native    interface{}
	}{
		{"uint32", false, &Arg{"uint32", "19480514"}, uint32(19480514)},
		{"uint32-fail", true, &Arg{"uint32", "bad text"}, uint32(0)},
		{"uint64", false, &Arg{"uint64", "19480514000000000"}, uint64(19480514000000000)},
		{"uint64-fail", true, &Arg{"uint64", "bad text"}, uint64(0)},
		{"string", false, &Arg{"string", "hello my name is ?"}, "hello my name is ?"},
		{"bytes", false, &Arg{"bytes", "ffee00eeff"}, []byte{0xff, 0xee, 0x00, 0xee, 0xff}},
		{"bytes-fail", true, &Arg{"bytes", "yyyy"}, []byte{}},
		{"bool-true", false, &Arg{"bool", "1"}, true},
		{"bool-false", false, &Arg{"bool", "0"}, false},
		{"bool-fail", true, &Arg{"bool", "2"}, false},
		{"bytes20", false, &Arg{"bytes20", "0011223344556677889900112233445566778899"}, [20]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}},
		{"bytes20-fail", true, &Arg{"bytes20", "yyyy"}, [20]byte{}},
		{"bytes20-fail-size", true, &Arg{"bytes20", "00112233445566778899001122334455667788"}, [20]byte{}},
		{"bytes32", false, &Arg{"bytes32", "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"}, [32]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
		{"bytes32-fail", true, &Arg{"bytes32", "yyyy"}, [32]byte{}},
		{"bytes32-fail-size", true, &Arg{"bytes32", "00112233445566778899aabbccddeeff001122334455667788"}, [32]byte{}},
		{"bigint", false, &Arg{"uint256", "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"}, expectedBigInt},
		{"bigint-fail", true, &Arg{"uint256", "yyyy"}, nil},
		{"bigint-fail-size", true, &Arg{"uint256", "00112233445566778899aabbccddeeff001122334455667788"}, nil},
		{"unknown type string", true, &Arg{"uint8", "19480514"}, uint32(0)},
	}

	for _, cTest := range tests {
		argList := []*Arg{cTest.arg}
		nativeList := []interface{}{cTest.native}

		resNativeList, err := UnmarshalArgs(argList, func(string) *RawKey { return nil })
		if cTest.shouldErr {
			require.Error(t, err, "unmarshal %s should fail", cTest.name)
		} else {
			require.NoError(t, err, "unmarshal %s should not fail", cTest.name)
			require.EqualValues(t, nativeList, resNativeList)
		}
	}
}

func TestArgumentsMarshaling_GoodFlow(t *testing.T) {
	aBigInt := big.NewInt(0)
	aBigInt.SetBytes([]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})

	tests := []struct {
		name   string
		arg    *Arg
		native interface{}
	}{
		{"uint32", &Arg{"uint32", "19480514"}, uint32(19480514)},
		{"uint64", &Arg{"uint64", "19480514000000000"}, uint64(19480514000000000)},
		{"string", &Arg{"string", "hello my name is ?"}, "hello my name is ?"},
		{"bytes", &Arg{"bytes", "0xffee00eeff"}, []byte{0xff, 0xee, 0x00, 0xee, 0xff}},
		{"bool-true", &Arg{"bool", "1"}, true},
		{"bool-false", &Arg{"bool", "0"}, false},
		{"bytes20", &Arg{"bytes20", "0x0011223344556677889900112233445566778899"}, [20]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}},
		{"bytes32", &Arg{"bytes32", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"}, [32]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
		{"bigint", &Arg{"uint256", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"}, aBigInt},
	}

	for _, cTest := range tests {
		argList := []*Arg{cTest.arg}
		nativeList := []interface{}{cTest.native}

		resArgList, err := MarshalArgs(nativeList)
		require.NoError(t, err, "unmarshal %s should not fail", cTest.name)
		require.EqualValues(t, argList, resArgList)
	}
}

func TestArgumentsMarshaling_BadFlow(t *testing.T) {
	nativeList := []interface{}{1.0}
	_, err := MarshalArgs(nativeList)
	require.Error(t, err, "unmarshal %s should fail")
}
