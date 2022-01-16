//go:generate go run github.com/golang/mock/mockgen -package matchersimpl -imports status=google.golang.org/grpc/status -destination mocks_test.go github.com/jamillosantos/gomega-grpc/matchersimpl StatusMatcher,ErrorInfoMatcher,BadRequestMatcher
//go:generate go run ../tools/replace_internal_status/main.go -- mocks_test.go
//go:generate go run github.com/golang/mock/mockgen -package matchersimpl -destination gomega_matcher_mock_test.go github.com/onsi/gomega/types GomegaMatcher

package matchersimpl

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewGRPCStatusMatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	sm := NewMockStatusMatcher(ctrl)
	m := NewGRPCStatusMatcher(sm)
	assert.Equal(t, sm, m.statusMatcher)
}

func TestGRPCStatusMatcher_Match(t *testing.T) {
	noop := func(m *MockStatusMatcher) {}

	wantErr := errors.New("want error")

	tests := []struct {
		name         string
		givenError   interface{}
		expectedCode codes.Code
		setupMock    func(m *MockStatusMatcher)
		wantSuccess  bool
		wantErr      error
	}{
		{"should match when StatusMatcher matches", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.InvalidArgument, func(m *MockStatusMatcher) {
			m.EXPECT().Match(status.New(codes.InvalidArgument, "invalid argument")).Return(true, nil)
		}, true, nil},
		{"should not match when StatusMatcher does not matches", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.InvalidArgument, func(m *MockStatusMatcher) {
			m.EXPECT().Match(gomock.Any()).Return(false, nil)
		}, false, nil},
		{"should fail StatusMatcher fails", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.InvalidArgument, func(m *MockStatusMatcher) {
			m.EXPECT().Match(gomock.Any()).Return(false, wantErr)
		}, false, wantErr},
		{"should fail when no error given", false, codes.OK, noop, false, errExpectedError},
		{"should fail when no status given", errors.New("non status error"), codes.OK, noop, false, errStatusErrExpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := NewMockStatusMatcher(ctrl)
			tt.setupMock(m)
			matcher := NewGRPCStatusMatcher(m)
			gotSuccess, err := matcher.Match(tt.givenError)
			assert.Equal(t, tt.wantSuccess, gotSuccess)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGRPCStatusMatcher_validateActual(t *testing.T) {
	tests := []struct {
		name         string
		givenError   interface{}
		expectedCode codes.Code
		wantMessage  string
		wantStatus   *status.Status
		wantResult   bool
	}{
		{"should fail when no error given", false, codes.OK, "is not an error", nil, false},
		{"should fail when no status given", errors.New("non status error"), codes.OK, "is not a grpc error", nil, false},
		{"should validate", status.New(codes.InvalidArgument, "invalid argument").Err(), codes.NotFound, "", status.New(codes.InvalidArgument, "invalid argument"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewGRPCStatusMatcher(nil)
			gotMessage, gotStatus, gotResult := matcher.validateActual(tt.givenError)
			assert.Contains(t, gotMessage, tt.wantMessage)
			assert.Equal(t, gotStatus, tt.wantStatus)
			assert.Equal(t, gotResult, tt.wantResult)
		})
	}
}

func TestGRPCStatusMatcher_FailureMessage(t *testing.T) {
	t.Run("should fail validation", func(t *testing.T) {
		matcher := NewGRPCStatusMatcher(nil)
		msg := matcher.FailureMessage(false)
		assert.Contains(t, msg, "is not an error")
	})

	t.Run("should call StatusMatcher FailureMessage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		m := NewMockStatusMatcher(ctrl)
		matcher := NewGRPCStatusMatcher(m)
		wantActual := status.New(codes.NotFound, "not found")
		m.EXPECT().FailureMessage(wantActual).Return(wantMessage)
		gotMessage := matcher.FailureMessage(wantActual.Err())
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestGRPCStatusMatcher_NegatedFailureMessage(t *testing.T) {
	t.Run("should fail validation", func(t *testing.T) {
		matcher := NewGRPCStatusMatcher(nil)
		msg := matcher.NegatedFailureMessage(false)
		assert.Contains(t, msg, "is not an error")
	})

	t.Run("should call StatusMatcher NegatedFailureMessage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		m := NewMockStatusMatcher(ctrl)
		matcher := NewGRPCStatusMatcher(m)
		wantActual := status.New(codes.NotFound, "not found")
		m.EXPECT().NegatedFailureMessage(wantActual).Return(wantMessage)
		gotMessage := matcher.NegatedFailureMessage(wantActual.Err())
		assert.Equal(t, wantMessage, gotMessage)
	})
}

func TestGRPCStatusPropMatcher_Match(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := NewMockGomegaMatcher(ctrl)

	wantST := status.New(codes.Internal, "message")

	matcher := GRPCStatusPropMatcher{
		PropMap: func(st *status.Status) interface{} {
			return st.Code()
		},
		Matcher: gm,
	}
	gm.EXPECT().Match(wantST.Code()).Return(true, nil)
	gotResult, err := matcher.Match(wantST)
	assert.NoError(t, err)
	assert.True(t, gotResult)
}

func TestGRPCStatusPropMatcher_FailureMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := NewMockGomegaMatcher(ctrl)

	wantST := status.New(codes.Internal, "message")

	matcher := GRPCStatusPropMatcher{
		PropMap: func(st *status.Status) interface{} {
			return st.Code()
		},
		Matcher: gm,
	}
	gm.EXPECT().FailureMessage(wantST.Code()).Return(wantMessage)
	gotMessage := matcher.FailureMessage(wantST)
	assert.Equal(t, wantMessage, gotMessage)
}

func TestGRPCStatusPropMatcher_NegatedFailureMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	gm := NewMockGomegaMatcher(ctrl)

	wantST := status.New(codes.Internal, "message")

	matcher := GRPCStatusPropMatcher{
		PropMap: func(st *status.Status) interface{} {
			return st.Code()
		},
		Matcher: gm,
	}
	gm.EXPECT().NegatedFailureMessage(wantST.Code()).Return(wantMessage)
	gotMessage := matcher.NegatedFailureMessage(wantST)
	assert.Equal(t, wantMessage, gotMessage)
}
