package wiplog

import (
	"encoding/json"
	"github.com/m25n/wip/wipfile"
	"time"
)

type WIPLog struct {
	wipfile wipfile.WIPFile
}

func New(wipfile wipfile.WIPFile) *WIPLog {
	return &WIPLog{wipfile: wipfile}
}

func (wl *WIPLog) Each(onPush func(time.Time, string), onPop func(time.Time)) error {
	return wl.wipfile.Lines(func(line []byte) error {
		var o op
		if err := json.Unmarshal(line, &o); err != nil {
			return err
		}
		if o.Push != nil {
			onPush(o.Push.At, o.Push.Item)
		}
		if o.Pop != nil {
			onPop(o.Pop.At)
		}
		return nil
	})
}

func (wl *WIPLog) Push(at time.Time, item string) error {
	return wl.writeOp(&op{Push: &pushOp{At: at, Item: item}})
}

func (wl *WIPLog) Pop(at time.Time) error {
	return wl.writeOp(&op{Pop: &popOp{At: at}})
}

func (wl *WIPLog) writeOp(o *op) error {
	line, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return wl.wipfile.AppendLine(line)
}
