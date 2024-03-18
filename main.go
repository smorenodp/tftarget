package main

import (
	"flag"
	"os"

	"github.com/smorenodp/tftarget/cli"
)

// func main() {
// 	var plan, apply bool
// 	var var_file string
// 	ifcColor := color.New(color.FgCyan, color.Bold)

// 	flag.StringVar(&var_file, "var-file", "", "The var file to use")
// 	flag.BoolVar(&plan, "plan", false, "If we want to make a terraform plan")
// 	flag.BoolVar(&apply, "apply", false, "If we want to make a terraform apply")
// 	flag.Parse()

// 	tf, err := terraform.NewClient()
// 	if err != nil {
// 		log.Fatalf("There was an error - %s\n", err)
// 	}
// 	list, err := tf.Plan()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	filteredList := list.Filter(func(resource terraform.ResourceChange) bool {
// 		action := resource.Change.Action
// 		if action == "update" || action == "delete" || action == "create" || action == "replace" {
// 			return true
// 		} else {
// 			return false
// 		}
// 	})

// 	ifcColor.Add(color.Bold).Println("Color legend")
// 	fmt.Println("")
// 	terraform.PrintColorLegend()
// 	fmt.Println("")
// 	ifcColor.Add(color.Bold).Println("Showing resources with changes...")
// 	fmt.Println("")
// 	for index, resource := range filteredList {
// 		fmt.Printf("%s %s\n", ifcColor.Sprintf("%d)", index), resource)
// 	}

// }

func main() {
	var plan, apply bool
	var var_file, dir string
	currentDir, _ := os.Getwd()

	flag.StringVar(&var_file, "var-file", "", "The var file to use")
	flag.BoolVar(&plan, "plan", false, "If we want to make a terraform plan")
	flag.BoolVar(&apply, "apply", false, "If we want to make a terraform apply")
	flag.StringVar(&dir, "dir", currentDir, "Directory to launch terraform")
	flag.Parse()

	// tf, err := terraform.NewClient(dir)
	// if err != nil {
	// 	log.Fatalf("There was an error - %s\n", err)
	// }
	// m := cli.New(tf)
	// m.Run()

	m := cli.New()
	m.Run()
}
