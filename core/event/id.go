package event

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func NewID() string {
	var random [4]byte

	if _, err := rand.Read(random[:]); err != nil {
		return fmt.Sprintf(
			"evt_%d",
			time.Now().UnixNano(),
		)
	}

	return fmt.Sprintf(
		"evt_%d_%s",
		time.Now().UnixMilli(),
		hex.EncodeToString(random[:]),
	)
}
