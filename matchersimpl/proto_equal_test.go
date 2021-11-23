package matchersimpl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestProtoEqualMatcher_Match(t *testing.T) {
	t.Run("should fail when expected and actual is nil", func(t *testing.T) {
		gotMatch, err := (&ProtoEqualMatcher{nil}).Match(nil)
		assert.False(t, gotMatch)
		assert.ErrorIs(t, err, errProtoEqualNil)
	})

	t.Run("should fail when actual is not a proto.Message", func(t *testing.T) {
		gotMatch, err := (&ProtoEqualMatcher{nil}).Match(false)
		assert.False(t, gotMatch)
		assert.ErrorIs(t, err, errProtoEqualActualNotMessage)
	})
}

func TestProtoEqualMatcher_FailureMessage(t *testing.T) {
	t.Run("when actual is a proto.Message", func(t *testing.T) {
		gotMessage := (&ProtoEqualMatcher{
			Expected: &helloworld.HelloReply{
				Message: "message 1",
			},
		}).FailureMessage(&helloworld.HelloReply{
			Message: "message 2",
		})
		fmt.Println(gotMessage)
		assert.Contains(t, gotMessage, "to equal")
		assert.Contains(t, gotMessage, "message 1")
		assert.Contains(t, gotMessage, "message 2")
	})
	t.Run("when actual is NOT a proto.Message", func(t *testing.T) {
		gotMessage := (&ProtoEqualMatcher{nil}).FailureMessage("string")
		assert.Contains(t, gotMessage, "to equal")
		assert.Contains(t, gotMessage, "string")
		assert.Contains(t, gotMessage, "nil")
	})
}

func TestProtoEqualMatcher_NegatedFailureMessage(t *testing.T) {
	gotMessage := (&ProtoEqualMatcher{nil}).NegatedFailureMessage("string")
	assert.Contains(t, gotMessage, "not to equal")
	assert.Contains(t, gotMessage, "string")
	assert.Contains(t, gotMessage, "nil")
}
