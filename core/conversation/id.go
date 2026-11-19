package conversation

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func NewConversationID() string {
	return "conv_" + uniqueSuffix()
}

func NewMessageID() string {
	return "msg_" + uniqueSuffix()
}

func uniqueSuffix() string {
	var random [4]byte

	if _, err := rand.Read(random[:]); err != nil {
		return fmt.Sprintf(
			"%d",
			time.Now().UnixNano(),
		)
	}

	return fmt.Sprintf(
		"%d_%s",
		time.Now().UnixMilli(),
		hex.EncodeToString(random[:]),
	)
}
