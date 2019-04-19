package fourletter

/*
 * fourletter_test.go
 * Tests for fourletter
 * By J. Stuart McMurray
 * Created 20190413
 * Last Modified 20190413
 */

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewEncoding(t *testing.T) {
	if _, err := NewEncoding("a"); nil == err {
		t.Fatalf("encoding made with short string")
	}
	if _, err := NewEncoding("abcdabcd12345678"); nil == err {
		t.Fatalf("encoding made with non-unique words")
	}
	e, err := NewEncoding("aaaabbbbccccdddd")
	if nil != err {
		t.Fatalf("error on valid encoding: %v", err)
	}
	for _, c := range []struct {
		have []byte
		want string
	}{
		{[]byte{0x00}, "aaaaaaaaaaaaaaaa"},
		{[]byte{0x01}, "bbbbaaaaaaaaaaaa"},
		{[]byte{0xFF}, "dddddddddddddddd"},
		{
			[]byte{0xFF, 0x00, 0x88},
			"ddddddddddddddddaaaaaaaaaaaaaaaaaaaaccccaaaacccc",
		},
		{
			[]byte("all your base are belong to us"),
			"bbbbaaaaccccbbbbaaaaddddccccbbbbaaaaddddccccbbbbaaaaaaaaccccaaaabbbbccccddddbbbbddddddddccccbbbbbbbbbbbbddddbbbbccccaaaaddddbbbbaaaaaaaaccccaaaaccccaaaaccccbbbbbbbbaaaaccccbbbbddddaaaaddddbbbbbbbbbbbbccccbbbbaaaaaaaaccccaaaabbbbaaaaccccbbbbccccaaaaddddbbbbbbbbbbbbccccbbbbaaaaaaaaccccaaaaccccaaaaccccbbbbbbbbbbbbccccbbbbaaaaddddccccbbbbddddddddccccbbbbccccddddccccbbbbddddbbbbccccbbbbaaaaaaaaccccaaaaaaaabbbbddddbbbbddddddddccccbbbbaaaaaaaaccccaaaabbbbbbbbddddbbbbddddaaaaddddbbbb",
		},
		{
			[]byte("ls -lart"),
			"aaaaddddccccbbbbddddaaaaddddbbbbaaaaaaaaccccaaaabbbbddddccccaaaaaaaaddddccccbbbbbbbbaaaaccccbbbbccccaaaaddddbbbbaaaabbbbddddbbbb",
		},
		{[]byte{}, ""},
	} {
		/* Make sure enc.Encode works */
		egot := make([]byte, 16*len(c.have))
		e.Encode(egot, c.have)
		if 0 != bytes.Compare(egot, []byte(c.want)) {
			t.Fatalf(
				"Encoding.Encode: have:%02x got:%q want:%q",
				c.have,
				egot,
				c.want,
			)
		}

		/* Make sure encoding works */
		got := e.EncodeToString(c.have)
		if c.want != got {
			t.Fatalf(
				"EncodeToString: have:%02X got:%v want:%v",
				c.have,
				got,
				c.want,
			)
		}
		/* Make sure decoding works */
		dec, err := e.DecodeString(got)
		if nil != err {
			t.Fatalf("DecodeString: %v", err)
		}
		if string(dec) != string(c.have) {
			t.Fatalf(
				"DecodeSTring: enc:%v dec:%02X have:%02X",
				got,
				dec,
				c.have,
			)
		}
	}
}

func ExampleEncoding_EncodeToString() {
	enc := MustNewEncoding("boatfeetbowlsoap")
	s := enc.EncodeToString([]byte("uname -a"))
	fmt.Printf("%v\n", s)

	// Output: feetfeetsoapfeetbowlsoapbowlfeetfeetboatbowlfeetfeetsoapbowlfeetfeetfeetbowlfeetboatboatbowlboatfeetsoapbowlboatfeetboatbowlfeet
}
