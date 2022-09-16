# Create certificate from [Let's Encrypt](https://letsencrypt.org/) for web server on Raspberry PI

- Install [Acme client from Let's Encrypt](https://github.com/go-acme/lego). Documentation [here](https://go-acme.github.io/lego/): 
```
# On Raspberry PI
go install github.com/go-acme/lego/v4/cmd/lego@latest
```
- Generate certificate using the built-in web server [see](https://go-acme.github.io/lego/usage/cli/obtain-a-certificate/)
```
# On Raspberry PI
# The port 4000 should not be bound
# The router is set up to forward the port 80 to the Raspberry PI port 4000
cd ~
mkdir lego
cd lego
lego --email="liviu@lm58.tplinkdns.com" --http.port ":4000" --domains="lm58.tplinkdns.com" --http run
```
- The generated files are in the .lego/certificates subdirectory. 
- The ***lm58.tplinkdns.com.crt*** is the server certificate (including the CA certificate)
- The ***tplinkdns.com.key*** is the private key needed for the server certificate.
- Transfer the files to MAC project

```
# On MAC in the tls subdirectory of the project
scp liviu@192.168.68.57:/home/liviu/lego/.lego/certificates/lm58.tplinkdns.com.key key.pem
scp liviu@192.168.68.57:/home/liviu/lego/.lego/certificates/lm58.tplinkdns.com.crt cert.pem
```
