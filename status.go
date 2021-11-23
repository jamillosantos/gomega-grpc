package grpcmatchers

import (
	"github.com/onsi/gomega/types"
	"google.golang.org/grpc/status"

	"github.com/jamillosantos/gomega-grpc/matchersimpl"
)

// HaveStatusCode will match the *status.Status Code property against the given matcher.
func HaveStatusCode(matcher types.GomegaMatcher) *matchersimpl.GRPCStatusMatcher {
	return matchersimpl.NewGRPCStatusMatcher(&matchersimpl.GRPCStatusPropMatcher{
		PropMap: func(st *status.Status) interface{} {
			return st.Code()
		},
		Matcher: matcher,
	})
}

// HaveStatusMessage will match the *status.Status Message property against the given matcher.
func HaveStatusMessage(matcher types.GomegaMatcher) *matchersimpl.GRPCStatusMatcher {
	return matchersimpl.NewGRPCStatusMatcher(&matchersimpl.GRPCStatusPropMatcher{
		PropMap: func(st *status.Status) interface{} {
			return st.Message()
		},
		Matcher: matcher,
	})
}
