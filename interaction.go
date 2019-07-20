package main

import (
	"github.com/c-bata/go-prompt"
)

// toplevel represents the top level choices in the interaction.
func toplevel(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "iam-roles", Description: "Select an AWS IAM role to explore"},
		{Text: "iam-policies", Description: "Select an AWS IAM policy to explore"},
		{Text: "k8s-sa", Description: "Select an Kubernetes service accounts to explore"},
		{Text: "k8s-secrets", Description: "Select a Kubernetes secret to explore"},
		{Text: "sync", Description: "Synchronize the local state with IAM and Kubernetes"},
		{Text: "help", Description: "Explain how it works and show available commands"},
		{Text: "quit", Description: "Terminate the interactive session and quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// selectRole allows user to select an IAM role by ARN.
func selectRole(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for rolearn := range ag.Roles {
		s = append(s, prompt.Suggest{Text: rolearn})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

// selectPolicy allows user to select an IAM policy by ARN.
func selectPolicy(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for policyarn := range ag.Policies {
		s = append(s, prompt.Suggest{Text: policyarn})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

// selectSA allows user to select an Kubernetes service account.
func selectSA(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for saname := range ag.ServiceAccounts {
		s = append(s, prompt.Suggest{Text: saname})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}

// selectSecret allows user to select an Kubernetes secret.
func selectSecret(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for secname := range ag.Secrets {
		s = append(s, prompt.Suggest{Text: secname})
	}
	return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
}
