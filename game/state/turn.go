package state

import (
	"github.com/kvnxiao/pictorio/model"
)

type SelectionIndex struct {
	User      model.User
	Timestamp int64
	Value     int
}

type Guess struct {
	User      model.User
	Timestamp int64
	Value     string
}
