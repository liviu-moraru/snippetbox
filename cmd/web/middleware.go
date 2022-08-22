package main

import (
	"fmt"
	"net/http"
)

func secureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability. You don't
		// need to do this in your own code.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (sr *StatusRecorder) WriteHeader(status int) {
	sr.Status = status
	sr.ResponseWriter.WriteHeader(status)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sr := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		rd := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(sr, r)
		app.InfoLog.Printf("%s Response status: %d", rd, sr.Status)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				//	Setting the
				//	header on the response acts as a trigger to make Go’s
				//Connection: Close
				//	HTTP server automatically close the current connection after a response has been sent. It
				//	also informs the user that the connection will be closed. Note: If the protocol being used is
				//	HTTP/2, Go will automatically strip the
				//	header from the response (so it
				//Connection: Close
				//	is not malformed) and send a
				//	GOAWAY frame.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
