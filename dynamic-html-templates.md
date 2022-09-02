# Chapter 5. Dynamic HTML templates

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

```
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

3. Jetbrains Goland: Associate *.tmpl files with Go templates

Preferences -> Editor -> File Types

Then for Go template, add *.tmpl in the list of associated file types.

# 5.2 Template actions and functions

List of template functions: [Template functions](https://pkg.go.dev/text/template#hdr-Functions)
