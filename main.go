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

// history keeps the selected items such as roles or service accounts around
var history []string

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("Can't load AWS config: %v", err.Error())
		os.Exit(1)
	}

	fmt.Println("Gathering info from IAM and Kubernetes. This may take a bit, please stand by ...")
	ag = NewAccessGraph(cfg)
	// fmt.Println(ag)
	tracecntr := 0
	tracemode := false
	cursel := "help" // make sure to first show the help to guide users what to do
	for {
		switch cursel {
		case "iam-user":
			presult(formatCaller(ag))
		case "iam-roles":
			targetrole := prompt.Input("  ↪ ", selectRole,
				prompt.OptionMaxSuggestion(30),
				prompt.OptionSuggestionBGColor(prompt.DarkBlue))
			if role, ok := ag.Roles[targetrole]; ok {
				presult(formatRole(&role))
				appendhist("IAM role", targetrole)
				if tracemode {
					tracecntr++
				}
			}
		case "iam-policies":
			targetpolicy := prompt.Input("  ↪ ", selectPolicy,
				prompt.OptionMaxSuggestion(30),
				prompt.OptionSuggestionBGColor(prompt.DarkBlue))
			if policy, ok := ag.Policies[targetpolicy]; ok {
				presult(formatPolicy(&policy))
				appendhist("IAM policy", targetpolicy)
				if tracemode {
					tracecntr++
				}
			}
		case "k8s-sa":
			targetsa := prompt.Input("  ↪ ", selectSA,
				prompt.OptionMaxSuggestion(30),
				prompt.OptionSuggestionBGColor(prompt.DarkBlue))
			if sa, ok := ag.ServiceAccounts[targetsa]; ok {
				presult(formatSA(&sa))
				appendhist("Kubernetes service account", targetsa)
				if tracemode {
					tracecntr++
				}
			}
		case "k8s-secrets":
			targetsec := prompt.Input("  ↪ ", selectSecret,
				prompt.OptionMaxSuggestion(30),
				prompt.OptionSuggestionBGColor(prompt.DarkBlue))
			if secret, ok := ag.Secrets[targetsec]; ok {
				presult(formatSecret(&secret))
				appendhist("Kubernetes secret", targetsec)
				if tracemode {
					tracecntr++
				}
			}
		case "k8s-pods":
			targetpod := prompt.Input("  ↪ ", selectPod,
				prompt.OptionMaxSuggestion(30),
				prompt.OptionSuggestionBGColor(prompt.DarkBlue))
			if pod, ok := ag.Pods[targetpod]; ok {
				presult(formatPod(&pod))
				appendhist("Kubernetes pod", targetpod)
				if tracemode {
					tracecntr++
				}
			}
		case "history":
			dumphist()
		case "sync":
			fmt.Println("Gathering info from IAM and Kubernetes. This may take a bit, please stand by ...")
			ag = NewAccessGraph(cfg)
		case "trace":
			tracemode = true
			tracecntr = 0
			presult("Starting to trace now. Use an 'export-xxx' command to stop tracing and export to one of the supported formats.\n")
		case "export-raw":
			tracemode = false
			fn, err := exportRaw(history[0:tracecntr], ag)
			if err != nil {
				pwarning(fmt.Sprintf("Can't export trace: %v\n", err))
				continue
			}
			presult(fmt.Sprintf("Raw trace exported to %v\n", fn))
		case "help":
			presult(fmt.Sprintf("\nThis is rbIAM in version %v\n\n", Version))
			presult(strings.Repeat("-", 80))
			presult("\nSelect one of the supported query commands:\n")
			presult("- iam-user … to look up the calling AWS IAM user\n")
			presult("- iam-roles … to look up an AWS IAM role by ARN\n")
			presult("- iam-policies … to look up an AWS IAM policy by ARN\n")
			presult("- k8s-sa … to look up an Kubernetes service account\n")
			presult("- k8s-secrets … to look up a Kubernetes secret\n")
			presult("- k8s-pods … to look up a Kubernetes pod\n")
			presult("- history … show history\n")
			presult("- sync … to refresh the local data\n")
			presult("- trace … start tracing\n")
			presult("- export-raw … stop tracing and export trace to JSON dump in current working directory\n")
			presult(strings.Repeat("-", 80))
			presult("\n\nNote: simply start typing and/or use the tab and cursor keys to select.\n")
			presult("CTRL+L clears the screen and if you're stuck type 'help' or 'quit' to leave.\n\n")
		case "quit":
			presult("bye!\n")
			os.Exit(0)
		default:
			presult("Not yet implemented, sorry\n")
		}
		cursel = prompt.Input("? ", toplevel,
			prompt.OptionMaxSuggestion(20),
			prompt.OptionSuggestionBGColor(prompt.DarkBlue),
			prompt.OptionSelectedDescriptionBGColor(prompt.DarkBlue))
	}
}

func appendhist(kind, entry string) {
	history = append([]string{fmt.Sprintf("[%v] %v", kind, entry)}, history...)
}

func dumphist() {
	for _, entry := range history {
		presult(fmt.Sprintf("%v\n", entry))
	}
}

// Below some helper functions for output. For available colors see:
// https://misc.flogisoft.com/bash/tip_colors_and_formatting

// presult writes msg in blue to stdout and note that you need to take
// care of newlines yourself.
func presult(msg string) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[34m%v\x1b[0m", msg)
}

// pwarning writes msg in red to stdout and note that you need to take
// care of newlines yourself.
func pwarning(msg string) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[31m%v\x1b[0m", msg)
}
