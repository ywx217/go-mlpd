package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/ianlancetaylor/demangle"

	"github.com/ywx217/d3-flame-server/flamewriter"

	"github.com/ywx217/go-mlpd/mlpd"
)

var (
	inputPath     string
	outputPath    string
	outputType    string
	cutThreshold  int
	includeNative bool
	splitThreads  bool
	skipIdle      bool
	debug         bool
)

var (
	idleSymbols = map[string]bool{
		"pthread_cond_wait":      true,
		"pthread_cond_timedwait": true,
	}
)

// CumulativeItem cumulative item for visualize
type CumulativeItem struct {
	Value int `json:"v"`
	Count int `json:"c"`
}

// TimeStatisticItem sample time statistic item
type TimeStatisticItem struct {
	Time     string `json:"date"`
	ThreadID uint64 `json:"tid"`
	Count    int    `json:"c"`
}

func debugOutputData(times map[string]map[uint64]int, record *flamewriter.Record) {
	rf, err := os.Create("debug.html")
	if err != nil {
		panic(err)
	}
	defer rf.Close()

	values := make([]CumulativeItem, 0)
	{
		m := make(map[int]int, 0)
		record.ValueStatisticInplace(m)
		for value, count := range m {
			values = append(values, CumulativeItem{
				Value: value,
				Count: count,
			})
		}
	}
	bsCumulative, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	timeStat := make([]TimeStatisticItem, 0)
	{
		for t, tc := range times {
			for tid, count := range tc {
				timeStat = append(timeStat, TimeStatisticItem{
					Time:     t,
					ThreadID: tid,
					Count:    count,
				})
			}
		}
	}
	bsTimes, err := json.Marshal(timeStat)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(rf, tmplDebugReport, string(bsCumulative), string(bsTimes))
}

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
	// start parsing flame data
	record := flamewriter.NewRecord("root", 0)
	stack := make([]string, 0, 100)
	times := make(map[string]map[uint64]int, 0)
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
		if debug {
			timeStr := d.Time().Format("Mon, 02 Jan 2006 15:04:05")
			tc, ok := times[timeStr]
			if !ok {
				tc = make(map[uint64]int, 0)
				times[timeStr] = tc
			}
			tc[d.ThreadID()]++
		}
		return nil
	}, includeNative))
	if err != nil {
		return err
	}
	writer := flamewriter.NewHTMLWriter(f)
	if cutThreshold > 1 {
		record.FixRootValue().CutoffInplace(cutThreshold).ReduceRoot()
	} else {
		record.FixRootValue().ReduceRoot()
	}
	if debug {
		debugOutputData(times, record)
	}
	return writer.Write(record)
}

func outputTimeline(r *mlpd.MlpdReader) error {
	return nil
}

func makeProgressReader(path string) (*os.File, *mlpd.MlpdReader, *pb.ProgressBar, error) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return nil, nil, nil, err
	}
	bar := pb.New64(fInfo.Size()).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	return f, mlpd.NewReader(bufio.NewReader(io.TeeReader(f, bar))), bar, nil
}

func main() {
	flag.StringVar(&inputPath, "i", "output.mlpd", "input file path in mlpd format")
	flag.StringVar(&outputPath, "o", "output.html", "output file path in html formatj")
	flag.StringVar(&outputType, "t", "flame", "output type, supported: flame, timeline")
	flag.IntVar(&cutThreshold, "cut", 1, "cut nodes less than the threshold")
	flag.BoolVar(&includeNative, "show-natives", false, "includes native only stacks")
	flag.BoolVar(&splitThreads, "split-threads", false, "split thread as different root nodes")
	flag.BoolVar(&skipIdle, "skip-idle", false, "skips idle symbols")
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()

	f, r, bar, err := makeProgressReader(inputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if outputType == "flame" {
		bar.Start()
		defer bar.Finish()
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
