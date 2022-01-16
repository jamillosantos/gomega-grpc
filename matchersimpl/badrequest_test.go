package matchersimpl

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jamillosantos/gomock-grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func Test_findBadRequest(t *testing.T) {
	t.Run("should find an error info", func(t *testing.T) {
		st := status.New(codes.Internal, "message")
		wantBadRequest := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "field1",
					Description: "description1",
				},
			},
		}
		st, err := st.WithDetails(wantBadRequest)
		require.NoError(t, err)
		gotBadRequest, ok := findBadRequest(st)
		assert.True(t, proto.Equal(wantBadRequest, gotBadRequest))
		assert.True(t, ok)
	})

	t.Run("should not find an error info", func(t *testing.T) {
		st := status.New(codes.Internal, "message")
		reqInfo := &errdetails.RequestInfo{
			RequestId: "rid",
		}
		st, err := st.WithDetails(reqInfo)
		require.NoError(t, err)
		gotBadRequest, ok := findBadRequest(st)
		assert.Nil(t, gotBadRequest)
		assert.False(t, ok)
	})
}

func TestGRPCBadRequestMatcher_Match(t *testing.T) {
	t.Run("should fail when there is not error info", func(t *testing.T) {
		matcher := NewGRPCMatchBadRequest(nil)
		gotResult, err := matcher.Match(status.New(codes.Internal, "random error").Err())
		assert.False(t, gotResult)
		assert.ErrorIs(t, err, errBadRequestNotFound)
	})

	t.Run("should return the match result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		errInfoMatcher := NewMockBadRequestMatcher(ctrl)
		matcher := NewGRPCMatchBadRequest(errInfoMatcher)

		wantErrInfo := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "field1",
					Description: "description1",
				},
			},
		}
		st, err := status.New(codes.Internal, "random error").WithDetails(wantErrInfo)
		require.NoError(t, err)

		errInfoMatcher.EXPECT().Match(gomockgrpc.ProtoEqual(wantErrInfo)).Return(true, nil)

		gotResult, err := matcher.Match(st.Err())
		assert.True(t, gotResult)
		assert.NoError(t, err, errBadRequestNotFound)
	})
}

func TestGRPCBadRequestMatcher_FailureMessage(t *testing.T) {
	t.Run("should fail finding BadRequest", func(t *testing.T) {
		matcher := &GRPCBadRequestMatcher{nil}
		gotMessage := matcher.FailureMessage(status.New(codes.Internal, "random error"))
		assert.Contains(t, gotMessage, "does not have an *errdetails.BadRequest in the details")
	})

	t.Run("should fail finding BadRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		errInfoMatcher := NewMockBadRequestMatcher(ctrl)
		matcher := &GRPCBadRequestMatcher{errInfoMatcher}

		wantErrInfo := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "field1",
					Description: "description1",
				},
			},
		}
		st, err := status.New(codes.Internal, "random error").WithDetails(wantErrInfo)
		require.NoError(t, err)

		errInfoMatcher.EXPECT().FailureMessage(gomockgrpc.ProtoEqual(wantErrInfo)).Return(wantMessage)

		gotMessage := matcher.FailureMessage(st)
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestGRPCBadRequestMatcher_NegatedFailureMessage(t *testing.T) {
	t.Run("should fail finding BadRequest", func(t *testing.T) {
		matcher := &GRPCBadRequestMatcher{nil}
		gotMessage := matcher.NegatedFailureMessage(status.New(codes.Internal, "random error"))
		assert.Contains(t, gotMessage, "does not have an *errdetails.BadRequest in the details")
	})

	t.Run("should fail finding BadRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		errInfoMatcher := NewMockBadRequestMatcher(ctrl)
		matcher := &GRPCBadRequestMatcher{errInfoMatcher}

		wantErrInfo := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "field1",
					Description: "description1",
				},
			},
		}
		st, err := status.New(codes.Internal, "random error").WithDetails(wantErrInfo)
		require.NoError(t, err)

		errInfoMatcher.EXPECT().NegatedFailureMessage(gomockgrpc.ProtoEqual(wantErrInfo)).Return(wantMessage)

		gotMessage := matcher.NegatedFailureMessage(st)
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestGRPCBadRequestFieldViolation_Match(t *testing.T) {
	violation := errdetails.BadRequest_FieldViolation{
		Field:       "field",
		Description: "description",
	}
	anotherViolation := errdetails.BadRequest_FieldViolation{
		Field:       "another field",
		Description: "description",
	}
	violations := []*errdetails.BadRequest_FieldViolation{
		&anotherViolation,
		&violation,
	}

	tests := []struct {
		name   string
		actual *errdetails.BadRequest
		want   bool
	}{
		{"should match when given a BadRequest", &errdetails.BadRequest{
			FieldViolations: violations,
		}, true},
		{"should not match when given a BadRequest", &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				&anotherViolation,
			},
		}, true},
		{"should not match when there are no BadRequest", &errdetails.BadRequest{
			FieldViolations: nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GRPCBadRequestFieldViolation{
				Field:       "field",
				Description: "description",
			}
			got, err := m.Match(tt.actual)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGRPCBadRequestFieldViolation_matchesFV(t *testing.T) {
	tests := []struct {
		name        string
		field       string
		description string
		want        bool
	}{
		{"should not match a given field does not match", "non matching field", "", false},
		{"should match when field matches and description is empty", "field", "", true},
		{"should not match when field matches and description does not match", "field", "non matching description", false},
		{"should match when both field and description match", "field", "description", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GRPCBadRequestFieldViolation{
				Field:       tt.field,
				Description: tt.description,
			}
			got := m.matchesFV(&errdetails.BadRequest_FieldViolation{
				Field:       "field",
				Description: "description",
			})
			assert.Equal(t, tt.want, got)
		})
	}
}
