package grpcmatchers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("BadRequest", func() {
	st, err := status.New(codes.NotFound, "random error").WithDetails(&errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       "field1",
				Description: "description1",
			},
			{
				Field:       "field2",
				Description: "description2",
			},
			{
				Field:       "field3",
				Description: "description3",
			},
		},
	})
	Expect(err).ToNot(HaveOccurred())
	err = st.Err()

	Describe("HaveFieldValidation", func() {
		It("should match a given field violation by field", func() {
			Expect(err).To(HaveFieldViolation("field1"))
		})
		It("should match a given field violation by field and description", func() {
			Expect(err).To(HaveFieldViolation("field2", "description2"))
		})
		It("should not match a given field violation by field and description", func() {
			Expect(err).ToNot(HaveFieldViolation("field_non_existing"))
		})
		It("should not match a given field violation by field and description", func() {
			Expect(err).ToNot(HaveFieldViolation("field2", "description3"))
		})
	})
})
