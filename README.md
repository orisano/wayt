# wayt
wayt is a command utility of a wait.

## Installation
```bash
go get github.com/orisano/wayt
```

## How to use
### wayt
```
$ wayt
wayt: subcommand is required:
Available SubCommands:
 - tcp
 - sql
 - http
 - file
 - sh
 - grpc

$ wayt -h
Usage of wayt:
  -i duration
    	interval (default 1s)
  -t duration
    	timeout (default 5m0s)
  -x	execute command
```

### wayt tcp
```
$ wayt tcp -h
Usage of tcp:
  -a string
    	target address (required)
```

### wayt sql
```
$ wayt sql -h
Usage of sql:
  -d string
    	driver (default "mysql")
  -dsn string
    	data source name (required)
  -env string
    	 (default "DB_URL")
  -q string
    	query (default "SELECT 1;")
  -url string
    	url
```

### wayt http
```
$ wayt http -h
Usage of http:
  -m string
    	method (default "GET")
  -u string
    	url (required)
```

### wayt file
```
$ wayt file -h
Usage of file:
  -p string
    	path (required)
```

### wayt sh
```
$ wayt sh -h
Usage of sh:
  -c string
    	command (required)
```

### wayt grpc
```
$ wayt grpc -h
Usage of grpc:
  -addr string
    	address (required)
  -service string
    	service name to check
  -tls
    	use TLS
  -tls-ca-cert string
    	trusted certificates for verifying server
  -tls-client-cert string

  -tls-client-key string

  -tls-no-verify
    	do not verify the certificate
  -tls-server-name string

```

## Author
Nao Yonashiro (@orisano)

## License
MIT
