package terraform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type TFClient struct {
	*tfexec.Terraform
	VarFile string
	ChDir   string
}

func NewClient(dir string) (*TFClient, error) {
	tf, err := tfexec.NewTerraform(dir, "terraform")
	var file string
	if err != nil {
		return nil, err
	}
	workspace, _ := tf.WorkspaceShow(context.Background())
	if _, err = os.Stat(fmt.Sprintf("%s/vars/%s.tfvars", dir, workspace)); err == nil {
		file = fmt.Sprintf("%s/vars/%s.tfvars", dir, workspace)
	}
	return &TFClient{tf, file, dir}, nil
}

func (tf *TFClient) Plan() (list ResourceChangeList, err error) {
	//Generate plan
	var buffer bytes.Buffer
	workspace, _ := tf.WorkspaceShow(context.Background())
	if _, err = os.Stat(fmt.Sprintf("%s/vars/%s.tfvars", tf.ChDir, workspace)); err != nil {
		_, err = tf.PlanJSON(context.Background(), &buffer)
	} else {
		_, err = tf.PlanJSON(context.Background(), &buffer, tfexec.VarFile(fmt.Sprintf("%s/vars/%s.tfvars", tf.ChDir, workspace)))
	}
	buffer.UnreadByte()
	lines := strings.Split(buffer.String(), "\n")

	for _, line := range lines {
		var output ResourceChange
		json.Unmarshal([]byte(line), &output)
		list = append(list, output)
	}
	if err != nil {
		errors := []string{}
		list = list.Filter(func(rc ResourceChange) bool {
			return strings.Contains(rc.Message, "error")
		})
		for _, e := range list {
			errors = append(errors, fmt.Sprintf(" > %s", e.Message))
		}
		return nil, fmt.Errorf("There were some errors:\n%s", strings.Join(errors, "\n"))
	}

	return list.Filter(func(rc ResourceChange) bool {
		return rc.Change.Action != ""
	}), nil
}
