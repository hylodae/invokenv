package envfile

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"

	"github.com/hylodae/invokenv/internal/providers"
)

const (
	EnvFilename = "invokenv.yml"
)

type EnvFileItem struct {
	Name     string                 `yaml:"name"`
	Provider string                 `yaml:"provider"`
	Params   map[string]interface{} `yaml:"params"`
}

type EnvFile struct {
	Items []*EnvFileItem
}

// Exists checks if the invokenv.yml file exists in the directory.
func Exists(dirPath string) bool {
	path := fmt.Sprintf("%s/%s", dirPath, EnvFilename)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func Parse(dirPath string) (*EnvFile, error) {
	path := fmt.Sprintf("%s/%s", dirPath, EnvFilename)

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var envFile EnvFile

	err = yaml.Unmarshal(yamlFile, &envFile.Items)
	if err != nil {
		return nil, err
	}

	return &envFile, nil
}

func (e *EnvFile) FetchValues() ([]string, error) {
	envVars := []string{}

	for _, item := range e.Items {
		p, err := providers.New(item.Provider, item.Params)
		if err != nil {
			return []string{}, err
		}

		v, err := p.FetchValue()
		if err != nil {
			return []string{}, fmt.Errorf("failed to fetch value for item `%s`: %s", item.Name, err)
		}

		envVars = append(envVars, fmt.Sprintf("%s=%s", item.Name, v))
	}

	return envVars, nil
}
