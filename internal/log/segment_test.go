package log

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	req "github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {
	dir, _ := ioutil.TempDir("", "segment-test")
	defer os.RemoveAll(dir)

	want := []byte("hello world")

	c := Config{}
	c.Segment.MaxLogBytes = 1024
	c.Segment.MaxIndexBytes = entWidth * 3

	s, err := newSegment(dir, 16, c)
	req.NoError(t, err)
	req.Equal(t, uint64(16), s.nextOffset, s.nextOffset)
	req.False(t, s.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := s.Append(want)
		req.NoError(t, err)
		req.Equal(t, 16+i, off)

		got, err := s.Read(off)
		req.NoError(t, err)
		req.Equal(t, want, got)
	}

	_, err = s.Append(want)
	req.Equal(t, io.EOF, err)

	// is maxed when index is max size
	req.True(t, s.IsMaxed())

	// is maxed when log is max size
	c.Segment.MaxLogBytes = uint64(len(want) * 3)
	c.Segment.MaxIndexBytes = 1024

	s, err = newSegment(dir, 16, c)
	req.NoError(t, err)
	req.True(t, s.IsMaxed())
}