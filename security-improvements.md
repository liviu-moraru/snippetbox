# Chapter 10. Security improvements

## Configuring HTTP settings

- Go supports a few elliptic curves, but as of Go 1.18 only ###tls.CurveP256 and tls.X25519 have assembly implementations. The others are very CPU intensive, so omitting them helps ensure that our server will remain performant under heavy loads.
- The full range of cipher suites that Go supports are defined in the [crypto/tls](https://pkg.go.dev/crypto/tls#pkg-constants) package constants.
- Mozilla’s [recommended configurations](https://wiki.mozilla.org/Security/Server_Side_TLS) for modern, intermediate and old browsers
- If a TLS 1.3 connection is negotiated, any field in your tls.Config. will be ignored.
- How to redirect http request to https:

```
main.go

// redirect every http request to https
	go http.ListenAndServe(":4000", http.HandlerFunc(httpRedirect))

	srv := &http.Server{
		Addr:      cfg.Addr,
		ErrorLog:  app.ErrorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	infoLog.Printf("Starting server on %s\n", cfg.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

handlers.go

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		"https://"+r.Host+r.URL.String(),
		http.StatusMovedPermanently)
}

```

  
  