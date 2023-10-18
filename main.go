package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		err := wiplog.Push(item)
		if err != nil {
			fmt.Fprintf(os.Stderr, `error pushing item "%s": %s\n`, item, err.Error())
			os.Exit(1)
		}
	case "pop":
		err := wiplog.Pop()
		if err != nil {
			fmt.Fprintf(os.Stderr, `error poping item: %s\n`, err.Error())
			os.Exit(1)
		}
	case "show":
		var stack []string
		err := wiplog.Each(func(pushedItem string) {
			stack = append(stack, pushedItem)
		}, func() {
			if len(stack) == 0 {
				return
			}
			stack = stack[:len(stack)-1]
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing wipfile: %s\n", err.Error())
			os.Exit(1)
		}
		for i := 0; i < len(stack); i++ {
			fmt.Printf("%d: %s\n", i, stack[len(stack)-1-i])
		}
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

func (wl *WIPLog) Each(onPush func(string), onPop func()) error {
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
			onPush(op.Push.Item)
		}
		if op.Pop != nil {
			onPop()
		}
	}
	return lines.Err()
}

func (wl *WIPLog) Push(item string) error {
	fh, err := wl.openWritable()
	if err != nil {
		return err
	}
	defer fh.Close()
	op := &Op{Push: &PushOp{Item: item}}
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(op)
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, buf)
	return err
}

func (wl *WIPLog) Pop() error {
	fh, err := wl.openWritable()
	if err != nil {
		return err
	}
	defer fh.Close()
	op := &Op{Pop: &PopOp{}}
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
	Item string `json:"item"`
}

type PopOp struct{}

func (wl *WIPLog) openWritable() (*os.File, error) {
	return os.OpenFile(wl.wipfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func (wl *WIPLog) openReadable() (*os.File, error) {
	return os.OpenFile(wl.wipfile, os.O_CREATE|os.O_RDONLY, 0600)
}

func NewWIPLog(wipfile string) *WIPLog {
	return &WIPLog{wipfile: wipfile}
}
