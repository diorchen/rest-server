// middleware for groceryItemStore

package middleware

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/diorchen/rest-server/internal/authdb"
)

// Logging information for each request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { // intercepts incoming request and performs desired operations before calling next handler
		start := time.Now() // captures start time
		next.ServeHTTP(w, req) // pass control to next handler in middleware chain, allowing it to continue processing request
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start)) // logs HTTP method, requestURI, and elapsed time since start time after next handler concludes
	})
}

// PanicRecovery from panics in 'next'
// returns a StatusInternalError to client
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {  // intercepts request
		defer func() {
			if err := recover(); err != nil { // if recover() called in deferred func, captures value passed to panic(), if no panic, then recover() returns nil
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // if panic, generates HTTP error response
				log.Println(string(debug.Stack())) // logs stack trace of goroutine that panicked (logs details)
			}
		}()
		next.ServeHTTP(w, req) // calls next handler
	})
}

// UserContextKey is the key in a request's context used to check if the request
// has an authenticated user. The middleware will set the value of this key to
// the username, if the user was properly authenticated with a password.
const UserContextKey = "user"

// BasicAuth is middleware that verifies the request has appropriate basic auth
// set up with a user:password pair verified by authdb.
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		if ok && authdb.VerifyUserPass(user, pass) {
			newctx := context.WithValue(req.Context(), UserContextKey, user)
			next.ServeHTTP(w, req.WithContext(newctx))
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}