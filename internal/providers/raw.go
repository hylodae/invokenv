package providers

import "fmt"

const (
	RawName = "raw"
)

type RawProviderParams struct {
	Value string `json:"value"`
}

type RawProvider struct {
	Params *RawProviderParams
}

func NewRawProvider(params map[string]interface{}) (*RawProvider, error) {
	var rawParams RawProviderParams
	err := parseParams(params, &rawParams)
	if err != nil {
		return nil, err
	}

	if rawParams.Value == "" {
		return nil, fmt.Errorf("parameter `value` is required in provider `raw`")
	}

	p := &RawProvider{
		Params: &rawParams,
	}

	return p, nil
}

func (p *RawProvider) FetchValue() (string, error) {
	return p.Params.Value, nil
}
