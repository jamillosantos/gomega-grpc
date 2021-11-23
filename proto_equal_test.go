package grpcmatchers

import (
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func ExampleProtoEqual() {
	errInfo := &errdetails.ErrorInfo{
		Reason: "some reason",
	}

	Expect(errInfo).To(ProtoEqual(&errdetails.ErrorInfo{
		Reason: "some reason",
	}))
}
