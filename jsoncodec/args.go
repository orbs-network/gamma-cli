// Copyright 2019 the gamma-cli authors
// This file is part of the gamma-cli library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package jsoncodec

import (
	"encoding/hex"
	"github.com/orbs-network/orbs-client-sdk-go/crypto/encoding"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
	"strings"
)

const supported = "Supported types are: uint32 uint64 uint256 bool string bytes bytes20 bytes32 gamma:address gamma:keys-file-address"

type Arg struct {
	Type  string
	Value string
}

func UnmarshalArgs(args []*Arg, getTestKeyFromFile func(string) *RawKey) ([]interface{}, error) {
	res := []interface{}{}
	for i, arg := range args {
		switch arg.Type {
		case "uint32":
			val, err := strconv.ParseUint(arg.Value, 10, 32)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the numeric value\n\nCurrent value: '%s'", i+1, arg.Value)
			}
			res = append(res, uint32(val))
		case "uint64":
			val, err := strconv.ParseUint(arg.Value, 10, 64)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the numeric value\n\nCurrent value: '%s'", i+1, arg.Value)
			}
			res = append(res, val)
		case "string":
			res = append(res, arg.Value)
		case "bytes":
			val, err := simpleDecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			res = append(res, val)
		case "bool":
			if arg.Value == "1" {
				res = append(res, true)
			} else if arg.Value == "0" {
				res = append(res, false)
			} else {
				return nil, errors.Errorf("Value of argument %d should be a string containing 1 or 0\n\nCurrent value: '%s'", i+1, arg.Value)
			}
		case "uint256":
			valBytes, err := simpleDecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the uint256 32 bytes in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			if len(valBytes) != 32 {
				return nil, errors.Errorf("Value of argument %d should be a string of 64 hex containing the uint256\n Actual size : %d", i+1, len(valBytes))
			}
			val := big.NewInt(0)
			val.SetBytes(valBytes)
			res = append(res, val)
		case "bytes20":
			valBytes, err := simpleDecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes20 in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			if len(valBytes) != 20 {
				return nil, errors.Errorf("Value of argument %d should be a string of 40 hex containing the bytes20\n Actual size : %d", i+1, len(valBytes))
			}
			var val [20]byte
			copy(val[:], valBytes)
			res = append(res, val)
		case "bytes32":
			valBytes, err := simpleDecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes32 in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			if len(valBytes) != 32 {
				return nil, errors.Errorf("Value of argument %d should be a string of 64 hex containing the bytes32\n Actual size : %d", i+1, len(valBytes))
			}
			var val [32]byte
			copy(val[:], valBytes)
			res = append(res, val)
		case "gamma:address":
			val, err := encoding.DecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			res = append(res, val)
		case "gamma:keys-file-address":
			key := getTestKeyFromFile(arg.Value)
			res = append(res, key.Address)
		default:
			return nil, errors.Errorf("Type of argument %d '%s' is unsupported\n\n%s", i+1, arg.Type, supported)
		}
	}
	return res, nil
}

func MarshalArgs(arguments []interface{}) ([]*Arg, error) {
	res := []*Arg{}
	for i, arg := range arguments {
		switch arg.(type) {
		case uint32:
			res = append(res, &Arg{"uint32", strconv.FormatUint(uint64(arg.(uint32)), 10)})
		case uint64:
			res = append(res, &Arg{"uint64", strconv.FormatUint(uint64(arg.(uint64)), 10)})
		case string:
			res = append(res, &Arg{"string", arg.(string)})
		case []byte:
			res = append(res, &Arg{"bytes", "0x" + hex.EncodeToString(arg.([]byte))})
		case bool:
			if arg.(bool) {
				res = append(res, &Arg{"bool", "1"})
			} else {
				res = append(res, &Arg{"bool", "0"})
			}
		case *big.Int:
			val := [32]byte{}
			b := arg.(*big.Int).Bytes()
			copy(val[32-len(b):], b)
			res = append(res, &Arg{"uint256", "0x" + hex.EncodeToString(val[:])})
		case [20]byte:
			val := arg.([20]byte)
			res = append(res, &Arg{"bytes20", "0x" + hex.EncodeToString(val[:])})
		case [32]byte:
			val := arg.([32]byte)
			res = append(res, &Arg{"bytes32", "0x" + hex.EncodeToString(val[:])})
		default:
			return nil, errors.Errorf("Type of argument %d '%T' is unsupported\n\n%s", i+1, arg, supported)
		}
	}
	return res, nil
}

func simpleDecodeHex(value string) ([]byte, error) {
	if strings.HasPrefix(value, "0x") {
		value = value[2:]
	}
	return hex.DecodeString(value)
}
