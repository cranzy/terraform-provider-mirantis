package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/msr/client"
)

type testStruct struct {
	server           *httptest.Server
	expectedResponse client.HealthResponse
	expectedErr      error
}

func TestMSRClientHealthy(t *testing.T) {
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"error": "", "healthy":true}`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{
			Error:   "",
			Healthy: true,
		},
		expectedErr: nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create test client")
	}
	ctx := context.Background()
	resp, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(tc.expectedResponse.Healthy, resp) {
		t.Errorf("expected (%v), got (%v)", tc.expectedResponse.Healthy, resp)
	}
	if err != tc.expectedErr {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestMSRClientUnhealthy(t *testing.T) {
	unhealthyRes := client.HealthResponse{Healthy: false}
	bodyRes, _ := json.Marshal(unhealthyRes)
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(bodyRes); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: unhealthyRes,
		expectedErr:      nil,
	}
	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "fakepass")
	if err != nil {
		t.Error("couldn't create client new client")
	}
	ctx := context.Background()
	isHealthy, err := testClient.IsHealthy(ctx)
	if !reflect.DeepEqual(unhealthyRes.Healthy, isHealthy) {
		t.Errorf("expected (%v), got (%v)", unhealthyRes.Healthy, isHealthy)
	}
	if err != tc.expectedErr {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestMSRClientBadrequest(t *testing.T) {
	resError := client.ResponseError{
		Errors: []client.Errors{
			{
				Code:    strconv.Itoa(http.StatusBadRequest),
				Message: "Bad request",
			},
		},
	}
	bodyRes, err := json.Marshal(resError)
	if err != nil {
		t.Errorf("couldn't marshal struct %+v", resError)
		return
	}
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(bodyRes); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{
			Healthy: false,
		},
		expectedErr: errors.New("MSR API Error: response status is: 400. Bad request"),
	}

	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "fakepass")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if err.Error() != tc.expectedErr.Error() {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestMSRClientUnauthorized(t *testing.T) {
	resError := client.ResponseError{
		Errors: []client.Errors{
			{
				Code:    strconv.Itoa(http.StatusUnauthorized),
				Message: "Bad creds",
			},
		},
	}
	bodyRes, err := json.Marshal(resError)
	if err != nil {
		t.Errorf("couldn't marshal struct %+v", resError)
		return
	}
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write(bodyRes); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      errors.New("MSR API Error: response status: 401. Unauthorized request"),
	}

	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "fakepass")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if err.Error() != tc.expectedErr.Error() {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestDoRequestWrongErrorStruct(t *testing.T) {
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{"msg": "wrong struct"}`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      errors.New("MSR API Error: response status is: 400. Wrong unmarshal struct for {\"msg\": \"wrong struct\"}"),
	}

	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "fakepass")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if err.Error() != tc.expectedErr.Error() {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestDoRequestWrongErrorStructField(t *testing.T) {
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{"errors":[{"code":true, "message":"lol"}]}`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      errors.New("MSR API Error: response status: 400. json: cannot unmarshal bool into Go struct field Errors.errors.code of type string"),
	}

	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "fakepass")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if err.Error() != tc.expectedErr.Error() {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestEmptyUsernameField(t *testing.T) {
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{ "msg": "ok"`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      errors.New("no username or password provided"),
	}

	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "", "fakepass")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if err.Error() != tc.expectedErr.Error() {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}

func TestEmptyPasswordField(t *testing.T) {
	tc := testStruct{
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte(`{ "msg": "ok"`)); err != nil {
				t.Error(err)
				return
			}
		})),
		expectedResponse: client.HealthResponse{},
		expectedErr:      errors.New("no username or password provided"),
	}

	defer tc.server.Close()
	testClient, err := client.NewClient(tc.server.URL, "fakeuser", "")
	if !reflect.DeepEqual(testClient, client.Client{}) {
		t.Errorf("expected (%v), got (%v)", client.Client{}, testClient)
	}
	if err.Error() != tc.expectedErr.Error() {
		t.Errorf("expected (%v), got (%v)", tc.expectedErr, err)
	}
}
