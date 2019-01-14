package jsoncodec

import "encoding/json"

type Key struct {
	PrivateKey string // hex string starting with 0x
	PublicKey  string // hex string starting with 0x
	Address    string // hex string starting with 0x
}

type RawKey struct {
	PrivateKey []byte
	PublicKey  []byte
	Address    []byte
}

func UnmarshalKeys(bytes []byte) (map[string]*Key, error) {
	keys := make(map[string]*Key)
	err := json.Unmarshal(bytes, &keys)
	return keys, err
}

func MarshalKeys(keys map[string]*Key) ([]byte, error) {
	return json.MarshalIndent(keys, "", "  ")
}
