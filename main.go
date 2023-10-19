package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m25n/wip/stack"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const Version = "0.2.0"

func main() {
	wipfile := GetWIPFile()
	wiplog := NewWIPLog(wipfile)
	args := os.Args
	command := "show"
	if len(args) > 1 {
		command = args[1]
	}

	switch command {
	case "push":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: wip push <item>\n")
			os.Exit(1)
		}
		item := strings.Join(args[2:], " ")
		err := wiplog.Push(time.Now(), item)
		if err != nil {
			fmt.Fprintf(os.Stderr, `error pushing item "%s": %s\n`, item, err.Error())
			os.Exit(1)
		}
	case "pop":
		err := wiplog.Pop(time.Now())
		if err != nil {
			fmt.Fprintf(os.Stderr, `error poping item: %s\n`, err.Error())
			os.Exit(1)
		}
	case "show":
		var items stack.Stack[string]
		err := wiplog.Each(func(_ time.Time, item string) {
			items = items.Push(item)
		}, func(_ time.Time) {
			items = items.Pop()
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing wipfile: %s\n", err.Error())
			os.Exit(1)
		}
		if items.Size() == 0 {
			fmt.Println("no items.")
			os.Exit(0)
		}
		i := 0
		for ; items.Size() > 0; items = items.Pop() {
			fmt.Printf("%d: %s\n", i, items.Top())
			i++
		}
	case "stats":
		var times stack.Stack[time.Time]
		var intervals []time.Duration
		err := wiplog.Each(func(at time.Time, _ string) {
			times = times.Push(at)
		}, func(end time.Time) {
			start := times.Top()
			times.Pop()
			intervals = append(intervals, end.Sub(start))
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing wipfile: %s\n", err.Error())
			os.Exit(1)
		}
		sort.Slice(intervals, func(i, j int) bool {
			return intervals[i] < intervals[j]
		})

		if len(intervals) == 0 {
			fmt.Println("no stats.")
			os.Exit(0)
		}
		medianLow := max((len(intervals)/2)-1, 0)
		medianHigh := ((len(intervals) + 1) / 2) - 1
		medianCompletion := (intervals[medianLow] + intervals[medianHigh]) / 2
		maxCompletion := intervals[len(intervals)-1]
		fmt.Printf("median completion time: %s\n", medianCompletion)
		fmt.Printf("max completion time: %s\n", maxCompletion)
	case "version":
		fmt.Println(Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s\n", args[1])
		os.Exit(1)
	}
}

func GetWIPFile() string {
	wipfile := os.Getenv("WIPFILE")
	if wipfile != "" {
		var err error
		wipfile, err = filepath.Abs(wipfile)
		if err != nil {
			panic(err)
		}
		return wipfile
	}
	return DefaultWIPFile()
}

func DefaultWIPFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".wip")
}

type WIPLog struct {
	wipfile string
}

func (wl *WIPLog) Each(onPush func(time.Time, string), onPop func(time.Time)) error {
	fh, err := wl.openReadable()
	if err != nil {
		return err
	}
	defer fh.Close()
	lines := bufio.NewScanner(fh)
	lines.Split(bufio.ScanLines)
	for lines.Scan() {
		line := lines.Bytes()
		var op Op
		err = json.Unmarshal(line, &op)
		if err != nil {
			return err
		}
		if op.Push != nil {
			onPush(op.Push.At, op.Push.Item)
		}
		if op.Pop != nil {
			onPop(op.Pop.At)
		}
	}
	return lines.Err()
}

func (wl *WIPLog) Push(at time.Time, item string) error {
	fh, err := wl.openWritable()
	if err != nil {
		return err
	}
	defer fh.Close()
	op := &Op{Push: &PushOp{At: at, Item: item}}
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(op)
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, buf)
	return err
}

func (wl *WIPLog) Pop(at time.Time) error {
	fh, err := wl.openWritable()
	if err != nil {
		return err
	}
	defer fh.Close()
	op := &Op{Pop: &PopOp{At: at}}
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(op)
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, buf)
	return err
}

type Op struct {
	Push *PushOp `json:"push,omitempty"`
	Pop  *PopOp  `json:"pop,omitempty"`
}

type PushOp struct {
	At   time.Time `json:"at"`
	Item string    `json:"item"`
}

type PopOp struct {
	At time.Time `json:"at"`
}

func (wl *WIPLog) openWritable() (*os.File, error) {
	return os.OpenFile(wl.wipfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func (wl *WIPLog) openReadable() (*os.File, error) {
	return os.OpenFile(wl.wipfile, os.O_CREATE|os.O_RDONLY, 0600)
}

func NewWIPLog(wipfile string) *WIPLog {
	return &WIPLog{wipfile: wipfile}
}
