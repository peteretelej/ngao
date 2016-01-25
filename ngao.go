/*
Package ngao is a reverse proxy that limits the maximum number of connections
to an upstream host.


*/
package ngao

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/sessions"
)

type Config struct {
	ListenAddr    string // your listen address e.g ":9010"
	Host          string // the backend server to reverseproxy
	Scheme        string // protocol scheme of host e.g. https, http
	TotalAllowed  int    // Total client sessions allowed
	ClearInterval int    // Interval to clear older sessions (secs)
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
