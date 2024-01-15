package terraform

import (
	"github.com/fatih/color"
)

var (
	delete  = color.New(color.FgRed)
	update  = color.New(color.FgYellow)
	replace = color.New(color.FgHiYellow)
	create  = color.New(color.FgGreen)
)

type Resource struct {
	Addr     string `json:"addr"`
	Resource string `json:"resource"`
}

type Change struct {
	Action   string   `json:"action"`
	Resource Resource `json:"resource"`
}

type ResourceChange struct {
	Message  string `json:"@message"`
	Change   Change `json:"change"`
	Type     string `json:"type"`
	Selected bool
	Hidden   bool
}

type resourceFilter func(ResourceChange) bool

type ResourceChangeList []ResourceChange

func (list ResourceChangeList) Filter(fn resourceFilter) (returnList ResourceChangeList) {
	for _, change := range list {
		if fn(change) {
			returnList = append(returnList, change)
		}
	}
	return
}

func (r *ResourceChange) Title() string       { return r.Change.Resource.Addr }
func (r *ResourceChange) Description() string { return r.Change.Action }
func (r *ResourceChange) FilterValue() string { return r.Change.Resource.Addr }
func (r *ResourceChange) IsSelected() bool    { return r.Selected }
