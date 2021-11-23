package matchersimpl

import (
	"errors"

	"github.com/onsi/gomega/format"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	errProtoEqualNil              = errors.New("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.") //nolint // Needs verbosity on the test output.
	errProtoEqualActualNotMessage = errors.New("The given actual value is not a proto.Message")                                                                                                                        // nolint // This is the gomega standard.
)

type ProtoEqualMatcher struct {
	Expected proto.Message
}

func (matcher *ProtoEqualMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, errProtoEqualNil
	}
	actualProtoMessage, ok := actual.(proto.Message)
	if !ok {
		return false, errProtoEqualActualNotMessage
	}
	return proto.Equal(actualProtoMessage, matcher.Expected), nil
}

func (matcher *ProtoEqualMatcher) FailureMessage(actual interface{}) (message string) {
	if pactual, ok := actual.(proto.Message); ok {
		return format.Message(protojson.Format(pactual), "to equal", protojson.Format(matcher.Expected))
	}
	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *ProtoEqualMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	if pactual, ok := actual.(proto.Message); ok {
		return format.Message(protojson.Format(pactual), "not to equal", protojson.Format(matcher.Expected))
	}
	return format.Message(actual, "not to equal", matcher.Expected)
}
