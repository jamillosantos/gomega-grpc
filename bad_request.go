package grpcmatchers

import (
	"github.com/jamillosantos/gomega-grpc/matchersimpl"
)

// HaveFieldValidation matches a given field and description (if informed) for a given errdetails.BadRequest or
// errdetails.BadRequest_FieldViolation[] or errdetails.BadRequest_FieldViolation.
func HaveFieldViolation(fieldAndDescription ...string) *matchersimpl.GRPCStatusMatcher {
	var field, description string
	if len(fieldAndDescription) > 0 {
		field = fieldAndDescription[0]
	}
	if len(fieldAndDescription) > 1 {
		description = fieldAndDescription[1]
	}
	return matchersimpl.NewGRPCMatchBadRequest(&matchersimpl.GRPCBadRequestFieldViolation{
		Field:       field,
		Description: description,
	})
}
