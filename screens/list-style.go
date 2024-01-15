package screens

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/smorenodp/tftarget/terraform"
)

const (
	ellipsis = "…"
)

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle      lipgloss.Style
	NormalDescCreate lipgloss.Style
	NormalDescUpdate lipgloss.Style
	NormalDescRemove lipgloss.Style

	// The selected item state.
	SelectedTitle      lipgloss.Style
	SelectedDescCreate lipgloss.Style
	SelectedDescUpdate lipgloss.Style
	SelectedDescRemove lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	TargetTitle      lipgloss.Style
	TargetDescCreate lipgloss.Style
	TargetDescUpdate lipgloss.Style
	TargetDescRemove lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#999999"}).
		Padding(0, 0, 0, 2)

	s.NormalDescCreate = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#b6d7a8", Dark: "#b6d7a8"})

	s.NormalDescUpdate = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#f9cb9c", Dark: "#f9cb9c"})

	s.NormalDescRemove = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#ea9999", Dark: "#ea9999"})

	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#ffffff"}).
		Padding(0, 0, 0, 1)

	s.SelectedDescCreate = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#6aa84f", Dark: "#6aa84f"})

	s.SelectedDescUpdate = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#f1c232", Dark: "#f1c232"})

	s.SelectedDescRemove = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#cc0000", Dark: "#cc0000"})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.TargetTitle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#70ea3b", Dark: "#70ea3b"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#999999", Dark: "#ffffff"}).
		Padding(0, 0, 0, 1)

	s.TargetDescCreate = s.TargetTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#6aa84f", Dark: "#6aa84f"})

	s.TargetDescUpdate = s.TargetTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#f1c232", Dark: "#f1c232"})

	s.TargetDescRemove = s.TargetTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#cc0000", Dark: "#cc0000"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

// DefaultItem describes an items designed to work with DefaultDelegate.
type DefaultItem interface {
	list.Item
	Title() string
	Description() string
}

// DefaultDelegate is a standard delegate designed to work in lists. It's
// styled by DefaultItemStyles, which can be customized as you like.
//
// The description line can be hidden by setting Description to false, which
// renders the list as single-line-items. The spacing between items can be set
// with the SetSpacing method.
//
// Setting UpdateFunc is optional. If it's set it will be called when the
// ItemDelegate called, which is called when the list's Update function is
// invoked.
//
// Settings ShortHelpFunc and FullHelpFunc is optional. They can be set to
// include items in the list's default short and full help menus.
type DefaultDelegate struct {
	ShowDescription bool
	Styles          DefaultItemStyles
	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
	ShortHelpFunc   func() []key.Binding
	FullHelpFunc    func() [][]key.Binding
	height          int
	spacing         int
}

// NewDefaultDelegate creates a new delegate with default styles.
func NewDefaultDelegate() DefaultDelegate {
	return DefaultDelegate{
		ShowDescription: true,
		Styles:          NewDefaultItemStyles(),
		height:          2,
		spacing:         1,
	}
}

// SetHeight sets delegate's preferred height.
func (d *DefaultDelegate) SetHeight(i int) {
	d.height = i
}

// Height returns the delegate's preferred height.
// This has effect only if ShowDescription is true,
// otherwise height is always 1.
func (d DefaultDelegate) Height() int {
	if d.ShowDescription {
		return d.height
	}
	return 1
}

// SetSpacing sets the delegate's spacing.
func (d *DefaultDelegate) SetSpacing(i int) {
	d.spacing = i
}

// Spacing returns the delegate's spacing.
func (d DefaultDelegate) Spacing() int {
	return d.spacing
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (d DefaultDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

// Render prints an item.
func (d DefaultDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc  string
		selected     bool
		matchedRunes []int
		s            = &d.Styles
	)

	if i, ok := item.(DefaultItem); ok {
		title = i.Title()
		desc = i.Description()
		selected = i.(*terraform.ResourceChange).IsSelected()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := uint(m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.height-1 {
				break
			}
			lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && index < len(m.VisibleItems()) {
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.Render(desc)
	} else if selected {
		title = s.TargetTitle.Render(title)
		if desc == "create" {
			desc = s.TargetDescCreate.Render(desc)
		} else if desc == "update" || desc == "replace" {
			desc = s.TargetDescUpdate.Render(desc)
		} else {
			desc = s.TargetDescRemove.Render(desc)
		}
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		if desc == "create" {
			desc = s.SelectedDescCreate.Render(desc)
		} else if desc == "update" || desc == "replace" {
			desc = s.SelectedDescUpdate.Render(desc)
		} else {
			desc = s.SelectedDescRemove.Render(desc)
		}
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		if desc == "create" {
			desc = s.NormalDescCreate.Render(desc)
		} else if desc == "update" || desc == "replace" {
			desc = s.NormalDescUpdate.Render(desc)
		} else {
			desc = s.NormalDescRemove.Render(desc)
		}
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc)
		return
	}
	fmt.Fprintf(w, "%s", title)
}

// ShortHelp returns the delegate's short help.
func (d DefaultDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

// FullHelp returns the delegate's full help.
func (d DefaultDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
