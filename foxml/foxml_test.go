package foxml

import (
	"os"
	"testing"
)

func TestDecodeFoxml(t *testing.T) {
	f, err := os.Open("info%3Afedora%2Fvecnet%3Ax920fw85p")
	if err != nil {
		t.Error(err)
	}
	do, err := DecodeFoxml(f)
	if err != nil {
		t.Error(err)
	}
	if do.Version != "1.1" {
		t.Error("got version", do.Version)
	}
	if do.Pid != "vecnet:x920fw85p" {
		t.Error("got PID", do.Pid)
	}
	f.Close()
}

func BenchmarkDecode(b *testing.B) {
	f, err := os.Open("info%3Afedora%2Fvecnet%3Ax920fw85p")
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Seek(0, 0)
		_, err := DecodeFoxml(f)
		if err != nil {
			b.Error(err)
		}
	}
	f.Close()
}
