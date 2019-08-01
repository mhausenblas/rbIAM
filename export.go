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
	lsa := g.Node("SERVICE ACCOUNT").Attr("style", "filled").Attr("fillcolor", "#33ff33").Attr("fontcolor", "#000000")
	lsecret := g.Node("SECRET").Attr("style", "filled").Attr("fillcolor", "#ff3399").Attr("fontcolor", "#000000")
	lpod := g.Node("POD").Attr("style", "filled").Attr("fillcolor", "#9900ff").Attr("fontcolor", "#000000")
	g.Edge(lpod, lsa)
	g.Edge(lsa, lsecret)

	for _, item := range trace {
		itype, ikey := extractTK(item)
		switch itype {
		case "IAM role":
		case "IAM policy":
		case "Kubernetes service account":
			g.Node(ikey).Attr("style", "filled").Attr("fillcolor", "#33ff33").Attr("fontcolor", "#000000")
		case "Kubernetes secret":
			g.Node(ikey).Attr("style", "filled").Attr("fillcolor", "#ff3399").Attr("fontcolor", "#000000")
		case "Kubernetes pod":
			g.Node(ikey).Attr("style", "filled").Attr("fillcolor", "#9900ff").Attr("fontcolor", "#000000")
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
