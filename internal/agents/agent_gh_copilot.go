package agents

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
)

type GhCopilotAgent struct {
	exec string
}

var GhCopilotAgentType = reflect.TypeOf(new(GhCopilotAgent))

var _ AgentTool = (*GhCopilotAgent)(nil)

func NewGhCopilotAgent(path string) *GhCopilotAgent {
	return &GhCopilotAgent{
		exec: path,
	}
}

func (a *GhCopilotAgent) Descriptor() string {
	return "GitHub Copilot CLI (" + a.exec + ")"
}

func (a *GhCopilotAgent) ExecuteCommand(command string) (string, error) {
	cmd := exec.Command(a.exec, "--allow-all-tools", "--prompt", command)

	var buf bytes.Buffer
	prefix := "  Copilot> "
	cmd.Stdout = NewPrefixWriter(os.Stdout, prefix)
	cmd.Stderr = NewPrefixWriter(os.Stderr, prefix)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
