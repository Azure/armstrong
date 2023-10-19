package azurerm

import (
	_ "embed"
	"encoding/json"
	"os"
)

type MappingJsonDependencyLoader struct {
	MappingJsonFilepath string
}

//go:embed mappings.json
var mappingsJson string

func (m MappingJsonDependencyLoader) Load() ([]Mapping, error) {
	var mappings []Mapping
	var data []byte
	var err error
	if len(m.MappingJsonFilepath) > 0 {
		data, err = os.ReadFile(m.MappingJsonFilepath)
		if err != nil {
			return nil, err
		}
	} else {
		data = []byte(mappingsJson)
	}
	err = json.Unmarshal(data, &mappings)
	if err != nil {
		return nil, err
	}
	return mappings, nil
}

var _ DependencyLoader = MappingJsonDependencyLoader{}
