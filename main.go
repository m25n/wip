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
		HandleErr(PushItem(wl, args[2:]))
	case "pop":
		HandleErr(wl.Pop(time.Now()))
	case "show":
		HandleErr(ShowStack(wl))
	case "stats":
		HandleErr(ShowStats(wl))
	case "version":
		fmt.Println(Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %s\n", args[1])
		os.Exit(1)
	}
}

func PushItem(wl *wiplog.WIPLog, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: wip push <item>\n")
	}
	item := strings.Join(args, " ")
	err := wl.Push(time.Now(), item)
	if err != nil {
		return fmt.Errorf(`error pushing item "%s": %s\n`, item, err.Error())
	}
	return nil
}

func HandleErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func ShowStats(wl *wiplog.WIPLog) error {
	completions, err := ComputeStats(wl)
	if err != nil {
		return err
	}
	if len(completions) == 0 {
		fmt.Println("no stats.")
		return nil
	}
	medianLow := max((len(completions)/2)-1, 0)
	medianHigh := ((len(completions) + 1) / 2) - 1
	medianCompletion := (completions[medianLow] + completions[medianHigh]) / 2
	maxCompletion := completions[len(completions)-1]
	fmt.Printf("median completion time: %s\n", medianCompletion.Truncate(time.Second))
	fmt.Printf("max completion time: %s\n", maxCompletion.Truncate(time.Second))
	return nil
}

func ComputeStats(wl *wiplog.WIPLog) ([]time.Duration, error) {
	var times stack.Stack[time.Time]
	var completions []time.Duration
	err := wl.Each(func(at time.Time, _ string) {
		times = times.Push(at)
	}, func(end time.Time) {
		start := times.Top()
		times = times.Pop()
		completions = append(completions, end.Sub(start))
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing wipfile: %s\n", err.Error())
	}
	sort.Slice(completions, func(i, j int) bool {
		return completions[i] < completions[j]
	})
	return completions, nil
}

func ShowStack(wl *wiplog.WIPLog) error {
	var items stack.Stack[string]
	err := wl.Each(func(_ time.Time, item string) {
		items = items.Push(item)
	}, func(_ time.Time) {
		items = items.Pop()
	})
	if err != nil {
		return fmt.Errorf("error parsing wipfile: %s\n", err.Error())
	}
	if items.Size() == 0 {
		fmt.Println("no items.")
		return nil
	}
	var i int
	for items.Size() > 0 {
		fmt.Printf("%d: %s\n", i, items.Top())
		items = items.Pop()
		i++
	}
	return nil
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
