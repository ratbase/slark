package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
)

var (
	vercel_token      string
	vercel_org_id     string
	vercel_project_id string
	telegram_token    string
	telegram_chat_id  string
)

type Project struct {
	project_name  string
	deploy_branch string
	build_folder  string
}

func main() {
	project := Project{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What is your Project Name?").
				Placeholder("slark_example").
				Value(&project.project_name).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("project name cannot be empty")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Oops, something went wrong:", err)
		return
	}

	fmt.Println("You entered: " + project.project_name)
}
