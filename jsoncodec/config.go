package jsoncodec

import "encoding/json"

type ConfFile struct {
	Environments map[string]*ConfEnv
}

type ConfEnv struct {
	VirtualChain uint32
	Endpoints    []string
	Experimental bool
}

func UnmarshalConfFile(bytes []byte) (*ConfFile, error) {
	var confFile *ConfFile
	err := json.Unmarshal(bytes, &confFile)
	return confFile, err
}
