package fourletter

/*
 * decoder.go
 * Streamingly decode four-byte words
 * By J. Stuart McMurray
 * Created 20190415
 * Last Modified 20190416
 */

/* Thanks to https://github.com/jrick */

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

/* decoder wraps an io.Reader to decode 4-byte words it gets from it */
type decoder struct {
	sync.Mutex
	enc *Encoding
	r   io.Reader
	buf []byte

	/* Holds leftover data between reads */
	leftover  []byte
	nleftover int
}

// NewDecoder constructs a new fourletter stream decoder.
func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	return newDecoder(enc, r)
}

func newDecoder(enc *Encoding, r io.Reader) *decoder {
	return &decoder{
		enc:      enc,
		r:        r,
		buf:      make([]byte, buflen),
		leftover: make([]byte, 16),
	}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	d.Lock()
	defer d.Unlock()

	/* Make sure we have enough buffer space to fill p */
	if 16*len(p) > len(d.buf) {
		d.buf = make([]byte, 16*len(p))
	}

	/* Use any leftover data */
	if 0 != d.nleftover {
		copy(d.buf, d.leftover[:d.nleftover])
	}

	/* Grab as much encoded data as we need to fill p */
	nr, err := d.r.Read(d.buf[d.nleftover : 16*len(p)])
	if 0 == nr {
		return 0, err
	}
	nr += d.nleftover

	/* Decode the rest */
	var db = -1 /* Current dst byte index */
	for i := 0; i < nr; i += 4 {
		/* Move to the next dst byte when we're ready */
		if 0 == i%16 {
			db++
			p[db] = 0
		}
		/* Source word */
		w := d.buf[i : i+4]
		var tb byte = 255
		/* Turn it into two bits. */
		for j := 0; j < len(d.enc.ws); j++ {
			if 0 == bytes.Compare(w, d.enc.ws[j][:]) {
				tb = byte(j)
				break
			}
		}
		if 255 == tb {
			return db, fmt.Errorf("invalid word %q", w)
		}
		p[db] >>= 2
		p[db] |= tb << 6
	}

	/* If we have any read left over, save it for next time. */
	if 0 != nr-db {
		copy(d.leftover, d.buf[db*16:nr])
	}

	return db, err
}

/* decodeInto calls d.Read and fills up dst until EOF or an error.  decodeInto
return nil on EOF. */
func (d *decoder) decodeInto(dst []byte) error {
	/* Read into the dst */
	var (
		nr  int
		tot int
		err error
	)
	for len(dst) > tot {
		nr, err = d.Read(dst[tot:])
		tot += nr
		if nil != err {
			break
		}
	}

	/* EOF is expected */
	if err == io.EOF {
		err = nil
	}
	return err
}
