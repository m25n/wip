package wiplog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"
)

type WIPLog struct {
	wipfile string
}

func New(wipfile string) *WIPLog {
	return &WIPLog{wipfile: wipfile}
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
		var o op
		err = json.Unmarshal(line, &o)
		if err != nil {
			return err
		}
		if o.Push != nil {
			onPush(o.Push.At, o.Push.Item)
		}
		if o.Pop != nil {
			onPop(o.Pop.At)
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
	o := &op{Push: &pushOp{At: at, Item: item}}
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(o)
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
	o := &op{Pop: &popOp{At: at}}
	buf := bytes.NewBuffer(nil)
	err = json.NewEncoder(buf).Encode(o)
	if err != nil {
		return err
	}
	_, err = io.Copy(fh, buf)
	return err
}

func (wl *WIPLog) openWritable() (*os.File, error) {
	return os.OpenFile(wl.wipfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
}

func (wl *WIPLog) openReadable() (*os.File, error) {
	return os.OpenFile(wl.wipfile, os.O_CREATE|os.O_RDONLY, 0600)
}
