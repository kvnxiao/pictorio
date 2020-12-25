package hint

import (
	"github.com/kvnxiao/pictorio/model"
)

type Hint struct {
	hints      []model.Hint
	hintsGiven int
	timings    []int
	maxToGive  int
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func NewHint(hints []model.Hint, timings []int) *Hint {
	return &Hint{
		hints:      hints,
		hintsGiven: 0,
		timings:    timings,
		maxToGive:  min(len(timings), len(hints)),
	}
}

func (h *Hint) NextHint(timeLeftSeconds int) (model.Hint, bool) {
	if h.hintsGiven < h.maxToGive && timeLeftSeconds == h.timings[0] {
		h.hintsGiven += 1

		// pop next hint
		nextHint := h.hints[0]
		h.hints = h.hints[1:]
		h.timings = h.timings[1:]

		return nextHint, true
	}
	return model.Hint{}, false
}
