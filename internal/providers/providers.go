package providers

import (
	"encoding/json/v2"
	"fmt"
)

type Provider struct{}

type ProviderInterface interface {
	FetchValue() (string, error)
}

func New(name string, params map[string]interface{}) (ProviderInterface, error) {
	switch name {
	case ExecName:
		return NewExecProvider(params)
	case FileName:
		return NewFileProvider(params)
	case RawName:
		return NewRawProvider(params)
	}
	return nil, fmt.Errorf("provider `%s` not supported", name)
}

func parseParams(input map[string]interface{}, output interface{}) (err error) {
	b, err := json.Marshal(input)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, output)
	if err != nil {
		return
	}

	return
}
