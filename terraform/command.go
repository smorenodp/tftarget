package terraform

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type Command struct {
	Type    string
	Targets []string
	ChDir   string
	VarFile string
	items   *[]list.Item
}

func NewCommand(dir string, varfile string) Command {
	return Command{"<type>", []string{}, dir, varfile, nil}
}

func (c *Command) SetItems(items []list.Item) {
	c.items = &items
}

func (c Command) String() string {
	var targets []string
	c.Targets = []string{}
	for _, i := range *c.items {
		if i.(*ResourceChange).IsSelected() {
			c.Targets = append(c.Targets, i.(*ResourceChange).Title())
		}
	}
	for _, t := range c.Targets {
		targets = append(targets, fmt.Sprintf("-target=%s", t))
	}
	return fmt.Sprintf("terraform -chdir=%s %s -var-file=%s %s", c.ChDir, c.Type, c.VarFile, strings.Join(targets, " "))
}
