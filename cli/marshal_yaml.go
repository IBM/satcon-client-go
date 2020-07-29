package cli

import (
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

func MarshalYAMLFromFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
}
