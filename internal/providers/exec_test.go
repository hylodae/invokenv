package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExecProviderOK(t *testing.T) {
	params := map[string]interface{}{
		"command": "echo foo",
	}

	p, err := NewExecProvider(params)

	assert.Nil(t, err)
	assert.Equal(t, "echo foo", p.Params.Command)
}

func TestNewExecProviderCommandEmpty(t *testing.T) {
	params := map[string]interface{}{
		"foo": "bar",
	}

	_, err := NewExecProvider(params)

	assert.NotNil(t, err)
}

func TestNewExecProviderCommandBadType(t *testing.T) {
	params := map[string]interface{}{
		"command": []string{
			"alpha",
			"beta",
		},
	}

	_, err := NewExecProvider(params)

	assert.NotNil(t, err)
}

func TestExecProviderFetchValueOK(t *testing.T) {
	params := map[string]interface{}{
		"command": "echo foo",
	}

	p, err := NewExecProvider(params)
	assert.Nil(t, err)

	value, err := p.FetchValue()

	assert.Nil(t, err)
	assert.Equal(t, "foo", value)
}
