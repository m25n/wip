package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/m25n/wip/stack"
	"github.com/m25n/wip/wiplog"
)

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

func PopItem(wl *wiplog.WIPLog) error {
	return wl.Pop(time.Now())
}

func MoveItem(wl *wiplog.WIPLog, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: wip move <from index> <to index>\n")
	}
	fromIndex, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf(`error parsing from index "%s": %s`, args[0], err.Error())
	}
	toIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf(`error parsing to index "%s": %s`, args[1], err.Error())
	}
	size, err := computeSize(wl, err)
	if err != nil {
		return err
	}
	err = checkIndexBounds("from", fromIndex, size)
	if err != nil {
		return err
	}
	err = checkIndexBounds("to", toIndex, size)
	if err != nil {
		return err
	}
	err = wl.Move(time.Now(), fromIndex, toIndex)
	if err != nil {
		return fmt.Errorf(`error moving item from %d to %d: %s\n`, fromIndex, toIndex, err.Error())
	}
	return nil
}

func computeSize(wl *wiplog.WIPLog, err error) (int, error) {
	var size int
	err = wl.Each(wiplog.Handlers{
		OnPush: func(_ time.Time, item string) {
			size++
		},
		OnPop: func(_ time.Time) {
			size--
		},
		OnMove: func(_ time.Time, _ int, _ int) {},
	})
	if err != nil {
		return 0, fmt.Errorf("error parsing wipfile: %s\n", err.Error())
	}
	return size, nil
}

func checkIndexBounds(indexType string, index int, size int) error {
	if index < 0 {
		return fmt.Errorf(`cannot have negative %s index (%d)`, indexType, index)
	}
	if index >= size {
		return fmt.Errorf(`%s index is too big (%d must be less than %d)`, indexType, index, size)
	}
	return nil
}

func ShowStack(wl *wiplog.WIPLog) error {
	var items stack.Stack[string]
	err := wl.Each(wiplog.Handlers{
		OnPush: func(_ time.Time, item string) {
			items = items.Push(item)
		},
		OnPop: func(_ time.Time) {
			items = items.Pop()
		},
		OnMove: func(t time.Time, from int, to int) {
			items = items.Move(from, to)
		},
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
	err := wl.Each(
		wiplog.Handlers{
			OnPush: func(at time.Time, _ string) {
				times = times.Push(at)
			},
			OnPop: func(end time.Time) {
				start := times.Top()
				times = times.Pop()
				completions = append(completions, end.Sub(start))
			},
			OnMove: func(t time.Time, from int, to int) {
				times = times.Move(from, to)
			},
		})
	if err != nil {
		return nil, fmt.Errorf("error parsing wipfile: %s\n", err.Error())
	}
	sort.Slice(completions, func(i, j int) bool {
		return completions[i] < completions[j]
	})
	return completions, nil
}
