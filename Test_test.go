package ip2tor

import (
	"testing"
	"time"
)

// test code for manual debugging
func Test1(t *testing.T) {
	Init(false, time.Hour*2, "tor-ips.txt")
}