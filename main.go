package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/c-bata/go-prompt"
)

var entity *Entity

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("Can't load config: %v", err.Error())
		os.Exit(1)
	}
	fmt.Println("Gathering info, this may take a bit ...")
	entity = NewEntity(cfg)
	// fmt.Println(entity)
	for {
		cursel := prompt.Input("? ", identities)
		switch cursel {
		case "iam-roles":
			targetrole := prompt.Input("  ↪ ", selectRole)
			if role, ok := entity.Roles[targetrole]; ok {
				presult(formatRole(&role))
			}
		case "k8s-sa":
			targetsa := prompt.Input("  ↪ ", selectSA)
			if sa, ok := entity.ServiceAccounts[targetsa]; ok {
				presult(formatSA(&sa))
			}
		case "help":
			fmt.Println("Select one of the supported query commands:")
			fmt.Println("- iam-roles … to look up an AWS IAM role by ARN")
			fmt.Println("- k8s-sa … to look up an Kubernetes service account")
			fmt.Println("\nNote: simply start typing and use tab and cursor keys to select. Also, CTRL+L clears the screen.")
		case "quit":
			os.Exit(0)
		default:
			fmt.Println("Not yet implemented, sorry")
		}
	}
}

func identities(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "iam-roles", Description: "Select an AWS IAM role to explore"},
		{Text: "k8s-sa", Description: "Select an Kubernetes service accounts to explore"},
		{Text: "help", Description: "Explain how it works and show available commands"},
		{Text: "quit", Description: "Terminate the interactive session and quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func selectRole(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for rolearn := range entity.Roles {
		s = append(s, prompt.Suggest{Text: rolearn})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

func selectSA(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for saname := range entity.ServiceAccounts {
		s = append(s, prompt.Suggest{Text: saname})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

// presult writes msg in light blue to stdout
// see also https://misc.flogisoft.com/bash/tip_colors_and_formatting
func presult(msg string) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[34m%v\x1b[0m\n", msg)
}
