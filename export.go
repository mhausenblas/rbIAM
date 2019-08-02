package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/emicklei/dot"
)

// exportRaw exports the trace as a raw dump in JSON format into a file
// in the current working directory with a name of 'rbiam-trace-NNNNNNNNNN' with
// the NNNNNNNNNN being the Unix timestamp of the creation time, for example:
// rbiam-trace-1564315687.json
func exportRaw(trace []string, ag *AccessGraph) (string, error) {
	dump := ""
	for _, item := range trace {
		itype, ikey := extractTK(item)
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

// exportGraph exports the trace as a graph in DOT format into a file
// in the current working directory with a name of 'rbiam-trace-NNNNNNNNNN' with
// the NNNNNNNNNN being the Unix timestamp of the creation time, for example:
// rbiam-trace-1564315687.dot
func exportGraph(trace []string, ag *AccessGraph) (string, error) {
	g := dot.NewGraph(dot.Directed)

	// legend:
	lsa := formatAsServiceAccount(g.Node("SERVICE ACCOUNT"))
	lsecret := formatAsSecret(g.Node("SECRET"))
	lpod := formatAsPod(g.Node("POD"))
	g.Edge(lpod, lsa)
	g.Edge(lsa, lsecret)

	for _, item := range trace {
		itype, ikey := extractTK(item)
		switch itype {
		case "IAM role":
		case "IAM policy":
		case "Kubernetes service account":
			formatAsServiceAccount(g.Node(ikey))
		case "Kubernetes secret":
			formatAsSecret(g.Node(ikey))
		case "Kubernetes pod":
			formatAsPod(g.Node(ikey))
		}
	}
	filename := fmt.Sprintf("rbiam-trace-%v.dot", time.Now().Unix())
	err := ioutil.WriteFile(filename, []byte(g.String()), 0644)
	if err != nil {
		return "", err
	}
	return filename, nil
}

// extractTK takes a history item in the form [TYPE] KEY
// and return t as the TYPE and k as the KEY, for example:
// [Kubernetes service account] default:s3-echoer ->
// t == Kubernetes service account
// k == default:s3-echoer
func extractTK(item string) (t, k string) {
	t = strings.TrimPrefix(strings.Split(item, "]")[0], "[")
	k = strings.TrimSpace(strings.Split(item, "]")[1])
	return
}

func formatAsServiceAccount(n dot.Node) dot.Node {
	return n.Attr("style", "filled").Attr("fillcolor", "#1BFF9F").Attr("fontcolor", "#000000").Attr("fontname", "Helvetica")
}

func formatAsSecret(n dot.Node) dot.Node {
	return n.Attr("style", "filled").Attr("fillcolor", "#F9ED49").Attr("fontcolor", "#000000").Attr("fontname", "Helvetica")
}

func formatAsPod(n dot.Node) dot.Node {
	return n.Attr("style", "filled").Attr("fillcolor", "#4260FA").Attr("fontcolor", "#f0f0f0").Attr("fontname", "Helvetica")
}
