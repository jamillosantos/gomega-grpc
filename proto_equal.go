package grpcmatchers

import (
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
)

// ProtoEqual is the matcher for comparing proto.Message values. This matcher is powered by the
// proto.Equal function, provided by the default GRPC implementation.
func ProtoEqual(m proto.Message) types.GomegaMatcher {
	return &matchers.EqualMatcher{
		Expected: m,
	}
}
