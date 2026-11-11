package providers

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	ExecName = "exec"
)

type ExecProviderParams struct {
	Command string `json:"command"`
}

type ExecProvider struct {
	Params *ExecProviderParams
}

func NewExecProvider(params map[string]interface{}) (*ExecProvider, error) {
	var execParams ExecProviderParams
	err := parseParams(params, &execParams)
	if err != nil {
		return nil, err
	}

	if execParams.Command == "" {
		return nil, fmt.Errorf("parameter `command` is required in provider `exec`")
	}

	p := &ExecProvider{
		Params: &execParams,
	}
	return p, nil
}

func (p *ExecProvider) FetchValue() (string, error) {
	args := strings.Split(p.Params.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
