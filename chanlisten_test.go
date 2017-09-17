package chanlisten

import (
	"io"
	"net"
	"testing"
)

func TestChanListener(t *testing.T) {
	var c net.Conn

	l := New(0)
	go func(list net.Listener) {
		for {
			_, err := list.Accept()
			if err != nil {
				if err != io.EOF {
					t.Fatalf("got error %q, want %q", err, io.EOF)
				}
				t.Logf("got io.EOF")
				break
			}
			t.Logf("got a conn")
		}
	}(l)
	if err := l.Queue(c); err != nil {
		t.Fatalf("got unexpected error sending conn: %q", err)
	}
	l.Close()
	if err := l.Queue(c); err == nil {
		t.Fatalf("got nil error, expected io.EOF")
	} else {
		if err != io.EOF {
			t.Fatalf("got err %q, want io.EOF", err)
		}
	}
}
