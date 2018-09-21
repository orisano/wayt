# wayt
 wayt is a command utility of a wait.

## Installation
```bash
go get github.com/orisano/wayt
```

## How to use
```
$ wayt
wayt: subcommand is required:
Available SubCommands:
 - tcp
 - sql
 - http
 - file
 - sh

$ wayt -h
Usage of wayt:
  -i duration
    	interval (default 1s)
  -t duration
    	timeout (default 5m0s)
```

## Author
Nao YONASHIRO (@orisano)

## License
MIT
