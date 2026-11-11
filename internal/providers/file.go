package providers

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
)

const (
	FileName = "file"
)

type FileProviderParams struct {
	Path   string `json:"path"`
	Format string `json:"format,omitempty"`
	Key    string `json:"key,omitempty"`
}

type FileProvider struct {
	Params *FileProviderParams
}

func NewFileProvider(params map[string]interface{}) (*FileProvider, error) {
	var fileParams FileProviderParams
	err := parseParams(params, &fileParams)
	if err != nil {
		return nil, err
	}

	if fileParams.Path == "" {
		return nil, fmt.Errorf("parameter `path` is required in provider `file`")
	}

	p := &FileProvider{
		Params: &fileParams,
	}
	return p, nil
}

func (p *FileProvider) FetchValue() (string, error) {
	switch strings.ToLower(p.Params.Format) {
	case "yaml":
		return readYamlFile(p.Params.Path, p.Params.Key)
	case "raw", "":
		return readRawFile(p.Params.Path)
	}

	return "", fmt.Errorf("unsupported format: %s", strings.ToLower(p.Params.Format))
}

func readRawFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func readYamlFile(path string, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("parameter `key` is required with `yaml` format")
	}

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return "", err
	}

	if _, ok := data[key]; !ok {
		return "", fmt.Errorf("top-level YAML key `%s` not found", key)
	}

	v, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("failed to cast top-level YAML key `%s` as string", key)
	}

	return v, nil
}
