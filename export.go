package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// exportRaw exports the trace as a raw dump in JSON format into a file
// in the current working directory with a name of 'rbiam-trace-NNNNNNNNNN' with
// the NNNNNNNNNN being the Unix timestamp of the creation time, for example:
// rbiam-trace-1564315687.json
func exportRaw(trace []string, ag *AccessGraph) (string, error) {
	dump := ""
	for _, item := range trace {
		itype := strings.TrimPrefix(strings.Split(item, "]")[0], "[")
		ikey := strings.TrimSpace(strings.Split(item, "]")[1])
		switch itype {
		case "IAM role":
			b, err := json.Marshal(ag.Roles[ikey])
			if err != nil {
				return "", err
			}
			dump = fmt.Sprintf("%v\n%v", dump, string(b))
		case "IAM policy":
			b, err := json.Marshal(ag.Policies[ikey])
			if err != nil {
				return "", err
			}
			dump = fmt.Sprintf("%v\n%v", dump, string(b))
		case "Kubernetes service account":
			b, err := json.Marshal(ag.ServiceAccounts[ikey])
			if err != nil {
				return "", err
			}
			dump = fmt.Sprintf("%v\n%v", dump, string(b))
		case "Kubernetes secret":
			b, err := json.Marshal(ag.Secrets[ikey])
			if err != nil {
				return "", err
			}
			dump = fmt.Sprintf("%v\n%v", dump, string(b))
		case "Kubernetes pod":
			b, err := json.Marshal(ag.Pods[ikey])
			if err != nil {
				return "", err
			}
			dump = fmt.Sprintf("%v\n%v", dump, string(b))
		}
	}

	filename := fmt.Sprintf("rbiam-trace-%v.json", time.Now().Unix())
	err := ioutil.WriteFile(filename, []byte(dump), 0644)
	if err != nil {
		return "", err
	}
	return filename, nil
}
