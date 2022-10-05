#!/usr/bin/env bash

docker exec -it mysql mysql -uroot -pmy-passw --verbose  -e "CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
docker exec -it mysql mysql -uroot -pmy-passw --verbose  -e "use snippetbox;"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "CREATE TABLE snippets ( id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, title VARCHAR(100) NOT NULL, content TEXT NOT NULL, created DATETIME NOT NULL, expires DATETIME NOT NULL );"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "CREATE INDEX idx_snippets_created ON snippets(created);"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "INSERT INTO snippets (title, content, created, expires) VALUES ( 'An old silent pond', 'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY) );"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "INSERT INTO snippets (title, content, created, expires) VALUES ( 'Over the wintry forest', 'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY) );"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "INSERT INTO snippets (title, content, created, expires) VALUES ( 'First autumn morning', 'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo', UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY) );"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "CREATE USER 'web'@'%';GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'%';ALTER USER 'web'@'%' IDENTIFIED BY 'pass';"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "CREATE TABLE sessions (token CHAR(43) PRIMARY KEY,data BLOB NOT NULL,expiry TIMESTAMP(6) NOT NULL);"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "CREATE INDEX sessions_expiry_idx ON sessions (expiry);"
docker exec -it mysql mysql -uroot -pmy-passw  snippetbox -e "CREATE TABLE users ( id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255) NOT NULL, email VARCHAR(255) NOT NULL, hashed_password CHAR(60) NOT NULL, created DATETIME NOT NULL ); ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);"
