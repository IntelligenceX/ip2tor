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
// Mode: 0 = disabled (no IP check), 1 = active (ban exit nodes only), 2 = active (ban all nodes), 3 = active, no fetching (only use file cache)
func Init(mode int, waitTime time.Duration, filename string) {
	if mode == 0 { // disabled?
		torIPs = make(map[string]struct{})
		return
	}

	var useCache bool
	torIPs, useCache = readCacheFile(filename)

	if mode == 3 { // only use file cache?
		startFileCacheFetcher(waitTime, filename)
		return
	}

	startDownloadDaemon(mode == 1, waitTime, filename, !useCache)
}

// IsTor checks if an IP address is listed as Tor IP
func IsTor(IP net.IP) bool {
	if IP == nil { // invalid input?
		return false
	}

	_, ok := torIPs[IP.String()]

	return ok
}
