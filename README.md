# 4.2 Installing a database driver

```shell
go get github.com/go-sql-driver/mysql@v1
# Remove package ( from module and from installed packages $GOPATH/pkg/mod)
go get github.com/go-sql-driver/mysql@none
```

# 4.3 Modules and reproducible builds

```shell
# Verify if the checksum s of the downloaded packages mathch the entries in go.sum
go mod verify
# Download the exact versions of all the packages in the project
go mod download
# Upgrade packages to the latest version
go get -u github.com/foo/bar
# Or alternatively, if you want to upgrade to a specific version
go get -u github.com/foo/bar@v2.0.0
# Removing unused packages
go get github.com/foo/bar@none
# go mod tidy will automatically remove any unused packages from your go.mod and go.sum files
# go mod tidy doesn't remove the modules from $GOPATH/bin/mod
go mod tidy -v # -v flag causes tidy to print information about removed modules to standard error.
```

# 4.6 Executing SQL statements

To test create request:

```shell
curl -iL -X POST http://localhost:4000/snippet/create

# Test inside container
# With container mysql running 
docker run -e MYSQL_ROOT_PASSWORD=my-passw --name mysql -p 3306:3306 -v mysql:/var/lib/mysql mysql
( or docker start mysql)

docker exec -it mysql mysql -uroot -p
# Insert password for root (my-passw)
# Inside mysql client REPL
use snippetbox; select * from snippets;
```

To test view request:

```shell
docker start mysql
curl -iL "http://127.0.0.1:4000/snippet/view?id=1"
```

# 4.9 Transactions and other details. Transaction and DB tests.

```shell
go test -v ./internal/models
# or
go test -v ./...
```

# 5.1 Displaying dynamic data

1. DB NULL values in templates

```
<strong>{{.Title.Value}}</strong>
or
 <strong>{{if .Newcol.Valid}}{{.Title.String}}{{else}}-{{end}}</strong>

```

2. Time fields

[Time package](https://pkg.go.dev/time#LoadLocation)

Create an other Time for a time zone

```go
t := time.Now()
tz, _ := time.LoadLocation("America/Toronto")
t = t.In(tz)
```

Ex. of time zone: See the file `$GOROOT/lib/time/zoneinfo.zip`

**Format date/time**
See [Format](https://pkg.go.dev/time#Time.Format)

In `format.go` (package time), see the constants:
```
const (
Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
ANSIC       = "Mon Jan _2 15:04:05 2006"
UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
RFC822      = "02 Jan 06 15:04 MST"
RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
RFC3339     = "2006-01-02T15:04:05Z07:00"
RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
Kitchen     = "3:04PM"
// Handy time stamps.
Stamp      = "Jan _2 15:04:05"
StampMilli = "Jan _2 15:04:05.000"
StampMicro = "Jan _2 15:04:05.000000"
StampNano  = "Jan _2 15:04:05.000000000"
)
```

Ex in templates:
```
 <time>Created: {{.Created.Format "02 Jan 06 15:04 -0700"}}</time>
```