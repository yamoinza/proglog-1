package log

import (
	"io/ioutil"
	"os"
	"testing"

	req "github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "log_test")
	req.NoError(t, err)
	defer os.Remove(f.Name())

	b := []byte("hello world")
	_, err = f.Write(b)
	req.NoError(t, err)

	l, err := newLog(f)
	req.NoError(t, err)

	size, err := l.Append(b)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Equal(t, size, uint64(len(b)*2))
}
