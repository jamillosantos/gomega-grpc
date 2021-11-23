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

var (
	wantMessage = "random message"
)

func Test_findErrorInfo(t *testing.T) {
	t.Run("should find an error info", func(t *testing.T) {
		st := status.New(codes.Internal, "message")
		wantErrorInfo := &errdetails.ErrorInfo{
			Reason: "reason",
		}
		st, err := st.WithDetails(wantErrorInfo)
		require.NoError(t, err)
		gotErrorInfo, ok := findErrorInfo(st)
		assert.True(t, proto.Equal(wantErrorInfo, gotErrorInfo))
		assert.True(t, ok)
	})

	t.Run("should not find an error info", func(t *testing.T) {
		st := status.New(codes.Internal, "message")
		reqInfo := &errdetails.RequestInfo{
			RequestId: "rid",
		}
		st, err := st.WithDetails(reqInfo)
		require.NoError(t, err)
		gotErrorInfo, ok := findErrorInfo(st)
		assert.Nil(t, gotErrorInfo)
		assert.False(t, ok)
	})
}

func TestGRPCErrorInfoMatcher_Match(t *testing.T) {
	t.Run("should fail when there is not error info", func(t *testing.T) {
		matcher := NewGRPCMatchErrorInfo(nil)
		gotResult, err := matcher.Match(status.New(codes.Internal, "random error").Err())
		assert.False(t, gotResult)
		assert.ErrorIs(t, err, errErrorInfoNotFound)
	})

	t.Run("should return the match result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		errInfoMatcher := NewMockErrorInfoMatcher(ctrl)
		matcher := NewGRPCMatchErrorInfo(errInfoMatcher)

		wantErrInfo := &errdetails.ErrorInfo{
			Reason: "reason",
		}
		st, err := status.New(codes.Internal, "random error").WithDetails(wantErrInfo)
		require.NoError(t, err)

		errInfoMatcher.EXPECT().Match(gomockgrpc.ProtoEqual(wantErrInfo)).Return(true, nil)

		gotResult, err := matcher.Match(st.Err())
		assert.True(t, gotResult)
		assert.NoError(t, err, errErrorInfoNotFound)
	})
}

func TestGRPCErrorInfoMatcher_FailureMessage(t *testing.T) {
	t.Run("should fail finding ErrorInfo", func(t *testing.T) {
		matcher := &GRPCErrorInfoMatcher{nil}
		gotMessage := matcher.FailureMessage(status.New(codes.Internal, "random error"))
		assert.Contains(t, gotMessage, "does not have an *errdetails.ErrorInfo in the details")
	})

	t.Run("should fail finding ErrorInfo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		errInfoMatcher := NewMockErrorInfoMatcher(ctrl)
		matcher := &GRPCErrorInfoMatcher{errInfoMatcher}

		wantErrInfo := &errdetails.ErrorInfo{
			Reason: "reason",
		}
		st, err := status.New(codes.Internal, "random error").WithDetails(wantErrInfo)
		require.NoError(t, err)

		errInfoMatcher.EXPECT().FailureMessage(gomockgrpc.ProtoEqual(wantErrInfo)).Return(wantMessage)

		gotMessage := matcher.FailureMessage(st)
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestGRPCErrorInfoMatcher_NegatedFailureMessage(t *testing.T) {
	t.Run("should fail finding ErrorInfo", func(t *testing.T) {
		matcher := &GRPCErrorInfoMatcher{nil}
		gotMessage := matcher.NegatedFailureMessage(status.New(codes.Internal, "random error"))
		assert.Contains(t, gotMessage, "does not have an *errdetails.ErrorInfo in the details")
	})

	t.Run("should fail finding ErrorInfo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		errInfoMatcher := NewMockErrorInfoMatcher(ctrl)
		matcher := &GRPCErrorInfoMatcher{errInfoMatcher}

		wantErrInfo := &errdetails.ErrorInfo{
			Reason: "reason",
		}
		st, err := status.New(codes.Internal, "random error").WithDetails(wantErrInfo)
		require.NoError(t, err)

		errInfoMatcher.EXPECT().NegatedFailureMessage(gomockgrpc.ProtoEqual(wantErrInfo)).Return(wantMessage)

		gotMessage := matcher.NegatedFailureMessage(st)
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestGRPCErrorInfoReasonMatcher_Match(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := NewMockGomegaMatcher(ctrl)

	wantErrorInfo := &errdetails.ErrorInfo{
		Reason: "reason",
	}

	matcher := GRPCErrorInfoReasonMatcher{
		PropMap: func(errInfo *errdetails.ErrorInfo) interface{} {
			return errInfo.GetReason()
		},
		Matcher: gm,
	}
	gm.EXPECT().Match(wantErrorInfo.Reason).Return(true, nil)
	gotResult, err := matcher.Match(wantErrorInfo)
	assert.NoError(t, err)
	assert.True(t, gotResult)
}

func TestGRPCErrorInfoReasonMatcher_FailureMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := NewMockGomegaMatcher(ctrl)

	wantErrorInfo := &errdetails.ErrorInfo{
		Reason: "reason",
	}

	matcher := GRPCErrorInfoReasonMatcher{
		PropMap: func(errInfo *errdetails.ErrorInfo) interface{} {
			return errInfo.GetReason()
		},
		Matcher: gm,
	}
	gm.EXPECT().FailureMessage(wantErrorInfo.Reason).Return(wantMessage)
	gotMessage := matcher.FailureMessage(wantErrorInfo)
	assert.Equal(t, wantMessage, gotMessage)
}

func TestGRPCErrorInfoReasonMatcher_NegatedFailureMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := NewMockGomegaMatcher(ctrl)

	wantErrorInfo := &errdetails.ErrorInfo{
		Reason: "reason",
	}

	matcher := GRPCErrorInfoReasonMatcher{
		PropMap: func(errInfo *errdetails.ErrorInfo) interface{} {
			return errInfo.GetReason()
		},
		Matcher: gm,
	}
	gm.EXPECT().NegatedFailureMessage(wantErrorInfo.Reason).Return(wantMessage)
	gotMessage := matcher.NegatedFailureMessage(wantErrorInfo)
	assert.Equal(t, wantMessage, gotMessage)
}
