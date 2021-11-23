package grpcmatchers

import (
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"

	"github.com/jamillosantos/gomega-grpc/matchersimpl"
)

// ProtoEqual is the matcher for comparing proto.Message values. This matcher is powered by the
// proto.Equal function, provided by the default GRPC implementation.
func ProtoEqual(m proto.Message) types.GomegaMatcher {
	return &matchersimpl.ProtoEqualMatcher{
		Expected: m,
	}
}
