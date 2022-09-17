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

## Connection timeouts

- IdleTimeout, ReadTimeout, WriteTimeout are server-wide settings which act on the underlying connection and apply to all requests irrespective of their handler or URL.

### IdleTimeout

- By default, keep-alive connections will be automatically closed after a couple of minutes (with the exact time depending on your [operating system](https://github.com/golang/go/issues/23459#issuecomment-374777402))
- [Keep-Alive Process](http://coryklein.com/tcp/2015/11/25/custom-configuration-of-tcp-socket-keep-alive-timeouts.html)

> There are three configurable properties that determine how Keep-Alives work. On Linux they are1:

> tcp_keepalive_time (default 7200 seconds)
>
> tcp_keepalive_probes (default 9 in Linux)
> 
> tcp_keepalive_intvl (default 75 seconds)
> 
> The process works like this:

> > Client opens TCP connection
> 
> > If the connection is silent for tcp_keepalive_time seconds, send a single empty ACK packet.1
> 
> > Did the server respond with a corresponding ACK of its own?
> 
> > No
> 
> > Wait tcp_keepalive_intvl seconds, then send another ACK
> 
> > Repeat until the number of ACK probes that have been sent equals tcp_keepalive_probes.
> 
> > If no response has been received at this point, send a RST and terminate the connection.
> 
> > Yes: Return to step 2
> 

### ReadTimeout

- In our code we’ve also set the ReadTimeout setting to 5 seconds. This means that if the
  request headers or body are still being read 5 seconds after the request is first accepted, then
  Go will close the underlying connection. Because this is a ‘hard’ closure on the connection,
  the user won’t receive any HTTP(S) response.
- If you set ReadTimeout but don’t set IdleTimeout, then IdleTimeout will
  default to using the same setting as ReadTimeout. For instance, if you set ReadTimeout
  to 3 seconds, then there is the side-effect that all keep-alive connections will also be
  closed after 3 seconds of inactivity. Generally, my recommendation is to avoid any
  ambiguity and always set an explicit IdleTimeout value for your server.

### WriteTimeout

- The WriteTimeout setting will close the underlying connection if our server attempts to write
  to the connection after a given period (in our code, 10 seconds). But this behaves slightly
  differently depending on the protocol being used.
- For HTTP connections, if some data is written to the connection more than 10 seconds
  after the read of the request header finished, Go will close the underlying connection
  instead of writing the data.
- For HTTPS connections, if some data is written to the connection more than 10 seconds
  after the request is first accepted, Go will close the underlying connection instead of
  writing the data. This means that if you’re using HTTPS (like we are) it’s sensible to set
  WriteTimeout to a value greater than ReadTimeout.
- It’s important to bear in mind that writes made by a handler are buffered and written to the
  connection as one when the handler returns. Therefore, the idea of WriteTimeout is generally
  not to prevent long-running handlers, but to prevent the data that the handler returns from
  taking too long to write.

See also in the book **ReadHeaderTimeout** and **MaxHeaderBytes** settings