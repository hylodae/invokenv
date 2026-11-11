package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRawProviderOK(t *testing.T) {
	params := map[string]interface{}{
		"value": "foo",
	}

	p, err := NewRawProvider(params)

	assert.Nil(t, err)
	assert.Equal(t, "foo", p.Params.Value)
}

func TestNewRawProviderEmptyValue(t *testing.T) {
	params := map[string]interface{}{
		"foo": "bar",
	}

	_, err := NewRawProvider(params)

	assert.NotNil(t, err)
}

func TestNewRawProviderValueBadType(t *testing.T) {
	params := map[string]interface{}{
		"foo": []string{
			"alpha",
			"beta",
		},
	}

	_, err := NewRawProvider(params)

	assert.NotNil(t, err)
}

func TestRawProviderFetchValueOK(t *testing.T) {
	params := map[string]interface{}{
		"value": "foo",
	}

	p, err := NewRawProvider(params)
	assert.Nil(t, err)

	value, err := p.FetchValue()

	assert.Nil(t, err)
	assert.Equal(t, "foo", value)
}
