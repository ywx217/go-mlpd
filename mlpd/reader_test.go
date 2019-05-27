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
	if err != nil {
		t.Error(err)
		return
	}
	if header == nil {
		t.Error(header)
		return
	}
}
