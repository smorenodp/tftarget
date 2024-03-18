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
	f, _ := os.Create("resources")
	defer f.Close()
	workspace, _ := tf.WorkspaceShow(context.Background())
	options := []tfexec.PlanOption{tfexec.Out("/tmp/.tftarget")}
	if _, err = os.Stat(fmt.Sprintf("%s/vars/%s.tfvars", tf.ChDir, workspace)); err == nil {
		options = append(options, tfexec.VarFile(fmt.Sprintf("%s/vars/%s.tfvars", tf.ChDir, workspace)))
	}
	_, err = tf.Terraform.Plan(context.Background(), options...)
	f.WriteString("--------------------------------------------------------\n")
	plan, _ := tf.Terraform.ShowPlanFile(context.Background(), "/tmp/.tftarget")
	for _, c := range plan.ResourceChanges {
		aux := NewPlanChange(c)
		f.WriteString("--------------------------------------------------------\n")
		f.WriteString(fmt.Sprintf("Resource %s", c.Address))
		f.WriteString(fmt.Sprintf("%s\n", aux))
		f.WriteString("--------------------------------------------------------\n\n")
	}
	return nil, err
}

func (tf *TFClient) PlanOld() (list ResourceChangeList, err error) {
	//Generate plan
	var buffer bytes.Buffer
	f, _ := os.Create("output")
	workspace, _ := tf.WorkspaceShow(context.Background())
	if _, err = os.Stat(fmt.Sprintf("%s/vars/%s.tfvars", tf.ChDir, workspace)); err != nil {
		_, err = tf.PlanJSON(context.Background(), &buffer)
	} else {
		_, err = tf.PlanJSON(context.Background(), &buffer, tfexec.VarFile(fmt.Sprintf("%s/vars/%s.tfvars", tf.ChDir, workspace)))
	}

	buffer.UnreadByte()
	output := buffer.String()
	f, _ := os.Create("output")
	f.WriteString(output)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		var output ResourceChange
		json.Unmarshal([]byte(line), &output)
		f.WriteString("-----------------------------------------------------")
		f.WriteString(fmt.Sprintf("The line is %s\n", line))
		f.WriteString(fmt.Sprintf("The resource %+v has action %s\n", output.Change.Resource, output.Change.Action))
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
		return (rc.Change.Action != "" && rc.Type != "resource_drift")
	}), nil
}
