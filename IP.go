/*
File Name:  IP.go
Copyright:  2020 Kleissner Investments s.r.o.
Author:     Peter Kleissner
*/

package ip2tor

import (
	"net"
	"time"
)

var torIPs map[string]struct{}

// Init starts the download daemon and optionally reads the cache file, if specified
func Init(ExitOnly bool, waitTime time.Duration, filename string) {
	var useCache bool
	torIPs, useCache = readCacheFile(filename)

	startDownloadDaemon(ExitOnly, waitTime, filename, !useCache)
}

// IsTor checks if an IP address is listed as Tor IP
func IsTor(IP net.IP) bool {
	_, ok := torIPs[IP.String()]

	return ok
}
