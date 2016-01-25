/*
Package ngao (shield) is a reverse proxy that limits the maximum number of connections
to an upstream host.

Usage:
Get the package
   go get github.com/etelej/ngao

Import into your application:

   import "github.com/etelej/ngao"

Set your configuration details, and run ngao:

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
   ngao.Run(c)

ngao was written by Peter Etelej <peter@etelej.com>
*/
package ngao

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

// Config allows configuring of ngao's options
type Config struct {
	ListenAddr    string
	Host          string
	Scheme        string
	TotalAllowed  int
	ClearInterval int
}

var store = sessions.NewCookieStore([]byte("ngao"))

var totalAllowed, clearInterval int

// Run launches ngao reverse proxy
func Run(c *Config) {
	log.Printf("Starting Ngao server: Listening on '%s', Serving: %s on scheme %s",
		c.ListenAddr, c.Host, c.Scheme)
	log.Printf("Max client sessions: %d. Clear older sessions interval: %d secs",
		c.TotalAllowed, c.ClearInterval)

	totalAllowed = c.TotalAllowed
	clearInterval = c.ClearInterval
	initConf()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		rProxy(c.ListenAddr, c.Host, c.Scheme)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		clearer()
		wg.Done()
	}()
	wg.Wait()
}
func initConf() {
	t := time.Now().Format("20060102150405")
	i, err := strconv.Atoi(t)
	if err != nil {
		log.Fatal(err.Error())
	}
	uniqgen = &uniqData{
		uniqID: i,
	}

	nconfig = &conf{
		AllowedMap: make(map[int]int),
		Allowed:    make([]int, 0),
		WaitingMap: make(map[int]string),
		Waiting:    make([]int, 0),
	}

}
