package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smorenodp/tftarget/terraform"
)

type SpinnerScreen struct {
	spinner      spinner.Model
	tf           *terraform.TFClient
	err          error
	message      string
	planFinished bool
}

type SpinnerMessageOut struct {
	Output planStatus
}

type planStatus terraform.ResourceChangeList

func (s SpinnerScreen) logicPlan() tea.Msg {
	list, err := s.tf.Plan()
	if err != nil {
		return err
	}
	return planStatus(list)
}

func NewSpinnerScreen(tf *terraform.TFClient) SpinnerScreen {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

	return SpinnerScreen{spinner: s, tf: tf}
}

func (s SpinnerScreen) Start() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.logicPlan)
}

func (s SpinnerScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		if s.planFinished {
			cmd = nil
		} else {
			s.message = fmt.Sprintf("Obtaining terraform resources with changes...press q to cancel")
			s.spinner, cmd = s.spinner.Update(msg)
		}
		return s, cmd
	case planStatus:
		s.planFinished = true
		return s, func() tea.Msg { return SpinnerMessageOut{Output: msg} }
	default:
		return s, nil
	}
}

func (s SpinnerScreen) View() string {
	var str string
	if s.err != nil {
		str = s.err.Error()
	} else {
		str = fmt.Sprintf("%s %s\n\n", s.spinner.View(), s.message)
	}
	return str
}
