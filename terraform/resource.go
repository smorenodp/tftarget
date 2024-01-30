package terraform

import (
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

var (
	change = `
resource {{ .Name }} {{ .ResourceType }} {
	{{ range .Changes }}
		{{ range .mapChanges }}
			Yes, it's a map {{ .Name }}
		{{ else }}
			{{ .Type }} {{.Name }} = {{ .Before }} -> {{ .After }}
		{{ end }}
	{{ end }}
}
	`
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

type MapInnerChange struct {
	Name       string
	Type       string
	Before     string
	After      string
	mapChanges []MapInnerChange // Only used for map values
}

type MapChange struct {
	Name         string
	ResourceType string
	ChangeType   string
	Changes      []MapInnerChange
}

type PlanChange struct {
	tfjson.ResourceChange
}

func stringChange(key string, value string, after map[string]interface{}, unknown map[string]interface{}, change MapInnerChange) MapInnerChange {
	if afterValue, ok := after[key]; ok {
		fmt.Printf("Changing in %s from %s to %s\n", key, value, afterValue)
		change.Type = "Update"
		change.Before = value
		change.After = afterValue.(string)
		return change
	} else if unknown[key].(bool) {
		fmt.Printf("Unknown value %s\n", key)
		change.Type = "Update"
		change.Before = value
		change.After = "(Known after apply)"
		return change
	} else {
		fmt.Sprintf("Removing %s\n", key)
		change.Type = "Removing"
		change.Before = value
		return change
	}
}

func NewPlanChange(change *tfjson.ResourceChange) PlanChange {
	return PlanChange{*change}
}

func stringMap(before map[string]interface{}, after map[string]interface{}, unknown map[string]interface{}) []MapInnerChange {
	var output []string = []string{}
	auxAfter := make(map[string]interface{})
	changes := []MapInnerChange{}
	for k, v := range after {
		auxAfter[k] = v
	}
	for key, value := range before {
		delete(auxAfter, key)
		change := MapInnerChange{Name: key}
		switch value.(type) {
		case map[string]interface{}:
			output = append(output, fmt.Sprintf("Printing the output for key %s", key))
			change.mapChanges = stringMap(value.(map[string]interface{}), after[key].(map[string]interface{}), unknown[key].(map[string]interface{}))
			output = append(output, fmt.Sprintf("Finishing printing the output for key %s", key))
		case string:
			change = stringChange(key, value.(string), after, unknown, change)
		case int:
			change = stringChange(key, fmt.Sprintf("%d", value.(int)), after, unknown, change)
		case bool:
			change = stringChange(key, fmt.Sprintf("%t", value.(bool)), after, unknown, change)
		default:
			fmt.Printf("Default option for key %s with value %v\n", key, value, change)
		}
		changes = append(changes, change)
	}
	for key, value := range auxAfter {
		// The chang to string will fail if it's not it
		changes = append(changes, MapInnerChange{Name: key, After: value.(string), Type: "Add"})
	}
	return changes
}

// We don't have the creating, only in after
func (p PlanChange) String() string {
	changes := MapChange{Name: p.Address}
	before := p.Change.Before.(map[string]interface{})
	after := p.Change.After.(map[string]interface{})
	unknown := p.Change.AfterUnknown.(map[string]interface{})

	changes.Changes = stringMap(before, after, unknown)
	return fmt.Sprintf("%+v\n", changes)
}

func (m MapChange) GenerateTemplate() {

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
