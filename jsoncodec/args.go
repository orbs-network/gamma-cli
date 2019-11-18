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
	"reflect"
	"strconv"
	"strings"
)

const supported = "Supported types are: uint32 uint64 uint256 bool string bytes bytes20 bytes32 uint32Array uint64Array uint256Array boolArray stringArray bytesArray bytes20Array bytes32Array gamma:address gamma:keys-file-address"

type Arg struct {
	Type  string
	Value interface{}
}

func isArgsInputStructureValid(args []*Arg) error {
	for i, arg := range args {
		rValue := reflect.TypeOf(arg.Value).String()
		if strings.HasSuffix(arg.Type, "Array") {
			if rValue != "[]interface {}" {
				return errors.Errorf("Argument %d's Type is marked as an Array and it's Value should contain an array of string\nCurrently %s\n", i+1, rValue)
			}
		} else if rValue != "string" {
			return errors.Errorf("Argument %d's Type is marked as a Scalar and it's Value should contain a string", i+1)
		}
	}
	return nil
}

func unmarshalScalar(argType, value string) (interface{}, error) {
	switch argType {
	case "uint32":
		val, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, errors.Errorf("a numeric value\nCurrent value: '%s'", value)
		}
		return uint32(val), nil
	case "uint64":
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, errors.Errorf("a numeric value\nCurrent value: '%s'", value)
		}
		return val, nil
	case "string":
		return value, nil
	case "bytes":
		val, err := simpleDecodeHex(value)
		if err != nil {
			return nil, errors.Errorf("bytes in hex format\nHex decoder returned error: %s\nCurrent value: '%s'", err.Error(), value)
		}
		return val, nil
	case "bool":
		if value == "1" {
			return true, nil
		} else if value == "0" {
			return false, nil
		} else {
			return nil, errors.Errorf("1 or 0\nCurrent value: '%s'", value)
		}
	case "uint256":
		valBytes, err := simpleDecodeHex(value)
		if err != nil {
			return nil, errors.Errorf("uint256 value in bytes in a hex format (64 hexes)\nHex decoder returned error: %s\nCurrent value: '%s'", err.Error(), value)
		}
		if len(valBytes) != 32 {
			return nil, errors.Errorf("uint256 value in bytes in a hex format (64 hexes)\n Actual size : %d", len(valBytes))
		}
		val := big.NewInt(0)
		val.SetBytes(valBytes)
		return val, nil
	case "bytes20":
		valBytes, err := simpleDecodeHex(value)
		if err != nil {
			return nil, errors.Errorf("bytes20 in a hex format (40 hexes)\nHex decoder returned error: %s\nCurrent value: '%s'", err.Error(), value)
		}
		if len(valBytes) != 20 {
			return nil, errors.Errorf("bytes20 in a hex format (40 hexes)\n Actual size : %d", len(valBytes))
		}
		var val [20]byte
		copy(val[:], valBytes)
		return val, nil
	case "bytes32":
		valBytes, err := simpleDecodeHex(value)
		if err != nil {
			return nil, errors.Errorf("bytes32 in a hex format (64 hexes)\nHex decoder returned error: %s\nCurrent value: '%s'", err.Error(), value)
		}
		if len(valBytes) != 32 {
			return nil, errors.Errorf("bytes32 in a hex format (64 hexes)\n Actual size : %d", len(valBytes))
		}
		var val [32]byte
		copy(val[:], valBytes)
		return val, nil
	default:
		return nil, errors.Errorf("a known type. '%s' is unsupported\n%s", argType, supported)
	}
}

func unmarshalArray(argType string, argValues []interface{}) (interface{}, error) {
	switch argType {
	case "uint32Array":
		var argArrayRes []uint32
		for j, argValue := range argValues {
			val, err := strconv.ParseUint(argValue.(string), 10, 32)
			if err != nil {
				return nil, errors.Errorf("element %d should be a string containing a numeric value\nCurrent value: '%s'", j+1, argValue)
			}
			argArrayRes = append(argArrayRes, uint32(val))
		}
		return argArrayRes, nil
	case "uint64Array":
		var argArrayRes []uint64
		for j, argValue := range argValues {
			val, err := strconv.ParseUint(argValue.(string), 10, 64)
			if err != nil {
				return nil, errors.Errorf("element %d should be a string containing a numeric value\nCurrent value: '%s'", j+1, argValue)
			}
			argArrayRes = append(argArrayRes, val)
		}
		return argArrayRes, nil
	case "stringArray":
		var argArrayRes []string
		for _, argValue := range argValues {
			argArrayRes = append(argArrayRes, argValue.(string))
		}
		return argArrayRes, nil
	case "bytesArray":
		var argArrayRes [][]byte
		for j, argValue := range argValues {
			val, err := simpleDecodeHex(argValue.(string))
			if err != nil {
				return nil, errors.Errorf("element %d should be a string containing bytes in hex format\nHex decoder returned error: %s\nCurrent value: '%s'", j+1, err.Error(), argValue)
			}
			argArrayRes = append(argArrayRes, val)
		}
		return argArrayRes, nil
	case "boolArray":
		var argArrayRes []bool
		for j, argValue := range argValues {
			s := argValue.(string)
			if s == "1" {
				argArrayRes = append(argArrayRes, true)
			} else if s == "0" {
				argArrayRes = append(argArrayRes, false)
			} else {
				return nil, errors.Errorf("element %d should be a string containing 1 or 0\nCurrent value: '%s'", j+1, argValue)
			}
		}
		return argArrayRes, nil
	case "uint256Array":
		var argArrayRes []*big.Int
		for j, argValue := range argValues {
			valBytes, err := simpleDecodeHex(argValue.(string))
			if err != nil {
				return nil, errors.Errorf("element %d should be a string containing uint256 in hex format (64 hexes)\nHex decoder returned error: %s\nCurrent value: '%s'", j+1, err.Error(), argValue)
			}
			if len(valBytes) != 32 {
				return nil, errors.Errorf("element %d should be a string containing uint256 in a hex format (64 hexes)\n Actual size : %d", j+1, len(valBytes))
			}
			val := big.NewInt(0)
			val.SetBytes(valBytes)
			argArrayRes = append(argArrayRes, val)
		}
		return argArrayRes, nil
	case "bytes20Array":
		var argArrayRes [][20]byte
		for j, argValue := range argValues {
			valBytes, err := simpleDecodeHex(argValue.(string))
			if err != nil {
				return nil, errors.Errorf("element %d should be a string containing bytes20 in hex format (40 hexes)\nHex decoder returned error: %s\nCurrent value: '%s'", j+1, err.Error(), argValue)
			}
			if len(valBytes) != 20 {
				return nil, errors.Errorf("element %d should be a string containing bytes20 in a hex format (40 hexes)\n Actual size : %d", j+1, len(valBytes))
			}
			var val [20]byte
			copy(val[:], valBytes)
			argArrayRes = append(argArrayRes, val)
		}
		return argArrayRes, nil
	case "bytes32Array":
		var argArrayRes [][32]byte
		for j, argValue := range argValues {
			valBytes, err := simpleDecodeHex(argValue.(string))
			if err != nil {
				return nil, errors.Errorf("element %d should be a string containing bytes32 in hex format (64 hexes)\nHex decoder returned error: %s\nCurrent value: '%s'", j+1, err.Error(), argValue)
			}
			if len(valBytes) != 32 {
				return nil, errors.Errorf("element %d should be a string containing bytes32 in a hex format (64 hexes)\n Actual size : %d", j+1, len(valBytes))
			}
			var val [32]byte
			copy(val[:], valBytes)
			argArrayRes = append(argArrayRes, val)
		}
		return argArrayRes, nil
	default:
		return nil, errors.Errorf("a known type. '%s' is unsupported\n%s", argType, supported)
	}
}

func UnmarshalArgs(args []*Arg, getTestKeyFromFile func(string) *RawKey) ([]interface{}, error) {
	if err := isArgsInputStructureValid(args); err != nil {
		return nil, err
	}
	var res []interface{}
	for i, arg := range args {
		if arg.Type == "gamma:address" {
			val, err := encoding.DecodeHex(arg.Value.(string))
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing the bytes in hex\nHex decoder returned error: %s\n\nCurrent value: '%s'", i+1, err.Error(), arg.Value)
			}
			res = append(res, val)
		} else if arg.Type == "gamma:keys-file-address" {
			key := getTestKeyFromFile(arg.Value.(string))
			res = append(res, key.Address)
		} else if strings.HasSuffix(arg.Type, "Array") {
			valArray, err := unmarshalArray(arg.Type, arg.Value.([]interface{}))
			if err != nil {
				return nil, errors.Errorf("Value of array argument %d, %s", i+1, err.Error())
			}
			res = append(res, valArray)
		} else {
			val, err := unmarshalScalar(arg.Type, arg.Value.(string))
			if err != nil {
				return nil, errors.Errorf("Value of argument %d should be a string containing %s", i+1, err.Error())
			}
			res = append(res, val)
		}
	}
	return res, nil
}

func MarshalArgs(arguments []interface{}) ([]*Arg, error) {
	var res []*Arg
	for i, arg := range arguments {
		if reflect.TypeOf(arg).Kind() == reflect.Slice { // all []Type including []byte
			var arrArguments []string
			switch arg := arg.(type) {
			case []byte:
				res = append(res, &Arg{"bytes", "0x" + hex.EncodeToString(arg)})
			case []uint32:
				for _, v := range arg {
					arrArguments = append(arrArguments, strconv.FormatUint(uint64(v), 10))
				}
				res = append(res, &Arg{"uint32Array", arrArguments})
			case []uint64:
				for _, v := range arg {
					arrArguments = append(arrArguments, strconv.FormatUint(v, 10))
				}
				res = append(res, &Arg{"uint64Array", arrArguments})
			case []string:
				for _, v := range arg {
					arrArguments = append(arrArguments, v)
				}
				res = append(res, &Arg{"stringArray", arrArguments})
			case [][]byte:
				for _, v := range arg {
					arrArguments = append(arrArguments, "0x"+hex.EncodeToString(v))
				}
				res = append(res, &Arg{"bytesArray", arrArguments})
			case []bool:
				for _, v := range arg {
					if v {
						arrArguments = append(arrArguments, "1")
					} else {
						arrArguments = append(arrArguments, "0")
					}
				}
				res = append(res, &Arg{"boolArray", arrArguments})
			case []*big.Int:
				val := [32]byte{}
				for _, v := range arg {
					b := v.Bytes()
					copy(val[32-len(b):], b)
					arrArguments = append(arrArguments, "0x"+hex.EncodeToString(val[:]))
				}
				res = append(res, &Arg{"uint256Array", arrArguments})
			case [][20]byte:
				for _, v := range arg {
					arrArguments = append(arrArguments, "0x"+hex.EncodeToString(v[:]))
				}
				res = append(res, &Arg{"bytes20Array", arrArguments})
			case [][32]byte:
				for _, v := range arg {
					arrArguments = append(arrArguments, "0x"+hex.EncodeToString(v[:]))
				}
				res = append(res, &Arg{"bytes32Array", arrArguments})
			default:
				return nil, errors.Errorf("Type of argument %d '%T' is unsupported\n\n%s", i+1, arg, supported)
			}
		} else {
			switch arg := arg.(type) {
			case uint32:
				res = append(res, &Arg{"uint32", strconv.FormatUint(uint64(arg), 10)})
			case uint64:
				res = append(res, &Arg{"uint64", strconv.FormatUint(arg, 10)})
			case string:
				res = append(res, &Arg{"string", arg})
			case bool:
				if arg {
					res = append(res, &Arg{"bool", "1"})
				} else {
					res = append(res, &Arg{"bool", "0"})
				}
			case *big.Int:
				val := [32]byte{}
				b := arg.Bytes()
				copy(val[32-len(b):], b)
				res = append(res, &Arg{"uint256", "0x" + hex.EncodeToString(val[:])})
			case [20]byte:
				res = append(res, &Arg{"bytes20", "0x" + hex.EncodeToString(arg[:])})
			case [32]byte:
				res = append(res, &Arg{"bytes32", "0x" + hex.EncodeToString(arg[:])})
			default:
				return nil, errors.Errorf("Type of argument %d '%T' is unsupported\n\n%s", i+1, arg, supported)
			}
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
