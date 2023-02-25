package grpcmatchers

import (
	"strconv"

	"github.com/onsi/gomega/format"
	"google.golang.org/grpc/status"
)

type errorFormat struct {
	Status     string
	StatusCode string
	Message    string
	Details    []interface{}
}

type errorInfo struct {
}

func grpcErrorFormatter(value interface{}) (string, bool) {
	err, ok := value.(error)
	if !ok {
		return "", false
	}
	st, ok := status.FromError(err)
	if !ok {
		return "", false
	}
	return format.Object(errorFormat{
		Status:     st.Code().String(),
		StatusCode: strconv.Itoa(int(st.Code())),
		Message:    st.Message(),
		Details:    st.Details(),
	}, 0), true
}

func init() {
	format.UseStringerRepresentation = true
	_ = format.RegisterCustomFormatter(grpcErrorFormatter)
}
