package client

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	URLTargetForClientBundle = "api/clientbundle"

	filenameCAPem      = "ca.pem"
	filenameCertPem    = "cert.pem"
	filenameKeyPem     = "key.pem"
	filenameCertPub    = "cert.pub"
	filenameKubeconfig = "kube.yml"
)

var (
	ErrFailedToRetrieveClientBundle = errors.New("failed to retrieve the client bundle from MKE")
)

// ApiClientBundle retrieve a client bundle(
func (c *Client) ApiClientBundle(ctx context.Context) (ClientBundle, error) {
	var cb ClientBundle

	req, err := c.RequestFromTargetAndBytesBody(ctx, http.MethodPost, URLTargetForClientBundle, []byte{})
	if err != nil {
		return cb, err
	}

	resp, err := c.doAuthorizedRequest(req)
	if err != nil {
		return cb, err
	}
	defer resp.Body.Close()

	zipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cb, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), resp.ContentLength)
	if err != nil {
		return cb, err
	}

	for _, f := range zipReader.File {
		switch f.Name {
		case filenameKubeconfig:
			fReader, _ := f.Open()
			k8bytes, _ := ioutil.ReadAll(fReader)
			fReader.Close()

			kube, err := NewClientBundleKubeFromKubeYml(k8bytes)
			if err != nil {
				return cb, err
			}
			cb.Kube = &kube
		}
	}

	return cb, nil
}
