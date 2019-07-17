package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("Can't load config: %v", err.Error())
		os.Exit(1)
	}
	fmt.Println("Gathering info, this may take a bit ...")
	entity := NewEntity(cfg)

	switch selectStartingPoint(entity) {
	case "myself":
		targetrole := selectRole(entity)
		fmt.Printf("%v", targetrole)
	case "service account":
		targetsa := selectSA(entity)
		fmt.Printf("%v", targetsa)
	default:
		fmt.Println("Not yet implemented, sorry")
	}
}

func selectStartingPoint(e *Entity) (selection string) {
	survey.AskOne(&survey.Select{
		Message: "What should I use as the starting point?",
		Options: []string{"myself", "service account"},
		Default: "myself",
		Help:    "Select 'myself' to use your IAM identity or 'service account' to pick a Kubernetes entity",
	}, &selection)
	return
}

func selectRole(e *Entity) (selection string) {
	roles := []string{}
	for _, role := range e.Roles {
		roles = append(roles, *role.RoleName)
	}
	survey.AskOne(&survey.Select{
		Message: "Which role would you like to use?",
		Options: roles,
		Help:    "Select an IAM role to explore. You can filter by start typing ...",
	}, &selection)
	return
}

func selectSA(e *Entity) (selection string) {
	// sas := []string{}

	return
}
