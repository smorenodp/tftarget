package terraform

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	tfjson "github.com/hashicorp/terraform-json"
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

type PlanChange struct {
	tfjson.ResourceChange
}

func stringChange(key string, value string, after map[string]interface{}) string {
	if afterValue, ok := after[key]; ok {
		return fmt.Sprintf("Changing in %s from %s to %s", key, value, afterValue)
	} else {
		return fmt.Sprintf("Removing %s", key)
	}
}

func NewPlanChange(change *tfjson.ResourceChange) PlanChange {
	return PlanChange{*change}
}

// We don't have the creating, only in after
func (p PlanChange) String() string {
	before := p.Change.Before.(map[string]interface{})
	var output []string = []string{}
	after := p.Change.After.(map[string]interface{})
	for key, value := range before {
		switch value.(type) {
		case map[string]string:
			output = append(output, fmt.Sprintf("The key %s contains map with only strings", key))
		case map[string]interface{}:
			output = append(output, fmt.Sprintf("The key %s contains map with more than strings", key))
		case string:
			output = append(output, stringChange(key, value.(string), after))
		case int:
			output = append(output, stringChange(key, fmt.Sprintf("%d", value.(int)), after))
		case bool:
			output = append(output, stringChange(key, fmt.Sprintf("%t", value.(bool)), after))
		default:
			output = append(output, fmt.Sprintf("Default option for key %s with value %v", key, value))
		}
	}
	return fmt.Sprintf(strings.Join(output, "\n"))
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
