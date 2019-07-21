package main

////////////////////////////////////////////////////////////////////////////////
// https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/types.go

// ObjectMeta is metadata that all persisted resources must have.
type ObjectMeta struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	ClusterName string            `json:"clusterName,omitempty"`
}

// ServiceAccountList is a list of service accounts.
type ServiceAccountList struct {
	Items []ServiceAccount `json:"items"`
}

// SecretList is a list of secrets.
type SecretList struct {
	Items []Secret `json:"items"`
}

// PodList is a list of pods.
type PodList struct {
	Items []Pod `json:"items"`
}

////////////////////////////////////////////////////////////////////////////////
// https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

// ObjectReference for inspecting the referred object.
type ObjectReference struct {
	Kind      string `json:"kind,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	FieldPath string `json:"fieldPath,omitempty"`
}

// LocalObjectReference for locating the referenced object inside a namespace.
type LocalObjectReference struct {
	Name string `json:"name,omitempty"`
}

// ServiceAccount represents a Kubernetes service account.
type ServiceAccount struct {
	ObjectMeta                   `json:"metadata,omitempty"`
	Secrets                      []ObjectReference      `json:"secrets,omitempty"`
	ImagePullSecrets             []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	AutomountServiceAccountToken *bool                  `json:"automountServiceAccountToken,omitempty"`
}

// Secret represents a Kubernetes secrect holding sensitive data.
type Secret struct {
	ObjectMeta `json:"metadata,omitempty"`
	Data       map[string][]byte `json:"data,omitempty" `
	StringData map[string]string `json:"stringData,omitempty"`
	Type       SecretType        `json:"type,omitempty"`
}

// SecretType is a custom data type facilitating programmatic handling of
// secret data.
type SecretType string

// Pod is a collection of containers that can run on a host.
type Pod struct {
	ObjectMeta `json:"metadata,omitempty"`
	Spec       PodSpec   `json:"spec,omitempty"`
	Status     PodStatus `json:"status,omitempty"`
}

// PodSpec is a description of a pod.
type PodSpec struct {
	Volumes            []Volume               `json:"volumes,omitempty"`
	InitContainers     []Container            `json:"initContainers,omitempty"`
	Containers         []Container            `json:"containers"`
	ServiceAccountName string                 `json:"serviceAccountName,omitempty"`
	ImagePullSecrets   []LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// Volume represents a named volume in a pod.
type Volume struct {
	Name string `json:"name"`
}

// Container is an application container running within a pod.
type Container struct {
	Name            string     `json:"name"`
	Image           string     `json:"image,omitempty"`
	Command         []string   `json:"command,omitempty"`
	Args            []string   `json:"args,omitempty"`
	Env             []EnvVar   `json:"env,omitempty"`
	ImagePullPolicy PullPolicy `json:"imagePullPolicy,omitempty"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	Name  string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
}

// PullPolicy describes a policy for if/when to pull a container image.
type PullPolicy string

// PodStatus represents information about the status of a pod.
type PodStatus struct {
	Phase             PodPhase       `json:"phase,omitempty"`
	Conditions        []PodCondition `json:"conditions,omitempty"`
	Message           string         `json:"message,omitempty"`
	NominatedNodeName string         `json:"nominatedNodeName,omitempty"`
	HostIP            string         `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
	PodIP             string         `json:"podIP,omitempty" protobuf:"bytes,6,opt,name=podIP"`
}

// PodPhase is a label for the condition of a pod at the current time.
type PodPhase string

// PodCondition contains details for the current condition of this pod.
type PodCondition struct {
	Type   PodConditionType `json:"type"`
	Status ConditionStatus  `json:"status"`
	Reason string           `json:"reason,omitempty"`
}

// PodConditionType is a valid value for PodCondition.Type
type PodConditionType string

// ConditionStatus provides the status.
type ConditionStatus string

////////////////////////////////////////////////////////////////////////////////
// https://github.com/kubernetes/client-go/blob/master/tools/clientcmd/api/v1/types.go

// Config is a lean v1.Config variant, representing a Kubernetes client config.
type Config struct {
	// Clusters is a map of referencable names to cluster configs
	Clusters []NamedCluster `json:"clusters"`
	// AuthInfos is a map of referencable names to user configs
	AuthInfos []NamedAuthInfo `json:"users"`
	// Contexts is a map of referencable names to context configs
	Contexts []NamedContext `json:"contexts"`
	// CurrentContext is the name of the context that you would like to use by default
	CurrentContext string `json:"current-context"`
}

// Cluster contains information about how to communicate with a cluster
type Cluster struct {
	// Server is the address of the kubernetes cluster (https://hostname:port)
	Server string `json:"server"`
	// InsecureSkipTLSVerify skips the validity check for the server's certificate
	InsecureSkipTLSVerify bool `json:"insecure-skip-tls-verify,omitempty"`
	// CertificateAuthority is the path to a cert file for the certificate authority
	CertificateAuthority string `json:"certificate-authority,omitempty"`
	// CertificateAuthorityData contains PEM-encoded certificate authority certificates
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
}

// AuthInfo contains information that describes identity information
type AuthInfo struct {
	// ClientCertificate is the path to a client cert file for TLS
	ClientCertificate string `json:"client-certificate,omitempty"`
	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS
	ClientCertificateData []byte `json:"client-certificate-data,omitempty"`
	// ClientKey is the path to a client key file for TLS.
	ClientKey string `json:"client-key,omitempty"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS
	ClientKeyData []byte `json:"client-key-data,omitempty"`
	// Token is the bearer token for authentication
	Token string `json:"token,omitempty"`
	// TokenFile is a pointer to a file that contains a bearer token.
	// If both Token and TokenFile are present, Token takes precedence.
	TokenFile string `json:"tokenFile,omitempty"`
	// Impersonate is the username to impersonate. The name matches the flag.
	Impersonate string `json:"as,omitempty"`
	// ImpersonateGroups is the groups to impersonate
	ImpersonateGroups []string `json:"as-groups,omitempty"`
	// ImpersonateUserExtra contains additional information for impersonated user
	ImpersonateUserExtra map[string][]string `json:"as-user-extra,omitempty"`
	// Username is the username for basic authentication
	Username string `json:"username,omitempty"`
	// Password is the password for basic authentication
	Password string `json:"password,omitempty"`
	// AuthProvider specifies a custom authentication plugin
	AuthProvider *AuthProviderConfig `json:"auth-provider,omitempty"`
	// Exec specifies a custom exec-based authentication plugin
	Exec *ExecConfig `json:"exec,omitempty"`
}

// Context is a tuple of references to a cluster:
// - how do I communicate with a kubernetes cluster
// - a user, that is, how do I identify myself
// - namespace, that is, what subset of resources do I want to work with
type Context struct {
	// Cluster is the name of the cluster for this context
	Cluster string `json:"cluster"`
	// AuthInfo is the name of the authInfo for this context
	AuthInfo string `json:"user"`
	// Namespace is the default namespace to use on unspecified requests
	Namespace string `json:"namespace,omitempty"`
}

// NamedCluster relates nicknames to cluster information
type NamedCluster struct {
	// Name is the nickname for this Cluster
	Name string `json:"name"`
	// Cluster holds the cluster information
	Cluster Cluster `json:"cluster"`
}

// NamedContext relates nicknames to context information
type NamedContext struct {
	// Name is the nickname for this Context
	Name string `json:"name"`
	// Context holds the context information
	Context Context `json:"context"`
}

// NamedAuthInfo relates nicknames to auth information
type NamedAuthInfo struct {
	// Name is the nickname for this AuthInfo
	Name string `json:"name"`
	// AuthInfo holds the auth information
	AuthInfo AuthInfo `json:"user"`
}

// AuthProviderConfig holds the configuration for a specified auth provider
type AuthProviderConfig struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

// ExecConfig specifies a command to provide client credentials
type ExecConfig struct {
	// Command to execute.
	Command string `json:"command"`
	// Arguments to pass to the command when executing it
	Args []string `json:"args"`
	// Env defines additional environment variables to expose to the process
	Env []ExecEnvVar `json:"env"`
	// Preferred input version of the ExecInfo
	APIVersion string `json:"apiVersion,omitempty"`
}

// ExecEnvVar is used for setting environment variables when executing an exec-based
// credential plugin.
type ExecEnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
