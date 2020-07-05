# IP2Tor

IP2Tor allows to check if an IP address is a Tor exit node. This can be used to identify and block traffic from the Tor network. IPv4 and IPv6 addresses are supported.

It will run a daemon to download the Tor list from https://www.dan.me.uk/tornodes and update it according to the time specified. It optionally uses a local file to cache the last results and allow continuous operation even after restart of your application.

## Usage

It is a Go package with no external dependencies.

```go
package main

import (
    "github.com/IntelligenceX/ip2tor"
)

func main() {
    // Only download exit nodes, refetch the list every 2 hours.
    // Cache it to the file "tor-ips.txt".
    ip2tor.Init(true, time.Hour*2, "tor-ips.txt")
    
    ip := net.ParseIP("1.2.3.4")

    ip2tor.IsTor(ip) // returns true or false
}
```