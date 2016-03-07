/*
Package ngao (shield) is a reverse proxy that limits the maximum number of connections
to an upstream host.

Usage:

Get the package

   go get github.com/peteretelej/ngao

Import into your application:

   import "github.com/peteretelej/ngao"

Set your configuration details, and run ngao:

   func main() {
      c := &ngao.Config{
         ListenAddr:    ":9015",          // your listen address e.g ":9010"
   	     Host:          "etelej.com", // the backend server to reverseproxy
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
	log.Print("Starting Ngao server")
	log.Printf("Listening on '%s', Serving: %s://%s",
		c.ListenAddr, c.Scheme, c.Host)
	log.Printf("Max sessions: %d. Older sessions cleared every: %d secs",
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
		// Launch session queues manager
		manager()
		wg.Done()
	}()
	wg.Wait()
}
func initConf() {
	t := time.Now().Format("20060102150405")
	i, _ := strconv.Atoi(t)
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
