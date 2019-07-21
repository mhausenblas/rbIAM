package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// AccessGraph represents the combined IAM and RBAC access control regime
// found in Kubernetes on AWS. It includes IAM users, roles and AWS service
// with temporary credentials (STS) as well as Kubernetes service accounts
// and users. See also:
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
// https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/
type AccessGraph struct {
	// Caller is the caller of an AWS service (STS).
	Caller *sts.GetCallerIdentityOutput
	// User is the AWS user (IAM).
	User *iam.User
	// KubeConfig is the Kubernetes client configuration.
	KubeConfig *Config
	// Roles is the collection of all IAM roles pertinent to user/caller.
	Roles map[string]iam.Role
	// Policies is the collection of all IAM policies pertinent to user/caller.
	Policies map[string]iam.Policy
	// ServiceAccounts is the collection of all service accounts in the
	// Kubernetes cluster.
	ServiceAccounts map[string]ServiceAccount
	// Secrets is the collection of all secrets in the Kubernetes cluster.
	Secrets map[string]Secret
	// Pods is the collection of all pods in the Kubernetes cluster.
	Pods map[string]Pod
}

// NewAccessGraph a new access graph for the currently authenticated AWS user,
// retrieving IAM-related as well as Kubernetes-related info. We try to be as
// graceful as possbile here but if the IAM queries fail, there's no point in
// continuing and we exit early.
func NewAccessGraph(cfg aws.Config) *AccessGraph {
	ag := &AccessGraph{}
	err := ag.user(cfg)
	if err != nil {
		fmt.Printf("Can't get user: %v", err.Error())
		os.Exit(2)
	}
	err = ag.callerIdentity(cfg)
	if err != nil {
		fmt.Printf("Can't get caller identity: %v", err.Error())
		os.Exit(2)
	}
	err = ag.roles(cfg)
	if err != nil {
		fmt.Printf("Can't get roles: %v", err.Error())
		os.Exit(2)
	}
	err = ag.policies(cfg)
	if err != nil {
		fmt.Printf("Can't get policies: %v", err.Error())
		os.Exit(2)
	}
	err = ag.kubeIdentity()
	if err != nil {
		fmt.Printf("Can't get Kubernetes identity: %v", err.Error())
	}
	err = ag.kubeServiceAccounts()
	if err != nil {
		fmt.Printf("Can't get Kubernetes service accounts: %v", err.Error())
	}
	err = ag.kubeSecrets()
	if err != nil {
		fmt.Printf("Can't get Kubernetes secrets: %v", err.Error())
	}
	err = ag.kubePods()
	if err != nil {
		fmt.Printf("Can't get Kubernetes pods: %v", err.Error())
	}
	return ag
}

// String provides a textual rendering of the access graph
func (ag *AccessGraph) String() string {
	return fmt.Sprintf(
		"User: %v\n"+
			"STS caller identity: %v\n"+
			"EKS roles: %v\n"+
			"EKS policies: %v\n"+
			"Kube context: %+v\n"+
			"Kube service accounts: %+v\n"+
			"Kube secrets: %+v\n"+
			"Kube pods: %+v\n",
		ag.User,
		ag.Caller,
		ag.Roles,
		ag.Policies,
		ag.KubeConfig.CurrentContext,
		ag.ServiceAccounts,
		ag.Secrets,
		ag.Pods,
	)
}
