package ngao

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/context"
)

// rProxy is a single host reverse proxy implementation
// See github.com/etelej/rproxy
func rProxy(listenAddr, remoteHost, scheme string) {
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "sess")
			if err != nil {
				// compromised cookie?
				http.Error(w, err.Error(), 500)
				return
			}
			var allow, resetsess bool
			var wait string
			sid, exists := session.Values["sessid"]
			if exists {
				if val, ok := sid.(int); ok {
					allow, err = isAllowed(val)
					if err != nil {
						// not allowed and not waiting
						delete(session.Values, "sessid")
						resetsess = true
					} else {
						if !allow {
							// client already waiting, get wait time
							wait = getWait(val)
						}
					}
				} else {
					http.Error(w, "User SessionID not int", 500)
					return
				}
			} else {
				// No session ID exists (new user)
				resetsess = true
			}
			// Reset user session, either new user or
			// User with and retired session
			if resetsess {
				id, allowed := handleNew()
				log.Printf("id: %v, allowed=%v", id, allowed)
				session.Values["sessid"] = id
				allow = allowed
				if !allowed {
					// get wait time if new session was set to wait
					wait = getWait(id)
				}
			}
			session.Save(r, w)

			// Handle waiting user, notify when they can access
			if !allow {
				w.Write([]byte("Please try again later at: " + wait))
				return
			}
			p := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: scheme,
				Host:   remoteHost,
			})
			p.ServeHTTP(w, r)

		})
	log.Fatal(http.ListenAndServe(listenAddr, context.ClearHandler(http.DefaultServeMux)))
}
