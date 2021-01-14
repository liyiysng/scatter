package encoding

import (
	"io"
	"strings"
)

// Codec defines the interface node rpc uses to encode and decode messages.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
}

var registeredCodecs = make(map[string]Codec)

// RegisterCodec registers the provided Codec for use with all node clients and
// servers.
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.  If multiple Compressors are
// registered with the same name, the one registered last will take effect.
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentSubtype string) Codec {
	return registeredCodecs[contentSubtype]
}

// Compressor is used for compressing and decompressing when sending or
// receiving messages.
type Compressor interface {
	// Compress writes the data written to wc to w after compressing it.  If an
	// error occurs while initializing the compressor, that error is returned
	// instead.
	Compress(w io.Writer) (io.WriteCloser, error)
	// Decompress reads data from r, decompresses it, and provides the
	// uncompressed data via the returned io.Reader.  If an error occurs while
	// initializing the decompressor, that error is returned instead.
	Decompress(r io.Reader) (io.Reader, error)
	// Name is the name of the compression codec and is used to set the content
	// coding header.  The result must be static; the result cannot change
	// between calls.
	Name() string
	// If a Compressor implements
	// DecompressedSize(compressedBytes []byte) int, rpc will call it
	// to determine the size of the buffer allocated for the result of decompression.
	// Return -1 to indicate unknown size.
}

var registeredCompressor = make(map[string]Compressor)

// RegisterCompressor registers the compressor with gRPC by its name.  It can
// be activated when sending an RPC via grpc.UseCompressor().  It will be
// automatically accessed when receiving a message based on the content coding
// header.  Servers also use it to send a response with the same encoding as
// the request.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.  If multiple Compressors are
// registered with the same name, the one registered last will take effect.
func RegisterCompressor(c Compressor) {
	registeredCompressor[c.Name()] = c
}

// GetCompressor returns Compressor for the given compressor name.
func GetCompressor(name string) Compressor {
	return registeredCompressor[name]
}
