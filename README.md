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
./httpcommander config-examples/test.conf
```

Run with TLS
-------------------------
```
./httpcommander config-examples/test-tls.conf
```

Run with x509 client auth
-------------------------

You can use the httpcommander with x509 client auth.
Client auth is automatically enabled as soon as you specify a ca-file to validate the client certs. Simply adjust the config accordingly.

To test this:
```
./httpcommander config-examples/test-client-auth.conf
```
and then you can do
```
curl -k --key x509/client-tester.key \
        --cert x509/client-tester.crt \
        https://localhost:8989/cmd/echotest
```

To create a CA, you can use https://github.com/r3ek0/cluster-ca

Build standalone binary for docker image
-------------------------------------------
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' .
docker build -t httpcommander:myversion .
docker run --rm --name httpcommander \
    -v $(pwd):/conf \
    httpcommander:myversion /conf/test.conf
```
