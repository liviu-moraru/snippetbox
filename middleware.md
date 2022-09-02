# Chapter 6. Middleware

# 6.3 Request logging

- How to log response code

See: [Logging the status code of a HTTP Handler in Go](https://dev.to/julienp/logging-the-status-code-of-a-http-handler-in-go-25aa)

```
// cmd/web/middleware.go

...
type StatusRecorder struct {
    http.ResponseWriter
    Status int
}

//Override the WriteHeader method of the embedded ResponseWriter
func (sr *StatusRecorder) WriteHeader(status int) {
    sr.Status = status
  
    // Without this, the Status Code of the response would not be set.
    sr.ResponseWriter.WriteHeader(status)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
sr := &StatusRecorder{
ResponseWriter: w,
Status:         200,
}

rd := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
next.ServeHTTP(sr, r) // the sr.Status will be set by the WriteHeader method
app.InfoLog.Printf("%s Response status: %d", rd, sr.Status)
})
}
```

```
// cmd/web/routes..go

func (app *Application) routes() http.Handler {
    mux := http.NewServeMux()
	....
    return app.logRequest(secureHeader(mux))
}
```

```
// cnd/web/main.go
...
srv := &http.Server{
    Addr:     cfg.Addr,
    ErrorLog: app.ErrorLog,
    Handler:  app.routes(),
}

err = srv.ListenAndServe()
errorLog.Fatal(err)

```

# 6.4 Panic recovery

1. Setting the **Connection: Close** header on the response acts as a trigger to make Go’ HTTP server automatically close the current connection after a response has been sent.I talso informs the user that the connection will be closed. Note: If the protocol being used is HTTP/2, Go will automatically strip the **Connection: Close** header from the response (so it is not malformed) and send a GOAWAY frame.
2. Go’s HTTP server assumes that the effect of any panic is isolated to the goroutine serving the active HTTP request. Specifically, following a panic our server will log a stack trace to the server error log, unwind
   the stack for the affected goroutine (calling any deferred functions along the way) and close
   the underlying HTTP connection.
3. A panic in a function called as a goroutine by a handler, will end the application (web server will crash)

```
func panicTest() {
	panic("Panic")
}

func (app *Application) SnippetViewHandler() http.Handler { 
...
go panicTest()
}
```

# 6.5 Composable middleware chains

Install alice package:

```
go get github.com/justinas/alice@v1
```