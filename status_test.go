package grpcmatchers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ErrorInfo", func() {
	st := status.New(codes.NotFound, "random error")
	err := st.Err()

	Describe("HaveStatusCode", func() {
		It("should match a status code", func() {
			Expect(err).To(HaveStatusCode(Equal(codes.NotFound)))
		})

		It("should not match a status code", func() {
			Expect(err).ToNot(HaveStatusCode(Equal(codes.Internal)))
		})
	})

	Describe("HaveStatusMessage", func() {
		It("should match an error info by domain", func() {
			Expect(err).To(HaveStatusMessage(Equal("random error")))
		})

		It("should not match an error info by domain", func() {
			Expect(err).ToNot(HaveStatusMessage(Equal("not random error")))
		})
	})
})

func ExampleHaveStatusCode() {
	err := status.New(codes.Internal, "message").Err()

	Expect(err).To(HaveStatusCode(Equal(codes.Internal)))
}

func ExampleHaveStatusMessage() {
	err := status.New(codes.Internal, "this is a message").Err()

	Expect(err).To(HaveStatusMessage(Equal("this is a message")))
	Expect(err).To(HaveStatusMessage(ContainSubstring("message")))
}
