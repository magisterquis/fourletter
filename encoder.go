package fourletter

/*
 * encoder.go
 * Streamingly encode bytes into four-byte words
 * By J. Stuart McMurray
 * Created 20190415
 * Last Modified 20190416
 */

/* Thanks to https://github.com/jrick */

import (
	"io"
	"sync"
)

type encoder struct {
	sync.Mutex

	w   io.Writer
	enc *Encoding
	buf []byte

	closed   bool
	closeErr error
}

// NewEncoder constructs a new fourletter stream decoder.
func NewEncoder(enc *Encoding, w io.Writer) io.Writer {
	return newEncoder(enc, w)
}

func newEncoder(enc *Encoding, w io.Writer) *encoder {
	return &encoder{
		w:   w,
		enc: enc,
	}
}

/* Write writes 16 encoded bytes to the underlying io.Writer for every
byte in p. */
func (e *encoder) Write(p []byte) (n int, err error) {
	e.Lock()
	defer e.Unlock()

	/* Make sure we have enough buffer room */
	if 16*len(p) > len(e.buf) {
		e.buf = make([]byte, 16*len(p))
	}

	/* Encode p into the buffer */
	var (
		bufi int /* Index into buf */
	)
	for _, b := range p {
		/* Encode the byte */
		for i := 0; i < 4; i++ {
			copy(e.buf[bufi:], e.enc.ws[0x03&b][:])
			b >>= 2
			bufi += 4
		}
	}

	/* Send it to the underlying encoder */
	_, err = e.w.Write(e.buf[:bufi])
	return bufi / 16, err
}
