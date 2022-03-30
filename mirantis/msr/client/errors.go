package client

import "errors"

var (
	ErrEmptyUsernamePass = errors.New("error - no username or password provided in MSR client")
	ErrRequestCreation   = errors.New("error creating request in MSR client")
	ErrMarshaling        = errors.New("error occured while marshalling struct in MSR client")
	ErrUnmarshaling      = errors.New("error occured while unmarshalling struct in MSR client")
	ErrEmptyResError     = errors.New("error - request returned empty ResponseError struct in MSR client")
	ErrResponseError     = errors.New("error - request returned ResponseError in MSR client")
	ErrUnauthorizedReq   = errors.New("unauthorized request in MSR client")
	ErrEmptyStruct       = errors.New("error - empty struct passed in MSR client")
	ErrInvalidFilter     = errors.New("error - passing invalid account retrieval filter in MSR client")
)
