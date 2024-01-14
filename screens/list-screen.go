package screens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/smorenodp/tftarget/terraform"
)

type listKeyMap struct {
	launchApply     key.Binding
	launchPlan      key.Binding
	toggleHelpMenu  key.Binding
	toggleTargetAll key.Binding
	toggleCommand   key.Binding
}

func (l *ListScreen) TargetAll() {
	items := l.List.Items()
	for index, i := range items {
		r := i.(terraform.ResourceChange)
		r.Selected = !l.targetedAll
		l.List.SetItem(index, r)
	}
	l.targetedAll = !l.targetedAll
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		launchApply: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Launch apply with targets"),
		),
		launchPlan: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "Launch plan with targets"),
		),
		// toggleHelpMenu: key.NewBinding(
		// 	key.WithKeys("h"),
		// 	key.WithHelp("h", "toggle help"),
		// ),
		toggleTargetAll: key.NewBinding(
			key.WithKeys("*"),
			key.WithHelp("*", "target all"),
		),
		// toggleCommand: key.NewBinding(
		// 	key.WithKeys("c"),
		// 	key.WithHelp("c", "Show command"),
		// ),
	}
}

func NewListScreen(resources []terraform.ResourceChange, width, height int) ListScreen {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Make initial list of items

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	test := []list.Item{}
	for _, r := range resources {
		test = append(test, r)
	}
	resourceList := list.New(test, delegate, 0, 0)
	resourceList.Title = "Resources"
	resourceList.Styles.Title = titleStyle
	resourceList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.launchApply,
			listKeys.launchPlan,
			listKeys.toggleHelpMenu,
			listKeys.toggleTargetAll,
			listKeys.toggleCommand,
		}
	}
	h, v := appStyle.GetFrameSize()
	resourceList.SetSize(width-h, height-v)

	return ListScreen{List: resourceList, keys: listKeys, delegateKeys: delegateKeys}
}

type ListScreen struct {
	List         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	targetedAll  bool
}

type ListCommandLaunched struct {
	Type string
}

func (l ListScreen) returnCommand(t string) tea.Cmd {
	return func() tea.Msg {
		return ListCommandLaunched{t}
	}
}

func (l ListScreen) Start() tea.Cmd {
	return nil
}

func (l ListScreen) Update(msg tea.Msg) (MainScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		l.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		if l.List.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, l.keys.toggleTargetAll):
			l.TargetAll()
		case key.Matches(msg, l.keys.launchApply):
			items := l.List.Items()
			var resources []terraform.ResourceChange = []terraform.ResourceChange{}
			for _, i := range items {
				if resource, ok := i.(terraform.ResourceChange); ok && resource.IsSelected() {
					resources = append(resources, resource)
				}
			}
			return l, l.returnCommand("apply")
		case key.Matches(msg, l.keys.launchPlan):
			items := l.List.Items()
			var resources []terraform.ResourceChange = []terraform.ResourceChange{}
			for _, i := range items {
				if resource, ok := i.(terraform.ResourceChange); ok && resource.IsSelected() {
					resources = append(resources, resource)
				}
			}
			return l, l.returnCommand("plan")

		case key.Matches(msg, l.keys.toggleHelpMenu):
			l.List.SetShowHelp(!l.List.ShowHelp())
			return l, nil
		}
	}
	newListModel, cmd := l.List.Update(msg)
	l.List = newListModel

	return l, cmd
}

func (l ListScreen) View() string {
	return appStyle.Render(l.List.View())
}
