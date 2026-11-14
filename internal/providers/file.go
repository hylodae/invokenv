package providers

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jmespath/go-jmespath"
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
		return readYAMLFile(p.Params.Path, p.Params.Key)
	case "json":
		return readJSONFile(p.Params.Path, p.Params.Key)
	case "toml":
		return readTOMLFile(p.Params.Path, p.Params.Key)
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

func searchValueFromKey(data interface{}, key string) (string, error) {
	result, err := jmespath.Search(key, data)
	if err != nil {
		return "", fmt.Errorf("failed to fetch value from key `%s`: %s", key, err)
	}

	switch v := result.(type) {
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("failed to cast value of key `%s`", key)
	}
}

func readYAMLFile(path string, key string) (string, error) {
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

	return searchValueFromKey(data, key)
}

func readJSONFile(path string, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("parameter `key` is required with `json` format")
	}

	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return "", err
	}

	return searchValueFromKey(data, key)
}

func readTOMLFile(path string, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("parameter `key` is required with `toml` format")
	}

	tomlFile, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = toml.Unmarshal(tomlFile, &data)
	if err != nil {
		return "", err
	}

	return searchValueFromKey(data, key)
}
