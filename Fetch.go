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

// URLs for downloading the Tor node list, full and exit-only.
// These URLs shall not be downloaded more often than 30 minutes according to the website, otherwise it risks being blocked.
const linkTorIPsFull = "https://www.dan.me.uk/torlist/"
const linkTorIPsExit = "https://www.dan.me.uk/torlist/?exit"

// startDownloadDaemon starts the download daemon according to the parameters
// filename defines the file to use as cache. Empty to disable.
func startDownloadDaemon(exitOnly bool, waitTime time.Duration, filename string, fetchImmediately bool) {
	dlLink := linkTorIPsFull
	if exitOnly {
		dlLink = linkTorIPsExit
	}

	if fetchImmediately {
		downloadTorList(dlLink, filename)
	}

	for {
		time.Sleep(waitTime)

		downloadTorList(dlLink, filename)
	}
}

func downloadTorList(link, filename string) (err error) {
	data, err := downloadLink(link)
	if err != nil {
		return err
	}

	// store as file?
	if filename != "" {
		ioutil.WriteFile(filename, data, 0644)
	}

	processTorList(data)

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

// processTorList processes a list of Tor nodes
func processTorList(data []byte) {
	records := strings.Fields(string(data))

	ipMap := make(map[string]struct{})

	for _, record := range records {
		ip := net.ParseIP(record)
		if ip == nil {
			continue
		}

		ipMap[ip.String()] = struct{}{}
	}

	updateMap(ipMap)
}
