package wiplog

import (
	"encoding/json"
	"time"

	"github.com/m25n/wip/wipfile"
)

type WIPLog struct {
	wipfile wipfile.WIPFile
}

func New(wipfile wipfile.WIPFile) *WIPLog {
	return &WIPLog{wipfile: wipfile}
}

type Handlers struct {
	OnPush func(time.Time, string)
	OnPop  func(time.Time)
	OnMove func(time.Time, int, int)
}

func (wl *WIPLog) Each(h Handlers) error {
	return wl.wipfile.Lines(func(line []byte) error {
		var o op
		if err := json.Unmarshal(line, &o); err != nil {
			return err
		}
		if o.Push != nil {
			h.OnPush(o.Push.At, o.Push.Item)
		}
		if o.Pop != nil {
			h.OnPop(o.Pop.At)
		}
		if o.Move != nil {
			h.OnMove(o.Move.At, o.Move.From, o.Move.To)
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

func (wl *WIPLog) Move(at time.Time, from, to int) error {
	return wl.writeOp(&op{Move: &moveOp{At: at, From: from, To: to}})
}

func (wl *WIPLog) writeOp(o *op) error {
	line, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return wl.wipfile.AppendLine(line)
}
