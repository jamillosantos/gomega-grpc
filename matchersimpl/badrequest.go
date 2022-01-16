package matchersimpl

import (
	"errors"
	"fmt"

	"github.com/onsi/gomega/format"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

var (
	errBadRequestNotFound = errors.New("no errdetails.BadRequest found in the details")
)

// BadRequestMatcher abstracts a matcher specialized on errdetails.BadRequest types.
type BadRequestMatcher interface {
	Match(actual *errdetails.BadRequest) (bool, error)
	FailureMessage(actual *errdetails.BadRequest) string
	NegatedFailureMessage(actual *errdetails.BadRequest) string
}

func NewGRPCMatchBadRequest(matcher BadRequestMatcher) *GRPCStatusMatcher {
	return NewGRPCStatusMatcher(&GRPCBadRequestMatcher{matcher})
}

// GRPCBadRequestMatcher implements a StatusMatcher that finds a errdetails.BadRequest instance on the given status.Status.
//
// This is a helper for dealing with BadRequest matcher.
type GRPCBadRequestMatcher struct {
	badRequestMatcher BadRequestMatcher
}

func (m *GRPCBadRequestMatcher) Match(st *status.Status) (bool, error) {
	errInfo, ok := findBadRequest(st)
	if !ok {
		return false, errBadRequestNotFound
	}
	return m.badRequestMatcher.Match(errInfo)
}

func (m *GRPCBadRequestMatcher) FailureMessage(st *status.Status) string {
	errInfo, ok := findBadRequest(st)
	if !ok {
		return format.Message(st, "does not have an *errdetails.BadRequest in the details")
	}
	return m.badRequestMatcher.FailureMessage(errInfo)
}

func (m *GRPCBadRequestMatcher) NegatedFailureMessage(st *status.Status) string {
	errInfo, ok := findBadRequest(st)
	if !ok {
		return format.Message(st, "does not have an *errdetails.BadRequest in the details")
	}
	return m.badRequestMatcher.NegatedFailureMessage(errInfo)
}

// GRPCBadRequestFieldViolation tries matching the matching the Field of a given errdetails.BadRequest_FieldViolation.
// If the given Description is not empty, this will only match if both Field and Description match.
type GRPCBadRequestFieldViolation struct {
	Field       string
	Description string
}

func (m *GRPCBadRequestFieldViolation) Match(actual *errdetails.BadRequest) (bool, error) {
	for _, fv := range actual.GetFieldViolations() {
		if m.matchesFV(fv) {
			return true, nil
		}
	}
	return false, nil
}

func (m *GRPCBadRequestFieldViolation) matchesFV(fv *errdetails.BadRequest_FieldViolation) bool {
	return fv.GetField() == m.Field && (m.Description == "" || fv.GetDescription() == m.Description)
}

func (m *GRPCBadRequestFieldViolation) FailureMessage(actual *errdetails.BadRequest) string {
	return format.Message(actual, "does not have", fmt.Sprintf("[%s: %s]", m.Field, m.Description))
}

func (m *GRPCBadRequestFieldViolation) NegatedFailureMessage(actual *errdetails.BadRequest) string {
	return format.Message(actual, "have", fmt.Sprintf("[%s: %s]", m.Field, m.Description))
}

// findBadRequest walks through the error details of the given st trying to find a errdetails.BadRequest instance. If it
// find any, returns the instance and true.
//
// Otherwise, it returns false.
func findBadRequest(st *status.Status) (*errdetails.BadRequest, bool) {
	for _, d := range st.Details() {
		errInfo, ok := d.(*errdetails.BadRequest)
		if ok {
			return errInfo, true
		}
	}
	return nil, false
}
