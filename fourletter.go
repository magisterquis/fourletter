// Package fourletter encodes bytes into four-byte words
//
// Its intended use is as a silly encoding to turn arbitrary text into cat
// noises for low-throughput DNS tunneling which will bypass defensive tools
// which alert on entropy.
package fourletter

/*
 * fourletter.go
 * Encode bytes into four-byte words
 * By J. Stuart McMurray
 * Created 20190412
 * Last Modified 20190416
 */

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

/* buflen is the length used for the io.{Read,Write}er buffers*/
const buflen = 1024

// DefaultEncoding is an encoding using the default cat noises
var DefaultEncoding = MustNewEncoding("meowmrowpurrmeww")

// An Encoding is a four-word encoding scheme which encodes each byte as four
// four-byte words.
type Encoding struct {
	ws [4][4]byte
}

// NewEncoding returns a new Encoding defined by the given alphabet, which must
// be a 16-byte string which will be interpreted as four 4-byte words.
func NewEncoding(encoder string) (*Encoding, error) {
	/* Make sure we have the right number of bytes */
	if 16 != len(encoder) {
		return nil, errors.New("encoder not 16 bytes")
	}

	var enc Encoding

	/* Grab words */
	copy(enc.ws[0][:], encoder[0:4])
	copy(enc.ws[1][:], encoder[4:8])
	copy(enc.ws[2][:], encoder[8:12])
	copy(enc.ws[3][:], encoder[12:16])

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

	return newDecoder(enc, bytes.NewReader(src)).decodeInto(dst)
}

// DecodeString returns the bytes represented by s.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	if 0 != len(s)%16 {
		return nil, errors.New("invalid source length")
	}

	dst := make([]byte, len(s)/16)

	return dst, newDecoder(enc, strings.NewReader(s)).decodeInto(dst)
}

// Encode places an src, encoded, into dst.
func (enc *Encoding) Encode(dst, src []byte) error {
	if len(src)*16 > len(dst) {
		return errors.New("destination buffer too small")
	}
	b := bytes.NewBuffer(dst[:0])
	_, err := newEncoder(enc, b).Write(src)
	return err
}

// EncodeToString returns a string containing src, encoded.
func (enc *Encoding) EncodeToString(src []byte) string {
	var b strings.Builder
	if _, err := newEncoder(enc, &b).Write(src); nil != err {
		/* Should never happen */
		panic(err)
	}
	return b.String()
}

// EncodeToWriter writes the encoded form of src to dst.  It stops on the first
// error returned by dst.Write.
func (enc *Encoding) EncodeToWriter(dst io.Writer, src []byte) (n int, err error) {
	return newEncoder(enc, dst).Write(src)
}
