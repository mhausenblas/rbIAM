package main

// ServiceAccountList is based on List from:
// https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/types.go
type ServiceAccountList struct {
	Items []ServiceAccount `json:"items"`
}

// ObjectMeta is metadata that all persisted resources must have.
type ObjectMeta struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	ClusterName string            `json:"clusterName,omitempty"`
}

// ServiceAccount is a lean v.1 based on:
// https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go
type ServiceAccount struct {
	ObjectMeta                   `json:"metadata,omitempty"`
	Secrets                      []ObjectReference      `json:"secrets,omitempty"`
	ImagePullSecrets             []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	AutomountServiceAccountToken *bool                  `json:"automountServiceAccountToken,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	Kind      string `json:"kind,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	FieldPath string `json:"fieldPath,omitempty"`
}

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
type LocalObjectReference struct {
	Name string `json:"name,omitempty"`
}

// Config is a lean v1.Config based on:
// https://github.com/kubernetes/client-go/blob/master/tools/clientcmd/api/v1/types.go
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
