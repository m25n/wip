package wiplog

import "time"

type op struct {
	Push *pushOp `json:"push,omitempty"`
	Pop  *popOp  `json:"pop,omitempty"`
}

type pushOp struct {
	At   time.Time `json:"at"`
	Item string    `json:"item"`
}

type popOp struct {
	At time.Time `json:"at"`
}
