package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smorenodp/tftarget/screens"
	"github.com/smorenodp/tftarget/terraform"
)

type Model struct {
	tfClient   *terraform.TFClient
	screen     screens.MainScreen
	width      int
	height     int
	quitting   bool
	err        error
	message    string
	Command    terraform.Command
	listScreen screens.ListScreen
}

func New(tf *terraform.TFClient) Model {

	return Model{screen: screens.NewSpinnerScreen(tf), tfClient: tf, Command: terraform.NewCommand(tf.ChDir, tf.VarFile)}
}

func (m Model) Run() error {
	var err error

	if _, err = tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}

func (m Model) Init() tea.Cmd {
	return m.screen.Start()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			screen, cmd := m.screen.Update(msg)
			m.screen = screen
			return m, cmd
		}
	case error:
		m.err = msg
		return m, nil
	case screens.SpinnerMessageOut:
		listScreen := screens.NewListScreen(msg.Output, m.width, m.height)
		//TODO: I don't quite like this thing
		m.Command.SetItems(listScreen.List.Items())
		m.screen = listScreen
		return m, nil
	// Maybe add one message for adding/removing resource target
	case screens.ListCommandLaunched:
		m.listScreen = m.screen.(screens.ListScreen)
		m.Command.Type = msg.Type
		command := strings.Split(m.Command.String(), " ")
		commandScreen := screens.NewCommandScreen(command[0], command[1:])
		m.screen = commandScreen
		return m, m.screen.Start()
	case screens.CommandContinue:
		m.screen = m.listScreen
		return m, nil
	default:
		screen, cmd := m.screen.Update(msg)
		m.screen = screen
		return m, cmd
	}
}

func (m Model) View() string {
	if !m.quitting {
		if m.message != "" {
			return m.message
		} else if m.err != nil {
			return fmt.Sprintf("Error: %s", m.err)
		} else {
			return m.screen.View()
		}
	}
	return ""
}
