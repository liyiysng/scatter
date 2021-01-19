package ratelimit

import "io"

type reader struct {
	r      io.Reader
	bucket *Bucket
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Reader(r io.Reader, bucket *Bucket) io.Reader {
	return &reader{
		r:      r,
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}
