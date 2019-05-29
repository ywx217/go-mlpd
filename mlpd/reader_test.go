package mlpd

import (
	"bufio"
	"os"
	"testing"
)

func TestMlpdReader_ReadHeader(t *testing.T) {
	f, err := os.Open("output.mlpd")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	r := MlpdReader{
		data: bufio.NewReader(f),
	}
	header, err := r.ReadHeader()
	if err != nil || header == nil {
		t.Error(err)
		return
	}

	// The file is composed by a header followed by 0 or more buffers.

	// Each buffer contains events that happened on a thread: for a given thread
	// buffers that appear later in the file are guaranteed to contain events
	// that happened later in time.
	err = r.ReadBuffer(func(bh *BufferHeader, ev *Event) error {
		t.Logf("Read event data=%v %+v", ev.data.Name(), ev.data)
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
}
