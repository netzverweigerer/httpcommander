# httpcommander
Simple program to execute predefined commands via http calls
It uses https://github.com/gin-gonic/gin 

Install
-------

Install go

```
export GOPATH=<path_to_go_dev>
go get github.com/r3ek0/httpcommander...
go build github.com/r3ek0/httpcommander
```

Run
---
```
./httpcommander httpcommander.conf
```

Run with x509 client auth
-------------------------

You can use the httpcommander with x509 client auth.
To create the CA, you can use https://github.com/r3ek0/cluster-ca

Client auth is automatically enabled as soon as you specify a ca-file to validate the client certs. Simply adjust the config accordingly.
