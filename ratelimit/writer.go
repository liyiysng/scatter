package ratelimit

import "io"

type writer struct {
	w      io.Writer
	bucket *Bucket
}

// Writer returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Writer(w io.Writer, bucket *Bucket) io.Writer {
	return &writer{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}
