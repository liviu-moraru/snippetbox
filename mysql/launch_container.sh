#!/usr/bin/env bash

docker run -e MYSQL_ROOT_PASSWORD=my-passw --name mysql -p 3306:3306 -v mysql:/var/lib/mysql mysql