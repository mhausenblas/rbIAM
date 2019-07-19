package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mhausenblas/kubecuddler"
)

// kubeIdentity queries the local Kube config to retrieve the configuration.
func (e *Entity) kubeIdentity() error {
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
	e.KubeConfig = kconf
	return nil
}

// kubeServiceAccounts retrieve the service accounts in the cluster
func (e *Entity) kubeServiceAccounts() error {
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
	e.ServiceAccounts = make(map[string]ServiceAccount)
	for _, sa := range sal.Items {
		e.ServiceAccounts[namespaceit(sa.Namespace, sa.Name)] = sa
	}
	return nil
}

func namespaceit(ns, name string) string {
	return fmt.Sprintf("%v:%v", ns, name)
}
