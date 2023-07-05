// middleware for groceryItemStore

package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
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