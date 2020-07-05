/*
File Name:  IP.go
Copyright:  2020 Kleissner Investments s.r.o.
Author:     Peter Kleissner
*/

package ip2tor

import (
	"io/ioutil"
	"net"
	"sync"
	"time"
)

var mapMutex sync.Mutex

var torIPs map[string]struct{}

func updateMap(newMap map[string]struct{}) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	torIPs = newMap
}

// Init starts the download daemon and optionally reads the cache file, if specified
func Init(ExitOnly bool, waitTime time.Duration, filename string) {
	torIPs = make(map[string]struct{})
	useCache := false

	if filename != "" {
		if data, err := ioutil.ReadFile(filename); err == nil {
			processTorList(data)
			useCache = true
		}
	}

	go startDownloadDaemon(ExitOnly, waitTime, filename, !useCache)
}

// IsTor checks if an IP address is listed as Tor IP
func IsTor(IP net.IP) bool {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	_, ok := torIPs[IP.String()]

	return ok
}
