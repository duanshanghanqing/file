package util

import (
	"time"
)

func CurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
