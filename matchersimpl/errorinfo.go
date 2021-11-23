package matchersimpl

import (
	"errors"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

var (
	errErrorInfoNotFound = errors.New("no errdetails.ErrorInfo found in the details")
)

// ErrorInfoMatcher abstracts a matcher specialized on errdetails.ErrorInfo types.
type ErrorInfoMatcher interface {
	Match(errInfo *errdetails.ErrorInfo) (bool, error)
	FailureMessage(errInfo *errdetails.ErrorInfo) string
	NegatedFailureMessage(errInfo *errdetails.ErrorInfo) string
}

func NewGRPCMatchErrorInfo(matcher ErrorInfoMatcher) *GRPCStatusMatcher {
	return NewGRPCStatusMatcher(&GRPCErrorInfoMatcher{matcher})
}

// GRPCErrorInfoMatcher implements a StatusMatcher that finds a errdetails.ErrorInfo instance on the given status.Status.
//
// This is a helper for dealing with ErrorInfo matcher.
type GRPCErrorInfoMatcher struct {
	errorInfoMatcher ErrorInfoMatcher
}

func (m *GRPCErrorInfoMatcher) Match(st *status.Status) (bool, error) {
	errInfo, ok := findErrorInfo(st)
	if !ok {
		return false, errErrorInfoNotFound
	}
	return m.errorInfoMatcher.Match(errInfo)
}

func (m *GRPCErrorInfoMatcher) FailureMessage(st *status.Status) string {
	errInfo, ok := findErrorInfo(st)
	if !ok {
		return format.Message(st, "does not have an *errdetails.ErrorInfo in the details")
	}
	return m.errorInfoMatcher.FailureMessage(errInfo)
}

func (m *GRPCErrorInfoMatcher) NegatedFailureMessage(st *status.Status) string {
	errInfo, ok := findErrorInfo(st)
	if !ok {
		return format.Message(st, "does not have an *errdetails.ErrorInfo in the details")
	}
	return m.errorInfoMatcher.NegatedFailureMessage(errInfo)
}

type GRPCErrorInfoReasonMatcher struct {
	PropMap func(errInfo *errdetails.ErrorInfo) interface{}
	Matcher types.GomegaMatcher
}

func (m *GRPCErrorInfoReasonMatcher) Match(errInfo *errdetails.ErrorInfo) (bool, error) {
	return m.Matcher.Match(m.PropMap(errInfo))
}

func (m *GRPCErrorInfoReasonMatcher) FailureMessage(errInfo *errdetails.ErrorInfo) string {
	return m.Matcher.FailureMessage(m.PropMap(errInfo))
}

func (m *GRPCErrorInfoReasonMatcher) NegatedFailureMessage(errInfo *errdetails.ErrorInfo) string {
	return m.Matcher.NegatedFailureMessage(m.PropMap(errInfo))
}

// findErrorInfo walks through the error details of the given st trying to find a errdetails.ErrorInfo instance. If it
// find any, returns the instance and true.
//
// Otherwise, it returns false.
func findErrorInfo(st *status.Status) (*errdetails.ErrorInfo, bool) {
	for _, d := range st.Details() {
		errInfo, ok := d.(*errdetails.ErrorInfo)
		if ok {
			return errInfo, true
		}
	}
	return nil, false
}
