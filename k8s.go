package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mhausenblas/kubecuddler"
)

// kubeIdentity queries the local Kube config to retrieve the configuration.
func (ag *AccessGraph) kubeIdentity() error {
	res, err := kubecuddler.Kubectl(false, false, "", "config", "view", "--minify", "--output", "json")
	if err != nil {
		return err
	}
	sr := strings.NewReader(res)
	decoder := json.NewDecoder(sr)
	kconf := &Config{}
	err = decoder.Decode(kconf)
	if err != nil {
		return err
	}
	ag.KubeConfig = kconf
	return nil
}

// kubeServiceAccounts retrieve the service accounts in the cluster
func (ag *AccessGraph) kubeServiceAccounts() error {
	res, err := kubecuddler.Kubectl(false, false, "", "get", "sa", "--all-namespaces", "--output", "json")
	if err != nil {
		return err
	}
	sr := strings.NewReader(res)
	decoder := json.NewDecoder(sr)
	sal := ServiceAccountList{}
	err = decoder.Decode(&sal)
	if err != nil {
		return err
	}
	ag.ServiceAccounts = make(map[string]ServiceAccount)
	for _, sa := range sal.Items {
		ag.ServiceAccounts[namespaceit(sa.Namespace, sa.Name)] = sa
	}
	return nil
}

// namespaceit joins Kubernetes namespace and name.
func namespaceit(ns, name string) string {
	return fmt.Sprintf("%v:%v", ns, name)
}

// format provides a textual rendering of the service account
func formatSA(sa *ServiceAccount) string {
	var secrets strings.Builder
	for _, sec := range sa.Secrets {
		secrets.WriteString(sec.Name + " ")
	}
	return fmt.Sprintf(
		"     Namespace: %v\n"+
			"     Name: %v\n"+
			"     Secrets: %v\n",
		sa.Namespace,
		sa.Name,
		secrets.String(),
	)
}
