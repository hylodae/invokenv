package providers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileProviderOK(t *testing.T) {
	params := map[string]interface{}{
		"path": "/foo/bar",
	}

	p, err := NewFileProvider(params)

	assert.Nil(t, err)
	assert.Equal(t, "/foo/bar", p.Params.Path)
	assert.Equal(t, "", p.Params.Format)
	assert.Equal(t, "", p.Params.Key)
}

func TestNewFileProviderPathEmpty(t *testing.T) {
	params := map[string]interface{}{
		"foo": "bar",
	}

	_, err := NewFileProvider(params)

	assert.NotNil(t, err)
}

func TestNewFileProviderPathBadType(t *testing.T) {
	params := map[string]interface{}{
		"foo": []string{
			"alpha",
			"beta",
		},
	}

	_, err := NewFileProvider(params)

	assert.NotNil(t, err)
}

func TestFileProviderFetchValueRawFileOK(t *testing.T) {
	file, _ := os.CreateTemp("", "invokenv")
	defer os.Remove(file.Name())

	file.WriteString("foo")
	file.Close()

	params := map[string]interface{}{
		"path": file.Name(),
	}

	p, err := NewFileProvider(params)
	assert.Nil(t, err)

	value, err := p.FetchValue()

	assert.Nil(t, err)
	assert.Equal(t, "foo", value)
}

func TestFileProviderFetchValueRawFileFormatOK(t *testing.T) {
	file, _ := os.CreateTemp("", "invokenv")
	defer os.Remove(file.Name())

	file.WriteString("foo")
	file.Close()

	params := map[string]interface{}{
		"path":   file.Name(),
		"format": "raw",
	}

	p, err := NewFileProvider(params)
	assert.Nil(t, err)

	value, err := p.FetchValue()

	assert.Nil(t, err)
	assert.Equal(t, "foo", value)
}

func TestFileProviderFetchValueYAMLFileOK(t *testing.T) {
	file, _ := os.CreateTemp("", "invokenv")
	defer os.Remove(file.Name())

	file.WriteString("---\nfoo: bar\n")
	file.Close()

	params := map[string]interface{}{
		"path":   file.Name(),
		"format": "yaml",
		"key":    "foo",
	}

	p, err := NewFileProvider(params)
	assert.Nil(t, err)

	value, err := p.FetchValue()

	assert.Nil(t, err)
	assert.Equal(t, "bar", value)
}

func TestFileProviderFetchValueYAMLFileNoKey(t *testing.T) {
	file, _ := os.CreateTemp("", "invokenv")
	defer os.Remove(file.Name())

	file.WriteString("---\nfoo: bar\n")
	file.Close()

	params := map[string]interface{}{
		"path":   file.Name(),
		"format": "yaml",
	}

	p, err := NewFileProvider(params)
	assert.Nil(t, err)

	_, err = p.FetchValue()

	assert.NotNil(t, err)
}

func TestFileProviderFetchValueYAMLFileNoKeyInFile(t *testing.T) {
	file, _ := os.CreateTemp("", "invokenv")
	defer os.Remove(file.Name())

	file.WriteString("---\ntoto: bar\n")
	file.Close()

	params := map[string]interface{}{
		"path":   file.Name(),
		"format": "yaml",
		"key":    "foo",
	}

	p, err := NewFileProvider(params)
	assert.Nil(t, err)

	_, err = p.FetchValue()

	assert.NotNil(t, err)
}
