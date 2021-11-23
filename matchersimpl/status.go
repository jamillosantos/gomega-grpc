package matchersimpl

import (
	"errors"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"google.golang.org/grpc/status"
)

var (
	errExpectedError     = errors.New("the given object is not an error")
	errStatusErrExpected = errors.New("the given error is not a GRPC error")
)

// StatusMatcher abstracts a matcher specialized on status.Status types.
type StatusMatcher interface {
	Match(st *status.Status) (bool, error)
	FailureMessage(st *status.Status) string
	NegatedFailureMessage(st *status.Status) string
}

// GRPCStatusMatcher implements a gomega.Matcher that makes creating status.Status checks easier by abstracting all the
// validation when receiving a Match call. Then, it uses the given functions (matchFunc, failureMessageFunc and
// negatedfailureMessageFunc) to extend its behaviour.
//
// For reference, you can check MatchErrorInfoReason.
type GRPCStatusMatcher struct {
	statusMatcher StatusMatcher
}

// NewGRPCStatusMatcher is the constructor for the GRPCStatusMatcher.
func NewGRPCStatusMatcher(statusMatcher StatusMatcher) *GRPCStatusMatcher {
	return &GRPCStatusMatcher{statusMatcher}
}

// Match validates if the given actual is a status.Status, if so it will call the matchFunc.
func (matcher *GRPCStatusMatcher) Match(actual interface{}) (success bool, err error) {
	actualErr, ok := actual.(error)
	if !ok {
		return false, errExpectedError
	}
	st, ok := status.FromError(actualErr)
	if !ok {
		return false, errStatusErrExpected
	}
	return matcher.statusMatcher.Match(st)
}

func (matcher *GRPCStatusMatcher) validateActual(actual interface{}) (string, *status.Status, bool) {
	actualErr, ok := actual.(error)
	if !ok {
		return format.Message(actual, "is not an error"), nil, false
	}
	st, ok := status.FromError(actualErr)
	if !ok {
		return format.Message(actual, "is not a grpc error"), nil, false
	}
	return "", st, true
}

func (matcher *GRPCStatusMatcher) FailureMessage(actual interface{}) (message string) {
	m, st, ok := matcher.validateActual(actual)
	if !ok {
		return m
	}
	return matcher.statusMatcher.FailureMessage(st)
}

func (matcher *GRPCStatusMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	m, st, ok := matcher.validateActual(actual)
	if !ok {
		return m
	}
	return matcher.statusMatcher.NegatedFailureMessage(st)
}

type GRPCStatusPropMatcher struct {
	PropMap func(status *status.Status) interface{}
	Matcher gomega.OmegaMatcher
}

func (matcher *GRPCStatusPropMatcher) Match(st *status.Status) (bool, error) {
	return matcher.Matcher.Match(matcher.PropMap(st))
}

func (matcher *GRPCStatusPropMatcher) FailureMessage(st *status.Status) string {
	return matcher.Matcher.FailureMessage(matcher.PropMap(st))
}

func (matcher *GRPCStatusPropMatcher) NegatedFailureMessage(st *status.Status) string {
	return matcher.Matcher.NegatedFailureMessage(matcher.PropMap(st))
}
