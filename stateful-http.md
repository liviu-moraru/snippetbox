# Chapter 9. Stateful HTTP

# 9.1 Choosing a session manager

1. Security considerations when working with sessions: [Session Management Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html)
2. Recommended packages: **gorilla/sessions** or **alexedwards/scs**
3. gorilla/sessions is the most established and well-known session management package
   for Go. It has a simple and easy-to-use API, and let’s you store session data client-side (in
   signed and encrypted cookies) or server-side (in a database like MySQL, PostgreSQL or
   Redis).
   However — importantly — it doesn’t provide a mechanism to renew session IDs (which is
   necessary to reduce risks associated with session fixation attacks if you’re using one of the
   server-side session stores).
4. alexedwards/scs let’s you store session data server-side only. It supports automatic
   loading and saving of session data via middleware, has a nice interface for type-safe
   manipulation of data, and does allow renewal of session IDs. Like gorilla/sessions, it
   also supports a variety of databases (including MySQL, PostgreSQL and Redis).
5. In summary, if you want to store session data client-side in a cookie then gorilla/sessions is a good choice, but otherwise alexedwards/scs is generally the better option due to the ability to renew session IDs.
6. In client-side sessions, the information is stored in the cookie. The advantage is that the server is entirely stateless. The disadvantage is that the user can see that data [Client Side Session vs Server Side Session](https://medium.com/@tiff.sage/client-side-session-vs-server-side-session-d506f5408e8c#:~:text=Server%2Dside%20sessions%20are%20mostly,use%20and%20smaller%20data%20size.)
7. Sesssion fixation attacks [Renew the Session ID After Any Privilege Level Change¶](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#renew-the-session-id-after-any-privilege-level-change)
8. Install packages:

```
go get github.com/alexedwards/scs/v2@v2
go get github.com/alexedwards/scs/mysqlstore
```

# 9.2 Setting up the session manager

1. The documentation for alexedwards/scs package: [documentation](https://github.com/alexedwards/scs) and [API reference](https://pkg.go.dev/github.com/alexedwards/scs/v2)
2. Create the sessions table in MySQL

```
CREATE TABLE sessions (
token CHAR(43) PRIMARY KEY,
data BLOB NOT NULL,
expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

- The token field will contain a unique, randomly-generated, identifier for each session.
- The data field will contain the actual session data that you want to share between HTTP
  requests. This is stored as binary data in a BLOB (binary large object) type.
- The expiry field will contain an expiry time for the session. The scs package will
  automatically delete expired sessions from the sessions table so that it doesn’t grow too
  large.

3. Establish a session manager and make it available through dependency injection (application struct)

```
main.go

import (
...
"database/sql"
_ "github.com/go-sql-driver/mysql" 
"github.com/alexedwards/scs/mysqlstore"
"github.com/alexedwards/scs/v2"
)
type Application struct {
   ....
   sessionManager *scs.SessionManager
}

db, err := openDB("web:pass@/snippetbox?parseTime=true")
if err != nil {
errorLog.Fatal(err)
}

defer db.Close()

sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
app := &Application{
    ...
    sessionManager: sessionManager,
}


func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
```

- [Other fields](https://pkg.go.dev/github.com/alexedwards/scs/v2#SessionManager) that can be configured in SessionManager

4. Wrap our application routes with the middleware provided by the SessionManager.LoadAndSave() method. This middleware automatically
   loads and saves session data with every HTTP request and response.
5. It doesn't make sense to add the session middleware to some routes (e.g. static route in the project). So, we must wrap only some middlewares with the session middleware.

```
routes.go
func (app *Application) routes() http.Handler {
...
    // Create a new middleware chain containing the middleware specific to our
    // dynamic application routes. For now, this chain will only contain the
    // LoadAndSave session middleware but we'll add more to it later.
    dynamic := alice.New(app.sessionManager.LoadAndSave)

    router.Handler(http.MethodGet, "/", dynamic.Then(app.HomeHandler()))
    router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.Then(app.SnippetViewHandler()))
    router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
    router.Handler(http.MethodPost, "/snippet/create", dynamic.Then(app.SnippetCreatePostHandler()))

    standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
    return standard.Then(router)
}
```
# 9.3 Working with session data

## Behind-the-scenes of session management

- The session manager set a session cookie to the browser.
- The session cookie contains the **session token** (session ID). This is a high-entropy random string.
- The Expires and MaxAge for the cookie that are sent back with the response are set from the value of the sessionManager.Lifetime
- By default, the sessionManager created by the scs.New() function, creates a Persistent cookie.
- The name of the cookie is by default session.
See [implementetion](https://github.com/alexedwards/scs/blob/12165213346f24f731a5b8f0e064dc41fea05935/session.go#L102)
```
// New returns a new session manager with the default options. It is safe for
// concurrent use.
func New() *SessionManager {
	s := &SessionManager{
		IdleTimeout: 0,
		Lifetime:    24 * time.Hour,
		Store:       memstore.New(),
		Codec:       GobCodec{},
		ErrorFunc:   defaultErrorFunc,
		contextKey:  generateContextKey(),
		Cookie: SessionCookie{
			Name:     "session",
			Domain:   "",
			HttpOnly: true,
			Path:     "/",
			Persist:  true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
	}
	return s
}

```

- We can make some changes: 

```
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)
sessionManager.Lifetime = 1 * time.Hour
cookie := &sessionManager.Cookie
cookie.Name = "mySecondSession"
cookie.Persist = false
```

- It’s important to emphasize that the session token is just a random string. In itself, it does’t
  carry or convey any session data (like the flash message that we set in this chapter).
- The data field in the database actuallu contains the session data. This field is a MySQL BLOB (binary large object) containing a **gob-encoded** representation of the session data.
- The [gob package](https://pkg.go.dev/encoding/gob) is from the standard library

> Package gob manages streams of gobs - binary values exchanged between an Encoder (transmitter) and a Decoder (receiver). 
> A typical use is transporting arguments and results of remote procedure calls (RPCs) such as those provided by package "net/rpc".
> 
> The implementation compiles a custom codec for each data type in the stream and is most efficient when a single Encoder is used to transmit a stream of values, amortizing the cost of compilation.

- Lastly, the final column in the database is the **expiry** time, after wich the session will no longer be considered valid.

> So, what happens in our application is that the **LoadAndSave** middleware checks each incoming request for a session cookie. 
> If a session cookie is present, it reads the session token and retrieves the corresponding session data from the database (while also checking that the
session hasn’t expired).
> It then adds the session data to the request context so it can be used
in your handlers.
> Any changes that you make to the session data in your handlers are updated in the request context, and then the LoadAndSave middleware updates the database with any changes to the session data before it returns.

- So, LoadAndSave middleware process the request before other middleware and the response after these middleware finished handling the request.
- When processing the response, LoadAndSave commits the changes made to the session data in the request context by other middlewares into the store.

