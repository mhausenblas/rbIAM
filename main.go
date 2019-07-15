package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Entity represents the caller of a service, could be a human on the CLI
// or another AWS service with temporary credentials (STS). See also:
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html
type Entity struct {
	Kind       string
	Caller     *sts.GetCallerIdentityOutput
	User       *iam.User
	KubeConfig *Config
	Roles      []iam.Role
	Policies   []iam.Policy
}

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Printf("Can't load config: %v", err.Error())
		os.Exit(1)
	}

	targetEntity := ""
	survey.AskOne(&survey.Select{
		Message: "What entity should I use as the starting point?",
		Options: []string{"myself", "a user", "a service"},
		Default: "myself",
	}, &targetEntity)
	switch targetEntity {
	case "myself":
		entity := NewEntity(targetEntity, cfg)
		fmt.Println(entity)
	default:
		fmt.Println("Not yet implemented, sorry")
	}
}

// NewEntity creates a new entity for the currently authenticated AWS user,
// retrieving both the IAM as well as the Kubernetes-related info. If either of
// the two queries fails, there's no point in continuing, hence we exit early.
func NewEntity(kind string, cfg aws.Config) *Entity {
	entity := &Entity{
		Kind: kind,
	}
	err := entity.user(cfg)
	if err != nil {
		fmt.Printf("Can't get user: %v", err.Error())
		os.Exit(2)
	}
	err = entity.callerIdentity(cfg)
	if err != nil {
		fmt.Printf("Can't get caller identity: %v", err.Error())
		os.Exit(2)
	}
	err = entity.roles(cfg)
	if err != nil {
		fmt.Printf("Can't get roles: %v", err.Error())
		os.Exit(2)
	}
	err = entity.policies(cfg)
	if err != nil {
		fmt.Printf("Can't get policies: %v", err.Error())
		os.Exit(2)
	}
	err = entity.kubeIdentity()
	if err != nil {
		fmt.Printf("Can't get Kubernetes identity: %v", err.Error())
		os.Exit(2)
	}
	return entity
}

// String provides a textual rendering of the entity
func (e *Entity) String() string {
	return fmt.Sprintf(
		"Kind: %v\n"+
			"User: %v\n"+
			"STS caller identity: %v\n"+
			"EKS roles: %v\n"+
			"EKS policies: %v\n"+
			"Kube context: %+v",
		e.Kind,
		e.User,
		e.Caller,
		e.Roles,
		e.Policies,
		e.KubeConfig.CurrentContext,
	)
}
