package client

import (
	"encoding/base64"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ClientBundle interpretation of the ClientBundle data in memory
type ClientBundle struct {
	ID         string            `json:"id"`
	PrivateKey string            `json:"private_key"`
	PublicKey  string            `json:"public_key"`
	Cert       string            `json:"cert"`
	CACert     string            `json:"ca_cert"`
	Kube       *ClientBundleKube `json:"kube"`
}

// ClientBundleKube Kubernetes parts of the client bundle
// primarily we are focused on satisfying requirements for a kubernetes provider
// such as https://github.com/hashicorp/terraform-provider-kubernetes/blob/main/kubernetes/provider.go
type ClientBundleKube struct {
	Config            string `json:"config"`
	Host              string `json:"host"`
	ClientKey         string `json:"client_key"`
	ClientCertificate string `json:"client_certificate"`
	CACertificate     string `json:"cluster_ca_certificate"`
	Insecure          string `json:"insecure"`
}

// ClientBundleRetrieveValue read a value
func ClientBundleRetrieveValue(val io.Reader) (string, error) {
	allBytes, err := ioutil.ReadAll(val)
	if err != nil {
		return "", err
	}
	return string(allBytes), nil
}

// ClientBundleDecodeBase64Value read a value and base64 decode it
func ClientBundleDecodeBase64Value(val io.Reader) (string, error) {
	allBytes, err := ioutil.ReadAll(val)
	if err != nil {
		return "", err
	}
	return helperStringBase64Decode(string(allBytes)), nil
}

// NewClientBundleKubeFromKubeYml ClientBundleKube constructor from byte list of a kubeconfig file
func NewClientBundleKubeFromKubeYml(val io.Reader) (ClientBundleKube, error) {
	k8bytes, _ := ioutil.ReadAll(val)

	var cbk ClientBundleKube

	// Struct representation of a kube config file.
	// see https://zhwt.github.io/yaml-to-go/
	var cbkHolder struct {
		APIVersion  string            `yaml:"apiVersion"`
		Kind        string            `yaml:"kind"`
		Preferences map[string]string `yaml:"preferences"`
		Clusters    []struct {
			Name    string `yaml:"name"`
			Cluster struct {
				CertificateAuthorityData string `yaml:"certificate-authority-data"`
				Server                   string `yaml:"server"`
			} `yaml:"cluster"`
		} `yaml:"clusters"`
		Contexts []struct {
			Name    string `yaml:"name"`
			Context struct {
				Cluster string `yaml:"cluster"`
				User    string `yaml:"user"`
			} `yaml:"context"`
		} `yaml:"contexts"`
		CurrentContext string `yaml:"current-context"`
		Users          []struct {
			Name string `yaml:"name"`
			User struct {
				ClientCertificateData string `yaml:"client-certificate-data"`
				ClientKeyData         string `yaml:"client-key-data"`
			} `yaml:"user"`
		} `yaml:"users"`
	}

	cbk.Config = string(k8bytes)

	if err := yaml.UnmarshalStrict(k8bytes, &cbkHolder); err != nil {
		return cbk, err
	}

	var contextName, clusterName, userName string

	contextName = cbkHolder.CurrentContext

	for _, context := range cbkHolder.Contexts {
		if context.Name == contextName {
			clusterName = context.Context.Cluster
			userName = context.Context.User
			break
		}
	}

	for _, cluster := range cbkHolder.Clusters {
		if cluster.Name == clusterName {
			cbk.Host = cluster.Cluster.Server
			cbk.CACertificate = helperStringBase64Decode(cluster.Cluster.CertificateAuthorityData)
			break
		}
	}

	for _, user := range cbkHolder.Users {
		if user.Name == userName {
			cbk.ClientKey = helperStringBase64Decode(user.User.ClientKeyData)
			cbk.ClientCertificate = helperStringBase64Decode(user.User.ClientCertificateData)
			break
		}
	}

	return cbk, nil
}

// this decodes some strings in the file that are base64 encoded
func helperStringBase64Decode(val string) string {
	valDecodedBytes, _ := base64.StdEncoding.DecodeString(val)
	return string(valDecodedBytes)
}
