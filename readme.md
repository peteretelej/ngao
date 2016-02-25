# Ngao 

Ngao **_(shield)_** - A Golang reverse proxy that limits number of client sessions connected to an upstream server.

Set a number of __maximum client sessions__, and the __clear interval__ (interval in which older sessions are cleared).

Ngao is useful if you have want to rate limit access to your application/backend server, for example the configuration below only allows 20 sessions at a time.

## Example

```go
package main

import "github.com/etelej/ngao"

func main() {
	c := &ngao.Config{
		ListenAddr:    ":9015",          // your listen address e.g ":9010"
		Host:          "etelej.com", // the backend server to reverseproxy
		Scheme:        "https",          // protocol scheme of backend host e.g. https, http
		TotalAllowed:  20,                // Maximum client sessions allowed
		ClearInterval: 60 * 5,           // Interval to clear older sessions (secs)
	}
	ngao.Run(c)
}
```


Contributions welcome & greatly appreciated.

