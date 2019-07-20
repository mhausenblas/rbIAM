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
	Caller          *sts.GetCallerIdentityOutput
	User            *iam.User
	KubeConfig      *Config
	Roles           map[string]iam.Role
	Policies        []iam.Policy
	ServiceAccounts map[string]ServiceAccount
}

// NewAccessGraph a new access graph for the currently authenticated AWS user,
// retrieving IAM-related as well as Kubernetes-related info. We try to be graceful
// here but if the IAM queries fail, there's no point in continuing and we exit early.
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
			"Kube service accounts: %+v\n",
		ag.User,
		ag.Caller,
		ag.Roles,
		ag.Policies,
		ag.KubeConfig.CurrentContext,
		ag.ServiceAccounts,
	)
}
