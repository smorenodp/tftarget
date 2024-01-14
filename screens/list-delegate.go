package screens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/smorenodp/tftarget/terraform"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

func newItemDelegate(keys *delegateKeyMap) DefaultDelegate {
	d := NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) (output tea.Cmd) {
		var title string
		var i terraform.ResourceChange
		var ok bool

		if i, ok = m.SelectedItem().(terraform.ResourceChange); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				i.Selected = !i.Selected
				if i.Selected {
					output = m.NewStatusMessage(statusMessageStyle("You added " + title))
				} else {
					output = m.NewStatusMessage(statusMessageStyle("You removed " + title))
				}
				m.SetItem(m.Index(), i)
			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				output = m.NewStatusMessage(statusMessageStyle("You erased " + title))
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("space/enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "remove"),
		),
	}
}
