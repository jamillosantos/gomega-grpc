package grpcmatchers

import (
	"github.com/onsi/gomega/types"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/jamillosantos/gomega-grpc/matchersimpl"
)

// HaveErrorInfoReason will match the *errdetails.ErrorInfo Reason property against the given matcher.
func HaveErrorInfoReason(matcher types.GomegaMatcher) *matchersimpl.GRPCStatusMatcher {
	return matchersimpl.NewGRPCMatchErrorInfo(&matchersimpl.GRPCErrorInfoReasonMatcher{
		PropMap: func(errInfo *errdetails.ErrorInfo) interface{} {
			return errInfo.GetReason()
		},
		Matcher: matcher,
	})
}

// HaveErrorInfoDomain will match the *errdetails.ErrorInfo Reason property against the given matcher.
func HaveErrorInfoDomain(matcher types.GomegaMatcher) *matchersimpl.GRPCStatusMatcher {
	return matchersimpl.NewGRPCMatchErrorInfo(&matchersimpl.GRPCErrorInfoReasonMatcher{
		PropMap: func(errInfo *errdetails.ErrorInfo) interface{} {
			return errInfo.GetDomain()
		},
		Matcher: matcher,
	})
}

// HaveErrorInfoMetadata will match the *errdetails.ErrorInfo Reason property against the given matcher.
func HaveErrorInfoMetadata(matcher types.GomegaMatcher) *matchersimpl.GRPCStatusMatcher {
	return matchersimpl.NewGRPCMatchErrorInfo(&matchersimpl.GRPCErrorInfoReasonMatcher{
		PropMap: func(errInfo *errdetails.ErrorInfo) interface{} {
			return errInfo.GetMetadata()
		},
		Matcher: matcher,
	})
}
