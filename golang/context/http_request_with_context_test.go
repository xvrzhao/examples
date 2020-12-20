package context

import (
	"testing"
	"time"
)

func TestHTTPRequestWithContext(t *testing.T) {
	// attempt to adjust ctxTimeout to less than 2 seconds
	HTTPRequestWithContext(9999, time.Second*3)
}
