package screens

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type CommandScreen struct {
	command *exec.Cmd
	cmd     string
	args    []string
}

type CommandFinished struct {
	err error
}

type CommandContinue struct {
}

func NewCommandScreen(command string, args []string) CommandScreen {

	c := exec.Command(command, args...)

	cmd := CommandScreen{command: c, cmd: command, args: args}
	return cmd
}

func (c CommandScreen) Start() tea.Cmd {
	cmd := tea.ExecProcess(c.command, func(err error) tea.Msg {
		return CommandFinished{err}
	})
	return tea.Batch(cmd, tea.ExitAltScreen)
}

func (c CommandScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return c, tea.Batch(tea.EnterAltScreen, func() tea.Msg { return CommandContinue{} })
		}
	}
	return c, nil
}

func (c CommandScreen) View() string {
	return "\n\n Command finished, press enter to return to list or q to exit..."
}
