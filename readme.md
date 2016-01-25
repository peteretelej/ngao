# Ngao 

Ngao **_(shield)_** - A Golang reverse proxy that limits number of client sessions connected to an upstream server.

Set a number of __maximum client sessions__, and the __clear interval__ (interval in which older sessions are cleared).

## Example

```go
package main

import "github.com/etelej/ngao"

func main() {
	c := &ngao.Config{
		ListenAddr:    ":9015",          // your listen address e.g ":9010"
		Host:          "dev.etelej.com", // the backend server to reverseproxy
		Scheme:        "https",          // protocol scheme of backend host e.g. https, http
		TotalAllowed:  4,                // Maximum client sessions allowed
		ClearInterval: 60 * 5,           // Interval to clear older sessions (secs)
	}
	ngao.Run(c)
}
```