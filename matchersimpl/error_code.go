package matchersimpl

import (
	"github.com/onsi/gomega/format"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCErrorCodeMatcher struct {
	expectedCode codes.Code
}

// MatchGRPCStatusCode expects an error to match the given code.
func MatchGRPCStatusCode(code codes.Code) *GRPCErrorCodeMatcher {
	return &GRPCErrorCodeMatcher{
		expectedCode: code,
	}
}

// Match checks if the given actual is an error and a status.Status. If so, it tries to match the actual status code
// with the expected one.
func (matcher *GRPCErrorCodeMatcher) Match(actual interface{}) (success bool, err error) {
	actualErr, ok := actual.(error)
	if !ok {
		return false, errExpectedError
	}
	st, ok := status.FromError(actualErr)
	if !ok {
		return false, errStatusErrExpected
	}

	return st.Code() == matcher.expectedCode, nil
}

// FailureMessage returns the error messages when this matcher does not receive an error, neither a status.Status or
// the status code does not match.
func (matcher *GRPCErrorCodeMatcher) FailureMessage(actual interface{}) (message string) {
	actualErr, ok := actual.(error)
	if !ok {
		return format.Message(actual, "is not an error")
	}
	st, ok := status.FromError(actualErr)
	if !ok {
		return format.Message(actual, "is not a grpc error")
	}

	return format.Message(st.Code(), "to match", matcher.expectedCode)
}

// NegatedFailureMessage returns the error messages when this matcher does not receive an error, neither a status.Status
// or the status code does not match negating the sentence.
func (matcher *GRPCErrorCodeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualErr, ok := actual.(error)
	if !ok {
		return format.Message(actual, "is not an error")
	}
	st, ok := status.FromError(actualErr)
	if !ok {
		return format.Message(actual, "is not a grpc error")
	}

	return format.Message(actual, "not to match", st.Code())
}
