package grpcmatchers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ErrorInfo", func() {
	st, err := status.New(codes.NotFound, "random error").WithDetails(&errdetails.ErrorInfo{
		Reason: "reason",
		Domain: "domain",
		Metadata: map[string]string{
			"key": "value",
		},
	})
	Expect(err).ToNot(HaveOccurred())
	err = st.Err()

	Describe("HaveErrorInfoReason", func() {
		It("should match an error info by reason", func() {
			Expect(err).To(HaveErrorInfoReason(Equal("reason")))
		})

		It("should not match an error info by reason", func() {
			Expect(err).ToNot(HaveErrorInfoReason(Equal("reason_nonexisting")))
		})
	})

	Describe("HaveErrorInfoDomain", func() {
		It("should match an error info by domain", func() {
			Expect(err).To(HaveErrorInfoDomain(Equal("domain")))
		})

		It("should not match an error info by domain", func() {
			Expect(err).ToNot(HaveErrorInfoDomain(Equal("domain_nonexisting")))
		})
	})

	Describe("HaveErrorInfoMetadata", func() {
		It("should match an error info by metadata", func() {
			Expect(err).To(HaveErrorInfoMetadata(HaveKey("key")))
		})

		It("should not match an non existing key", func() {
			Expect(err).ToNot(HaveErrorInfoMetadata(HaveKey("key_nonexisting")))
		})
	})
})

func ExampleHaveErrorInfoReason() {
	st, err := status.New(codes.Internal, "message").WithDetails(&errdetails.ErrorInfo{
		Reason: "some reason",
	})
	Expect(err).ToNot(HaveOccurred())

	Expect(st).To(HaveErrorInfoReason(Equal("some reason")))
}

func ExampleHaveErrorInfoDomain() {
	st, err := status.New(codes.Internal, "message").WithDetails(&errdetails.ErrorInfo{
		Domain: "some domain",
	})
	Expect(err).ToNot(HaveOccurred())

	Expect(st).To(HaveErrorInfoDomain(Equal("some domain")))
}

func ExampleHaveErrorInfoMetadata() {
	st, err := status.New(codes.Internal, "message").WithDetails(&errdetails.ErrorInfo{
		Metadata: map[string]string{
			"key1": "value1",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	Expect(st).To(HaveErrorInfoMetadata(HaveKey("key1")))
}
