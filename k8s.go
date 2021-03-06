package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mhausenblas/kubecuddler"
)

// namespaceit joins Kubernetes namespace and name.
func namespaceit(ns, name string) string {
	return fmt.Sprintf("%v:%v", ns, name)
}

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

// kubeServiceAccounts retrieves the service accounts in the cluster.
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

// formatSA provides a textual rendering of the service account.
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

// kubeSecrets retrieves the secrets in the cluster.
func (ag *AccessGraph) kubeSecrets() error {
	res, err := kubecuddler.Kubectl(false, false, "", "get", "secrets", "--all-namespaces", "--output", "json")
	if err != nil {
		return err
	}
	sr := strings.NewReader(res)
	decoder := json.NewDecoder(sr)
	secl := SecretList{}
	err = decoder.Decode(&secl)
	if err != nil {
		return err
	}
	ag.Secrets = make(map[string]Secret)
	for _, secret := range secl.Items {
		ag.Secrets[namespaceit(secret.Namespace, secret.Name)] = secret
	}
	return nil
}

// formatSecret provides a textual rendering of the secret.
func formatSecret(secret *Secret) string {
	strdatamap := ""
	// case of string data present:
	if len(secret.StringData) != 0 {
		for k, v := range secret.StringData {
			strdatamap += fmt.Sprintf("%v: %v ", k, v)
		}
	}
	datamap := ""
	// case of data present:
	if len(secret.Data) != 0 {
		for k, v := range secret.Data {
			datamap += fmt.Sprintf("\n      %v: %v ", k, string(v))
		}
	}
	return fmt.Sprintf(
		"     Namespace: %v\n"+
			"     Name: %v\n"+
			"     Type: %v\n"+
			"     String data: %v\n"+
			"     Data: %v\n",
		secret.Namespace,
		secret.Name,
		secret.Type,
		strdatamap,
		datamap,
	)
}

// kubePods retrieves the pods in the cluster.
func (ag *AccessGraph) kubePods() error {
	res, err := kubecuddler.Kubectl(false, false, "", "get", "pods", "--all-namespaces", "--output", "json")
	if err != nil {
		return err
	}
	sr := strings.NewReader(res)
	decoder := json.NewDecoder(sr)
	podl := PodList{}
	err = decoder.Decode(&podl)
	if err != nil {
		return err
	}
	ag.Pods = make(map[string]Pod)
	for _, pod := range podl.Items {
		ag.Pods[namespaceit(pod.Namespace, pod.Name)] = pod
	}
	return nil
}

// formatPod provides a textual rendering of the pod.
func formatPod(pod *Pod) string {
	strcontainers := ""
	for _, container := range pod.Spec.Containers {
		strcontainers += fmt.Sprintf(
			"      Name: %v\n"+
				"       Image pull policy: %v\n"+
				"       Image: %v\n"+
				"       Command: %v\n"+
				"       Args: %v\n"+
				"       Env: %v",
			container.Name,
			container.ImagePullPolicy,
			container.Image,
			container.Command,
			container.Args,
			container.Env,
		)
	}
	return fmt.Sprintf(
		"     Namespace: %v\n"+
			"     Name: %v\n"+
			"     Service account: %v\n"+
			"     Image pull secrets: %v\n"+
			"     Volumes: %v\n"+
			"     Containers:\n %v\n"+
			"     Host IP: %v\n"+
			"     Pod IP: %v\n"+
			"     Phase: %v\n",
		pod.Namespace,
		pod.Name,
		pod.Spec.ServiceAccountName,
		pod.Spec.ImagePullSecrets,
		pod.Spec.Volumes,
		strcontainers,
		pod.Status.HostIP,
		pod.Status.PodIP,
		pod.Status.Phase,
	)
}
