package client

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	URLTargetForClientBundle = "api/clientbundle"

	filenameCAPem      = "ca.pem"
	filenameCertPem    = "cert.pem"
	filenamePrivKeyPem = "key.pem"
	filenamePubKeyPem  = "cert.pub"
	filenameKubeconfig = "kube.yml"
)

var (
	ErrFailedToRetrieveClientBundle = errors.New("failed to retrieve the client bundle from MKE")
)

// ApiClientBundle retrieve a client bundle(
func (c *Client) ApiClientBundleCreate(ctx context.Context) (ClientBundle, error) {
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

	cb.ID = zipReader.Comment

	errs := []error{}

	for _, f := range zipReader.File {
		switch f.Name {
		case filenameCAPem:
			fReader, _ := f.Open()
			capem, err := ClientBundleRetrieveValue(fReader)
			fReader.Close()

			if err != nil {
				errs = append(errs, err)
			} else {
				cb.CACert = capem
			}
		case filenameCertPem:
			fReader, _ := f.Open()
			cert, err := ClientBundleRetrieveValue(fReader)
			fReader.Close()

			if err != nil {
				errs = append(errs, err)
			} else {
				cb.Certs = append(cb.Certs, cert)
			}
		case filenamePrivKeyPem:
			fReader, _ := f.Open()
			capem, err := ClientBundleRetrieveValue(fReader)
			fReader.Close()

			if err != nil {
				errs = append(errs, err)
			} else {
				cb.PrivateKey = capem
			}
		case filenamePubKeyPem:
			fReader, _ := f.Open()
			capem, err := ClientBundleRetrieveValue(fReader)
			fReader.Close()

			if err != nil {
				errs = append(errs, err)
			} else {
				cb.PublicKey = capem
			}
		case filenameKubeconfig:
			fReader, _ := f.Open()
			kube, err := NewClientBundleKubeFromKubeYml(fReader)
			fReader.Close()

			if err != nil {
				errs = append(errs, err)
			} else {
				cb.Kube = &kube
			}

		}
	}

	if len(errs) > 0 {
		errString := ""

		for _, err := range errs {
			errString = fmt.Sprintf("%s, %s", errString, err)
		}

		return cb, fmt.Errorf("%w; %s", ErrFailedToRetrieveClientBundle, errString)
	}

	return cb, nil
}
