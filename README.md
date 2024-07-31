# gomega-grpc

When comparing protobuf objects in Gomega could be a bit complicated as unmarshaled structs cannot be directly compared
by the tranditional gomega `Equal`. This library provides a set of matchers to compare protobuf objects.

## Usage

### ProtoEqual

```go
import (
    . "github.com/onsi/gomega"
    . "github.com/onsi/gomega-grpc"
)

var _ = Describe("MyTest", func() {
    It("should compare protobuf objects", func() {
        Expect(Message{
			Message: "message 1",
        }).To ProtoEqual(Message{
            Message: "message 2",
        })
    })
})
```

The above will produce:

```
  Expected
    <string>: {
      "message": "message 2"
    }
  to equal
    <string>: {
      "message": "message 1"
    }

  Diff:
   {
  -  "message": "message 1"
  +  "message": "message 2"
   }
```

