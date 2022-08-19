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
