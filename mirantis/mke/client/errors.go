package client

import "errors"

var (
	ErrEmptyUsernamePass = errors.New("No username or password provided in MKE client")
	ErrEmptyEndpoint     = errors.New("No endpoint provided in MKE client")
	ErrRequestCreation   = errors.New("Error creating request in MKE client")
	ErrMarshaling        = errors.New("Error occured while marshalling struct in MKE client")
	ErrUnmarshaling      = errors.New("Error occured while unmarshalling struct in MKE client")
	ErrEmptyResError     = errors.New("Request returned empty ResponseError struct in MKE client")
	ErrResponseError     = errors.New("Request returned ResponseError in MKE client")
	ErrUnauthorizedReq   = errors.New("Unauthorized request in MKE client")
	ErrUnknownTarget     = errors.New("Unknown API target")
	ErrServerError       = errors.New("Server error occured")
	ErrEmptyStruct       = errors.New("Empty struct passed in MKE client")
	ErrInvalidFilter     = errors.New("Passing invalid account retrieval filter in MKE client")
)
