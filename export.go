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
	g.Edge(lpod, lsa, "uses").Attr("fontname", "Helvetica")
	g.Edge(lsa, lsecret, "has").Attr("fontname", "Helvetica")

	// first let's draw the nodes and remember the
	// graph entry points for traversals to later draw
	// the edges (Kubernetes pods and IAM roles)
	pods := make(map[string]dot.Node)
	sas := make(map[string]dot.Node)
	for _, item := range trace {
		itype, ikey := extractTK(item)
		switch itype {
		case "IAM role":
		case "IAM policy":
		case "Kubernetes service account":
			sas[ikey] = formatAsServiceAccount(g.Node(ikey))
		case "Kubernetes secret":
			formatAsSecret(g.Node(ikey))
		case "Kubernetes pod":
			pods[ikey] = formatAsPod(g.Node(ikey))
		}
	}

	// next, we draw the edges, using K8s pods and IAM roles
	// as the entry points into the graph
	for podname, node := range pods {
		fmt.Printf("Looking at pod %v in node %v\n", podname, node)
		for _, item := range trace {
			itype, ikey := extractTK(item)
			if itype == "Kubernetes service account" {
				podsa := namespaceit(ag.Pods[podname].Namespace, ag.Pods[podname].Spec.ServiceAccountName)
				fmt.Printf("Checking service account %v against the pod's service account %v\n", ikey, podsa)
				if podsa == ikey {
					fmt.Printf("Connecting %v to %v", node, sas[ikey])
					g.Edge(node, sas[ikey])
				}
			}
		}
	}

	// now we can write out the graph into a file in DOT format:
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
