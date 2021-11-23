package matchersimpl

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMatchGRPCStatusCode(t *testing.T) {
	wantCode := codes.Internal
	matcher := MatchGRPCStatusCode(wantCode)
	assert.Equal(t, wantCode, matcher.expectedCode)
}

func TestGRPCErrorCodeMatcher_Match(t *testing.T) {
	tests := []struct {
		name         string
		givenError   interface{}
		expectedCode codes.Code
		wantSuccess  bool
		wantErr      error
	}{
		{"should match", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.InvalidArgument, true, nil},
		{"should not match", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.NotFound, false, nil},
		{"should fail when no error given", false, codes.OK, false, errExpectedError},
		{"should fail when no status given", errors.New("non status error"), codes.OK, false, errStatusErrExpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &GRPCErrorCodeMatcher{
				expectedCode: tt.expectedCode,
			}
			gotSuccess, err := matcher.Match(tt.givenError)
			assert.Equal(t, tt.wantSuccess, gotSuccess)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGRPCErrorCodeMatcher_FailureMessage(t *testing.T) {
	tests := []struct {
		name         string
		givenError   interface{}
		expectedCode codes.Code
		wantMessage  string
	}{
		{"should fail when no error given", false, codes.OK, "is not an error"},
		{"should fail when no status given", errors.New("non status error"), codes.OK, "is not a grpc error"},
		{"should not match", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.NotFound, "to match"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &GRPCErrorCodeMatcher{
				expectedCode: tt.expectedCode,
			}
			gotMessage := matcher.FailureMessage(tt.givenError)
			assert.Contains(t, gotMessage, tt.wantMessage)
		})
	}
}

func TestGRPCErrorCodeMatcher_NegatedFailureMessage(t *testing.T) {
	tests := []struct {
		name         string
		givenError   interface{}
		expectedCode codes.Code
		wantMessage  string
	}{
		{"should fail when no error given", false, codes.OK, "is not an error"},
		{"should fail when no status given", errors.New("non status error"), codes.OK, "is not a grpc error"},
		{"should not match", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.NotFound, "not to match"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &GRPCErrorCodeMatcher{
				expectedCode: tt.expectedCode,
			}
			gotMessage := matcher.NegatedFailureMessage(tt.givenError)
			assert.Contains(t, gotMessage, tt.wantMessage)
		})
	}
}
