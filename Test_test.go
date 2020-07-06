package ip2tor

import (
	"testing"
	"time"
)

// test code for manual debugging
func Test1(t *testing.T) {
	Init(1, time.Hour*2, "tor-ips.txt")

	select {}
}
