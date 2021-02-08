package json

import (
	"encoding/json"

	"github.com/liyiysng/scatter/encoding"
)

// Name is the name registered for the proto compressor.
const Name = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
