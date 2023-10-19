package main

import (
	"fmt"
	"github.com/m25n/wip/stack"
	"github.com/m25n/wip/wiplog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const Version = "0.2.0"

func main() {
	wipfile := GetWIPFile()
	wl := wiplog.New(wipfile)
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
		err := wl.Push(time.Now(), item)
		if err != nil {
			fmt.Fprintf(os.Stderr, `error pushing item "%s": %s\n`, item, err.Error())
			os.Exit(1)
		}
	case "pop":
		err := wl.Pop(time.Now())
		if err != nil {
			fmt.Fprintf(os.Stderr, `error poping item: %s\n`, err.Error())
			os.Exit(1)
		}
	case "show":
		var items stack.Stack[string]
		err := wl.Each(func(_ time.Time, item string) {
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
		err := wl.Each(func(at time.Time, _ string) {
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
