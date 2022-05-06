package client

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	ErrFailedToRetrieveClientBundle         = errors.New("failed to retrieve the client bundle from MKE")
	ErrFailedToFindClientBundleMKEPublicKey = errors.New("no MKE Public key was found that matches the client bundle")
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
				cb.Cert = cert
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

// ApiClientBundleGetPublicKey retrieve a client bundle by finding the matching public key
// There isn't really a great way of doing this.
func (c *Client) ApiClientBundleGetPublicKey(ctx context.Context, cb ClientBundle) (AccountPublicKey, error) {
	var k AccountPublicKey

	account := c.Username()

	keys, err := c.ApiPublicKeyList(ctx, account)
	if err != nil {
		return k, err
	}

	foundKeys := []string{}
	cbpk := strings.TrimSpace(cb.PublicKey)
	for _, key := range keys {
		pk := strings.TrimSpace(key.PublicKey)
		if pk == cbpk {
			return key, nil
		}
		foundKeys = append(foundKeys, pk)
	}

	return k, fmt.Errorf("%w; Could not match key: \n%s\n in \n%s", ErrFailedToFindClientBundleMKEPublicKey, cb.PublicKey, strings.Join(foundKeys, "\n"))
}

// ApiClientBundleDelete delete a client bundle by finding and deleting the matching public key
// There isn't really a great way of doing this.
func (c *Client) ApiClientBundleDelete(ctx context.Context, cb ClientBundle) error {
	account := c.Username()

	key, err := c.ApiClientBundleGetPublicKey(ctx, cb)
	if err != nil {
		return err
	}

	return c.ApiPublicKeyDelete(ctx, account, key.ID)
}
