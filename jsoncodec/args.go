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
	"strconv"
	"strings"
)

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
			res = append(res, uint64(val))
		case "string":
			res = append(res, string(arg.Value))
		case "bytes":
			val, err := simpleDecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			res = append(res, []byte(val))
		case "gamma:address":
			val, err := encoding.DecodeHex(arg.Value)
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			res = append(res, []byte(val))
		case "gamma:keys-file-address":
			key := getTestKeyFromFile(arg.Value)
			res = append(res, []byte(key.Address))
		default:
			supported := "Supported types are: uint32 uint64 string bytes gamma:keys-file-address"
			return nil, errors.Errorf("Type of argument %d '%s' is unsupported\n\n%s", i+1, arg.Type, supported)
		}
	}
	return res, nil
}

func MarshalArgs(arguments []interface{}) []*Arg {
	res := []*Arg{}
	for _, arg := range arguments {
		switch arg.(type) {
		case uint32:
			res = append(res, &Arg{"uint32", strconv.FormatUint(uint64(arg.(uint32)), 10)})
		case uint64:
			res = append(res, &Arg{"uint64", strconv.FormatUint(uint64(arg.(uint64)), 10)})
		case string:
			res = append(res, &Arg{"string", arg.(string)})
		case []byte:
			res = append(res, &Arg{"bytes", "0x" + hex.EncodeToString(arg.([]byte))})
		default:
			panic("unsupported type in json marshal of method arguments")
		}
	}
	return res
}

func simpleDecodeHex(value string) ([]byte, error) {
	if strings.HasPrefix(value, "0x") {
		value = value[2:]
	}
	return hex.DecodeString(value)
}
