package screens

import tea "github.com/charmbracelet/bubbletea"

type MainScreen interface {
	Start() tea.Cmd
	Update(msg tea.Msg) (MainScreen, tea.Cmd)
	View() string
}
