package main

import (
	"github.com/c-bata/go-prompt"
)

// toplevel represents the top level choices in the interaction
func toplevel(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "iam-roles", Description: "Select an AWS IAM role to explore"},
		{Text: "k8s-sa", Description: "Select an Kubernetes service accounts to explore"},
		{Text: "help", Description: "Explain how it works and show available commands"},
		{Text: "quit", Description: "Terminate the interactive session and quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// selectRole allows user to select an IAM role by ARN
func selectRole(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for rolearn := range ag.Roles {
		s = append(s, prompt.Suggest{Text: rolearn})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

// selectSA allows user to select an Kubernetes service account
func selectSA(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for saname := range ag.ServiceAccounts {
		s = append(s, prompt.Suggest{Text: saname})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}
