# Setup MYSQL (MariaDB) Database on Kali Linux Raspberry PI

- Kali Linux comes with MariaDB preinstalled
- Enable MariaDB on Kali Linux

```
sudo systemctl enable --now mariadb
```
- [Secure MariaDB server](https://computingforgeeks.com/how-to-install-mariadb-on-kali-linux/)

```
sudo mysql_secure_installation

Questions and answers:
Enter current password for root (enter for none):
    Answer: none
Switch to unix_socket authentication [Y/n]
    Answer: n
Change the root password? [Y/n]
    Answer: y
New password: 
    Answer: my-passw
Remove anonymous users?
    Answer: y
Disallow root login remotely?
    Answer: y
Remove test database and access to it?
    Answer: y
Reload privilege tables now?
    Answer: y
```

- Connect to MySQL server

```
mysql -u root -p
Enter password: (my-passw)
```

- Run the sql commands from file mysql/initialize_data.sh

- [Opening a port on Linux to allow TCP connections](https://www.digitalocean.com/community/tutorials/opening-a-port-on-linux) (e.g. port 4000)

```
iptables -A INPUT -p tcp --dport 4000 -j ACCEPT
```

