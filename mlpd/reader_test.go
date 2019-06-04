package mlpd

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestMlpdReader_Read(t *testing.T) {
	f, err := os.Open("output.unity.mlpd")
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
	err = r.ReadBuffer(func(bh *BufferHeader, ev *Event, ver byte) error {
		fmt.Printf("Read event data=%v %+v\n", ev.data.Name(), ev.data)
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMlpdReader_ReadSamples(t *testing.T) {
	f, err := os.Open("output.unity.mlpd")
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
	err = r.ReadBuffer(MakeEventIterFromSampleIter(func(d *SampleData) error {
		fmt.Printf(
			"# %4d-%02d-%02d %02d:%02d:%02d.%06d\n",
			d.time.Year(), d.time.Month(), d.time.Day(),
			d.time.Hour(), d.time.Minute(), d.time.Second(),
			d.time.Nanosecond()/1000,
		)
		for _, node := range d.methods {
			if node == nil {
				fmt.Println("Undefined")
			} else {
				fmt.Println(node.name)
			}
		}
		fmt.Printf("thread-0x%x\n", d.threadID)
		fmt.Printf("-- end of sample\n\n")
		return nil
	}, true))
	if err != nil {
		t.Error(err)
		return
	}
}
