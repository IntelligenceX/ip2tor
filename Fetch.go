/*
File Name:  Fetch.go
Copyright:  2020 Kleissner Investments s.r.o.
Author:     Peter Kleissner
*/

package ip2tor

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// URL for the official Tor exit node list. This list seems to match only about 72% IPs of the dan.me.uk one.
const linkTorIPsOfficial = "https://check.torproject.org/torbulkexitlist"

// URLs for downloading the Tor node list, full and exit-only.
// These URLs shall not be downloaded more often than 30 minutes according to the website, otherwise it risks being blocked.
const linkTorIPsFull = "https://www.dan.me.uk/torlist/"
const linkTorIPsExit = "https://www.dan.me.uk/torlist/?exit"

// startDownloadDaemon starts the download daemon according to the parameters
// filename defines the file to use as cache. Empty to disable.
func startDownloadDaemon(exitOnly bool, waitTime time.Duration, filename string, fetchImmediately bool) {
	if fetchImmediately {
		fetchTorLists(exitOnly, filename)
	}

	go func() {
		for {
			time.Sleep(waitTime)

			fetchTorLists(exitOnly, filename)
		}
	}()
}

func fetchTorLists(exitOnly bool, filename string) {
	ipMap := make(map[string]struct{})

	// first download the official one
	err1 := downloadTorList(&ipMap, linkTorIPsOfficial)

	// second source
	dlLink := linkTorIPsFull
	if exitOnly {
		dlLink = linkTorIPsExit
	}
	err2 := downloadTorList(&ipMap, dlLink)

	// in case any of the sources fail, re-use the old list
	if err1 != nil || err2 != nil {
		for ipA := range torIPs {
			ipMap[ipA] = struct{}{}
		}
	}

	// update the live map
	torIPs = ipMap

	// write out cache file
	storeCacheFile(ipMap, filename)
}

// downloadTorList downloads the IP list and applies it to the map.
func downloadTorList(ipMap *map[string]struct{}, link string) (err error) {
	data, err := downloadLink(link)
	if err != nil {
		return err
	}

	records := strings.Fields(string(data))

	for _, record := range records {
		ip := net.ParseIP(record)
		if ip == nil {
			continue
		}

		(*ipMap)[ip.String()] = struct{}{}
	}

	return nil
}

func downloadLink(link string) (buffer []byte, err error) {
	r, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, errors.New("Download failed")
	}

	return ioutil.ReadAll(r.Body)
}

// storeCacheFile stores the input map to the file specified
func storeCacheFile(ipMap map[string]struct{}, filename string) {
	if filename == "" {
		return
	}

	var data string

	for ipA := range ipMap {
		data += ipA + "\n"
	}

	ioutil.WriteFile(filename, []byte(data), 0644)
}

// readCacheFile reads the cache file and processes it
func readCacheFile(filename string) (ipMap map[string]struct{}, valid bool) {
	ipMap = make(map[string]struct{})

	if filename == "" {
		return ipMap, false
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ipMap, false
	}

	records := strings.Fields(string(data))

	for _, record := range records {
		ip := net.ParseIP(record)
		if ip == nil {
			continue
		}

		ipMap[ip.String()] = struct{}{}
	}

	return ipMap, true
}

// startFileCacheFetcher starts a Go routine to continuously reload the file cache
func startFileCacheFetcher(waitTime time.Duration, filename string) {
	if filename == "" {
		return
	}

	go func() {
		for {
			time.Sleep(waitTime)

			if ipMap, valid := readCacheFile(filename); valid {
				torIPs = ipMap
			}
		}
	}()
}
