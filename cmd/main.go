package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/ianlancetaylor/demangle"

	"github.com/ywx217/d3-flame-server/flamewriter"

	"github.com/ywx217/go-mlpd/mlpd"
)

var (
	inputPath     string
	outputPath    string
	outputType    string
	includeNative bool
	splitThreads  bool
	skipIdle      bool
)

var (
	idleSymbols = map[string]bool{
		"pthread_cond_wait":      true,
		"pthread_cond_timedwait": true,
	}
)

func outputFlame(r *mlpd.MlpdReader) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = r.ReadHeader()
	if err != nil {
		return err
	}
	record := flamewriter.NewRecord("root", 0)
	stack := make([]string, 0, 100)
	err = r.ReadBuffer(mlpd.MakeEventIterFromSampleIter(func(d *mlpd.SampleData) error {
		stack = stack[:0]
		// thread name
		if splitThreads {
			stack = append(stack, fmt.Sprintf("thread-0x%x", d.ThreadID()))
		}
		// method from bottom to top
		methods := d.Methods()
		for i := len(methods) - 1; i >= 0; i-- {
			name := methods[i].Name()
			if i == 0 && name[:2] == "_Z" {
				if nn, err := demangle.ToString(name); err == nil {
					name = nn
				}
			}
			stack = append(stack, name)
		}
		if skipIdle && len(stack) > 0 {
			if _, exist := idleSymbols[stack[len(stack)-1]]; exist {
				return nil
			}
		}
		// add to flame record
		record.Add(stack, 1)
		return nil
	}, includeNative))
	if err != nil {
		return err
	}
	writer := flamewriter.NewHTMLWriter(f)
	return writer.Write(record.FixRootValue().ReduceRoot())
}

func outputTimeline(r *mlpd.MlpdReader) error {
	return nil
}

func main() {
	flag.StringVar(&inputPath, "i", "output.mlpd", "input file path in mlpd format")
	flag.StringVar(&outputPath, "o", "output.html", "output file path in html formatj")
	flag.StringVar(&outputType, "t", "flame", "output type, supported: flame, timeline")
	flag.BoolVar(&includeNative, "n", false, "includes native only stacks")
	flag.BoolVar(&splitThreads, "s", false, "split thread as different root nodes")
	flag.BoolVar(&skipIdle, "skip-idle", false, "skips idle symbols")
	flag.Parse()

	f, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := mlpd.NewReader(bufio.NewReader(f))
	if outputType == "flame" {
		err = outputFlame(r)
	} else if outputType == "timeline" {
		err = outputTimeline(r)
	} else {
		panic(fmt.Sprintf("unsupported output type: %s", outputType))
	}

	if err != nil {
		panic(err)
	}
}