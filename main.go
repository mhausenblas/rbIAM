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
	// fmt.Println(entity)
	for {
		switch selectStartingPoint(entity) {
		case "IAM role":
			targetrole := selectRole(entity)
			fmt.Println(entity.Roles[targetrole])
		case "Kubernetes service account":
			targetsa := selectSA(entity)
			fmt.Println(entity.ServiceAccounts[targetsa])
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Not yet implemented, sorry")
		}
	}
}

func selectStartingPoint(e *Entity) (selection string) {
	survey.AskOne(&survey.Select{
		Message: "What should I use as the starting point?",
		Options: []string{"IAM role", "Kubernetes service account", "exit"},
		Default: "IAM role",
		Help:    "Select from which side to start exploring, either using an IAM role or a Kubernetes service account",
	}, &selection)
	return
}

func selectRole(e *Entity) (selection string) {
	roles := []string{}
	for rolearn := range e.Roles {
		roles = append(roles, rolearn)
	}
	survey.AskOne(&survey.Select{
		Message: "Which IAM role would you like to use?",
		Options: roles,
		Help:    "Select an IAM role to explore. You can filter by start typing.",
	}, &selection)
	return
}

func selectSA(e *Entity) (selection string) {
	sas := []string{}
	for saname := range e.ServiceAccounts {
		sas = append(sas, saname)
	}
	survey.AskOne(&survey.Select{
		Message: "Which Kubernetes service account would you like to use?",
		Options: sas,
		Help:    "Select an Kubernetes service account to explore. You can filter by start typing.",
	}, &selection)
	return
}
