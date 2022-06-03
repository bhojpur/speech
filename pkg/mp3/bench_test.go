package mp3

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	buf, err := ioutil.ReadFile("./classic.mp3")
	if err != nil {
		b.Fatal(err)
	}
	src := bytes.NewReader(buf)
	for i := 0; i < b.N; i++ {
		if _, err := src.Seek(0, io.SeekStart); err != nil {
			b.Fatal(err)
		}
		d, err := NewDecoder(src)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := ioutil.ReadAll(d); err != nil {
			b.Fatal(err)
		}
	}
}
