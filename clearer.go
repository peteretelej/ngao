package ngao

import (
	"errors"
	"sync"
	"time"
)

// conf defines properties used to handle clearing sessions
type conf struct {
	AllowedMap map[int]int
	Allowed    []int
	WaitingMap map[int]string
	Waiting    []int
	sync.RWMutex
}

var nconfig *conf

// Type that handles generation of unique IDs
type uniqData struct {
	uniqID int
	sync.Mutex
}

var uniqgen *uniqData

// getUniqID returns a unique ID for identifying
func getUniqID() (i int) {
	uniqgen.Lock()
	defer uniqgen.Unlock()
	i = uniqgen.uniqID + 1
	uniqgen.uniqID = i
	return i
}

// clearer handles adding waiting users to allowed list
// and displacing previously allowed users
func clearer() {
	for {
		nconfig.RLock()
		ca := nconfig.Allowed
		cw := nconfig.Waiting
		nconfig.RUnlock()
		var half int
		half = (totalAllowed / 2)
		if len(ca) > half {
			todelete := ca[:half]
			for _, val := range todelete {
				nconfig.Lock()
				nconfig.deleteAllowed(val)
				nconfig.Unlock()
			}
			addtotal := 0
			if len(cw) >= half {
				addtotal = half
			} else {
				addtotal = len(cw)
			}
			toadd := cw[:addtotal]
			for _, val := range toadd {
				nconfig.Lock()
				nconfig.deleteWaiting(val)
				nconfig.addAllowed(val)
				nconfig.Unlock()
			}
		}
		time.Sleep(time.Second * time.Duration(clearInterval))
	}
}

// Adds session to waiting queue
func (c *conf) addWaiting(n int) {
	l := len(c.Waiting)
	c.Waiting = append(c.Waiting, n)
	var waittime string
	var t int
	var half int
	half = totalAllowed / 2
	if l < half {
		t = clearInterval
	} else {
		t = (l / half) * clearInterval
	}
	waittime = time.Now().Add(time.Second * time.Duration(t)).Format("3:04PM (Jan 2 2006)")
	c.WaitingMap[n] = waittime
}
func (c *conf) deleteWaiting(n int) {
	delete(c.WaitingMap, n)
	c.Waiting = nil
	for key := range c.WaitingMap {
		c.Waiting = append(c.Waiting, key)
	}
}

func (c *conf) addAllowed(n int) {
	c.Allowed = append(c.Allowed, n)
	c.AllowedMap[n] = 1
}

func (c *conf) deleteAllowed(n int) {
	delete(c.AllowedMap, n)
	c.Allowed = nil
	for key := range c.AllowedMap {
		c.Allowed = append(c.Allowed, key)
	}
}

func isAllowed(n int) (allowed bool, err error) {
	nconfig.Lock()
	defer nconfig.Unlock()

	_, exists := nconfig.AllowedMap[n]
	if exists {
		return true, nil
	}
	_, waiting := nconfig.WaitingMap[n]
	if waiting {
		if len(nconfig.Allowed) < totalAllowed {
			nconfig.AllowedMap[n] = 1
			nconfig.Allowed = append(nconfig.Allowed, n)
			delete(nconfig.WaitingMap, n)
			return true, nil
		}
		return false, nil

	}
	return false, errors.New("Session not waiting")
}
func getWait(n int) string {
	nconfig.RLock()
	defer nconfig.RUnlock()
	return nconfig.WaitingMap[n]
}

func handleNew() (sessid int, allow bool) {
	nconfig.Lock()
	defer nconfig.Unlock()
	n := getUniqID()
	if len(nconfig.Allowed) < totalAllowed {
		nconfig.AllowedMap[n] = 1
		nconfig.Allowed = append(nconfig.Allowed, n)
		return n, true
	}
	nconfig.addWaiting(n)
	return n, false

}
