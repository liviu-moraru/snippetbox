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