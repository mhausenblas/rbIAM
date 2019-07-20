package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// user queries IAM to retrieve info on the user issuing the request.
func (ag *AccessGraph) user(cfg aws.Config) error {
	svc := iam.New(cfg)
	req := svc.GetUserRequest(&iam.GetUserInput{})
	res, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	ag.User = res.User
	return nil
}

// formatCaller provides a textual rendering of the combined IAM user and caller information.
func formatCaller(ag *AccessGraph) string {
	user := ag.User
	caller := ag.Caller
	return fmt.Sprintf(
		"     Account ID: %v\n"+
			"     User name: %v\n"+
			"     User ID: %v\n"+
			"     Caller ID: %v\n"+
			"     Path: %v\n"+
			"     Created at: %v\n"+
			"     Tags: %v\n",
		*caller.Account,
		*user.UserName,
		*user.UserId,
		*caller.UserId,
		*user.Path,
		user.CreateDate,
		user.Tags,
	)
}

// callerIdentity queries STS to retrieve the identity of the caller.
func (ag *AccessGraph) callerIdentity(cfg aws.Config) error {
	svc := sts.New(cfg)
	req := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	res, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	ag.Caller = res
	return nil
}

// roles queries IAM for roles in use, related to EKS.
// This is done simply by checking if the role ARN contains EKS or eks.
func (ag *AccessGraph) roles(cfg aws.Config) error {
	svc := iam.New(cfg)
	req := svc.ListRolesRequest(&iam.ListRolesInput{})
	res, err := req.Send(context.TODO())
	if err != nil {
		return err
	}
	ag.Roles = make(map[string]iam.Role)
	for _, role := range res.Roles {
		rolearn := *role.Arn
		ag.Roles[rolearn] = role
	}
	return nil
}

// formatRole provides a textual rendering of a role
func formatRole(role *iam.Role) string {
	arpd := ""
	u, err := url.QueryUnescape(*role.AssumeRolePolicyDocument)
	if err == nil {
		arpd = u
	}
	return fmt.Sprintf(
		"     Name: %v\n"+
			"     ID: %v\n"+
			"     Path: %v\n"+
			"     Assume role by: %v\n"+
			"     Maximum session duration: %v sec\n"+
			"     Created at: %v\n"+
			"     Tags: %v\n",
		*role.RoleName,
		*role.RoleId,
		*role.Path,
		arpd,
		*role.MaxSessionDuration,
		role.CreateDate,
		role.Tags,
	)
}

// policies queries IAM for attached policies, related to EKS.
// This is done simply by checking if the policy ARN contains EKS or eks.
func (ag *AccessGraph) policies(cfg aws.Config) error {
	svc := iam.New(cfg)
	req := svc.ListPoliciesRequest(&iam.ListPoliciesInput{OnlyAttached: aws.Bool(true)})
	res, err := req.Send(context.TODO())
	if err != nil {
		return err
	}
	ag.Policies = make(map[string]iam.Policy)
	for _, policy := range res.Policies {
		policyarn := *policy.Arn
		ag.Policies[policyarn] = policy
	}
	return nil
}

// formatPolicy provides a textual rendering of a policy.
func formatPolicy(policy *iam.Policy) string {
	return fmt.Sprintf(
		"     Name: %v\n"+
			"     ID: %v\n"+
			"     Path: %v\n"+
			"     Number of entities the policy is attached: %v\n"+
			"     Created at: %v\n"+
			"     Updated at: %v\n",
		*policy.PolicyName,
		*policy.PolicyId,
		*policy.Path,
		*policy.AttachmentCount,
		*policy.CreateDate,
		*policy.UpdateDate,
	)
}
