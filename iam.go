package main

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// user queries IAM to retrieve info on the user issuing the request.
func (e *Entity) user(cfg aws.Config) error {
	svc := iam.New(cfg)
	req := svc.GetUserRequest(&iam.GetUserInput{})
	res, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	e.User = res.User
	return nil
}

// callerIdentity queries STS to retrieve the identity of the caller.
func (e *Entity) callerIdentity(cfg aws.Config) error {
	svc := sts.New(cfg)
	req := svc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	res, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	e.Caller = res
	return nil
}

// roles queries IAM for roles in use, related to EKS.
// This is done simply by checking if the role ARN contains EKS or eks.
func (e *Entity) roles(cfg aws.Config) error {
	svc := iam.New(cfg)
	req := svc.ListRolesRequest(&iam.ListRolesInput{})
	res, err := req.Send(context.TODO())
	if err != nil {
		return err
	}
	for _, role := range res.Roles {
		if strings.Contains(*role.Arn, "EKS") ||
			strings.Contains(*role.Arn, "eks") {
			e.Roles = append(e.Roles, role)
		}
	}
	return nil
}

// policies queries IAM for attached policies, related to EKS.
// This is done simply by checking if the policy ARN contains EKS or eks.
func (e *Entity) policies(cfg aws.Config) error {
	svc := iam.New(cfg)
	req := svc.ListPoliciesRequest(&iam.ListPoliciesInput{OnlyAttached: aws.Bool(true)})
	res, err := req.Send(context.TODO())
	if err != nil {
		return err
	}
	for _, policy := range res.Policies {
		if strings.Contains(*policy.Arn, "EKS") ||
			strings.Contains(*policy.Arn, "eks") {
			e.Policies = append(e.Policies, policy)
		}
	}
	return nil
}
