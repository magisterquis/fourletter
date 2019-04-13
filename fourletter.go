// Package fourletter encodes bytes into four-byte words
//
// Its intended use is as a silly encoding to turn arbitrary text into cat
// noises for low-throughput DNS tunneling which will bypass defensive tools
// which alert on entropy.
package fourletter

import (
	"errors"
	"fmt"
)

/*
 * fourletter.go
 * Encode bytes into four-byte words
 * By J. Stuart McMurray
 * Created 20190412
 * Last Modified 20190413
 */

// DefaultEncoding is an encoding using the default cat noises
var DefaultEncoding = MustNewEncoding("meowmrowpurrmeww")

// An Encoding is a four-word encoding scheme which encodes each byte as four
// four-byte words.
type Encoding struct {
	ws []string
}

// NewEncoding returns a new Encoding defined by the given alphabet, which must
// be a 16-byte string which will be interpreted as four 4-byte words.
func NewEncoding(encoder string) (*Encoding, error) {
	/* Make sure we have the right number of bytes */
	if 16 != len(encoder) {
		return nil, errors.New("encoder not 16 bytes")
	}

	var enc Encoding
	enc.ws = make([]string, 4)

	/* Grab words */
	enc.ws[0] = encoder[0:4]
	enc.ws[1] = encoder[4:8]
	enc.ws[2] = encoder[8:12]
	enc.ws[3] = encoder[12:16]

	/* Make sure the words are unique */
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			if enc.ws[i] == enc.ws[j] {
				return nil, errors.New("words not unique")
			}
		}
	}

	return &enc, nil
}

// MustNewEncoding is like NewEncoding but panics if encoder isn't 16 bytes
// long.
func MustNewEncoding(encoder string) *Encoding {
	enc, err := NewEncoding(encoder)
	if nil != err {
		panic(err)
	}

	return enc
}

// Decode decodes src into dst, which must be large enough to hold the decoded
// bytes.
func (enc *Encoding) Decode(dst, src []byte) error {
	/* Check the buffer sizes */
	if 0 != len(src)%16 {
		return errors.New("invalid source length")
	}
	if len(dst) < len(src)/16 {
		return errors.New("destination buffer too small")
	}

	/* Decode it all */
	db := -1 /* Current dst byte */
	for i := 0; i < len(src); i += 4 {
		/* Move to the next dst byte when we're ready */
		if 0 == i%16 {
			db++
			dst[db] = 0
		}
		/* Source word */
		w := src[i : i+4]
		var tb byte = 255
		/* Turn it into two bits. */
		for j := 0; j < 4; j++ {
			if string(w) == enc.ws[j] {
				tb = byte(j)
				break
			}
		}
		if 255 == tb {
			return fmt.Errorf("invalid word %q", w)
		}
		dst[db] >>= 2
		dst[db] |= tb << 6
	}

	return nil
}

// DecodeString returns the bytes represented by s.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	o := make([]byte, len(s)/16)
	if err := enc.Decode(o, []byte(s)); nil != err {
		return nil, err
	}
	return o, nil
}

// Encode places an src, encoded, into dst.
func (enc *Encoding) Encode(dst, src []byte) error {
	if len(dst) < len(src)*16 {
		return errors.New("destination buffer too small")
	}
	cur := 0
	for _, v := range src {
		for i := 0; i < 4; i++ {
			copy(dst[cur:], []byte(enc.ws[0x03&v]))
			v >>= 2
			cur += 4
		}
	}

	return nil
}

// EncodeToString returns a string containing src, encoded.
func (enc *Encoding) EncodeToString(src []byte) string {
	dst := make([]byte, len(src)*16)
	if err := enc.Encode(dst, src); nil != err {
		panic(err)
	}
	return string(dst)
}
