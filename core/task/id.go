package task

import (
	"fmt"
	"sync/atomic"
	"time"
)

var sequence atomic.Uint64

func NewID() string {
	now := time.Now()

	n := sequence.Add(1)

	return fmt.Sprintf(
		"task_%s_%03d",
		now.Format("20060102"),
		n,
	)
}
