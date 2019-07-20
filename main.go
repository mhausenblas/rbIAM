package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/c-bata/go-prompt"
)

// Version is the CLI tool version, provided by the build process, see Makefile
var Version string

// ag is the global access graph, mostly read-only besides the init phase when
// all the pertinent information is gathered from IAM and Kubernetes via NewAccessGraph()
var ag *AccessGraph

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("Can't load AWS config: %v", err.Error())
		os.Exit(1)
	}

	fmt.Println("Gathering info from IAM and Kubernetes. This may take a bit, please stand by ...")
	ag = NewAccessGraph(cfg)
	// fmt.Println(ag)

	cursel := "help" // make sure to first show the help to guide users what to do
	for {
		switch cursel {
		case "iam-roles":
			targetrole := prompt.Input("  ↪ ", selectRole)
			if role, ok := ag.Roles[targetrole]; ok {
				presult(formatRole(&role))
			}
		case "k8s-sa":
			targetsa := prompt.Input("  ↪ ", selectSA)
			if sa, ok := ag.ServiceAccounts[targetsa]; ok {
				presult(formatSA(&sa))
			}
		case "help":
			presult(fmt.Sprintf("\nThis is rbIAM in version %v\n\n", Version))
			presult(strings.Repeat("-", 80))
			presult("\nSelect one of the supported query commands:\n")
			presult("- iam-roles … to look up an AWS IAM role by ARN\n")
			presult("- k8s-sa … to look up an Kubernetes service account\n")
			presult(strings.Repeat("-", 80))
			presult("\n\nNote: simply start typing and/or use the tab and cursor keys to select.\n")
			presult("CTRL+L clears the screen and if you're stuck type 'help'.\n\n")
		case "quit":
			presult("bye!\n")
			os.Exit(0)
		default:
			presult("Not yet implemented, sorry\n")
		}
		cursel = prompt.Input("? ", toplevel)
	}
}

// presult writes msg in  blue to stdout. you need to take care of
// newlines yourself and for available colors see:
// https://misc.flogisoft.com/bash/tip_colors_and_formatting
func presult(msg string) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[34m%v\x1b[0m", msg)
}
