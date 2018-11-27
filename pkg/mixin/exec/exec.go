package exec

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/deislabs/porter/pkg/context"
	"gopkg.in/yaml.v2"
)

// Exec is the logic behind the exec mixin
type Mixin struct {
	*context.Context

	instruction Instruction
}

type Instruction struct {
	Name       string            `yaml:"name"`
	Command    string            `yaml:"command"`
	Arguments  []string          `yaml:"arguments"`
	Parameters map[string]string `yaml:"parameters"`
}

// New exec mixin client, initialized with useful defaults.
func New() *Mixin {
	return &Mixin{
		Context: context.New(),
	}
}

func (m *Mixin) LoadInstruction(commandFile string) error {
	contents, err := m.getCommandFile(commandFile, m.Out)
	if err != nil {
		return fmt.Errorf("there was an error getting commands: %s", err)
	}
	return yaml.Unmarshal(contents, &m.instruction)
}

func (m *Mixin) Execute() error {
	cmd := exec.Command(m.instruction.Command, m.instruction.Arguments...)
	cmd.Stdout = m.Out
	cmd.Stderr = m.Err

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start...%s", err)
	}

	return cmd.Wait()
}

func (m *Mixin) getCommandFile(commandFile string, w io.Writer) ([]byte, error) {
	if commandFile == "" {
		reader := bufio.NewReader(m.In)
		return ioutil.ReadAll(reader)
	}
	return ioutil.ReadFile(commandFile)
}
