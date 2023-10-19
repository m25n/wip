package main

import (
	"fmt"
	"github.com/m25n/wip/wipfile"
	"github.com/m25n/wip/wiplog"
	"os"
)

const Version = "0.2.0"

func main() {
	wf, err := wipfile.FromEnv()
	HandleErr(err)
	wl := wiplog.New(wf)
	args := os.Args
	command := "show"
	if len(args) > 1 {
		command = args[1]
	}

	switch command {
	case "push":
		HandleErr(PushItem(wl, args[2:]))
	case "pop":
		HandleErr(PopItem(wl))
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

func HandleErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
